package ws

import (
	"net/http"

	"github.com/cngamesdk/live-chat/api/middleware"
	"github.com/gin-gonic/gin"
	gorillaWs "github.com/gorilla/websocket"
)

var wsUpgrader = gorillaWs.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpgradeHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "missing token"})
		return
	}

	claims, err := middleware.ParseSessionToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		SessionID: claims.SessionID,
		ProductID: claims.ProductID,
		UserID:    claims.UserID,
		UserType:  "user",
		TraceID:   middleware.GetTraceID(c),
		Conn:      conn,
		Send:      make(chan []byte, 256),
	}

	GlobalHub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
