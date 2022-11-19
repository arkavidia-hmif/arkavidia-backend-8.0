package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/services/database"
	"arkavidia-backend-8.0/competition/services/storage"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[ERROR] No .env file found!")
	}
}

func main() {
	r := gin.Default()

	// Setup services
	database.GetDB()
	storage.GetClient()

	// Middlewares
	r.Use(middlewares.CORSMiddleware())

	r.Run()
}
