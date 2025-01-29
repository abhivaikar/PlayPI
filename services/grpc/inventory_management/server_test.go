package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/abhivaikar/playpi/services/grpc/inventory_management/pb"
)

func setupTestServer() *server {
	s := &server{}
	s.loadMockData() // Load mock data into the inventory
	return s
}

func TestGetItem(t *testing.T) {
	s := setupTestServer()

	t.Run("Retrieve Existing Item", func(t *testing.T) {
		resp, err := s.GetItem(context.Background(), &pb.GetItemRequest{Id: 1})
		require.NoError(t, err)
		require.Equal(t, int32(1), resp.Item.Id)
		require.Equal(t, "Laptop", resp.Item.Name)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		_, err := s.GetItem(context.Background(), &pb.GetItemRequest{Id: 999})
		require.Error(t, err)
		require.Equal(t, "item not found", err.Error())
	})
}

func TestListItems(t *testing.T) {
	s := setupTestServer()

	t.Run("List All Items", func(t *testing.T) {
		resp, err := s.ListItems(context.Background(), &pb.ListItemsRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Items, 10) // Mock data has 10 items
	})
}

func TestAddItem(t *testing.T) {
	s := setupTestServer()

	t.Run("Successful Addition", func(t *testing.T) {
		resp, err := s.AddItem(context.Background(), &pb.AddItemRequest{
			Name:        "New Item",
			Description: "A new item description",
			Price:       100,
			Quantity:    5,
		})
		require.NoError(t, err)
		require.Equal(t, "New Item", resp.Item.Name)
		require.Equal(t, 11, len(s.inventory)) // Ensure the item was added
	})

	t.Run("Validation Error - Invalid Name", func(t *testing.T) {
		_, err := s.AddItem(context.Background(), &pb.AddItemRequest{
			Name:        "A", // Invalid name
			Description: "Valid Description",
			Price:       50,
			Quantity:    10,
		})
		require.Error(t, err)
		require.Equal(t, "name must be between 3 and 50 characters", err.Error())
	})

	t.Run("Validation Error - Invalid Description", func(t *testing.T) {
		_, err := s.AddItem(context.Background(), &pb.AddItemRequest{
			Name: "Valid Name",
			Description: "This description exceeds the character limit of 200 characters. " +
				"This description exceeds the character limit of 200 characters. " +
				"This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters. ",
			Price:    50,
			Quantity: 10,
		})
		require.Error(t, err)
		require.Equal(t, "description cannot exceed 200 characters", err.Error())
	})

	t.Run("Validation Error - Invalid Price", func(t *testing.T) {
		_, err := s.AddItem(context.Background(), &pb.AddItemRequest{
			Name:        "Valid Name",
			Description: "Valid Description",
			Price:       -10, // Invalid price
			Quantity:    10,
		})
		require.Error(t, err)
		require.Equal(t, "price must be a positive number not exceeding 10,000", err.Error())
	})

	t.Run("Validation Error - Invalid Quantity", func(t *testing.T) {
		_, err := s.AddItem(context.Background(), &pb.AddItemRequest{
			Name:        "Valid Name",
			Description: "Valid Description",
			Price:       50,
			Quantity:    -1, // Invalid quantity
		})
		require.Error(t, err)
		require.Equal(t, "quantity must be at least 0", err.Error())
	})
}

func TestUpdateItem(t *testing.T) {
	s := setupTestServer()

	t.Run("Successful Update", func(t *testing.T) {
		resp, err := s.UpdateItem(context.Background(), &pb.UpdateItemRequest{
			Id:          1,
			Name:        "Updated Laptop",
			Description: "Updated Description",
			Price:       1200,
			Quantity:    0,
		})
		require.NoError(t, err)
		require.Equal(t, "Updated Laptop", resp.Item.Name)
		require.Equal(t, "Updated Description", resp.Item.Description)
		require.Equal(t, float32(1200), resp.Item.Price)
	})

	t.Run("Validation Error - Invalid Name", func(t *testing.T) {
		_, err := s.UpdateItem(context.Background(), &pb.UpdateItemRequest{
			Id:          1,
			Name:        "A", // Invalid name
			Description: "Valid Description",
			Price:       50,
			Quantity:    10,
		})
		require.Error(t, err)
		require.Equal(t, "name must be between 3 and 50 characters", err.Error())
	})

	t.Run("Validation Error - Invalid Description", func(t *testing.T) {
		_, err := s.UpdateItem(context.Background(), &pb.UpdateItemRequest{
			Id:   1,
			Name: "Valid Name",
			Description: "This description exceeds the character limit of 200 characters. " +
				"This description exceeds the character limit of 200 characters. " +
				"This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters. ",
			Price:    50,
			Quantity: 10,
		})
		require.Error(t, err)
		require.Equal(t, "description cannot exceed 200 characters", err.Error())
	})

	t.Run("Validation Error - Invalid Price", func(t *testing.T) {
		_, err := s.UpdateItem(context.Background(), &pb.UpdateItemRequest{
			Id:          1,
			Name:        "Valid Name",
			Description: "Valid Description",
			Price:       15000, // Invalid price
			Quantity:    10,
		})
		require.Error(t, err)
		require.Equal(t, "price must be a positive number not exceeding 10,000", err.Error())
	})

	t.Run("Validation Error - Invalid Quantity", func(t *testing.T) {
		_, err := s.UpdateItem(context.Background(), &pb.UpdateItemRequest{
			Id:          1,
			Name:        "Valid Name",
			Description: "Valid Description",
			Price:       50,
			Quantity:    -1, // Invalid quantity
		})
		require.Error(t, err)
		require.Equal(t, "quantity must be at least 0", err.Error())
	})

	t.Run("Item Not Found", func(t *testing.T) {
		_, err := s.UpdateItem(context.Background(), &pb.UpdateItemRequest{
			Id:          999, // Non-existent item ID
			Name:        "Valid Name",
			Description: "Valid Description",
			Price:       50,
			Quantity:    10,
		})
		require.Error(t, err)
		require.Equal(t, "item not found", err.Error())
	})
}

func TestDeleteItem(t *testing.T) {
	s := setupTestServer()

	t.Run("Successful Deletion", func(t *testing.T) {
		s.inventory[0].Quantity = 0 // Ensure quantity is 0
		resp, err := s.DeleteItem(context.Background(), &pb.DeleteItemRequest{Id: 1})
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, 9, len(s.inventory)) // Ensure the item was removed
	})

	t.Run("Cannot Delete Item with Stock Remaining", func(t *testing.T) {
		_, err := s.DeleteItem(context.Background(), &pb.DeleteItemRequest{Id: 2})
		require.Error(t, err)
		require.Equal(t, "cannot delete an item with stock remaining", err.Error())
	})

	t.Run("Item Not Found", func(t *testing.T) {
		_, err := s.DeleteItem(context.Background(), &pb.DeleteItemRequest{Id: 999})
		require.Error(t, err)
		require.Equal(t, "item not found", err.Error())
	})
}
