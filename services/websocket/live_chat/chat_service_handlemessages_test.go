package live_chat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure MockWebSocketConn satisfies the WebSocketConn interface
var _ WebSocketConn = (*MockWebSocketConn)(nil)

func TestHandleMessage(t *testing.T) {
	service := NewChatService(5)

	t.Run("Public Message", func(t *testing.T) {
		mockConn1 := &MockWebSocketConn{}
		mockConn2 := &MockWebSocketConn{}

		// Add users
		service.users["user1"] = mockConn1
		service.users["user2"] = mockConn2

		msg := ChatMessage{
			Type:     "chat",
			Username: "user1",
			Message:  "Hello everyone!",
		}

		err := service.HandleMessage(msg, "user1")
		require.NoError(t, err)

		// Validate message broadcast
		require.Equal(t, msg, mockConn2.LastMessage) // Other users should receive the message
		require.Nil(t, mockConn1.LastMessage)        // Sender should not receive their own message
	})

	t.Run("Private Message", func(t *testing.T) {
		t.Run("Send Private Message to Recipient", func(t *testing.T) {
			mockConn1 := &MockWebSocketConn{}
			mockConn2 := &MockWebSocketConn{}

			// Add users
			service.users["user1"] = mockConn1
			service.users["user2"] = mockConn2

			msg := ChatMessage{
				Type:     "private",
				Username: "user1",
				Message:  "Hello user2!",
				To:       "user2",
			}

			err := service.HandleMessage(msg, "user1")
			require.NoError(t, err)

			// Validate private message delivery
			require.Nil(t, mockConn1.LastMessage)        // Sender shouldn't receive their own private message
			require.Equal(t, msg, mockConn2.LastMessage) // Recipient should receive the private message
		})

		t.Run("Recipient Not Found", func(t *testing.T) {
			mockConn1 := &MockWebSocketConn{}

			// Add sender only
			service.users["user1"] = mockConn1

			msg := ChatMessage{
				Type:     "private",
				Username: "user1",
				Message:  "Hello user3!",
				To:       "user3", // Non-existent recipient
			}

			err := service.HandleMessage(msg, "user1")
			require.Error(t, err)
			require.Equal(t, "recipient does not exist or is not online", err.Error()) // Updated expected error message
		})

	})

	t.Run("Message Validation", func(t *testing.T) {
		t.Run("Empty Message", func(t *testing.T) {
			mockConn1 := &MockWebSocketConn{}

			// Add sender
			service.users["user1"] = mockConn1

			msg := ChatMessage{
				Type:     "chat",
				Username: "user1",
				Message:  "", // Empty message
			}

			err := service.HandleMessage(msg, "user1")
			require.Error(t, err)
			require.Equal(t, "message cannot be empty", err.Error())
		})

		t.Run("Empty Username", func(t *testing.T) {
			msg := ChatMessage{
				Type:     "chat",
				Message:  "Hello everyone!",
				Username: "", // Empty username
			}

			err := service.HandleMessage(msg, "")
			require.Error(t, err)
			require.Equal(t, "username cannot be empty", err.Error())
		})

		t.Run("Message Exceeding Maximum Length", func(t *testing.T) {
			longMessage := strings.Repeat("A", 501) // 501 characters
			msg := ChatMessage{
				Type:     "chat",
				Username: "user1",
				Message:  longMessage,
			}

			err := service.HandleMessage(msg, "user1")
			require.Error(t, err)
			require.Equal(t, "message exceeds maximum length of 500 characters", err.Error())
		})

		t.Run("Invalid Message Type", func(t *testing.T) {
			mockConn1 := &MockWebSocketConn{}

			// Add sender
			service.users["user1"] = mockConn1

			msg := ChatMessage{
				Type:     "unknown",
				Username: "user1",
				Message:  "This message type is invalid",
			}

			err := service.HandleMessage(msg, "user1")
			require.Error(t, err)
			require.Equal(t, "invalid message type", err.Error())
		})
	})

}
