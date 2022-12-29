package main

import (
	"runtime"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	messageConfig "arkavidia-backend-8.0/competition/config/message"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/routes"
	"arkavidia-backend-8.0/competition/utils/mail"
)

// TODO: Gunakan gzip untuk mengkompresi size HTTP Response
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-http-gzip-compression.html
// ASSIGNED TO: @rayhankinan
// STATUS: DONE

// TODO: Tambahkan secure middleware untuk menambah security
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-secure-middleware.html
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-https-tls.html
// NOTES DARI GARE: Pake openssl aja buat generate certificate
// ASSIGNED TO: @confusionhill
// STATUS: IN PROGRESS

// TODO: Tambahkan route render photo dan submission untuk menghindari akses ke google cloud storage secara langsung
// REFERENCE: https://zetcode.com/golang/http-serve-image/
// REFERENCE: https://stackoverflow.com/questions/26744814/serve-image-in-go-that-was-just-created
// REFERENCE: https://freshman.tech/snippets/go/file-content-type/
// REFERENCE: https://stackoverflow.com/questions/51209439/mime-type-checking-of-files-uploaded-golang
// ASSIGNED TO: @patrickamadeus
// STATUS: IN PROGRESS

// TODO: Gunakan GormValuerInterface untuk mengautomatisasi enkripsi bcrypt password
// REFERENCE: https://gorm.io/docs/data_types.html#GormValuerInterface
// ASSIGNED TO: @graceclaudia19
// STATUS: DONE

// TODO: Tambahkan hash pada semua ID di model untuk mencegah terjadinya IDOR
// REFERENCE: https://www.securecoding.com/blog/how-to-prevent-idor-attacks/
// ASSIGNED TO: @akbarmridho
// STATUS: IN PROGRESS

// TODO: Gunakan syntax iota untuk membuat tipe enum untuk memperkecil size penyimpanan pada basis data (string menjadi integer)
// REFERENCE: https://levelup.gitconnected.com/implementing-enums-in-golang-9537c433d6e2
// ASSIGNED TO: @rayhankinan
// STATUS: IN PROGRESS

func main() {
	// Configure runtime
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Gin Framework
	engine := gin.Default()

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
