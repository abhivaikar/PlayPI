package graphql

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestAddItemScenarios(t *testing.T) {
	ts := httptest.NewServer(setupGraphQLServer())
	defer ts.Close()

	t.Run("Successful Addition", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "Test Item", description: "A test item", price: 50.5, quantity: 10) {
					id
					name
					description
					price
					quantity
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		data := result["data"].(map[string]interface{})
		addItem := data["addItem"].(map[string]interface{})

		if addItem["name"] != "Test Item" {
			t.Errorf("Expected name 'Test Item', got %v", addItem["name"])
		}
	})

	t.Run("Invalid Name - Too Short", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "Ti", description: "Invalid name", price: 20.5, quantity: 5) {
					id
					name
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		errors := result["errors"].([]interface{})
		if len(errors) == 0 || errors[0].(map[string]interface{})["message"] != "name must be between 3 and 50 characters" {
			t.Errorf("Expected validation error for short name, got %v", errors)
		}
	})

	t.Run("Invalid Name - Too Long", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "` + longString(51) + `", description: "invalid name length", price: 20.5, quantity: 5) {
					id
					name
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		errors := result["errors"].([]interface{})
		if len(errors) == 0 || errors[0].(map[string]interface{})["message"] != "name must be between 3 and 50 characters" {
			t.Errorf("Expected validation error for long name, got %v", errors)
		}
	})

	t.Run("Invalid Price - Negative", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "Negative Price", description: "Invalid price", price: -10.0, quantity: 5) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		errors := result["errors"].([]interface{})
		if len(errors) == 0 || errors[0].(map[string]interface{})["message"] != "price must be a positive number not exceeding 10,000" {
			t.Errorf("Expected validation error for negative price, got %v", errors)
		}
	})

	t.Run("Duplicate Name", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "Test Item", description: "Duplicate name test", price: 20.5, quantity: 5) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		errors := result["errors"].([]interface{})
		if len(errors) == 0 || errors[0].(map[string]interface{})["message"] != "an item with this name already exists" {
			t.Errorf("Expected validation error for duplicate name, got %v", errors)
		}
	})

	t.Run("Price Exceeding 10,000", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "Expensive Item", description: "Price too high", price: 15000, quantity: 10) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		errors := result["errors"].([]interface{})
		if len(errors) == 0 || errors[0].(map[string]interface{})["message"] != "price must be a positive number not exceeding 10,000" {
			t.Errorf("Expected validation error for price exceeding 10,000, got %v", errors)
		}
	})

	t.Run("Quantity Less Than 1", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "Invalid Quantity", description: "Quantity too low", price: 100, quantity: 0) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		errors := result["errors"].([]interface{})
		if len(errors) == 0 || errors[0].(map[string]interface{})["message"] != "quantity must be at least 1" {
			t.Errorf("Expected validation error for quantity less than 1, got %v", errors)
		}
	})

	t.Run("Description Exceeding Character Limit", func(t *testing.T) {
		query := `
			mutation {
				addItem(name: "Invalid Description", description: "` + longString(201) + `", price: 100, quantity: 5) {
					id
				}
			}
		`
		result := sendGraphQLRequest(t, ts, query)
		errors := result["errors"].([]interface{})
		if len(errors) == 0 || errors[0].(map[string]interface{})["message"] != "description cannot exceed 200 characters" {
			t.Errorf("Expected validation error for description exceeding 200 characters, got %v", errors)
		}
	})
}

// Utility to generate long strings for testing
func longString(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += "a"
	}
	return result
}

func setupGraphQLServer() http.Handler {
	// Reuse the global Schema variable from schema.go
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		body := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query, ok := body["query"].(string)
		if !ok {
			http.Error(w, "Query not provided", http.StatusBadRequest)
			return
		}

		// Execute the GraphQL query using the existing Schema
		params := graphql.Params{
			Schema:        Schema, // Reusing the global Schema variable
			RequestString: query,
		}
		result := graphql.Do(params)

		// Return the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
}

func sendGraphQLRequest(t *testing.T, ts *httptest.Server, query string) map[string]interface{} {
	payload := map[string]interface{}{
		"query": query,
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	log.Println("GraphQL Response:", result)
	return result
}
