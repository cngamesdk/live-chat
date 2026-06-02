package router

import (
	"github.com/cngamesdk/live-chat/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(r *gin.RouterGroup) {
	// Public routes (no auth)
	publicGroup := r.Group("")
	{
		publicGroup.GET("/product/:code/info", ApiGroupApp.ProductApi.GetProductInfo)
		publicGroup.POST("/chat/init", ApiGroupApp.ChatApi.InitSession)
	}

	// Session-authenticated routes
	sessionGroup := r.Group("").Use(middleware.SessionAuth())
	{
		sessionGroup.GET("/chat/faq", ApiGroupApp.FaqApi.QueryFaq)
		sessionGroup.GET("/chat/history", ApiGroupApp.ChatApi.GetHistory)
		sessionGroup.POST("/chat/mark-read", ApiGroupApp.ChatApi.MarkMessagesAsRead)
		sessionGroup.GET("/chat/unread-count", ApiGroupApp.ChatApi.GetUnreadCount)
		sessionGroup.POST("/chat/close", ApiGroupApp.ChatApi.CloseSession)
		sessionGroup.POST("/upload", ApiGroupApp.UploadApi.Upload)
	}

	// Admin routes (for background system to call)
	adminGroup := r.Group("")
	{
		adminGroup.POST("/admin/agent-reply", ApiGroupApp.AdminApi.AgentReply)
		adminGroup.POST("/admin/close-session", ApiGroupApp.SessionApi.CloseSession)
	}
}
