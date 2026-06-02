package model

import (
	"time"
)

type ChatMessage struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionID     int64     `json:"session_id" gorm:"column:session_id;type:bigint;index;not null;comment:会话ID"`
	SenderType    string    `json:"sender_type" gorm:"column:sender_type;type:varchar(16);not null;comment:发送者类型 user/agent/system"`
	SenderID      string    `json:"sender_id" gorm:"column:sender_id;type:varchar(128);comment:发送者标识"`
	SenderName    string    `json:"sender_name" gorm:"column:sender_name;type:varchar(128);comment:发送者名称"`
	Content       string    `json:"content" gorm:"column:content;type:text;comment:消息内容"`
	MsgType       string    `json:"msg_type" gorm:"column:msg_type;type:varchar(16);default:text;comment:消息类型 text/image/video/system"`
	AttachmentURL string    `json:"attachment_url" gorm:"column:attachment_url;type:varchar(512);comment:附件URL"`
	IsFaqReply    bool      `json:"is_faq_reply" gorm:"column:is_faq_reply;type:tinyint(1);default:0;comment:是否FAQ自动回复"`
	FaqID         *int64    `json:"faq_id" gorm:"column:faq_id;type:bigint;comment:匹配的FAQ ID"`
	IsRead        bool      `json:"is_read" gorm:"column:is_read;type:tinyint(1);default:0;index;comment:是否已读"`
	ReadAt        *time.Time `json:"read_at" gorm:"column:read_at;type:datetime(0);comment:阅读时间"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;type:datetime(0);autoCreateTime"`
}

func (ChatMessage) TableName() string { return "lc_chat_message" }
