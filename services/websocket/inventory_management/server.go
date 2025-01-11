package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	upgrader  websocket.Upgrader
	mu        sync.Mutex // Protect inventory
	inventory []InventoryItem
}

type InventoryItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections
			},
		},
		inventory: []InventoryItem{},
	}
}

func (server *WebSocketServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket connection failed: %v", err)
		return
	}
	defer conn.Close()

	server.clients[conn] = true
	log.Println("New client connected")

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(server.clients, conn)
			break
		}
		server.handleMessage(conn, msg)
	}
}

func (server *WebSocketServer) handleMessages() {
	for {
		msg := <-server.broadcast
		for client := range server.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				client.Close()
				delete(server.clients, client)
			}
		}
	}
}

func (server *WebSocketServer) handleMessage(conn *websocket.Conn, msg Message) {
	switch msg.Type {
	case "get_all_items":
		server.mu.Lock()
		items := server.inventory
		server.mu.Unlock()
		// Broadcast all items to all clients
		server.broadcast <- Message{Type: "all_items", Payload: items}

	case "get_item":
		id := int(msg.Payload.(map[string]interface{})["id"].(float64))
		server.mu.Lock()
		var item *InventoryItem
		for _, inv := range server.inventory {
			if inv.ID == id {
				item = &inv
				break
			}
		}
		server.mu.Unlock()
		conn.WriteJSON(Message{Type: "item", Payload: item})

	case "add_item":
		newItem := msg.Payload.(map[string]interface{})
		server.mu.Lock()
		item := InventoryItem{
			ID:          len(server.inventory) + 1,
			Name:        newItem["name"].(string),
			Description: newItem["description"].(string),
			Price:       newItem["price"].(float64),
			Quantity:    int(newItem["quantity"].(float64)),
		}
		server.inventory = append(server.inventory, item)
		server.mu.Unlock()
		server.broadcast <- Message{Type: "item_added", Payload: item}

	case "update_item":
		data := msg.Payload.(map[string]interface{})
		id := int(data["id"].(float64))
		server.mu.Lock()
		for i, item := range server.inventory {
			if item.ID == id {
				if name, ok := data["name"].(string); ok {
					server.inventory[i].Name = name
				}
				if description, ok := data["description"].(string); ok {
					server.inventory[i].Description = description
				}
				if price, ok := data["price"].(float64); ok {
					server.inventory[i].Price = price
				}
				if quantity, ok := data["quantity"].(float64); ok {
					server.inventory[i].Quantity = int(quantity)
				}
				server.broadcast <- Message{Type: "item_updated", Payload: server.inventory[i]}
				break
			}
		}
		server.mu.Unlock()

	case "delete_item":
		id := int(msg.Payload.(map[string]interface{})["id"].(float64))
		server.mu.Lock()
		for i, item := range server.inventory {
			if item.ID == id {
				server.inventory = append(server.inventory[:i], server.inventory[i+1:]...)
				server.broadcast <- Message{Type: "item_deleted", Payload: item}
				break
			}
		}
		server.mu.Unlock()
	}
}

func (server *WebSocketServer) StartServer() {
	http.HandleFunc("/ws", server.handleConnections)

	go server.handleMessages()

	log.Println("WebSocket server is running on port 8083")
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		log.Fatalf("WebSocket server failed: %v", err)
	}
}
