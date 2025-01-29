package graphql

import (
	"errors"

	"github.com/graphql-go/graphql"
)

// Define the InventoryItem type
var inventoryItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "InventoryItem",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.Int},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"price":       &graphql.Field{Type: graphql.Float},
		"quantity":    &graphql.Field{Type: graphql.Int},
	},
})

// Inventory data (mock data)
var inventory = []map[string]interface{}{
	{"id": 1, "name": "Laptop", "description": "High-performance laptop", "price": 1500.99, "quantity": 10},
	{"id": 2, "name": "Smartphone", "description": "Latest model smartphone", "price": 899.99, "quantity": 20},
	{"id": 3, "name": "Tablet", "description": "Portable and powerful tablet", "price": 499.99, "quantity": 15},
	{"id": 4, "name": "Smartwatch", "description": "Stylish smartwatch", "price": 199.99, "quantity": 25},
	{"id": 5, "name": "Headphones", "description": "Noise-canceling headphones", "price": 99.99, "quantity": 30},
	{"id": 6, "name": "Monitor", "description": "4K Ultra HD monitor", "price": 399.99, "quantity": 12},
	{"id": 7, "name": "Keyboard", "description": "Mechanical keyboard", "price": 79.99, "quantity": 50},
	{"id": 8, "name": "Mouse", "description": "Wireless ergonomic mouse", "price": 29.99, "quantity": 40},
	{"id": 9, "name": "Printer", "description": "Multi-function printer", "price": 249.99, "quantity": 8},
	{"id": 10, "name": "Camera", "description": "DSLR camera with lens kit", "price": 1199.99, "quantity": 5},
	{"id": 11, "name": "External Hard Drive", "description": "1TB external storage", "price": 59.99, "quantity": 30},
	{"id": 12, "name": "Gaming Console", "description": "Next-gen gaming console", "price": 499.99, "quantity": 10},
	{"id": 13, "name": "Router", "description": "High-speed wireless router", "price": 89.99, "quantity": 25},
	{"id": 14, "name": "Speaker", "description": "Bluetooth portable speaker", "price": 49.99, "quantity": 35},
	{"id": 15, "name": "Power Bank", "description": "Fast-charging power bank", "price": 19.99, "quantity": 50},
	{"id": 16, "name": "Projector", "description": "1080p home theater projector", "price": 299.99, "quantity": 7},
	{"id": 17, "name": "Smart Bulb", "description": "Color-changing smart bulb", "price": 14.99, "quantity": 60},
	{"id": 18, "name": "Fitness Tracker", "description": "Waterproof fitness tracker", "price": 49.99, "quantity": 20},
	{"id": 19, "name": "Electric Scooter", "description": "Lightweight and portable", "price": 699.99, "quantity": 3},
	{"id": 20, "name": "Drone", "description": "Quadcopter with camera", "price": 999.99, "quantity": 5},
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"items": &graphql.Field{
			Type: graphql.NewList(inventoryItemType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Return the current state of the inventory
				return inventory, nil
			},
		},
		"item": &graphql.Field{
			Type: inventoryItemType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Validate if the ID argument is provided
				id, ok := p.Args["id"].(int)
				if !ok {
					return nil, errors.New("invalid or missing 'id' argument. 'id' must be an integer")
				}

				// Search for the item with the given ID
				for _, item := range inventory {
					if item["id"] == id {
						return item, nil
					}
				}

				// Return an error if the item is not found
				return nil, errors.New("item not found")
			},
		},
	},
})

// Define the root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"addItem": &graphql.Field{
			Type: inventoryItemType, // Return the newly created item
			Args: graphql.FieldConfigArgument{
				"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
				"price":       &graphql.ArgumentConfig{Type: graphql.Float},
				"quantity":    &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				name := p.Args["name"].(string)
				if len(name) < 3 || len(name) > 50 {
					return nil, errors.New("name must be between 3 and 50 characters")
				}
				price := p.Args["price"].(float64)
				if price < 0 || price > 10000 {
					return nil, errors.New("price must be a positive number not exceeding 10,000")
				}
				quantity := p.Args["quantity"].(int)
				if quantity < 1 {
					return nil, errors.New("quantity must be at least 1")
				}
				description := p.Args["description"].(string)
				if len(description) > 200 {
					return nil, errors.New("description cannot exceed 200 characters")
				}

				// Prevent duplicate names
				for _, item := range inventory {
					if item["name"] == name {
						return nil, errors.New("an item with this name already exists")
					}
				}

				// Add the new item
				newItem := map[string]interface{}{
					"id":          len(inventory) + 1,
					"name":        name,
					"description": description,
					"price":       price,
					"quantity":    quantity,
				}
				inventory = append(inventory, newItem)
				return newItem, nil
			},
		},
		"updateItem": &graphql.Field{
			Type: inventoryItemType, // Return the updated item
			Args: graphql.FieldConfigArgument{
				"id":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":        &graphql.ArgumentConfig{Type: graphql.String},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
				"price":       &graphql.ArgumentConfig{Type: graphql.Float},
				"quantity":    &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				var item map[string]interface{}
				var found bool
				for _, i := range inventory {
					if i["id"] == id {
						item = i
						found = true
						break
					}
				}

				if !found {
					return nil, errors.New("item not found")
				}

				name := p.Args["name"].(string)
				if len(name) < 3 || len(name) > 50 {
					return nil, errors.New("name must be between 3 and 50 characters")
				}
				price := p.Args["price"].(float64)
				if price < 0 || price > 10000 {
					return nil, errors.New("price must be a positive number not exceeding 10,000")
				}
				quantity := p.Args["quantity"].(int)
				if quantity < 0 {
					return nil, errors.New("quantity cannot be negative")
				}
				description := p.Args["description"].(string)
				if len(description) > 200 {
					return nil, errors.New("description cannot exceed 200 characters")
				}

				// Update the item
				item["name"] = name
				item["price"] = price
				item["quantity"] = quantity
				item["description"] = description

				return item, nil
			},
		},
		"deleteItem": &graphql.Field{
			Type: graphql.Boolean, // Return true if successful
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				var index int
				var found bool
				for i, item := range inventory {
					if item["id"] == id {
						if item["quantity"].(int) > 0 {
							return nil, errors.New("cannot delete an item with stock remaining")
						}
						index = i
						found = true
						break
					}
				}

				if !found {
					return nil, errors.New("item not found")
				}

				inventory = append(inventory[:index], inventory[index+1:]...)
				return "Item deleted successfully", nil
			},
		},
	},
})

// Define the schema
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
