package initialize

import (
	"fmt"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/logger"
	"github.com/cngamesdk/live-chat/api/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Gorm() *gorm.DB {
	cfg := global.GVA_CONFIG.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.Username, cfg.Password, cfg.Path, cfg.Port, cfg.DbName, cfg.Config)

	logLevel := gormlogger.Info
	if global.GVA_CONFIG.Server.Mode == "release" {
		logLevel = gormlogger.Warn
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		logger.Error(nil, "failed to connect database", zap.Error(err))
		panic(fmt.Sprintf("database connect error: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("get sql.DB error: %v", err))
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	return db
}

func RegisterTables() {
	if global.GVA_DB == nil {
		return
	}
	global.GVA_DB.AutoMigrate(
		&model.Product{},
		&model.Faq{},
		&model.ChatSession{},
		&model.ChatMessage{},
		&model.Agent{},
		&model.Upload{},
		&model.DailyReport{},
	)
}
