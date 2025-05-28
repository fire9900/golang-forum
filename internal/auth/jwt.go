package auth

import (
	"fmt"
	"time"

	"github.com/fire9900/golang-forum/internal/model"
	"github.com/golang-jwt/jwt"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 24 * time.Hour * 30
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenManager interface {
	GenerateTokens(user *model.User) (*Tokens, error)
	ValidateAccessToken(tokenString string) (int64, error)
	ValidateRefreshToken(tokenString string) (int64, error)
	RefreshTokens(refreshToken string) (*Tokens, *model.User, error)
}

type Manager struct {
	signingKey string
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Type   string `json:"type"`
	jwt.StandardClaims
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("пустой ключ подписи")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) GenerateTokens(user *model.User) (*Tokens, error) {
	// Генерация Access токена
	accessToken, err := m.generateToken(user.ID, "access", accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации access токена: %w", err)
	}

	// Генерация Refresh токена
	refreshToken, err := m.generateToken(user.ID, "refresh", refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации refresh токена: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *Manager) generateToken(userID int64, tokenType string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ValidateAccessToken(tokenString string) (int64, error) {
	return m.validateToken(tokenString, "access")
}

func (m *Manager) ValidateRefreshToken(tokenString string) (int64, error) {
	return m.validateToken(tokenString, "refresh")
}

func (m *Manager) validateToken(tokenString, expectedType string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("ошибка валидации токена: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, fmt.Errorf("невалидная структура токена")
	}

	if !token.Valid {
		return 0, fmt.Errorf("токен недействителен")
	}

	if claims.Type != expectedType {
		return 0, fmt.Errorf("неверный тип токена")
	}

	return claims.UserID, nil
}

func (m *Manager) RefreshTokens(refreshToken string) (*Tokens, *model.User, error) {
	// Валидация refresh токена
	userID, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка валидации refresh токена: %w", err)
	}

	// TODO: Здесь должен быть запрос к сервису пользователей для получения актуальных данных
	user := &model.User{
		ID: userID,
	}

	// Генерация новой пары токенов
	tokens, err := m.GenerateTokens(user)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка генерации новых токенов: %w", err)
	}

	return tokens, user, nil
}
