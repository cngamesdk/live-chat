package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cngamesdk/live-chat/api/config"
	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/initialize"
	"github.com/cngamesdk/live-chat/api/logger"
	"github.com/cngamesdk/live-chat/api/ws"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// -------- 启动参数 --------
	configPath := flag.String("c", "config.yaml", "配置文件路径 (yaml)")
	serverPort := flag.Int("p", 0, "服务端口 (覆盖配置文件)")
	serverMode := flag.String("m", "", "运行模式 debug/release (覆盖配置文件)")
	logLevel := flag.String("l", "", "日志级别 debug/info/warn/error (覆盖配置文件)")
	dbHost := flag.String("db-host", "", "数据库地址 (覆盖配置文件)")
	dbPort := flag.Int("db-port", 0, "数据库端口 (覆盖配置文件)")
	dbName := flag.String("db-name", "", "数据库名称 (覆盖配置文件)")
	dbUser := flag.String("db-user", "", "数据库用户 (覆盖配置文件)")
	dbPass := flag.String("db-pass", "", "数据库密码 (覆盖配置文件)")
	redisAddr := flag.String("redis-addr", "", "Redis地址 (覆盖配置文件)")
	flag.Parse()

	// -------- 加载配置文件 --------
	v := viper.New()
	v.SetConfigFile(*configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file %s: %v", *configPath, err)
	}

	v.SetEnvPrefix("LIVE_CHAT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	// -------- 命令行参数覆盖 --------
	if *serverPort > 0 {
		cfg.Server.Port = *serverPort
	}
	if *serverMode != "" {
		cfg.Server.Mode = *serverMode
	}
	if *logLevel != "" {
		cfg.Log.Level = *logLevel
	}
	if *dbHost != "" {
		cfg.Mysql.Path = *dbHost
	}
	if *dbPort > 0 {
		cfg.Mysql.Port = *dbPort
	}
	if *dbName != "" {
		cfg.Mysql.DbName = *dbName
	}
	if *dbUser != "" {
		cfg.Mysql.Username = *dbUser
	}
	if *dbPass != "" {
		cfg.Mysql.Password = *dbPass
	}
	if *redisAddr != "" {
		cfg.Redis.Addr = *redisAddr
	}

	global.GVA_VP = v
	global.GVA_CONFIG = cfg

	// -------- 初始化日志 --------
	zapLogger := initialize.InitLogger()
	logger.Init(zapLogger)
	global.GVA_LOG = zapLogger // keep backward compatibility

	logger.Info(nil, "logger initialized",
		zap.String("level", cfg.Log.Level),
		zap.String("file", cfg.Log.FilePath),
	)

	// -------- 初始化数据库 --------
	global.GVA_DB = initialize.Gorm()
	if global.GVA_DB != nil {
		initialize.RegisterTables()
		logger.Info(nil, "database connected and tables migrated")
	}

	// -------- 初始化 Redis --------
	global.GVA_REDIS = initialize.Redis()

	// -------- 初始化 WebSocket Hub --------
	ws.GlobalHub = ws.NewHub()
	go ws.GlobalHub.Run()

	// -------- 启动服务 --------
	r := initialize.Routers()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logger.Info(nil, "server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(nil, "server error", zap.Error(err))
		}
	}()

	// -------- 优雅关闭 --------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(nil, "shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(nil, "server forced shutdown", zap.Error(err))
	}
	logger.Info(nil, "server exited")
}
