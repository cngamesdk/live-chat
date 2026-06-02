package model

import "time"

type Agent struct {
	BaseModel
	UserID          int64      `json:"user_id" gorm:"column:user_id;type:bigint;index;not null;comment:系统用户ID"`
	ProductID       int64      `json:"product_id" gorm:"column:product_id;type:bigint;index;not null;comment:产品ID"`
	Status          string     `json:"status" gorm:"column:status;type:varchar(16);default:offline;comment:状态 online/offline/busy"`
	MaxConcurrent   int        `json:"max_concurrent" gorm:"column:max_concurrent;type:int;default:5;comment:最大并发会话数"`
	CurrentSessions int        `json:"current_sessions" gorm:"column:current_sessions;type:int;default:0;comment:当前活跃会话数"`
	TotalServed     int64      `json:"total_served" gorm:"column:total_served;type:bigint;default:0;comment:总服务次数"`
	LastOnlineAt    *time.Time `json:"last_online_at" gorm:"column:last_online_at;type:datetime(0);comment:最后上线时间"`
}

func (Agent) TableName() string { return "lc_agent" }
