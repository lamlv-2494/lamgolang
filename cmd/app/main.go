package main

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/utils/constants"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Can't found env file")
	}

	_ = configs.ConnectDB()

	r := gin.Default()

	port := os.Getenv(constants.Port)
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
