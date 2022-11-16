package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"arkavidia-backend-8.0/competition/services/database"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[ERROR] No .env file found!")
	}
}

func main() {
	r := gin.Default()
	database.GetDB()

	r.Run()
}
