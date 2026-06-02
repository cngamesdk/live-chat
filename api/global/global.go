package global

import (
	"github.com/cngamesdk/live-chat/api/config"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GVA_DB     *gorm.DB
	GVA_REDIS  redis.UniversalClient
	GVA_CONFIG config.Config
	GVA_LOG    *zap.Logger
	GVA_VP     *viper.Viper
)
