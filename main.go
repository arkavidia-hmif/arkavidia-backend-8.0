package main

import (
	"runtime"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	messageConfig "arkavidia-backend-8.0/competition/config/message"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/routes"
	"arkavidia-backend-8.0/competition/utils/mail"
	"arkavidia-backend-8.0/competition/utils/validation"
)

// TODO: Gunakan gzip untuk mengkompresi size HTTP Response
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-http-gzip-compression.html
// ASSIGNED TO: @rayhankinan
// STATUS: DONE

// TODO: Tambahkan secure middleware untuk menambah security
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-secure-middleware.html
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-https-tls.html
// ASSIGNED TO: @rayhankinan

func main() {
	// Configure runtime
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Gin Framework
	engine := gin.Default()

	// Setup Validator
	binding.Validator = validation.Validator

	// Middlewares
	engine.Use(middlewares.CORSMiddleware())
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	// Routes
	routes.AdminRoute(engine)
	routes.TeamRoute(engine)
	routes.ParticipantRoute(engine)
	routes.SubmissionRoute(engine)
	routes.PhotoRoute(engine)
	routes.NotFoundRoute(engine)

	// Goroutine Worker
	configMessage := messageConfig.Config.GetMetadata()
	go mail.Broker.RunMailWorker(configMessage.WorkerSize)

	// Run App
	engine.Run()
}
