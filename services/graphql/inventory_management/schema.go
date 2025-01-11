package graphql

import (
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
	{"id": 2, "name": "Phone", "description": "Smartphone with excellent camera", "price": 799.99, "quantity": 25},
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
				id, ok := p.Args["id"].(int)
				if ok {
					// Find the item with the given ID
					for _, item := range inventory {
						if item["id"] == id {
							return item, nil
						}
					}
				}
				return nil, nil
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
				newItem := map[string]interface{}{
					"id":          len(inventory) + 1,
					"name":        p.Args["name"].(string),
					"description": p.Args["description"],
					"price":       p.Args["price"],
					"quantity":    p.Args["quantity"],
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
				for i, item := range inventory {
					if item["id"] == id {
						if name, ok := p.Args["name"].(string); ok {
							inventory[i]["name"] = name
						}
						if description, ok := p.Args["description"].(string); ok {
							inventory[i]["description"] = description
						}
						if price, ok := p.Args["price"].(float64); ok {
							inventory[i]["price"] = price
						}
						if quantity, ok := p.Args["quantity"].(int); ok {
							inventory[i]["quantity"] = quantity
						}
						return inventory[i], nil
					}
				}
				return nil, nil
			},
		},
		"deleteItem": &graphql.Field{
			Type: graphql.Boolean, // Return true if successful
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				for i, item := range inventory {
					if item["id"] == id {
						inventory = append(inventory[:i], inventory[i+1:]...)
						return true, nil
					}
				}
				return false, nil
			},
		},
	},
})

// Define the schema
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
