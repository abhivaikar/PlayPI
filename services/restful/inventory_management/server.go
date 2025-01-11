package restful

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InventoryItem represents an item in the inventory
type InventoryItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// In-memory store for inventory items
var inventory = []InventoryItem{}

func StartServer() {
	r := gin.Default()

	// List all items
	r.GET("/items", func(c *gin.Context) {
		c.JSON(http.StatusOK, inventory)
	})

	// Add a new item
	r.POST("/items", func(c *gin.Context) {
		var newItem InventoryItem
		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newItem.ID = len(inventory) + 1 // Auto-increment ID
		inventory = append(inventory, newItem)
		c.JSON(http.StatusCreated, newItem)
	})

	// Update an existing item
	r.PUT("/items/:id", func(c *gin.Context) {
		// Extract the ID from the URL
		idParam := c.Param("id")
		var id int
		if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Bind the JSON payload to the updatedItem struct
		var updatedData InventoryItem
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Locate the item in the inventory
		for i, item := range inventory {
			if item.ID == id {
				// Update the item with new data
				inventory[i].Name = updatedData.Name
				inventory[i].Description = updatedData.Description
				inventory[i].Price = updatedData.Price
				inventory[i].Quantity = updatedData.Quantity

				c.JSON(http.StatusOK, inventory[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
	})

	r.PATCH("/items/:id", func(c *gin.Context) {
		// Extract the ID from the URL
		idParam := c.Param("id")
		var id int
		if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Bind the JSON payload as a map for flexible updates
		var updatedData map[string]interface{}
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Locate the item in the inventory
		for i, item := range inventory {
			if item.ID == id {
				// Update only the fields provided in the payload
				if name, exists := updatedData["name"]; exists {
					inventory[i].Name = name.(string)
				}
				if description, exists := updatedData["description"]; exists {
					inventory[i].Description = description.(string)
				}
				if price, exists := updatedData["price"]; exists {
					inventory[i].Price = price.(float64)
				}
				if quantity, exists := updatedData["quantity"]; exists {
					inventory[i].Quantity = int(quantity.(float64))
				}

				c.JSON(http.StatusOK, inventory[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
	})

	// Delete an item
	r.DELETE("/items/:id", func(c *gin.Context) {
		var id int
		if err := c.ShouldBindJSON(&id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for i, item := range inventory {
			if item.ID == id {
				inventory = append(inventory[:i], inventory[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
	})

	fmt.Println("RESTful API is running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
