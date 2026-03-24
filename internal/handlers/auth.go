package handlers

import (
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register регистрация пользователя
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			utils.EmailAlreadyExists(c)
			return
		}
		utils.InternalServerError(c, "Ошибка при регистрации: "+err.Error())
		return
	}

	utils.SuccessResponse(c, resp)
}

// Login вход пользователя
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		utils.InvalidCredentials(c)
		return
	}

	utils.SuccessResponse(c, resp)
}
