package utils

import (
	"net/http"

	appErrors "razum-backend/internal/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, err *appErrors.AppError) {
	c.JSON(err.Status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    err.Code,
			Message: err.Message,
		},
	})
}

// Обертки для быстрого использования
func Unauthorized(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrUnauthorized)
}

func Forbidden(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrForbidden)
}

func NotFound(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrNotFound)
}

func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, &appErrors.AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  400,
	})
}

func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, &appErrors.AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  500,
	})
}

// Специфичные ошибки
func InvalidCredentials(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrInvalidCredentials)
}

func EmailAlreadyExists(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrEmailAlreadyExists)
}

func EventNotFound(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrEventNotFound)
}

func AlreadyRegistered(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrAlreadyRegistered)
}

func InvalidQRCode(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrInvalidQRCode)
}

func QRCodeAlreadyUsed(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrQRCodeAlreadyUsed)
}

func NotOrganizer(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrNotOrganizer)
}

func OnlyOrganizers(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrOnlyOrganizers)
}

func RegistrationDeadlinePassed(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrRegistrationDeadlinePassed)
}

func EventCancelled(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrEventCancelled)
}

func CannotCancelConfirmed(c *gin.Context) {
	ErrorResponse(c, appErrors.ErrCannotCancelConfirmed)
}
