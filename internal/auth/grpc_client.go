package auth

import (
	"context"
	"fmt"

	pb "github.com/fire9900/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcAuthClient struct {
	client pb.AuthServiceClient
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	Tokens *TokenPair `json:"tokens"`
	User   *User      `json:"user"`
	Error  string     `json:"error,omitempty"`
}

func NewGrpcAuthClient(address string) (*GrpcAuthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к серверу авторизации: %w", err)
	}

	return &GrpcAuthClient{
		client: pb.NewAuthServiceClient(conn),
	}, nil
}

func (c *GrpcAuthClient) ValidateToken(token string) (int64, error) {
	resp, err := c.client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return 0, fmt.Errorf("ошибка валидации токена: %w", err)
	}

	if !resp.Valid {
		return 0, fmt.Errorf("недействительный токен: %s", resp.Error)
	}

	return resp.UserId, nil
}

func (c *GrpcAuthClient) Register(username, email, password string) (*AuthResponse, error) {
	resp, err := c.client.Register(context.Background(), &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка регистрации: %w", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return convertAuthResponse(resp), nil
}

func (c *GrpcAuthClient) Login(email, password string) (*AuthResponse, error) {
	resp, err := c.client.Login(context.Background(), &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка входа: %w", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return convertAuthResponse(resp), nil
}

func (c *GrpcAuthClient) RefreshTokens(refreshToken string) (*AuthResponse, error) {
	resp, err := c.client.RefreshTokens(context.Background(), &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления токенов: %w", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return convertAuthResponse(resp), nil
}

func convertAuthResponse(resp *pb.AuthResponse) *AuthResponse {
	return &AuthResponse{
		Tokens: &TokenPair{
			AccessToken:  resp.Tokens.AccessToken,
			RefreshToken: resp.Tokens.RefreshToken,
		},
		User: &User{
			ID:       resp.User.Id,
			Username: resp.User.Username,
			Email:    resp.User.Email,
		},
	}
}
