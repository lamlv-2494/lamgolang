package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	err = db.AutoMigrate(
		&entities.User{},
		&entities.Category{},
		&entities.Product{},
		&entities.Rating{},
		&entities.CartItem{},
		&entities.Order{},
		&entities.OrderItem{},
		&entities.Suggestion{},
	)
	if err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	return db
}
