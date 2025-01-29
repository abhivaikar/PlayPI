package user_registration

import (
	"context"
	"errors"
	"log"
	"net"
	"regexp"
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

	// Validate required fields
	if req.User.Username == "" || len(req.User.Username) < 3 || len(req.User.Username) > 50 {
		return nil, errors.New("username must be between 3 and 50 characters")
	}
	if req.User.Password == "" || len(req.User.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	if req.User.FullName == "" {
		return nil, errors.New("fullname is required")
	}
	if req.User.Address != "" && len(req.User.Address) > 100 {
		return nil, errors.New("address cannot exceed 100 characters")
	}
	if !isValidEmail(req.User.Email) {
		return nil, errors.New("invalid email format")
	}
	if !isValidPhoneNumber(req.User.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// Check if username already exists
	if _, exists := s.users[req.User.Username]; exists {
		return nil, errors.New("username already exists")
	}

	// Register the user
	s.users[req.User.Username] = *req.User
	return &pb.RegisterUserResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

func isValidEmail(email string) bool {
	// Basic regex for email validation
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func isValidPhoneNumber(phone string) bool {
	// Check if the phone number is numeric and between 10-15 digits
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (s *server) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// Check username and password
	user, exists := s.users[req.Username]
	if !exists || user.Password != req.Password {
		return nil, errors.New("invalid username or password")
	}

	// Generate a token for the session
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

	// Validate token
	username, exists := s.tokens[req.Token]
	if !exists {
		return nil, errors.New("invalid token")
	}

	// Retrieve user profile
	user := s.users[username]
	return &pb.GetProfileResponse{User: &user}, nil
}

func (s *server) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate token
	username, exists := s.tokens[req.Token]
	if !exists {
		return nil, errors.New("invalid token")
	}

	// Get the current user profile
	user := s.users[username]

	// Validate and update username
	if req.User.Username != "" {
		if len(req.User.Username) < 3 || len(req.User.Username) > 50 {
			return nil, errors.New("username must be between 3 and 50 characters")
		}
		if _, exists := s.users[req.User.Username]; exists && req.User.Username != username {
			return nil, errors.New("username already exists")
		}
		user.Username = req.User.Username
	}

	// Validate and update password
	if req.User.Password != "" {
		if len(req.User.Password) < 8 {
			return nil, errors.New("password must be at least 8 characters long")
		}
		user.Password = req.User.Password
	}

	// Validate and update fullname
	if req.User.FullName != "" {
		if len(req.User.FullName) == 0 {
			return nil, errors.New("fullname is required")
		}
		user.FullName = req.User.FullName
	}

	// Validate and update email
	if req.User.Email != "" {
		if !isValidEmail(req.User.Email) {
			return nil, errors.New("invalid email format")
		}
		user.Email = req.User.Email
	}

	// Validate and update phone
	if req.User.Phone != "" {
		if !isValidPhoneNumber(req.User.Phone) {
			return nil, errors.New("invalid phone number format")
		}
		user.Phone = req.User.Phone
	}

	// Validate and update address
	if req.User.Address != "" {
		if len(req.User.Address) > 100 {
			return nil, errors.New("address cannot exceed 100 characters")
		}
		user.Address = req.User.Address
	}

	// Save the updated user profile
	s.users[username] = user

	return &pb.UpdateProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
	}, nil
}

func (s *server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate token
	username, exists := s.tokens[req.Token]
	if !exists {
		return nil, errors.New("invalid token")
	}

	// Delete the user
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
