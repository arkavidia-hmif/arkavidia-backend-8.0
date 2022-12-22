package main

import (
	"runtime"

	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/routes"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	storageService "arkavidia-backend-8.0/competition/services/storage"
	"arkavidia-backend-8.0/competition/utils/worker"
)

func main() {
	// Configure runtime
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Gin Framework
	r := gin.Default()

	// Setup Services Check
	databaseService.GetDB()
	storageService.GetClient()

	// Middlewares
	r.Use(middlewares.CORSMiddleware())

	// Routes
	routes.AdminRoute(r)
	routes.TeamRoute(r)
	routes.ParticipantRoute(r)
	routes.SubmissionRoute(r)
	routes.PhotoRoute(r)
	routes.NotFoundRoute(r)

	// Goroutine Worker
	go worker.MailRun()

	// RUn App
	r.Run()
}
