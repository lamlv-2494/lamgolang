package configs

import (
	"food_delivery/internal/models/entities"
	"food_delivery/internal/utils/constants"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	dsn := os.Getenv(constants.DBDSN)
	if dsn == "" {
		log.Fatal("Error: Not config DB_DSN in .env file")
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error: Failed to connect to database:", err)
	}

	err = database.AutoMigrate(
		&entities.User{},
		&entities.Category{},
		&entities.Product{},
		&entities.CartItem{},
		&entities.Order{},
		&entities.OrderItem{},
		&entities.Rating{},
		&entities.Suggestion{},
	)

	if err != nil {
		log.Fatal("Error: Failed to migrate database:", err)
	}

	log.Println("✅ Connect to MySQL success!")
	return database
}
