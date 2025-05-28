package controllers

import (
	"net/http"

	"github.com/fire9900/golang-forum/internal/auth"
	"github.com/fire9900/golang-forum/internal/model"
	"github.com/fire9900/golang-forum/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service      *service.UserService
	tokenManager auth.TokenManager
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type updateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type authResponse struct {
	Tokens *auth.Tokens `json:"tokens"`
	User   *model.User  `json:"user"`
}

func NewUserController(service *service.UserService, tokenManager auth.TokenManager) *UserController {
	return &UserController{
		service:      service,
		tokenManager: tokenManager,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка валидации",
			"details": map[string]string{
				"username": "Имя пользователя обязательно и должно быть от 3 до 50 символов",
				"email":    "Email обязателен и должен быть корректным",
				"password": "Пароль обязателен и должен быть от 6 до 50 символов",
			},
		})
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := c.service.Create(user); err != nil {
		if err == service.ErrEmailAlreadyTaken {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Этот email уже занят"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}

	// Генерация токенов
	tokens, err := c.tokenManager.GenerateTokens(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токенов"})
		return
	}

	ctx.JSON(http.StatusCreated, authResponse{
		Tokens: tokens,
		User:   user,
	})
}

func (c *UserController) Login(ctx *gin.Context) {
	var login model.UserLogin
	if err := ctx.ShouldBindJSON(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка валидации",
			"details": map[string]string{
				"email":    "Email обязателен и должен быть корректным",
				"password": "Пароль обязателен и должен быть от 6 символов",
			},
		})
		return
	}

	user, err := c.service.Login(&login)
	if err != nil {
		switch err {
		case service.ErrUserNotFound, service.ErrInvalidPassword:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	// Генерация токенов
	tokens, err := c.tokenManager.GenerateTokens(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токенов"})
		return
	}

	ctx.JSON(http.StatusOK, authResponse{
		Tokens: tokens,
		User:   user,
	})
}

func (c *UserController) RefreshTokens(ctx *gin.Context) {
	var req refreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка валидации",
			"details": map[string]string{
				"refresh_token": "Refresh токен обязателен",
			},
		})
		return
	}

	tokens, user, err := c.tokenManager.RefreshTokens(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный refresh токен"})
		return
	}

	// Получаем актуальные данные пользователя
	user, err = c.service.GetByID(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, authResponse{
		Tokens: tokens,
		User:   user,
	})
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	user, err := c.service.GetByID(userID.(int64))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка валидации",
			"details": map[string]string{
				"username": "Имя пользователя обязательно и должно быть от 3 до 50 символов",
				"email":    "Email обязателен и должен быть корректным",
			},
		})
		return
	}

	user := &model.User{
		ID:       userID.(int64),
		Username: req.Username,
		Email:    req.Email,
	}

	if err := c.service.Update(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) UpdatePassword(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required,min=6"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка валидации",
			"details": map[string]string{
				"old_password": "Старый пароль обязателен и должен быть от 6 символов",
				"new_password": "Новый пароль обязателен и должен быть от 6 до 50 символов",
			},
		})
		return
	}

	if err := c.service.UpdatePassword(userID.(int64), req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case service.ErrInvalidPassword:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный старый пароль"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	ctx.Status(http.StatusOK)
}
