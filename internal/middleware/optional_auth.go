package middleware

import (
	"strings"

	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// OptionalAuthMiddleware пытается получить пользователя из токена, но не требует его
func OptionalAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Убираем "Bearer " если есть
		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(token, secret)
		if err != nil {
			// Токен невалидный, но мы не блокируем запрос
			c.Next()
			return
		}

		// Токен валидный — сохраняем user_id и role
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
