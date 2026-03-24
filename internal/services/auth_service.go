package services

import (
	"errors"
	"fmt"

	"razum-backend/internal/config"
	"razum-backend/internal/models"
	"razum-backend/internal/repository"
	"razum-backend/internal/utils"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// RegisterRequest структура для регистрации
type RegisterRequest struct {
	Email     string            `json:"email" binding:"required,email"`
	Password  string            `json:"password" binding:"required,min=6"`
	FullName  string            `json:"full_name" binding:"required"`
	Role      models.UserRole   `json:"role" binding:"required,oneof=organizer participant observer admin"`
	City      *string           `json:"city"`
	Age       *int              `json:"age"`
	Direction *models.Direction `json:"direction"`
}

// RegisterResponse структура ответа при регистрации
type RegisterResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Создаем пользователя (пароль в открытом виде)
	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  req.Password,
		FullName:  req.FullName,
		Role:      req.Role,
		City:      req.City,
		Age:       req.Age,
		Direction: req.Direction,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Генерируем JWT токен
	token, err := utils.GenerateJWT(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RegisterResponse{
		User:  user,
		Token: token,
	}, nil
}

// LoginRequest структура для входа
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse структура ответа при входе
type LoginResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// Login авторизует пользователя
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Ищем пользователя по email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("ошибка в поиске: %w", err)
	}
	if user == nil {
		return nil, errors.New("неверный пароль или логин")
	}

	// Проверяем пароль (прямое сравнение)
	if user.Password != req.Password {
		return nil, errors.New("неверный пароль или логин")
	}

	// Генерируем JWT токен
	token, err := utils.GenerateJWT(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, fmt.Errorf("ошибка в генерации токена: %w", err)
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}
