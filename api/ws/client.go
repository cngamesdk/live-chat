package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cngamesdk/live-chat/api/logger"
	gorillaWs "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	SessionID int64
	ProductID int64
	UserID    string
	UserType  string // "user" or "agent"
	TraceID   string
	Conn      *gorillaWs.Conn
	Send      chan []byte
}

func (c *Client) ctx() context.Context {
	return logger.WithTraceID(context.Background(), c.TraceID)
}

func (c *Client) ReadPump() {
	defer func() {
		GlobalHub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if gorillaWs.IsUnexpectedCloseError(err, gorillaWs.CloseGoingAway, gorillaWs.CloseAbnormalClosure) {
				logger.Error(c.ctx(), "ws read error", zap.Error(err))
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}

		msg.Timestamp = time.Now().UnixMilli()
		msg.SessionID = c.SessionID

		// Handle message from client
		HandleClientMessage(c, &msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(gorillaWs.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(gorillaWs.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(gorillaWs.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
