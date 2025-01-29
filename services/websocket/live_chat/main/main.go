package main

import (
	"fmt"

	websocketLiveChat "github.com/abhivaikar/playpi/services/websocket/live_chat"
)

func main() {
	fmt.Println("Starting WebSocket Playground for live chat...")
	wsLiveChatServer := websocketLiveChat.NewWebSocketServer()
	wsLiveChatServer.StartServer()
}
