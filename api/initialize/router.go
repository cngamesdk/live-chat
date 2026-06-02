package initialize

import (
	"net/http"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/middleware"
	"github.com/cngamesdk/live-chat/api/router"
	"github.com/cngamesdk/live-chat/api/ws"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	mode := global.GVA_CONFIG.Server.Mode
	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.TraceID())    // 注入 trace_id
	r.Use(middleware.RequestLog()) // 请求日志
	r.Use(middleware.Cors())

	// Static file server for uploads
	uploadPath := global.GVA_CONFIG.Upload.Path
	r.Static("/uploads", uploadPath)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	apiGroup := r.Group("/api/v1")
	router.RegisterChatRoutes(apiGroup)

	// WebSocket route
	r.GET("/ws/chat", ws.UpgradeHandler)

	return r
}
