# 管理后台消息实时推送配置说明

## 问题
管理后台发送的消息只保存到数据库，但没有通过 WebSocket 实时推送给用户。

## 解决方案
管理后台在发送消息后，会调用 live-chat 服务的 HTTP 接口，由 live-chat 服务通过 WebSocket 广播消息给用户。

## 配置步骤

### 1. 在 background 项目的配置文件中添加 live-chat 服务地址

编辑 `background/server/config.yaml`，在 `system` 部分添加：

```yaml
system:
  # ... 其他配置
  live-chat-url: "http://localhost:8890"  # live-chat 服务的地址
```

### 2. 确保 live-chat 服务正常运行

live-chat 服务需要在 8890 端口运行（或配置文件中指定的端口）。

### 3. 新增的 API 接口

live-chat 服务新增了管理后台专用接口：

- `POST /api/v1/admin/agent-reply` - 接收管理后台的消息并广播

请求参数：
```json
{
  "session_id": 123,
  "agent_id": 1,
  "agent_name": "客服小王",
  "content": "您好，有什么可以帮您？",
  "msg_type": "text"
}
```

## 工作流程

1. 管理后台客服发送消息
2. 消息保存到数据库（lc_chat_message 表）
3. 异步调用 live-chat 服务的 `/api/v1/admin/agent-reply` 接口
4. live-chat 服务通过 WebSocket 广播消息给该会话的所有在线用户
5. H5 用户实时收到消息

## 注意事项

- 如果 live-chat 服务不可用，消息仍会保存到数据库，但不会实时推送
- 用户刷新页面后可以从历史记录中看到消息
- 建议在生产环境中配置正确的 live-chat 服务地址
