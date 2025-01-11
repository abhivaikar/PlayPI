package live_chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ChatMessage represents a single chat message
type ChatMessage struct {
	Type     string `json:"type"` // "chat", "system", "private"
	Username string `json:"username"`
	Message  string `json:"message"`
	To       string `json:"to,omitempty"` // For private messages
}

// WebSocketServer manages chat connections and broadcasts
type WebSocketServer struct {
	clients   map[*websocket.Conn]bool
	users     map[string]*websocket.Conn // Map usernames to connections
	broadcast chan ChatMessage
	upgrader  websocket.Upgrader
}

// NewWebSocketServer initializes the WebSocketServer
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:   make(map[*websocket.Conn]bool),
		users:     make(map[string]*websocket.Conn),
		broadcast: make(chan ChatMessage),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections
			},
		},
	}
}

// HandleConnections manages incoming WebSocket connections
func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to WebSocket
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket connection failed: %v", err)
		return
	}
	defer conn.Close()

	// Prompt user to send their username as the first message
	var username string
	if err := conn.ReadJSON(&username); err != nil {
		log.Printf("Failed to read username: %v", err)
		return
	}

	// Track the connection and username
	server.clients[conn] = true
	server.users[username] = conn
	defer delete(server.clients, conn)   // Remove the connection on disconnect
	defer delete(server.users, username) // Remove the username on disconnect

	log.Printf("User %s connected", username)

	// Broadcast join notification
	joinMessage := ChatMessage{
		Type:     "system",
		Username: "System",
		Message:  username + " has joined the chat.",
	}
	server.broadcast <- joinMessage // Broadcast the join message to all users

	// Handle messages from this connection
	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from %s: %v", username, err)

			// Notify others about the user leaving
			leaveMessage := ChatMessage{
				Type:     "system",
				Username: "System",
				Message:  username + " has left the chat.",
			}
			server.broadcast <- leaveMessage // Broadcast the leave message
			break
		}

		// Handle the received message (broadcast or private)
		server.broadcast <- msg
	}
}

// HandleMessages listens for incoming messages and broadcasts them
func (server *WebSocketServer) HandleMessages() {
	for {
		msg := <-server.broadcast
		if msg.Type == "private" && msg.To != "" {
			// Send private message to the intended recipient
			if recipientConn, exists := server.users[msg.To]; exists {
				err := recipientConn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending private message: %v", err)
					recipientConn.Close()
					delete(server.clients, recipientConn)
					delete(server.users, msg.To)
				}
			} else {
				log.Printf("User %s not found for private message", msg.To)
			}
		} else {
			// Broadcast message to all clients
			for client := range server.clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Error broadcasting message: %v", err)
					client.Close()
					delete(server.clients, client)
				}
			}
		}
	}
}

// StartServer starts the WebSocket server
func (server *WebSocketServer) StartServer() {
	http.HandleFunc("/ws", server.HandleConnections)

	go server.HandleMessages()

	log.Println("WebSocket Live Chat server is running on port 8086")
	err := http.ListenAndServe(":8086", nil)
	if err != nil {
		log.Fatalf("WebSocket server failed: %v", err)
	}
}
