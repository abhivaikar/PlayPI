/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	graphqlInventory "github.com/abhivaikar/playpi/services/graphql/inventory_management"
	grpcInventory "github.com/abhivaikar/playpi/services/grpc/inventory_management"
	grpcUserRegistration "github.com/abhivaikar/playpi/services/grpc/user_registration"
	restfulInventory "github.com/abhivaikar/playpi/services/restful/inventory_management"
	restfulTaskManagement "github.com/abhivaikar/playpi/services/restful/task_management"
	websocketInventory "github.com/abhivaikar/playpi/services/websocket/inventory_management"
	websocketLiveChat "github.com/abhivaikar/playpi/services/websocket/live_chat"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the PlayPI playground",
	Long: `Start the PlayPI playground. Choose an API type to test:
- RESTful API
- GraphQL API
- gRPC API
- WebSocket API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to PlayPI!")
		fmt.Println("Please choose an API type:")
		fmt.Println("1. RESTful API - Inventory Management")
		fmt.Println("2. GraphQL API - Inventory Management")
		fmt.Println("3. gRPC API - Inventory Management")
		fmt.Println("4. WebSocket API - Inventory Management")
		fmt.Println("5. RESTful API - Task Management")
		fmt.Println("6. gRPC API - User Registration and Sign In")
		fmt.Println("7. WebSocket API - Live Chat")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("Starting RESTful API Playground...")
			restfulInventory.StartServer()
		case 2:
			fmt.Println("Starting GraphQL API Playground...")
			graphqlInventory.StartServer()
		case 3:
			fmt.Println("Starting gRPC API Playground...")
			grpcInventory.StartServer()
		case 4:
			fmt.Println("Starting WebSocket Playground...")
			wsServer := websocketInventory.NewWebSocketServer()
			wsServer.StartServer()
		case 5:
			fmt.Println("Starting RESTful API PLayground...")
			restfulTaskManagement.StartServer()
		case 6:
			fmt.Println("Starting gRPC API Playground...")
			grpcUserRegistration.StartServer()
		case 7:
			fmt.Println("Starting WebSocket Playground...")
			wsLiveChatServer := websocketLiveChat.NewWebSocketServer()
			wsLiveChatServer.StartServer()
		default:
			fmt.Println("Feature not yet implemented.")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
