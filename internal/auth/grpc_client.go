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
