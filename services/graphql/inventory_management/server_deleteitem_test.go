package graphql

import (
	"net/http/httptest"
	"testing"
)

func TestDeleteItem(t *testing.T) {
	ts := httptest.NewServer(setupGraphQLServer())
	defer ts.Close()

	t.Run("Successful Deletion", func(t *testing.T) {
		// Set up an item with quantity = 0 for this test
		inventory = append(inventory, map[string]interface{}{
			"id":          99,
			"name":        "Item to Delete",
			"description": "This item will be deleted",
			"price":       50.0,
			"quantity":    0,
		})

		query := `
			mutation {
				deleteItem(id: 99)
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		data, ok := result["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to extract 'data' from response: %v", result)
		}

		success, ok := data["deleteItem"].(bool)
		if !ok || !success {
			t.Errorf("Expected deleteItem to return true, got %v", success)
		}
	})

	t.Run("Item Not Found", func(t *testing.T) {
		query := `
			mutation {
				deleteItem(id: 999)
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		errors, ok := result["errors"].([]interface{})
		if !ok || len(errors) == 0 {
			t.Fatalf("Expected error, got none")
		}

		errorMessage := errors[0].(map[string]interface{})["message"].(string)
		if errorMessage != "item not found" {
			t.Errorf("Expected error message 'item not found', got '%s'", errorMessage)
		}
	})

	t.Run("Cannot Delete Item with Stock Remaining", func(t *testing.T) {
		// Set up an item with quantity > 0 for this test
		inventory = append(inventory, map[string]interface{}{
			"id":          100,
			"name":        "Item with Stock",
			"description": "This item has stock remaining",
			"price":       30.0,
			"quantity":    10,
		})

		query := `
			mutation {
				deleteItem(id: 100)
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		errors, ok := result["errors"].([]interface{})
		if !ok || len(errors) == 0 {
			t.Fatalf("Expected error, got none")
		}

		errorMessage := errors[0].(map[string]interface{})["message"].(string)
		if errorMessage != "cannot delete an item with stock remaining" {
			t.Errorf("Expected error message 'cannot delete an item with stock remaining', got '%s'", errorMessage)
		}
	})
}
