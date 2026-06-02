# 优化总结：消息红点和会话关闭提示

## 1. H5消息红点自动消失 ✅

### 问题
用户查看消息后，红点没有自动消失。

### 解决方案

**自动标记已读逻辑：**
- 当收到客服消息时，延迟1秒后自动调用 `markAsRead` API
- 初始化加载历史消息时，自动标记未读的客服消息为已读（延迟1.5秒）
- 标记已读后，更新本地消息状态并刷新未读计数

**修改文件：**
- `h5/src/composables/useChat.js`
  - 在 `handleWsMessage` 中添加自动标记逻辑
  - 在 `init` 方法中加载历史消息后自动标记未读消息
  - 优化 `markAsRead` 方法，标记后刷新未读计数

**工作流程：**
```
用户打开聊天窗口
    ↓
加载历史消息
    ↓
识别未读的客服消息
    ↓
延迟1.5秒后自动标记为已读
    ↓
红点消失，未读计数更新
```

## 2. 会话关闭提示 ✅

### 问题
用户关闭和管理后台关闭会话时，都需要显示提示信息。

### 解决方案

#### A. 用户主动关闭会话

**H5 前端：**
- 在 `ChatHeader` 组件添加关闭按钮
- 点击关闭按钮时弹出确认对话框
- 确认后调用 `closeSession` API
- 显示"您已结束会话，感谢您的咨询"系统消息
- 禁用输入框，防止继续发送消息

**后端 API：**
- 新增 `POST /api/v1/chat/close` 接口
- 用户可以主动关闭自己的会话

**修改文件：**
- `h5/src/components/ChatHeader.vue` - 添加关闭按钮
- `h5/src/views/chat/index.vue` - 添加关闭确认逻辑
- `h5/src/composables/useChat.js` - 添加 `closeSession` 方法
- `h5/src/api/chat.js` - 添加关闭会话 API
- `api/api/v1/chat.go` - 添加 `CloseSession` 处理器
- `api/router/chat.go` - 注册关闭会话路由

#### B. 管理后台关闭会话

**实时通知机制：**
- 管理后台关闭会话时，调用 live-chat 的 `/api/v1/admin/close-session` 接口
- live-chat 服务通过 WebSocket 广播关闭通知
- H5 用户实时收到"客服已关闭会话"系统消息
- 自动禁用输入框

**修改文件：**
- `api/api/v1/session.go` - 新增管理后台关闭会话接口
- `api/api/v1/enter.go` - 注册 SessionApi
- `api/router/chat.go` - 添加管理后台关闭路由
- `background/server/service/live_chat/chat.go` - 修改 `CloseSession` 方法，调用 live-chat 服务

**工作流程：**
```
管理后台点击关闭
    ↓
保存到数据库（status = closed）
    ↓
异步调用 live-chat 的 /api/v1/admin/close-session
    ↓
live-chat 通过 WebSocket 广播 "closed" 消息
    ↓
H5 用户收到系统提示："客服已关闭会话"
    ↓
输入框自动禁用
```

## 3. 新增 API 接口

### H5 用户接口
- `POST /api/v1/chat/close` - 用户关闭会话

### 管理后台接口
- `POST /api/v1/admin/close-session` - 管理后台关闭会话并通知用户

**请求参数：**
```json
{
  "session_id": 123,
  "reason": "客服已关闭会话"  // 可选，自定义关闭原因
}
```

## 4. 用户体验优化

### 消息红点
- ✅ 自动消失，无需手动操作
- ✅ 实时更新未读计数
- ✅ 页面刷新后正确显示未读状态

### 会话关闭
- ✅ 用户可以主动结束会话
- ✅ 关闭前有确认提示，防止误操作
- ✅ 管理后台关闭时用户实时收到通知
- ✅ 会话关闭后自动禁用输入框
- ✅ 显示友好的系统提示消息

## 5. 配置说明

确保 background 的配置文件中已添加 live-chat 服务地址：

```yaml
system:
  live-chat-url: "http://localhost:8890"
```

## 6. 测试验证

### 测试消息红点
1. 管理后台发送消息给用户
2. H5 用户打开聊天窗口
3. 验证消息显示红点
4. 等待1秒后，红点自动消失

### 测试用户关闭
1. H5 用户点击右上角关闭按钮
2. 确认关闭对话框
3. 验证显示"您已结束会话"提示
4. 验证输入框已禁用

### 测试管理后台关闭
1. 管理后台点击关闭会话
2. H5 用户实时收到"客服已关闭会话"提示
3. 验证输入框已禁用
4. 验证无法继续发送消息
