package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:datetime(0);autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:datetime(0);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Product struct {
	BaseModel
	ProductCode    string `json:"product_code" gorm:"column:product_code;type:varchar(64);uniqueIndex;not null;comment:产品编码"`
	Name           string `json:"name" gorm:"column:name;type:varchar(128);not null;comment:产品名称"`
	Logo           string `json:"logo" gorm:"column:logo;type:varchar(512);comment:Logo URL"`
	WelcomeTitle   string `json:"welcome_title" gorm:"column:welcome_title;type:varchar(256);comment:欢迎标题"`
	WelcomeMessage string `json:"welcome_message" gorm:"column:welcome_message;type:text;comment:欢迎消息"`
	Status         string `json:"status" gorm:"column:status;type:varchar(16);default:normal;comment:状态 normal/disabled"`
}

func (Product) TableName() string { return "lc_product" }
