package middleware

import (
	"net/http"
	"strings"

	"github.com/fire9900/golang-forum/internal/auth"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user_id"
)

func AuthMiddleware(authClient *auth.GrpcAuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пустой заголовок авторизации"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			return
		}

		if len(headerParts[1]) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Токен пуст"})
			return
		}

		userID, err := authClient.ValidateToken(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный access токен",
				"code":  "invalid_access_token",
			})
			return
		}

		c.Set(userCtx, userID)
		c.Next()
	}
}
