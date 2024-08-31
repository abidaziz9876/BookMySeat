package main

import (
	"net/http"
	"time"

	"github.com/abidaziz9876/payment-service/config"
	"github.com/abidaziz9876/payment-service/services"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	config.LoadEnv()
	config.SetupRabbitMQ()
	go services.ConsumePaymentRequests()
	 

	corsConfig := cors.Config{
		Origins:         "*",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		Methods:         "GET, POST, PUT,DELETE",
		Credentials:     false,
		ValidateHeaders: false,
		MaxAge:          1 * time.Minute,
	}
	router := gin.Default()
	router.Use(cors.Middleware(corsConfig))

	// // routes.PaymentRoutes(router)

	log.Infof("Server listening on http://localhost:9999/")
	if err := http.ListenAndServe("0.0.0.0:9999", router); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
