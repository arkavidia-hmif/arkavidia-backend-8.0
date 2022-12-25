package main

import (
	"runtime"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/routes"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	storageService "arkavidia-backend-8.0/competition/services/storage"
	"arkavidia-backend-8.0/competition/utils/validation"
	"arkavidia-backend-8.0/competition/utils/worker"
)

// TODO: Gunakan gzip untuk mengkompresi size HTTP Handler
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-http-gzip-compression.html
// ASSIGNED TO: @rayhankinan
// STATUS: DONE

// TODO: Tambahkan validasi payload request dengan menggunakan validator
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-http-request-payload-validation.html
// ASSIGNED TO: @rayhankinan

// TODO: Tambahkan secure middleware untuk menambah security
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-secure-middleware.html
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-https-tls.html
// ASSIGNED TO: @rayhankinan

func main() {
	// Configure runtime
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Gin Framework
	r := gin.Default()

	// Setup Services Check
	databaseService.GetDB()
	storageService.GetClient()

	// Setup Validator
	binding.Validator = validation.GetValidator()

	// Middlewares
	r.Use(middlewares.CORSMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

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
