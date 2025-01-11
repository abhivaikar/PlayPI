package grpc

import (
	"context"
	"log"
	"net"

	pb "github.com/abhivaikar/playpi/services/grpc/inventory_management/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedInventoryServiceServer
	inventory []pb.Item
}

// GetItem fetches an item by ID
func (s *server) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	for i := range s.inventory {
		if s.inventory[i].Id == req.Id {
			return &pb.GetItemResponse{Item: &s.inventory[i]}, nil
		}
	}
	return nil, nil
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
	return nil, nil
}

// DeleteItem removes an item by ID
func (s *server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	for i := range s.inventory {
		if s.inventory[i].Id == req.Id {
			s.inventory = append(s.inventory[:i], s.inventory[i+1:]...)
			return &pb.DeleteItemResponse{Success: true}, nil
		}
	}
	return &pb.DeleteItemResponse{Success: false}, nil
}

func StartServer() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, &server{})

	log.Println("gRPC server is running on port 8082")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
