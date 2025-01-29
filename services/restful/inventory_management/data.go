package restful

// GetMockInventory returns the initial mock inventory items
func GetMockInventory() []InventoryItem {
	return []InventoryItem{
		{ID: 1, Name: "Laptop", Description: "High-performance laptop", Price: 1500.0, Quantity: 10},
		{ID: 2, Name: "Smartphone", Description: "Latest model smartphone", Price: 800.0, Quantity: 20},
		{ID: 3, Name: "Tablet", Description: "Portable tablet", Price: 600.0, Quantity: 15},
		{ID: 4, Name: "Headphones", Description: "Noise-cancelling headphones", Price: 200.0, Quantity: 25},
		{ID: 5, Name: "Smartwatch", Description: "Fitness tracking smartwatch", Price: 300.0, Quantity: 12},
		{ID: 6, Name: "Gaming Console", Description: "Next-gen gaming console", Price: 500.0, Quantity: 5},
		{ID: 7, Name: "Monitor", Description: "4K Ultra HD monitor", Price: 400.0, Quantity: 8},
		{ID: 8, Name: "Keyboard", Description: "Mechanical keyboard", Price: 100.0, Quantity: 30},
		{ID: 9, Name: "Mouse", Description: "Wireless ergonomic mouse", Price: 50.0, Quantity: 40},
		{ID: 10, Name: "External Hard Drive", Description: "1TB external hard drive", Price: 120.0, Quantity: 18},
		{ID: 11, Name: "Webcam", Description: "1080p HD webcam", Price: 80.0, Quantity: 22},
		{ID: 12, Name: "Microphone", Description: "Studio-quality microphone", Price: 150.0, Quantity: 7},
		{ID: 13, Name: "Router", Description: "Dual-band Wi-Fi router", Price: 90.0, Quantity: 10},
		{ID: 14, Name: "Printer", Description: "All-in-one printer", Price: 250.0, Quantity: 6},
		{ID: 15, Name: "Projector", Description: "Portable mini projector", Price: 350.0, Quantity: 4},
		{ID: 16, Name: "Power Bank", Description: "20,000mAh power bank", Price: 50.0, Quantity: 35},
		{ID: 17, Name: "Drone", Description: "4K camera drone", Price: 800.0, Quantity: 3},
		{ID: 18, Name: "VR Headset", Description: "Virtual reality headset", Price: 600.0, Quantity: 5},
		{ID: 19, Name: "Smart Home Hub", Description: "Voice-controlled smart home hub", Price: 150.0, Quantity: 20},
		{ID: 20, Name: "Fitness Tracker", Description: "Health and fitness tracker", Price: 100.0, Quantity: 25},
	}
}
