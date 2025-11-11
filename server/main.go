package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/Vadim-Makhnev/grpc/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	proto.UnimplementedUserServiceServer
	mu    sync.RWMutex
	users map[string]*proto.UserResponse
}

func NewUserService() *userService {
	return &userService{
		users: make(map[string]*proto.UserResponse),
	}
}

func (s *userService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	if req.Name == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "name and email are required")
	}

	if req.Age <= 0 {
		return nil, status.Error(codes.InvalidArgument, "age must be greater than 0")
	}

	user := &proto.UserResponse{
		Id:    uuid.NewString(),
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	s.mu.Lock()
	s.users[user.Id] = user
	s.mu.Unlock()

	log.Printf("Created user: %s (%s)", user.Name, user.Id)
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	s.mu.RLock()
	user, exists := s.users[req.Id]
	s.mu.RUnlock()

	if !exists {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]*proto.UserResponse, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return &proto.ListUsersResponse{Users: users}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("listener: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, NewUserService())

	log.Println("gRPC server listening on :5000")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server: %v", err)
	}
}
