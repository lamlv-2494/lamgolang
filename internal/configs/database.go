package configs

import (
	"food_delivery/internal/models/entities"
	"food_delivery/internal/utils/constants"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
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
	seedAdmin(database)
	return database
}

// seedAdmin checks if an admin user exists and creates one if not
func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&entities.User{}).Where("role = ?", constants.AdminRole).Count(&count)

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("❌ Failed to hash default admin password:", err)
			return
		}

		admin := entities.User{
			Username: "admin",
			Email:    "admin@fooddelivery.com",
			Password: string(hashedPassword),
			Role:     constants.AdminRole,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Println("❌ Failed to seed default admin user:", err)
		} else {
			log.Println("👑 Seeded default admin account successfully! (Email: admin@fooddelivery.com | Pass: admin123456)")
		}
	}
}
