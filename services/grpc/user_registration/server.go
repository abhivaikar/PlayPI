package user_registration

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	pb "github.com/abhivaikar/playpi/services/grpc/user_registration/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	users  map[string]pb.User // username as key
	mu     sync.Mutex
	tokens map[string]string // token -> username
}

func NewServer() *server {
	return &server{
		users:  make(map[string]pb.User),
		tokens: make(map[string]string),
	}
}

func (s *server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[req.User.Username]; exists {
		return &pb.RegisterUserResponse{
			Success: false,
			Message: "Username already exists",
		}, nil
	}

	s.users[req.User.Username] = *req.User
	return &pb.RegisterUserResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

func (s *server) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.Username]
	if !exists || user.Password != req.Password {
		return &pb.SignInResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	token := "token_" + req.Username
	s.tokens[token] = req.Username
	return &pb.SignInResponse{
		Success: true,
		Message: "Signed in successfully",
		Token:   token,
	}, nil
}

func (s *server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	username, exists := s.tokens[req.Token]
	if !exists {
		return nil, errors.New("invalid token")
	}

	user := s.users[username]
	return &pb.GetProfileResponse{User: &user}, nil
}

func (s *server) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	username, exists := s.tokens[req.Token]
	if !exists {
		return &pb.UpdateProfileResponse{
			Success: false,
			Message: "Invalid token",
		}, nil
	}

	s.users[username] = *req.User
	return &pb.UpdateProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
	}, nil
}

func (s *server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	username, exists := s.tokens[req.Token]
	if !exists {
		return &pb.DeleteAccountResponse{
			Success: false,
			Message: "Invalid token",
		}, nil
	}

	delete(s.users, username)
	delete(s.tokens, req.Token)
	return &pb.DeleteAccountResponse{
		Success: true,
		Message: "Account deleted successfully",
	}, nil
}

func StartServer() {
	// Listen on the specified port
	listener, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("Failed to listen on port 8084: %v", err)
	}

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	// Register the UserService with the gRPC server
	pb.RegisterUserServiceServer(grpcServer, NewServer())

	log.Println("Starting User Registration Service on port 8084...")

	// Start the gRPC server and serve requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
