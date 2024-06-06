package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"
	pb "order-microservice/internal/proto"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &AuthClient{client: pb.NewAuthServiceClient(conn)}, nil
}

func (c *AuthClient) ValidateToken(token string) (bool, int64, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.ValidateTokenRequest{Token: token}
	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return false, 0, "", err
	}
	return resp.Valid, resp.UserId, resp.Scope, nil
}
