# 修复：管理后台发送的消息用户未收到

## 问题原因
管理后台（background）和客服系统（live-chat）是两个独立的服务。管理后台发送消息时只保存到数据库，但没有通过 WebSocket 通知在线用户。

## 解决方案

### 架构设计
```
管理后台客服发送消息
    ↓
保存到数据库 (lc_chat_message)
    ↓
HTTP 调用 live-chat 服务
    ↓
live-chat 通过 WebSocket 广播
    ↓
H5 用户实时收到消息
```

### 代码修改

#### 1. live-chat 服务新增管理后台接口

**新增文件：** `api/api/v1/admin.go`
- 提供 `/api/v1/admin/agent-reply` 接口
- 接收管理后台的消息并通过 WebSocket 广播

**修改文件：**
- `api/api/v1/enter.go` - 注册 AdminApi
- `api/router/chat.go` - 添加管理后台路由
- `api/service/chat.go` - 添加 `SaveAgentMessage` 方法

#### 2. background 服务调用 live-chat

**修改文件：** `background/server/service/live_chat/chat.go`
- `AgentReply` 方法：保存消息后异步调用 live-chat 服务
- `notifyLiveChatService` 方法：发送 HTTP 请求通知 live-chat

**修改文件：** `background/server/config/system.go`
- 添加 `LiveChatURL` 配置字段

## 配置说明

### 在 background 的配置文件中添加

编辑 `background/server/config.yaml`：

```yaml
system:
  # ... 其他配置
  live-chat-url: "http://localhost:8890"  # live-chat 服务地址
```

### API 接口说明

**接口：** `POST /api/v1/admin/agent-reply`

**请求参数：**
```json
{
  "session_id": 123,
  "agent_id": 1,
  "agent_name": "客服小王",
  "content": "您好，有什么可以帮您？",
  "msg_type": "text"
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "message_id": 456
  },
  "msg": "操作成功"
}
```

## 部署步骤

1. **更新 live-chat 服务代码**
   - 部署新的 API 接口
   - 确保服务正常运行

2. **更新 background 服务代码**
   - 部署新的消息发送逻辑
   - 配置 live-chat 服务地址

3. **配置文件**
   - 在 background 的 config.yaml 中添加 `live-chat-url` 配置

4. **重启服务**
   - 先重启 live-chat 服务
   - 再重启 background 服务

## 容错机制

- 如果 live-chat 服务不可用，消息仍会保存到数据库
- 用户刷新页面后可以从历史记录中看到消息
- HTTP 调用采用异步方式，不影响管理后台的响应速度
- 调用失败会记录日志，便于排查问题

## 测试验证

1. 启动 live-chat 服务（端口 8890）
2. 启动 background 服务
3. H5 用户打开聊天窗口
4. 管理后台客服发送消息
5. 验证 H5 用户是否实时收到消息

## 注意事项

- 确保两个服务之间网络可达
- 生产环境需要配置正确的服务地址（不是 localhost）
- 建议配置服务健康检查
- 可以考虑使用消息队列（如 Redis、RabbitMQ）来解耦服务
