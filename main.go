package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/panjf2000/gnet/v2"
)

// chatServer 实现 gnet 的事件处理逻辑
type chatServer struct {
	gnet.BuiltinEventEngine                      // 嵌入默认事件引擎（可选）
	conns                   map[gnet.Conn]string // 所有已连接客户端映射：连接对象 -> 用户名（默认是远程地址）
	mu                      sync.RWMutex         // 并发读写保护 conns
}

// OnBoot 服务器启动初始化
func (cs *chatServer) OnBoot(eng gnet.Engine) gnet.Action {
	log.Println("🚀 Chatroom server is running.")
	cs.conns = make(map[gnet.Conn]string)
	return gnet.None
}

// OnOpen 有新连接建立时触发
func (cs *chatServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	cs.mu.Lock()
	name := c.RemoteAddr().String()
	cs.conns[c] = name
	cs.mu.Unlock()

	fmt.Printf("✅ New connection: %s (Total: %d)\n", name, len(cs.conns))

	// 向其他客户端广播加入信息
	msg := fmt.Sprintf("📢 %s has entered the chatroom.\n", name)
	cs.broadcast([]byte(msg), c)

	// 返回欢迎语
	return []byte("🎉 Welcome to the chatroom.\n"), gnet.None
}

// OnClose 某连接关闭时触发
func (cs *chatServer) OnClose(c gnet.Conn, err error) gnet.Action {
	cs.mu.Lock()
	name := cs.conns[c]
	delete(cs.conns, c)
	cs.mu.Unlock()

	msg := fmt.Sprintf("👋 %s has left the chatroom.\n", name)
	cs.broadcast([]byte(msg), c)
	return gnet.None
}

// OnTraffic 收到客户端数据时触发
func (cs *chatServer) OnTraffic(c gnet.Conn) gnet.Action {
	// 读取当前连接的全部数据
	data, _ := c.Next(-1)
	if len(data) == 0 {
		return gnet.None // 忽略空消息
	}

	// 获取当前发送者昵称
	cs.mu.RLock()
	name := cs.conns[c]
	cs.mu.RUnlock()

	// 格式化广播消息
	msg := fmt.Sprintf("💬 %s: %s\n", name, data)

	// 向其他客户端广播消息
	cs.broadcast([]byte(msg), c)

	// 服务端打印日志
	fmt.Printf("📨 Received from %s: %s\n", name, data)
	return gnet.None
}

// broadcast 广播消息给所有客户端（排除发送者）
func (cs *chatServer) broadcast(msg []byte, sender gnet.Conn) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	for conn := range cs.conns {
		if conn == sender {
			continue // 跳过发送者
		}
		// 复制一份 buffer 避免数据覆盖
		buf := append([]byte(nil), msg...)
		_ = conn.AsyncWrite(buf, nil)
	}
}

// main 启动服务器入口
func main() {
	server := &chatServer{}
	log.Fatal(gnet.Run(server, "tcp://:9000", gnet.WithMulticore(true)))
}
