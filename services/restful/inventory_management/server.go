package restful

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartServer initializes and starts the RESTful API server
func StartServer() {
	inventory = GetMockInventory() // Use mock inventory for production
	nextID = len(inventory) + 1
	r := setupRouter()
	fmt.Println("RESTful API is running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

// StartServerForTesting initializes the router for testing
func StartServerForTesting() *gin.Engine {
	inventory = GetMockInventory() // Use mock inventory for production
	nextID = len(inventory) + 1
	return setupRouter()
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	// GET /items - Get all items
	r.GET("/items", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetAllItems())
	})

	// POST /items - Add a new item
	r.POST("/items", func(c *gin.Context) {
		var newItem InventoryItem
		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
			return
		}

		item, err := AddItem(newItem)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	// PUT /items/:id - Update an existing item
	r.PUT("/items/:id", func(c *gin.Context) {
		var updatedData InventoryItem
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
			return
		}

		id, err := parseIDParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}

		item, err := UpdateItem(id, updatedData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	// PATCH /items/:id - Partially update an item
	r.PATCH("/items/:id", func(c *gin.Context) {
		id, err := parseIDParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
			return
		}

		item, err := PatchItem(id, updates)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	// DELETE /items/:id - Delete an item
	r.DELETE("/items/:id", func(c *gin.Context) {
		id, err := parseIDParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}

		err = DeleteItem(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
	})

	return r
}

// parseIDParam extracts the ID parameter from the request
func parseIDParam(c *gin.Context) (int, error) {
	var id int
	_, err := fmt.Sscanf(c.Param("id"), "%d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
