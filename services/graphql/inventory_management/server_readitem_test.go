package graphql

import (
	"log"
	"net/http/httptest"
	"testing"
)

func TestQueryItems(t *testing.T) {
	ts := httptest.NewServer(setupGraphQLServer())
	defer ts.Close()

	t.Run("Successful Retrieval of All Items", func(t *testing.T) {
		query := `
			query {
				items {
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
		items := data["items"].([]interface{})
		log.Printf("Items: %+v", items)

		if len(items) == 0 {
			t.Errorf("Expected items to be returned, got empty list")
		}
	})
}

func TestQueryItem(t *testing.T) {
	ts := httptest.NewServer(setupGraphQLServer())
	defer ts.Close()

	t.Run("Successful Retrieval of an Item by ID", func(t *testing.T) {
		query := `
			query {
				item(id: 1) {
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
		item := data["item"].(map[string]interface{})
		log.Printf("Item: %+v", item)

		if item["id"].(float64) != 1 {
			t.Errorf("Expected item with ID 1, got %v", item["id"])
		}
	})

	t.Run("Item Not Found (Invalid ID)", func(t *testing.T) {
		query := `
			query {
				item(id: 999) {
					id
					name
					description
					price
					quantity
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)

		// Validate response data
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to extract 'data' from response: %v", result)
		}
		item := data["item"]

		if item != nil {
			t.Errorf("Expected null, got %v", item)
		}

		// Validate error messages
		errors, ok := result["errors"].([]interface{})
		if !ok || len(errors) == 0 {
			t.Fatalf("Expected error, got none")
		}

		errorMessage := errors[0].(map[string]interface{})["message"].(string)
		if errorMessage != "item not found" {
			t.Errorf("Expected error message 'Item not found', got '%s'", errorMessage)
		}
	})
}
