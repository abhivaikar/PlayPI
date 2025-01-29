package restful

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupTestServer() *gin.Engine {
	return StartServerForTesting() // Use your updated testing setup
}

func TestGetItems(t *testing.T) {
	r := setupTestServer()

	t.Run("Retrieve all items", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/items", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var items []InventoryItem
		err := json.Unmarshal(resp.Body.Bytes(), &items)
		require.NoError(t, err)
		require.Len(t, items, 20)
		require.Equal(t, "Laptop", items[0].Name) // Validate first item
	})
}

func TestAddItem(t *testing.T) {
	r := setupTestServer()

	t.Run("Add a valid item", func(t *testing.T) {
		payload := `{
			"name": "Wireless Charger",
			"description": "Fast wireless charging pad",
			"price": 40.0,
			"quantity": 10
		}`
		req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusCreated, resp.Code)

		var item InventoryItem
		err := json.Unmarshal(resp.Body.Bytes(), &item)
		require.NoError(t, err)
		require.Equal(t, "Wireless Charger", item.Name)
		require.Equal(t, "Fast wireless charging pad", item.Description)
		require.Equal(t, 40.0, item.Price)
		require.Equal(t, 10, item.Quantity)
	})

	// Validation Error Cases
	t.Run("Validation Error - Missing Name", func(t *testing.T) {
		payload := `{
			"description": "Fast wireless charging pad",
			"price": 40.0,
			"quantity": 10
		}`
		req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "name must be between 3 and 50 characters", response["error"])
	})

	t.Run("Validation Error - Name Too Short", func(t *testing.T) {
		payload := `{
			"name": "Wi",
			"description": "Fast wireless charging pad",
			"price": 40.0,
			"quantity": 10
		}`
		req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "name must be between 3 and 50 characters", response["error"])
	})

	t.Run("Validation Error - Price Out of Range", func(t *testing.T) {
		payload := `{
			"name": "Wireless Charger",
			"description": "Fast wireless charging pad",
			"price": 20000.0,
			"quantity": 10
		}`
		req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "price must be a positive number not exceeding 10,000", response["error"])
	})

	t.Run("Validation Error - Negative Quantity", func(t *testing.T) {
		payload := `{
			"name": "Wireless Charger",
			"description": "Fast wireless charging pad",
			"price": 40.0,
			"quantity": -5
		}`
		req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "quantity must be at least 0", response["error"])
	})

	t.Run("Validation Error - Description Too Long", func(t *testing.T) {
		payload := `{
			"name": "Wireless Charger",
			"description": "This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters.",
			"price": 40.0,
			"quantity": 10
		}`
		req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "description cannot exceed 200 characters", response["error"])
	})
}

func TestUpdateItem(t *testing.T) {
	r := setupTestServer()

	t.Run("Update a valid item", func(t *testing.T) {
		payload := `{
			"name": "Updated Laptop",
			"description": "Updated high-performance laptop",
			"price": 1600.0,
			"quantity": 8
		}`
		req, _ := http.NewRequest(http.MethodPut, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var item InventoryItem
		err := json.Unmarshal(resp.Body.Bytes(), &item)
		require.NoError(t, err)
		require.Equal(t, "Updated Laptop", item.Name)
		require.Equal(t, "Updated high-performance laptop", item.Description)
		require.Equal(t, 1600.0, item.Price)
		require.Equal(t, 8, item.Quantity)
	})

	// Validation Error Cases
	t.Run("Validation Error - Name Too Short", func(t *testing.T) {
		payload := `{
			"name": "Up",
			"description": "Updated high-performance laptop",
			"price": 1600.0,
			"quantity": 8
		}`
		req, _ := http.NewRequest(http.MethodPut, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "name must be between 3 and 50 characters", response["error"])
	})

	t.Run("Validation Error - Description Too Long", func(t *testing.T) {
		payload := `{
			"name": "Updated Laptop",
			"description": "This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters.",
			"price": 1600.0,
			"quantity": 8
		}`
		req, _ := http.NewRequest(http.MethodPut, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "description cannot exceed 200 characters", response["error"])
	})

	t.Run("Validation Error - Negative Quantity", func(t *testing.T) {
		payload := `{
			"name": "Updated Laptop",
			"description": "Updated high-performance laptop",
			"price": 1600.0,
			"quantity": -5
		}`
		req, _ := http.NewRequest(http.MethodPut, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "quantity must be at least 0", response["error"])
	})
}

func TestPatchItem(t *testing.T) {
	r := setupTestServer()

	t.Run("Partial update with valid data", func(t *testing.T) {
		payload := `{
			"price": 1200.0
		}`
		req, _ := http.NewRequest(http.MethodPatch, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var item InventoryItem
		err := json.Unmarshal(resp.Body.Bytes(), &item)
		require.NoError(t, err)
		require.Equal(t, 1200.0, item.Price)
	})

	// Validation Error Cases
	t.Run("Validation Error - Name Too Short", func(t *testing.T) {
		payload := `{
			"name": "Wi"
		}`
		req, _ := http.NewRequest(http.MethodPatch, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "name must be between 3 and 50 characters", response["error"])
	})

	t.Run("Validation Error - Description Too Long", func(t *testing.T) {
		payload := `{
			"description": "This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters.This description exceeds the character limit of 200 characters. This description exceeds the character limit of 200 characters."
		}`
		req, _ := http.NewRequest(http.MethodPatch, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "description cannot exceed 200 characters", response["error"])
	})

	t.Run("Validation Error - Negative Quantity", func(t *testing.T) {
		payload := `{
			"quantity": -10
		}`
		req, _ := http.NewRequest(http.MethodPatch, "/items/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "quantity must be at least 0", response["error"])
	})
}

func TestDeleteItem(t *testing.T) {
	r := setupTestServer()

	t.Run("Delete an existing item", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/items/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "item deleted", response["message"])
	})

	t.Run("Item not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/items/999", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "item not found", response["error"])
	})
}
