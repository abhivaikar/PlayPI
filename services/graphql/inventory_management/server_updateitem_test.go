package graphql

import (
	"net/http/httptest"
	"testing"
)

func TestUpdateItem(t *testing.T) {
	ts := httptest.NewServer(setupGraphQLServer())
	defer ts.Close()

	t.Run("Successful Update", func(t *testing.T) {
		query := `
			mutation {
				updateItem(id: 1, name: "Updated Name", description: "Updated Description", price: 99.99, quantity: 20) {
					id
					name
					description
					price
					quantity
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		data, ok := result["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to extract 'data' from response: %v", result)
		}
		updatedItem := data["updateItem"].(map[string]interface{})
		if updatedItem["name"] != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %v", updatedItem["name"])
		}
		if updatedItem["price"] != 99.99 {
			t.Errorf("Expected price 99.99, got %v", updatedItem["price"])
		}
		if int(updatedItem["quantity"].(float64)) != 20 {
			t.Errorf("Expected quantity 20, got %v", updatedItem["quantity"])
		}
		if updatedItem["description"] != "Updated Description" {
			t.Errorf("Expected description 'Updated Description', got %v", updatedItem["description"])
		}
	})

	t.Run("Item Not Found", func(t *testing.T) {
		query := `
			mutation {
				updateItem(id: 999, name: "Non-Existent", description: "This item does not exist", price: 50, quantity: 10) {
					id
				}
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

	t.Run("Invalid Name - Too Short", func(t *testing.T) {
		query := `
			mutation {
				updateItem(id: 1, name: "Up", description: "Updated Description", price: 50, quantity: 10) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		errors, ok := result["errors"].([]interface{})
		if !ok || len(errors) == 0 {
			t.Fatalf("Expected error, got none")
		}
		errorMessage := errors[0].(map[string]interface{})["message"].(string)
		if errorMessage != "name must be between 3 and 50 characters" {
			t.Errorf("Expected error message 'name must be between 3 and 50 characters', got '%s'", errorMessage)
		}
	})

	t.Run("Invalid Price - Negative", func(t *testing.T) {
		query := `
			mutation {
				updateItem(id: 1, name: "Valid Name", description: "Valid Description", price: -10, quantity: 10) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		errors, ok := result["errors"].([]interface{})
		if !ok || len(errors) == 0 {
			t.Fatalf("Expected error, got none")
		}
		errorMessage := errors[0].(map[string]interface{})["message"].(string)
		if errorMessage != "price must be a positive number not exceeding 10,000" {
			t.Errorf("Expected error message 'price must be a positive number not exceeding 10,000', got '%s'", errorMessage)
		}
	})

	t.Run("Invalid Quantity - Negative", func(t *testing.T) {
		query := `
			mutation {
				updateItem(id: 1, name: "Valid Name", description: "Valid Description", price: 50, quantity: -5) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		errors, ok := result["errors"].([]interface{})
		if !ok || len(errors) == 0 {
			t.Fatalf("Expected error, got none")
		}
		errorMessage := errors[0].(map[string]interface{})["message"].(string)
		if errorMessage != "quantity cannot be negative" {
			t.Errorf("Expected error message 'quantity cannot be negative', got '%s'", errorMessage)
		}
	})
}
