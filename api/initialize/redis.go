package initialize

import (
	"context"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/logger"
	"github.com/redis/go-redis/v9"
)

func Redis() redis.UniversalClient {
	cfg := global.GVA_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Warn(nil, "redis connect failed, running without redis")
		return nil
	}
	logger.Info(nil, "redis connected successfully")
	return client
}
