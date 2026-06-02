package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cngamesdk/live-chat/api/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SessionClaims struct {
	SessionID int64  `json:"session_id"`
	ProductID int64  `json:"product_id"`
	UserID    string `json:"user_id"`
	jwt.RegisteredClaims
}

var sessionSecret = []byte("live-chat-session-secret-key-2026")

func SetSessionSecret(secret string) {
	if secret != "" {
		sessionSecret = []byte(secret)
	}
}

func GenerateSessionToken(sessionID, productID int64, userID string) (string, error) {
	claims := SessionClaims{
		SessionID: sessionID,
		ProductID: productID,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sessionSecret)
}

func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Session-Token")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			model.FailWithCode(model.CodeUnauthorized, "missing session token", c)
			c.Abort()
			return
		}

		claims := &SessionClaims{}
		parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return sessionSecret, nil
		})
		if err != nil || !parsed.Valid {
			model.FailWithCode(model.CodeUnauthorized, "invalid session token", c)
			c.Abort()
			return
		}

		c.Set("session_id", claims.SessionID)
		c.Set("product_id", claims.ProductID)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func SignUserToken(userID string, productCode string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%s:%s:%d", userID, productCode, time.Now().Unix()/300)))
	return hex.EncodeToString(h.Sum(nil))
}

func ParseSessionToken(tokenString string) (*SessionClaims, error) {
	claims := &SessionClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return sessionSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, err
	}
	_ = parsed
	return claims, nil
}

func VerifyUserToken(userID, productCode, token, secret string) bool {
	expected := SignUserToken(userID, productCode, secret)
	return hmac.Equal([]byte(token), []byte(expected))
}
