package model

import "time"

type Upload struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string `json:"user_id" gorm:"column:user_id;type:varchar(128);index;comment:上传者"`
	FileName  string `json:"file_name" gorm:"column:file_name;type:varchar(256);comment:原始文件名"`
	FileURL   string `json:"file_url" gorm:"column:file_url;type:varchar(512);comment:访问URL"`
	FileType  string `json:"file_type" gorm:"column:file_type;type:varchar(16);comment:文件类型 image/video"`
	FileSize  int64  `json:"file_size" gorm:"column:file_size;type:bigint;comment:文件大小"`
	MimeType  string `json:"mime_type" gorm:"column:mime_type;type:varchar(64);comment:MIME类型"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(0);autoCreateTime"`
}

func (Upload) TableName() string { return "lc_upload" }
