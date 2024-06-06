package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
	pb "github.com/nurtikaga/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

func (s *server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	tokenHash := sha256.Sum256([]byte(req.Token))
	var userID int64
	var expiry time.Time
	var scope string

	query := `SELECT user_id, expiry, scope FROM tokens WHERE hash = $1`
	err := s.db.QueryRow(query, tokenHash[:]).Scan(
		&userID,
		&expiry,
		&scope,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Invalid token: %s", req.Token)
			return &pb.ValidateTokenResponse{Valid: false}, nil
		}
		log.Printf("Error querying token: %v", err)
		return nil, err
	}

	if time.Now().After(expiry) {
		log.Printf("Token expired: %s", req.Token)
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{Valid: true, UserId: userID, Scope: scope}, nil
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("dsnAuth"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{db: db})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
