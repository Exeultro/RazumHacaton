package middleware

import (
	"log"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		log.Printf("RequireRole: role=%v, exists=%v, allowed=%v", role, exists, allowedRoles)

		if !exists {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				log.Printf("RequireRole: allowed %s", roleStr)
				c.Next()
				return
			}
		}

		log.Printf("RequireRole: forbidden for %s", roleStr)
		utils.Forbidden(c)
		c.Abort()
	}
}

// RequireOrganizer проверяет, что пользователь - организатор ИЛИ админ
func RequireOrganizer() gin.HandlerFunc {
	return RequireRole("organizer", "admin")
}

// RequireParticipant проверяет, что пользователь - участник
func RequireParticipant() gin.HandlerFunc {
	return RequireRole("participant", "organizer", "observer")
}

// RequireObserver проверяет, что пользователь - наблюдатель
func RequireObserver() gin.HandlerFunc {
	return RequireRole("observer")
}
