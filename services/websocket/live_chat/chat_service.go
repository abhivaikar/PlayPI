package live_chat

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// NewChatService initializes a new ChatService.
func NewChatService(maxClients int) *ChatService {
	service := &ChatService{
		users:      make(map[string]WebSocketConn),
		broadcast:  make(chan ChatMessage, 100), // Buffered channel for broadcast
		maxClients: maxClients,
	}
	service.StartBroadcastProcessor()
	return service
}

func (s *ChatService) StartBroadcastProcessor() {
	go func() {
		for msg := range s.broadcast {
			s.mu.Lock()
			for username, conn := range s.users {
				if username != msg.Username { // Exclude the sender
					s.writeMu.Lock()
					if err := conn.WriteJSON(msg); err != nil {
						fmt.Printf("Failed to send message to user %s: %v\n", username, err)
					} else {
						fmt.Printf("Message sent to user %s: %s\n", username, msg.Message)
					}
					s.writeMu.Unlock()
				}
			}
			s.mu.Unlock()
		}
	}()
}

func generateRandomUsername() string {
	// List of template names
	templateNames := []string{"Apple", "Banana", "Cherry", "Date", "Elderberry", "Fig", "Grape", "Honeydew"}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	templateName := templateNames[r.Intn(len(templateNames))]
	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%s%d", templateName, timestamp)
}

// RegisterUserWithUsername registers a new user with a random username to the chat service.
func (s *ChatService) RegisterUserWithUsername(conn WebSocketConn) (string, error) {
	username := generateRandomUsername()

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.users) >= s.maxClients {
		s.writeMu.Lock()
		conn.WriteJSON(map[string]string{"error": "server is full, please try again later"})
		s.writeMu.Unlock()
		return "", errors.New("server is full, please try again later")
	}

	s.users[username] = conn
	s.BroadcastSystemMessage(username+" has joined the chat.", username)
	return username, nil
}

// RemoveUser removes a user from the chat service and broadcasts a leave message.
func (s *ChatService) RemoveUser(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, exists := s.users[username]; exists {
		conn.Close()
		delete(s.users, username)

		// Broadcast leave message
		s.broadcast <- ChatMessage{
			Type:     "system",
			Username: "System",
			Message:  fmt.Sprintf("%s has left the chat.", username),
		}
	}
}

// BroadcastSystemMessage broadcasts a system message to all users except the sender.
func (s *ChatService) BroadcastSystemMessage(message string, sender string) {
	msg := ChatMessage{
		Type:     "system",
		Username: sender,
		Message:  message,
	}

	s.broadcast <- msg
}

// HandleMessage validates and processes incoming chat messages.
func (s *ChatService) HandleMessage(msg ChatMessage, sender string) error {
	// Validate the message
	if err := s.validateMessage(msg, sender); err != nil {
		return err
	}

	// Handle private messages
	if msg.Type == "private" {
		recipientConn, exists := s.users[msg.To]
		if !exists {
			return errors.New("recipient does not exist or is not online")
		}
		s.writeMu.Lock()
		err := recipientConn.WriteJSON(msg)
		s.writeMu.Unlock()
		return err
	}

	// Handle public messages
	for username, conn := range s.users {
		if username != sender { // Exclude the sender
			s.writeMu.Lock()
			err := conn.WriteJSON(msg)
			s.writeMu.Unlock()
			if err != nil {
				fmt.Printf("Failed to send message to user %s: %v\n", username, err)
			}
		}
	}
	return nil
}

// validateMessage validates a chat message based on its type.
func (s *ChatService) validateMessage(msg ChatMessage, sender string) error {
	if strings.TrimSpace(msg.Message) == "" {
		return errors.New("message cannot be empty")
	}

	if strings.TrimSpace(msg.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(msg.Message) > 500 {
		return errors.New("message exceeds maximum length of 500 characters")
	}
	if msg.Type != "chat" && msg.Type != "private" {
		return errors.New("invalid message type")
	}
	if msg.Username != sender {
		return errors.New("invalid username: you are not registered as " + msg.Username)
	}
	if msg.Type == "private" {
		if msg.To == "" {
			return errors.New("recipient cannot be empty for private messages")
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, exists := s.users[msg.To]; !exists {
			return errors.New("recipient does not exist or is not online")
		}
	}
	msg.Username = sender
	return nil
}
