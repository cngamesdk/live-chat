package model

import "time"

type ChatSession struct {
	BaseModel
	ProductID int64  `json:"product_id" gorm:"column:product_id;type:bigint;index;not null;comment:产品ID"`
	UserID    string `json:"user_id" gorm:"column:user_id;type:varchar(128);index;not null;comment:用户标识"`
	UserName  string `json:"user_name" gorm:"column:user_name;type:varchar(128);comment:用户名称"`
	AgentID   *int64 `json:"agent_id" gorm:"column:agent_id;type:bigint;index;comment:分配的客服ID"`
	Status    string `json:"status" gorm:"column:status;type:varchar(16);default:waiting;index;comment:状态 waiting/active/closed"`
	Source    string `json:"source" gorm:"column:source;type:varchar(32);default:h5;comment:来源"`
	UserIP    string `json:"user_ip" gorm:"column:user_ip;type:varchar(64);comment:客户端IP"`
	UserAgent string `json:"user_agent" gorm:"column:user_agent;type:varchar(512);comment:客户端UA"`
	ClosedAt  *time.Time `json:"closed_at" gorm:"column:closed_at;type:datetime(0);comment:关闭时间"`
}

func (ChatSession) TableName() string { return "lc_chat_session" }
