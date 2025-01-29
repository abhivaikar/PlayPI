package grpc

import (
	"context"
	"errors"
	"log"
	"net"

	pb "github.com/abhivaikar/playpi/services/grpc/inventory_management/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedInventoryServiceServer
	inventory []pb.Item
}

// Load mock data into the inventory
func (s *server) loadMockData() {
	s.inventory = []pb.Item{
		{Id: 1, Name: "Laptop", Description: "A high-performance laptop", Price: 999.99, Quantity: 10},
		{Id: 2, Name: "Smartphone", Description: "A powerful smartphone", Price: 699.99, Quantity: 20},
		{Id: 3, Name: "Headphones", Description: "Noise-cancelling headphones", Price: 199.99, Quantity: 15},
		{Id: 4, Name: "Monitor", Description: "A 24-inch full HD monitor", Price: 149.99, Quantity: 8},
		{Id: 5, Name: "Keyboard", Description: "Mechanical keyboard with RGB lighting", Price: 49.99, Quantity: 50},
		{Id: 6, Name: "Mouse", Description: "Wireless optical mouse", Price: 29.99, Quantity: 30},
		{Id: 7, Name: "Printer", Description: "All-in-one color printer", Price: 299.99, Quantity: 5},
		{Id: 8, Name: "Webcam", Description: "1080p HD webcam", Price: 89.99, Quantity: 12},
		{Id: 9, Name: "External Hard Drive", Description: "1TB portable external hard drive", Price: 79.99, Quantity: 25},
		{Id: 10, Name: "USB Hub", Description: "4-port USB 3.0 hub", Price: 19.99, Quantity: 40},
	}
}

// GetItem fetches an item by ID
func (s *server) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	for i := range s.inventory {
		if s.inventory[i].Id == req.Id {
			return &pb.GetItemResponse{Item: &s.inventory[i]}, nil
		}
	}
	return nil, errors.New("item not found")
}

// ListItems returns all items in the inventory
func (s *server) ListItems(ctx context.Context, req *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	var items []*pb.Item
	for i := range s.inventory {
		items = append(items, &s.inventory[i])
	}
	return &pb.ListItemsResponse{Items: items}, nil
}

// AddItem adds a new item to the inventory
func (s *server) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	if len(req.Name) < 3 || len(req.Name) > 50 {
		return nil, errors.New("name must be between 3 and 50 characters")
	}
	if len(req.Description) > 200 {
		return nil, errors.New("description cannot exceed 200 characters")
	}
	if req.Price < 0 || req.Price > 10000 {
		return nil, errors.New("price must be a positive number not exceeding 10,000")
	}
	if req.Quantity < 0 {
		return nil, errors.New("quantity must be at least 0")
	}

	newItem := pb.Item{
		Id:          int32(len(s.inventory) + 1),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}
	s.inventory = append(s.inventory, newItem)
	return &pb.AddItemResponse{Item: &newItem}, nil
}

// UpdateItem updates an existing item by ID
func (s *server) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	for i := range s.inventory {
		if s.inventory[i].Id == req.Id {
			if req.Name != "" && (len(req.Name) < 3 || len(req.Name) > 50) {
				return nil, errors.New("name must be between 3 and 50 characters")
			}
			if req.Description != "" && len(req.Description) > 200 {
				return nil, errors.New("description cannot exceed 200 characters")
			}
			if req.Price < 0 || req.Price > 10000 {
				return nil, errors.New("price must be a positive number not exceeding 10,000")
			}
			if req.Quantity < 0 {
				return nil, errors.New("quantity must be at least 0")
			}

			// Update the item
			if req.Name != "" {
				s.inventory[i].Name = req.Name
			}
			if req.Description != "" {
				s.inventory[i].Description = req.Description
			}
			if req.Price != 0 {
				s.inventory[i].Price = req.Price
			}
			if req.Quantity != 0 {
				s.inventory[i].Quantity = req.Quantity
			}
			return &pb.UpdateItemResponse{Item: &s.inventory[i]}, nil
		}
	}
	return nil, errors.New("item not found")
}

// DeleteItem removes an item by ID
func (s *server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	for i := range s.inventory {
		if s.inventory[i].Id == req.Id {
			if s.inventory[i].Quantity > 0 {
				return nil, errors.New("cannot delete an item with stock remaining")
			}
			s.inventory = append(s.inventory[:i], s.inventory[i+1:]...)
			return &pb.DeleteItemResponse{Success: true}, nil
		}
	}
	return nil, errors.New("item not found")
}

func StartServer() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := &server{}
	s.loadMockData() // Load mock data into the inventory

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, s)

	log.Println("gRPC server is running on port 8082")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
