package main

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/routes"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	storageService "arkavidia-backend-8.0/competition/services/storage"
)

func main() {
	r := gin.Default()

	// Setup services
	databaseService.GetDB()
	storageService.GetClient()

	// Middlewares
	r.Use(middlewares.CORSMiddleware())

	// Routes
	routes.NotFoundRoute(r)
	routes.TeamRoute(r)
	routes.ParticipantRoute(r)
	routes.SubmissionRoute(r)
	routes.PhotoRoute(r)

	r.Run()
}
