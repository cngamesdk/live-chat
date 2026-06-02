# 在线客服系统 - 对接文档

## 一、系统架构

```
┌──────────────┐        ┌───────────────────┐        ┌──────────┐
│  H5 客户端    │ ◄────► │  live-chat API    │ ◄────► │  MySQL   │
│  (Vue 3)     │  HTTP  │  (Gin :8890)      │        │  + Redis │
└──────────────┘  + WS  └───────────────────┘        └──────────┘
                                          ▲
┌──────────────┐        ┌───────────────────┐
│  管理后台     │ ◄────► │  background 服务   │
│  (Vue 3)     │  HTTP  │  (Gin :8888)      │
└──────────────┘        └───────────────────┘
```

- **live-chat API** (端口 8890)：处理 H5 客户端的聊天请求、WebSocket 连接
- **background 服务** (端口 8888)：管理后台，产品配置、客服管理、数据报表
- 两个服务共享同一 MySQL 数据库和 Redis

---

## 二、H5 集成指南

### 2.1 URL 参数接入

H5 聊天页面通过 URL 参数初始化：

```
https://your-domain.com/chat?product_code=DEMO&user_id=user123&user_name=张三&user_token=xxx
```

| 参数 | 必填 | 说明 |
|------|------|------|
| `product_code` | 是 | 产品编码，需先在管理后台创建 |
| `user_id` | 是 | 用户唯一标识 |
| `user_name` | 否 | 用户展示名称，不传则显示 user_id |
| `user_token` | 否 | 用户身份签名（见 2.2） |

### 2.2 用户身份签名（可选）

为防止伪造请求，可使用 HMAC-SHA256 签名：

```javascript
// 生成签名
const crypto = require('crypto')

function signUserToken(userId, productCode, secret) {
  const ts = Math.floor(Date.now() / 1000 / 300) // 5分钟窗口
  const data = `${userId}:${productCode}:${ts}`
  return crypto.createHmac('sha256', secret).update(data).digest('hex')
}
```

`secret` 为产品密钥，需与管理后台中配置的一致（待扩展字段）。

### 2.3 iframe 嵌入

```html
<iframe
  src="https://your-domain.com/chat?product_code=DEMO&user_id=user123&user_name=张三"
  width="100%"
  height="600px"
  frameborder="0"
></iframe>
```

建议宽度≥360px，高度≥500px。

---

## 三、API 接口文档

### 基础 URL

- 开发环境：`http://127.0.0.1:8890/api/v1`
- 生产环境：`https://api.example.com/api/v1`

### 3.1 获取产品信息

```
GET /api/v1/product/:code/info
```

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "product_code": "DEMO",
    "name": "示例游戏",
    "logo": "https://...",
    "welcome_title": "欢迎咨询",
    "welcome_message": "您好！请问有什么可以帮助您的？"
  }
}
```

### 3.2 初始化会话

```
POST /api/v1/chat/init
Content-Type: application/json

{
  "product_code": "DEMO",
  "user_id": "user123",
  "user_name": "张三",
  "user_token": "",
  "source": "h5"
}
```

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "session_id": 1,
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": "user123"
  }
}
```

### 3.3 获取聊天历史

```
GET /api/v1/chat/history?page=1&pageSize=50
Header: X-Session-Token: <token>
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "session_id": 1,
        "sender_type": "system",
        "sender_name": "系统",
        "content": "您好，欢迎咨询！",
        "msg_type": "system",
        "created_at": "2026-06-01T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 50
  }
}
```

### 3.4 查询 FAQ

```
GET /api/v1/chat/faq?keyword=充值&page=1&pageSize=10
Header: X-Session-Token: <token>
```

### 3.5 上传文件

```
POST /api/v1/upload
Content-Type: multipart/form-data
Header: X-Session-Token: <token>

file: <binary>
```

**支持类型**：image/jpeg, image/png, image/gif, image/webp, video/mp4
**大小限制**：10MB（可配置）

**响应**：
```json
{
  "code": 0,
  "data": {
    "url": "/uploads/image/1717200000000000001.jpg",
    "file_type": "image",
    "file_name": "photo.jpg",
    "file_size": 102400
  }
}
```

---

## 四、WebSocket 协议

### 4.1 连接

```
ws://127.0.0.1:8890/ws/chat?token=<session_token>
```

### 4.2 消息格式

#### 客户端 → 服务端

```json
// 发送文本消息
{"type": "message", "content": "如何充值？", "msg_type": "text"}

// 发送图片
{"type": "message", "content": "[图片]", "msg_type": "image", "attachment_url": "/uploads/image/xxx.jpg"}

// 发送视频
{"type": "message", "content": "[视频]", "msg_type": "video", "attachment_url": "/uploads/video/xxx.mp4"}

// 心跳
{"type": "ping"}
```

#### 服务端 → 客户端

```json
// 普通消息
{"type": "message", "content": "请点击充值按钮...", "sender_type": "system", "sender_name": "智能客服", "msg_type": "text", "is_faq_reply": true}

// 排队通知
{"type": "queue", "queue_position": 3}

// 客服接入
{"type": "connected", "agent_name": "客服001"}

// 会话关闭
{"type": "closed"}

// 心跳响应
{"type": "pong"}
```

### 4.3 消息流程

1. 客户端通过 `POST /chat/init` 获取 `session_token`
2. 使用 `session_token` 建立 WebSocket 连接
3. 发送消息 → 服务端进行 FAQ 匹配
4. 匹配成功 → 自动回复 FAQ 答案
5. 匹配失败 → 分配人工客服 / 加入排队队列
6. 客服接入 → 收到 `connected` 消息

---

## 五、FAQ 导入格式

批量导入 FAQ 的 JSON 格式：

```json
{
  "product_id": 1,
  "items": [
    {
      "category": "充值问题",
      "question": "如何充值？",
      "answer": "请点击游戏内充值按钮，选择支付方式进行充值。",
      "keywords": "充值,支付,购买,氪金",
      "priority": 10
    },
    {
      "category": "账号问题",
      "question": "忘记密码怎么办？",
      "answer": "请点击登录页的'忘记密码'，通过手机号或邮箱找回。",
      "keywords": "密码,忘记,找回,登录",
      "priority": 8
    }
  ]
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| product_id | int | 产品 ID |
| category | string | 分类名称 |
| question | string | 问题标题 |
| answer | string | 答案内容（支持纯文本） |
| keywords | string | 搜索关键词，逗号分隔 |
| priority | int | 优先级，越大越先匹配，默认 0 |

**导入 API**：`POST /liveChat/faq/import`

---

## 六、配置说明

### 6.1 API 服务配置 (`config.yaml`)

```yaml
server:
  port: 8890              # 监听端口
  mode: debug             # debug / release

mysql:
  path: 192.168.60.219    # MySQL 地址
  port: 3306
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: live_chat      # 数据库名
  username: root
  password: ""
  max-idle-conns: 10
  max-open-conns: 100

redis:
  db: 0
  addr: 127.0.0.1:6379
  password: ""

chat:
  faq_match_threshold: 0.6     # FAQ 匹配阈值 (0.0-1.0)，低于此值转人工
  max_upload_size: 10485760    # 上传文件最大字节数 (10MB)
  allowed_upload_types:        # 允许的上传文件类型
    - image/jpeg
    - image/png
    - image/gif
    - image/webp
    - video/mp4
  session_timeout: 1800        # 会话超时自动关闭（秒），默认 30 分钟
  agent_assignment: round_robin # 客服分配策略: round_robin / least_loaded

upload:
  path: ./uploads              # 上传文件存储路径
  max_size: 10485760
```

### 6.2 FAQ 匹配引擎

- **匹配方式**：关键词 + 问题文本分词匹配
- **评分机制**：精确匹配=1.0，包含匹配=0.9-0.95，分词匹配=命中词数/总词数
- **阈值**：`faq_match_threshold`（默认 0.6），低于此值转人工客服
- **优先匹配**：高优先级 FAQ 优先

---

## 七、部署指南

### 7.1 后端部署

```bash
# 1. 编译 API 服务
cd live-chat/api
go mod tidy
go build -o live-chat-api .

# 2. 配置 config.yaml（修改数据库连接等）

# 3. 创建 uploads 目录
mkdir -p uploads/image uploads/video

# 4. 启动
./live-chat-api
```

### 7.2 H5 部署

```bash
# 1. 安装依赖
cd live-chat/h5
npm install

# 2. 开发运行
npm run dev

# 3. 构建生产版本
npm run build
# 将 dist/ 部署到 Nginx 或 CDN
```

### 7.3 管理后台

管理后台复用上级目录的 `background` 项目，新增页面:
- 产品管理：配置产品 Logo、名称、欢迎语
- FAQ 问题库：管理自动回复知识库
- 客服管理：客服上下线、并发配置
- 会话管理：实时查看、接单、回复
- 数据报表：客服数据分析

### 7.4 Nginx 配置示例

```nginx
# H5 前端
server {
    listen 80;
    server_name chat.example.com;
    root /var/www/live-chat/dist;
    index index.html;
    location / { try_files $uri $uri/ /index.html; }
}

# API 反向代理
server {
    listen 80;
    server_name api.example.com;
    location /api { proxy_pass http://127.0.0.1:8890; }
    location /ws {
        proxy_pass http://127.0.0.1:8890;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    location /uploads { proxy_pass http://127.0.0.1:8890; }
}
```

---

## 八、扩展开发

### 8.1 消息类型扩展

添加新消息类型只需实现客户端的渲染逻辑：

```javascript
// H5 端新增卡片消息渲染
// views/chat/index.vue 的 handleWsMessage 添加:
case 'card':
  messages.value.push({
    id: Date.now(),
    msgType: 'card',
    senderType: 'system',
    content: msg.data
  })
```

### 8.2 客服分配策略

当前支持 `round_robin`（轮询）和 `least_loaded`（最少负载），扩展方式：

```go
// service/agent.go - 在 AssignAgent 方法中添加新策略
case "skill_based":
    // 根据用户问题关键词匹配有相应技能的客服
```

### 8.3 Webhook 事件

可在 `service/chat.go` 的会话生命周期中添加 Webhook 回调：

- 会话创建 → 通知业务系统
- FAQ 自动回复 → 记录到业务数据库
- 客服接入 → 记录客服服务数据
- 会话关闭 → 同步满意度评价

---

## 九、数据库表说明

| 表名 | 说明 |
|------|------|
| `lc_product` | 产品配置（Logo、名称、欢迎语等） |
| `lc_faq` | FAQ 知识库 |
| `lc_chat_session` | 聊天会话记录 |
| `lc_chat_message` | 聊天消息记录 |
| `lc_agent` | 客服人员配置 |
| `lc_upload` | 文件上传记录 |
| `lc_daily_report` | 每日统计汇总 |

---

## 十、性能指标

- WebSocket 单机并发连接：≥5000
- 消息持久化：异步写入，不阻塞 WebSocket 处理
- FAQ 匹配：内存级关键词搜索，<1ms 响应
- 文件上传：支持断点续传（待实现）
- 数据库：索引覆盖所有常用查询字段

---

## 十一、常见问题

**Q: H5 页面白屏怎么办？**
A: 检查 `product_code` 是否在管理后台配置，以及 API 服务是否正常启动。

**Q: WebSocket 断开后消息会丢失吗？**
A: 所有消息都会持久化到数据库，重新连接后可通过 `/chat/history` 获取历史消息。

**Q: FAQ 匹配不准确怎么办？**
A: 调整 `config.yaml` 中的 `faq_match_threshold`，越低越容易匹配但可能不准确；增加 FAQ 的关键词和优先级。

**Q: 如何添加新的客服人员？**
A: 在 `sys_users` 中创建用户，然后在客服管理页面选择产品和人员上线即可。

**Q: 上传文件存在哪里？**
A: 默认存储在服务器的 `uploads/` 目录，可通过 Nginx 反向代理访问。
