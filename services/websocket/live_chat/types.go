package live_chat

import "sync"

// WebSocketConn defines an interface for WebSocket connections.
// This is useful for mocking WebSocket behavior in unit tests.
type WebSocketConn interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

// ChatMessage represents a message in the chat.
type ChatMessage struct {
	Type     string `json:"type"`     // "chat" or "private"
	Username string `json:"username"` // Sender's username
	Message  string `json:"message"`  // Message content
	To       string `json:"to"`       // Recipient username for private messages
}

// ChatService manages users and messages in the chat.
type ChatService struct {
	users      map[string]WebSocketConn
	mu         sync.Mutex
	broadcast  chan ChatMessage
	maxClients int
	writeMu    sync.Mutex // Mutex to protect write operations
}
