package auth

import (
	"context"
	"order-microservice/internal/proto"
	"time"

	"google.golang.org/grpc"
)

type AuthClient struct {
	client proto.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &AuthClient{client: proto.NewAuthServiceClient(conn)}, nil
}

func (c *AuthClient) ValidateToken(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &proto.ValidateTokenRequest{Token: token}
	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.Valid, nil
}
