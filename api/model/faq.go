package model

type Faq struct {
	BaseModel
	ProductID  int64  `json:"product_id" gorm:"column:product_id;type:bigint;index;not null;comment:产品ID"`
	Category   string `json:"category" gorm:"column:category;type:varchar(64);comment:分类"`
	Question   string `json:"question" gorm:"column:question;type:varchar(512);not null;comment:问题"`
	Answer     string `json:"answer" gorm:"column:answer;type:text;not null;comment:答案"`
	Keywords   string `json:"keywords" gorm:"column:keywords;type:varchar(512);comment:关键词"`
	Priority   int    `json:"priority" gorm:"column:priority;type:int;default:0;comment:优先级"`
	MatchCount int64  `json:"match_count" gorm:"column:match_count;type:bigint;default:0;comment:匹配次数"`
	Status     string `json:"status" gorm:"column:status;type:varchar(16);default:normal;comment:状态 normal/disabled"`
}

func (Faq) TableName() string { return "lc_faq" }
