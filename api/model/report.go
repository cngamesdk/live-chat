package model

import "time"

type DailyReport struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID       int64     `json:"product_id" gorm:"column:product_id;type:bigint;index;not null;comment:产品ID"`
	ReportDate      string    `json:"report_date" gorm:"column:report_date;type:date;index;not null;comment:报表日期"`
	TotalSessions   int       `json:"total_sessions" gorm:"column:total_sessions;type:int;default:0;comment:总会话数"`
	FaqResolved     int       `json:"faq_resolved" gorm:"column:faq_resolved;type:int;default:0;comment:FAQ解决数"`
	AgentResolved   int       `json:"agent_resolved" gorm:"column:agent_resolved;type:int;default:0;comment:人工解决数"`
	AvgResponseTime int       `json:"avg_response_time" gorm:"column:avg_response_time;type:int;default:0;comment:平均响应时间(秒)"`
	AvgSessionTime  int       `json:"avg_session_time" gorm:"column:avg_session_time;type:int;default:0;comment:平均会话时长(秒)"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime(0);autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:datetime(0);autoUpdateTime"`
}

func (DailyReport) TableName() string { return "lc_daily_report" }
