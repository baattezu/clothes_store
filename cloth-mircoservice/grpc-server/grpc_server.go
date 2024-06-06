package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"

	pb "cloth-service/proto"

	"google.golang.org/grpc"
)

type clothServiceServer struct {
	pb.UnimplementedClothServiceServer
	db *sql.DB // Database connection
}

// NewClothServiceServer creates a new clothServiceServer with the given database connection.
func NewClothServiceServer(db *sql.DB) *clothServiceServer {
	return &clothServiceServer{
		db: db,
	}
}

func (s *clothServiceServer) CreateCloth(ctx context.Context, req *pb.CreateClothRequest) (*pb.CreateClothResponse, error) {
	query := `INSERT INTO clothes_info (cloth_name, cloth_cost, cloth_size) VALUES ($1, $2, $3) RETURNING id`
	var id int32
	err := s.db.QueryRowContext(ctx, query, req.ClothName, req.ClothCost, req.ClothSize).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.CreateClothResponse{Id: id}, nil
}

func (s *clothServiceServer) GetCloth(ctx context.Context, req *pb.GetClothRequest) (*pb.GetClothResponse, error) {
	query := `SELECT id, cloth_name, cloth_cost, cloth_size, created_at, updated_at, version FROM clothes_info WHERE id = $1`
	var cloth pb.Cloth
	err := s.db.QueryRowContext(ctx, query, req.Id).Scan(&cloth.Id, &cloth.ClothName, &cloth.ClothCost, &cloth.ClothSize, &cloth.CreatedAt, &cloth.UpdatedAt, &cloth.Version)
	if err != nil {
		return nil, err
	}
	return &pb.GetClothResponse{Cloth: &cloth}, nil
}

func (s *clothServiceServer) EditCloth(ctx context.Context, req *pb.EditClothRequest) (*pb.EditClothResponse, error) {
	query := `UPDATE clothes_info SET cloth_name = $1, cloth_cost = $2, cloth_size = $3 WHERE id = $4`
	_, err := s.db.ExecContext(ctx, query, req.ClothName, req.ClothCost, req.ClothSize, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.EditClothResponse{Success: true}, nil
}

func (s *clothServiceServer) DeleteCloth(ctx context.Context, req *pb.DeleteClothRequest) (*pb.DeleteClothResponse, error) {
	query := `DELETE FROM clothes_info WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, req.Id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return &pb.DeleteClothResponse{Success: true}, nil
}

func (s *clothServiceServer) ListClothes(ctx context.Context, req *pb.ListClothesRequest) (*pb.ListClothesResponse, error) {
	// Initialize slice to store clothes
	var clothes []*pb.Cloth

	// Construct the SQL query based on the request parameters
	query := "SELECT id, cloth_name, cloth_cost, cloth_size, created_at, updated_at, version FROM clothes_info WHERE 1=1"
	args := []interface{}{}

	// Add WHERE clauses based on the request parameters
	if req.ClothName != "" {
		query += " AND cloth_name LIKE $1"
		args = append(args, "%"+req.ClothName+"%")
	}
	if req.ClothSize != "" {
		query += " AND cloth_size = $2"
		args = append(args, req.ClothSize)
	}
	// You can add more conditions based on other request parameters

	// Execute the query
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the clothes slice
	for rows.Next() {
		var cloth pb.Cloth
		if err := rows.Scan(&cloth.Id, &cloth.ClothName, &cloth.ClothCost, &cloth.ClothSize, &cloth.CreatedAt, &cloth.UpdatedAt, &cloth.Version); err != nil {
			return nil, err
		}
		clothes = append(clothes, &cloth)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of clothes
	return &pb.ListClothesResponse{Clothes: clothes}, nil
}

func main() {
	// Initialize your database connection
	db, err := sql.Open("postgres", "postgres://nurtileu:root@localhost/a.nurtileuDB?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create gRPC server
	server := grpc.NewServer()
	pb.RegisterClothServiceServer(server, NewClothServiceServer(db))

	// Listen on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start the gRPC server
	log.Println("gRPC server listening on port 50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
