package live_chat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure MockWebSocketConn satisfies the WebSocketConn interface
var _ WebSocketConn = (*MockWebSocketConn)(nil)

func TestRemoveUser(t *testing.T) {
	service := NewChatService(5)
	mockConn := &MockWebSocketConn{}

	// Add a user
	service.users["user1"] = mockConn

	t.Run("Remove Valid User", func(t *testing.T) {
		service.RemoveUser("user1")
		_, exists := service.users["user1"]
		require.False(t, exists, "User should be removed from the map")
	})

	t.Run("Broadcast Leave Message", func(t *testing.T) {
		mockConn := &MockWebSocketConn{}
		service.broadcast = make(chan ChatMessage, 1) // Create a buffered channel
		service.users["user1"] = mockConn

		service.RemoveUser("user1")

		select {
		case msg := <-service.broadcast:
			require.Equal(t, "user1 has left the chat.", msg.Message)
			require.Equal(t, "System", msg.Username)
		default:
			t.Error("Leave message was not broadcasted")
		}
	})
}
