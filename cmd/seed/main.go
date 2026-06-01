package main

import (
	"log"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/entities"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Cannot find .env file, reading from environment variables")
	}

	// 2. Connect to Database (reuses existing ConnectDB function)
	db := configs.ConnectDB()
	log.Println("🚀 Starting to inject Seed Data...")

	// 3. Declare sample Categories
	categories := []entities.Category{
		{Name: "Drinks", Description: "All soft drinks, juices, and coffees"},
		{Name: "Desserts", Description: "Sweet treats, delicious cakes, and pastries"},
		{Name: "Fast Food", Description: "Juicy burgers, crispy snacks, and quick bites"},
		{Name: "Main Dishes", Description: "Hearty Asian and Western main courses"},
		{Name: "Healthy Salads", Description: "Fresh, low-calorie, and nutrient-dense options"},
	}

	// Store Category ID map to safely assign to Products below
	categoryMap := make(map[string]uint)

	for i, cat := range categories {
		var existing entities.Category
		// Avoid duplicate data if seed is executed multiple times
		if err := db.Where("name = ?", cat.Name).First(&existing).Error; err != nil {
			db.Create(&categories[i])
			categoryMap[cat.Name] = categories[i].ID
			log.Printf("✅ Created Category: %s", cat.Name)
		} else {
			categoryMap[cat.Name] = existing.ID
			log.Printf("⏭️  Skipped Category (Already exists): %s", cat.Name)
		}
	}

	// 4. Declare sample Products
	products := []entities.Product{
		{Name: "Iced Matcha Latte", Description: "Premium Japanese matcha whisked with fresh milk and ice", Price: 4.50, Image: "https://images.unsplash.com/photo-1536256263959-770b48d82b0a?w=500", CategoryID: categoryMap["Drinks"]},
		{Name: "Fresh Orange Juice", Description: "100% freshly squeezed organic oranges, cold and refreshing", Price: 3.99, Image: "https://images.unsplash.com/photo-1621506289937-a8e4df240d0b?w=500", CategoryID: categoryMap["Drinks"]},
		{Name: "Chocolate Lava Cake", Description: "Rich chocolate cake with a warm, gooey molten chocolate center", Price: 5.50, Image: "https://images.unsplash.com/photo-1606313564200-e75d5e30476c?w=500", CategoryID: categoryMap["Desserts"]},
		{Name: "Classic Beef Cheeseburger", Description: "Juicy beef patty, melted cheddar cheese, lettuce, tomato, and secret sauce", Price: 6.50, Image: "https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=500", CategoryID: categoryMap["Fast Food"]},
		{Name: "Crispy Chicken Tenders", Description: "5 pieces of golden, crispy fried chicken breast served with honey mustard", Price: 4.99, Image: "https://images.unsplash.com/photo-1562967914-608f82629710?w=500", CategoryID: categoryMap["Fast Food"]},
		{Name: "Spaghetti Carbonara", Description: "Classic Italian pasta with creamy sauce, crispy bacon, and parmesan cheese", Price: 8.50, Image: "https://images.unsplash.com/photo-1612874742237-6526221588e3?w=500", CategoryID: categoryMap["Main Dishes"]},
		{Name: "Grilled Chicken Caesar Salad", Description: "Crisp romaine lettuce, grilled chicken breast, croutons, and Caesar dressing", Price: 6.25, Image: "https://images.unsplash.com/photo-1550304943-4f24f54ddde9?w=500", CategoryID: categoryMap["Healthy Salads"]},
	}

	for _, prod := range products {
		var existing entities.Product
		if err := db.Where("name = ?", prod.Name).First(&existing).Error; err != nil {
			db.Create(&prod)
			log.Printf("🍔 Created Product: %s", prod.Name)
		} else {
			log.Printf("⏭️  Skipped Product (Already exists): %s", prod.Name)
		}
	}

	log.Println("🎉 SEED DATA INJECTION COMPLETED!")
}
