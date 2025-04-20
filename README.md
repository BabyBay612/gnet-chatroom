# gnet 聊天室服务端

一个基于gnet(https://github.com/panjf2000/gnet)构建的高性能终端聊天室服务端，采用事件驱动、非阻塞的TCP框架，支持多用户并发通信。

---

## 项目特点

- 基于gnet v2，事件驱动 + 多核并发
- 支持多个客户端同时连接
- 实现消息实时广播（排除发送者）
- 使用`sync.RWMutex`实现线程安全的连接管理
- 结构清晰，代码注释规范，便于扩展和学习

---

## 使用方式

### 1. 安装Go和gnet（v2）

```bash
go install github.com/panjf2000/gnet/v2@latest
```

---

### 2. 启动服务端

```bash
go run main.go
```

默认会监听本机`9000`端口。

---

## 测试方式

### 使用`ncat`

1. 打开两个终端窗口
2. 在每个窗口中输入以下命令连接服务器：

   ```bash
   ncat localhost 9000
   ```

3. 在任意一端输入消息，另一端即可收到实时广播。

---

## 示例输出

```bash
Chatroom server is running.
New connection: 127.0.0.1:51925 (Total: 1)
New connection: 127.0.0.1:51926 (Total: 2)
Received from 127.0.0.1: Hello everyone!
Broadcasting message to 1 client(s)...
```

---

欢迎 Fork、Star、二次开发