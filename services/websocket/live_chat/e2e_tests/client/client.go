package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// Message represents the payload structure for the WebSocket messages.
type Message struct {
	Type     string `json:"type"`     // "public" or "private"
	Username string `json:"username"` // Sender's username
	Message  string `json:"message"`  // Message content
	To       string `json:"to"`       // Recipient username (for private messages)
}

// Client represents the WebSocket connection and its state.
type Client struct {
	Conn       *websocket.Conn
	MessageLog []map[string]interface{}
	Lock       sync.Mutex
}

// Globals
var client *Client

// Initialize WebSocket connection
// Initialize WebSocket connection
func connectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the WebSocket server URL from the request body
	var body struct {
		WebSocketServerURL string `json:"websocket_server_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Establish a WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(body.WebSocketServerURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to WebSocket server: %v", err), http.StatusInternalServerError)
		return
	}

	client = &Client{
		Conn:       conn,
		MessageLog: []map[string]interface{}{}, // Initialize as a slice of map[string]interface{}
	}

	// Start listening for incoming messages
	go listenForMessages()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Connected successfully"})
}

// Send a message to the WebSocket server
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the message payload
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	log.Printf("Parsed message payload: %+v", msg) // Add logging

	client.Lock.Lock()
	defer client.Lock.Unlock()

	// Send the message to the WebSocket server
	if err := client.Conn.WriteJSON(msg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
}

// Read the latest message received by the client
func readMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	client.Lock.Lock()
	defer client.Lock.Unlock()

	if len(client.MessageLog) == 0 {
		http.Error(w, "No messages received yet", http.StatusNotFound)
		return
	}

	// Return the latest message from the message log
	latestMessage := client.MessageLog[len(client.MessageLog)-1]
	response, err := json.Marshal(latestMessage)
	if err != nil {
		http.Error(w, "Failed to serialize message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Disconnect from the WebSocket server
func disconnectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	client.Lock.Lock()
	defer client.Lock.Unlock()

	// Close the WebSocket connection
	if err := client.Conn.Close(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to disconnect: %v", err), http.StatusInternalServerError)
		return
	}

	client.Conn = nil
	client.MessageLog = nil

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Disconnected successfully"})
}

// Listen for incoming messages from the WebSocket server
func listenForMessages() {
	for {
		var msg map[string]interface{}
		if err := client.Conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		client.Lock.Lock()
		client.MessageLog = append(client.MessageLog, msg)
		client.Lock.Unlock()

		log.Printf("Received message: %+v", msg)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}

	port := 0
	if _, err := fmt.Sscanf(os.Args[1], "%d", &port); err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/send", sendMessageHandler)
	http.HandleFunc("/read", readMessageHandler)
	http.HandleFunc("/disconnect", disconnectHandler)

	log.Printf("RESTful client API server is running on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
