package httputil

import "github.com/gin-gonic/gin"

// Error представляет базовую ошибку
type Error struct {
	Message string `json:"message" example:"Произошла ошибка"`
}

// ErrorResponse представляет расширенный ответ с ошибкой
type ErrorResponse struct {
	Error   string            `json:"error" example:"Ошибка валидации"`
	Details map[string]string `json:"details,omitempty" example:"username:Имя пользователя обязательно"`
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Error   string            `json:"error" example:"Ошибка валидации"`
	Details map[string]string `json:"details" example:"field:описание ошибки"`
}

// AuthError представляет ошибку аутентификации
type AuthError struct {
	Message string `json:"message" example:"Неверные учетные данные"`
}

// NotFoundError представляет ошибку отсутствия ресурса
type NotFoundError struct {
	Message string `json:"message" example:"Ресурс не найден"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// NewError создает новую ошибку
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// NewErrorResponse создает новый ответ с ошибкой
func NewErrorResponse(err string, details map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Error:   err,
		Details: details,
	}
}

// NewValidationError создает новую ошибку валидации
func NewValidationError(details map[string]string) *ValidationError {
	return &ValidationError{
		Error:   "Ошибка валидации",
		Details: details,
	}
}

// NewAuthError создает новую ошибку аутентификации
func NewAuthError(message string) *AuthError {
	return &AuthError{
		Message: message,
	}
}

// NewNotFoundError создает новую ошибку отсутствия ресурса
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}
