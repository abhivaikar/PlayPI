package user_registration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	pb "github.com/abhivaikar/playpi/services/grpc/user_registration/pb"
)

func setupTestServer() *server {
	return NewServer()
}

func TestRegisterUser(t *testing.T) {
	s := setupTestServer()

	t.Run("Successful Registration", func(t *testing.T) {
		resp, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "testuser",
				Password: "securepassword",
				FullName: "Test User",
				Email:    "testuser@example.com",
				Phone:    "1234567890",
				Address:  "123 Test Street",
			},
		})
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "User registered successfully", resp.Message)
	})

	t.Run("Validation Error - Empty Username", func(t *testing.T) {
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "",
				Password: "securepassword",
				FullName: "Test User",
				Email:    "testuser@example.com",
				Phone:    "1234567890",
				Address:  "123 Test Street",
			},
		})
		require.Error(t, err)
		require.Equal(t, "username must be between 3 and 50 characters", err.Error())
	})

	t.Run("Validation Error - Invalid Password", func(t *testing.T) {
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "validusername",
				Password: "short", // Too short
				FullName: "Test User",
				Email:    "testuser@example.com",
				Phone:    "1234567890",
				Address:  "123 Test Street",
			},
		})
		require.Error(t, err)
		require.Equal(t, "password must be at least 8 characters long", err.Error())
	})

	t.Run("Validation Error - Invalid Email", func(t *testing.T) {
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "validusername",
				Password: "securepassword",
				FullName: "Test User",
				Email:    "invalid-email", // Invalid email format
				Phone:    "1234567890",
				Address:  "123 Test Street",
			},
		})
		require.Error(t, err)
		require.Equal(t, "invalid email format", err.Error())
	})

	t.Run("Validation Error - Invalid Phone Number", func(t *testing.T) {
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "validusername",
				Password: "securepassword",
				FullName: "Test User",
				Email:    "testuser@example.com",
				Phone:    "123ABC7890", // Non-numeric phone number
				Address:  "123 Test Street",
			},
		})
		require.Error(t, err)
		require.Equal(t, "invalid phone number format", err.Error())
	})

	t.Run("Validation Error - Address Too Long", func(t *testing.T) {
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "validusername",
				Password: "securepassword",
				FullName: "Test User",
				Email:    "testuser@example.com",
				Phone:    "1234567890",
				Address:  "This address is far too long and exceeds the limit of 100 characters. This is invalid and should not be accepted.",
			},
		})
		require.Error(t, err)
		require.Equal(t, "address cannot exceed 100 characters", err.Error())
	})

	t.Run("Validation Error - Missing Fullname", func(t *testing.T) {
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "validusername",
				Password: "securepassword",
				FullName: "", // Missing fullname
				Email:    "testuser@example.com",
				Phone:    "1234567890",
				Address:  "123 Test Street",
			},
		})
		require.Error(t, err)
		require.Equal(t, "fullname is required", err.Error())
	})

	t.Run("Duplicate Username", func(t *testing.T) {
		// First registration
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "duplicateuser",
				Password: "securepassword",
				FullName: "Test User",
				Email:    "duplicateuser@example.com",
				Phone:    "1234567890",
				Address:  "123 Test Street",
			},
		})
		require.NoError(t, err)

		// Second registration with the same username
		_, err = s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: "duplicateuser",
				Password: "securepassword",
				FullName: "Test User",
				Email:    "duplicateuser@example.com",
				Phone:    "1234567890",
				Address:  "123 Test Street",
			},
		})
		require.Error(t, err)
		require.Equal(t, "username already exists", err.Error())
	})
}

func TestSignIn(t *testing.T) {
	s := setupTestServer()

	// Pre-register a user for testing
	_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		User: &pb.User{
			Username: "signinuser",
			Password: "mypassword",
			FullName: "Sign In User",
			Email:    "signinuser@example.com",
			Phone:    "1234567890",
			Address:  "123 Sign In Street",
		},
	})
	require.NoError(t, err)

	t.Run("Successful Sign In", func(t *testing.T) {
		resp, err := s.SignIn(context.Background(), &pb.SignInRequest{
			Username: "signinuser",
			Password: "mypassword",
		})
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Signed in successfully", resp.Message)
		require.NotEmpty(t, resp.Token)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		_, err := s.SignIn(context.Background(), &pb.SignInRequest{
			Username: "signinuser",
			Password: "wrongpassword",
		})
		require.Error(t, err)
		require.Equal(t, "invalid username or password", err.Error())
	})
}

func TestGetProfile(t *testing.T) {
	s := setupTestServer()

	// Pre-register a user and generate a token
	_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		User: &pb.User{
			Username: "profileuser",
			Password: "mypassword",
			FullName: "Profile User",
			Email:    "profileuser@example.com",
			Phone:    "1234567890",
			Address:  "123 Profile Street",
		},
	})
	require.NoError(t, err)

	signInResp, err := s.SignIn(context.Background(), &pb.SignInRequest{
		Username: "profileuser",
		Password: "mypassword",
	})
	require.NoError(t, err)
	token := signInResp.Token

	t.Run("Successful GetProfile", func(t *testing.T) {
		resp, err := s.GetProfile(context.Background(), &pb.GetProfileRequest{
			Token: token,
		})
		require.NoError(t, err)
		require.Equal(t, "profileuser", resp.User.Username)
		require.Equal(t, "Profile User", resp.User.FullName)
		require.Equal(t, "profileuser@example.com", resp.User.Email)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := s.GetProfile(context.Background(), &pb.GetProfileRequest{
			Token: "invalidtoken",
		})
		require.Error(t, err)
		require.Equal(t, "invalid token", err.Error())
	})
}

func TestDeleteAccount(t *testing.T) {
	s := setupTestServer()

	t.Run("Successful Deletion", func(t *testing.T) {
		// Register a unique account
		username := "deleteuser_" + generateUniqueID()
		_, err := s.RegisterUser(context.Background(), &pb.RegisterUserRequest{
			User: &pb.User{
				Username: username,
				Password: "mypassword",
				FullName: "Delete User",
				Email:    username + "@example.com",
				Phone:    "1234567890",
				Address:  "123 Delete Street",
			},
		})
		require.NoError(t, err)

		// Sign in and get the token
		signInResp, err := s.SignIn(context.Background(), &pb.SignInRequest{
			Username: username,
			Password: "mypassword",
		})
		require.NoError(t, err)
		token := signInResp.Token

		// Delete the account
		resp, err := s.DeleteAccount(context.Background(), &pb.DeleteAccountRequest{
			Token: token,
		})
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Account deleted successfully", resp.Message)

		// Verify the account no longer exists
		_, err = s.SignIn(context.Background(), &pb.SignInRequest{
			Username: username,
			Password: "mypassword",
		})
		require.Error(t, err)
		require.Equal(t, "invalid username or password", err.Error())
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Attempt to delete an account with an invalid token
		_, err := s.DeleteAccount(context.Background(), &pb.DeleteAccountRequest{
			Token: "invalidtoken",
		})
		require.Error(t, err)
		require.Equal(t, "invalid token", err.Error())
	})
}

func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
