package live_chat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure MockWebSocketConn satisfies the WebSocketConn interface
var _ WebSocketConn = (*MockWebSocketConn)(nil)

func TestRegisterUserWithUsername(t *testing.T) {
	service := NewChatService(5)

	t.Run("Register User", func(t *testing.T) {
		mockConn := &MockWebSocketConn{}

		username, err := service.RegisterUserWithUsername(mockConn)
		require.NoError(t, err)
		require.NotEmpty(t, username)
		require.NotNil(t, service.users[username])
	})

	t.Run("Max Clients Reached", func(t *testing.T) {
		service := NewChatService(1) // Set max clients to 1
		mockConn1 := &MockWebSocketConn{}
		mockConn2 := &MockWebSocketConn{}

		service.RegisterUserWithUsername(mockConn1) // Register first user
		_, err := service.RegisterUserWithUsername(mockConn2)
		require.Error(t, err)
		require.Equal(t, "server is full, please try again later", err.Error())
	})
}
