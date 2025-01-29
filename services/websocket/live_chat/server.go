package live_chat

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	service *ChatService
	server  *http.Server
}

func NewWebSocketServer() *WebSocketServer {
	server := &WebSocketServer{
		service: NewChatService(5), // Initialize with a max of 5 users
	}
	server.service.StartBroadcastProcessor()
	return server
}

func (s *WebSocketServer) StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.HandleConnections)

	s.server = &http.Server{
		Addr:    ":8086",
		Handler: mux,
	}

	log.Println("WebSocket server started on port 8086...")
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func (s *WebSocketServer) StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("WebSocket server stopped gracefully")
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// Register the user with a random username
	username, err := s.service.RegisterUserWithUsername(conn)
	if err != nil {
		log.Println("Error registering user:", err)
		conn.Close()
		return
	}
	defer func() {
		s.service.RemoveUser(username)
		conn.Close()
	}()

	// Inform the client of their assigned username
	if err := s.sendJSON(conn, map[string]string{"message": "You have connected as " + username}); err != nil {
		log.Println("Error sending welcome message:", err)
		return
	}

	// Listen for incoming messages
	for {
		var msg ChatMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("Connection closed:", err)
			} else {
				log.Println("Error reading message:", err)
			}
			break // Exit the loop for any error
		}

		// Delegate message handling to ChatService
		if err := s.service.HandleMessage(msg, username); err != nil {
			if sendErr := s.sendJSON(conn, map[string]string{"error": err.Error()}); sendErr != nil {
				log.Println("Error sending error message:", sendErr)
			}
		}
	}
}

// sendJSON safely sends a JSON message over a WebSocket connection
func (s *WebSocketServer) sendJSON(conn *websocket.Conn, v interface{}) error {
	s.service.writeMu.Lock()
	defer s.service.writeMu.Unlock()
	return conn.WriteJSON(v)
}
