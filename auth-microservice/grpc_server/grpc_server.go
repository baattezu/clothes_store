package main

import (
	"context"
	"log"
	"net"

	"github.com/nurtikaga/proto" // Import your proto package here
	"google.golang.org/grpc"
)

// server is used to implement auth_microservice.AuthServiceServer.
type server struct {
	proto.UnimplementedAuthServiceServer
}

// RegisterUser implements auth_microservice.AuthServiceServer.RegisterUser
func (s *server) RegisterUser(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	// Implement your logic here
	return &proto.RegisterResponse{}, nil
}

// LoginUser implements auth_microservice.AuthServiceServer.LoginUser
func (s *server) LoginUser(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Implement your logic here
	return &proto.LoginResponse{}, nil
}

// GetUserByEmail implements auth_microservice.AuthServiceServer.GetUserByEmail
func (s *server) GetUserByEmail(ctx context.Context, in *proto.GetUserByEmailRequest) (*proto.GetUserByEmailResponse, error) {
	// Implement your logic here
	return &proto.GetUserByEmailResponse{}, nil
}

// GetUserById implements auth_microservice.AuthServiceServer.GetUserById
func (s *server) GetUserById(ctx context.Context, in *proto.GetUserByIdRequest) (*proto.GetUserByIdResponse, error) {
	// Implement your logic here
	return &proto.GetUserByIdResponse{}, nil
}

// UpdateUser implements auth_microservice.AuthServiceServer.UpdateUser
func (s *server) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	// Implement your logic here
	return &proto.UpdateUserResponse{}, nil
}

// DeleteUser implements auth_microservice.AuthServiceServer.DeleteUser
func (s *server) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	// Implement your logic here
	return &proto.DeleteUserResponse{}, nil
}

// GetAllUsers implements auth_microservice.AuthServiceServer.GetAllUsers
func (s *server) GetAllUsers(ctx context.Context, in *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	// Implement your logic here
	return &proto.GetAllUsersResponse{}, nil
}

func main() {
	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Create a new gRPC server
	s := grpc.NewServer()
	// Register the AuthService server
	proto.RegisterAuthServiceServer(s, &server{})
	// Start the server
	log.Println("Server started on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
