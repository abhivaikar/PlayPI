package restful

import (
	"errors"
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
var nextID = 1 // Auto-increment ID

// Validation functions
func validateItem(item InventoryItem) error {
	if len(item.Name) < 3 || len(item.Name) > 50 {
		return errors.New("name must be between 3 and 50 characters")
	}
	if len(item.Description) > 200 {
		return errors.New("description cannot exceed 200 characters")
	}
	if item.Price < 0 || item.Price > 10000 {
		return errors.New("price must be a positive number not exceeding 10,000")
	}
	if item.Quantity < 0 {
		return errors.New("quantity must be at least 0")
	}
	return nil
}

// Service functions
func GetAllItems() []InventoryItem {
	return inventory
}

func GetItemByID(id int) (*InventoryItem, error) {
	for _, item := range inventory {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, errors.New("item not found")
}

func AddItem(newItem InventoryItem) (*InventoryItem, error) {
	if err := validateItem(newItem); err != nil {
		return nil, err
	}

	newItem.ID = nextID
	nextID++
	inventory = append(inventory, newItem)

	return &newItem, nil
}

func UpdateItem(id int, updatedData InventoryItem) (*InventoryItem, error) {
	for i, item := range inventory {
		if item.ID == id {
			if err := validateItem(updatedData); err != nil {
				return nil, err
			}
			updatedData.ID = id // Preserve the original ID
			inventory[i] = updatedData
			return &updatedData, nil
		}
	}
	return nil, errors.New("item not found")
}

func PatchItem(id int, updates map[string]interface{}) (*InventoryItem, error) {
	for i, item := range inventory {
		if item.ID == id {
			if name, exists := updates["name"]; exists {
				nameStr, ok := name.(string)
				if !ok || len(nameStr) < 3 || len(nameStr) > 50 {
					return nil, errors.New("name must be between 3 and 50 characters")
				}
				item.Name = nameStr
			}
			if description, exists := updates["description"]; exists {
				descriptionStr, ok := description.(string)
				if !ok || len(descriptionStr) > 200 {
					return nil, errors.New("description cannot exceed 200 characters")
				}
				item.Description = descriptionStr
			}
			if price, exists := updates["price"]; exists {
				priceFloat, ok := price.(float64)
				if !ok || priceFloat < 0 || priceFloat > 10000 {
					return nil, errors.New("price must be a positive number not exceeding 10,000")
				}
				item.Price = priceFloat
			}
			if quantity, exists := updates["quantity"]; exists {
				quantityInt, ok := quantity.(float64)
				if !ok || int(quantityInt) < 0 {
					return nil, errors.New("quantity must be at least 0")
				}
				item.Quantity = int(quantityInt)
			}

			inventory[i] = item
			return &item, nil
		}
	}
	return nil, errors.New("item not found")
}

func DeleteItem(id int) error {
	for i, item := range inventory {
		if item.ID == id {
			inventory = append(inventory[:i], inventory[i+1:]...)
			return nil
		}
	}
	return errors.New("item not found")
}
