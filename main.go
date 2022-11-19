package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"arkavidia-backend-8.0/competition/middlewares"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	storageService "arkavidia-backend-8.0/competition/services/storage"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()

	// Setup services
	databaseService.GetDB()
	storageService.GetClient()

	// Middlewares
	r.Use(middlewares.CORSMiddleware())

	r.Run()
}
