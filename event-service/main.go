package main

import (
	"net/http"
	"time"

	"github.com/abidaziz9876/event-service/config"
	"github.com/abidaziz9876/event-service/routes"
	"github.com/abidaziz9876/event-service/services"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	config.LoadEnv()
	config.ConnectPostGresDB()
	config.SetupPostgresSchemaAndTables()
	config.SetupRabbitMQ()
	go services.ConsumeBookingRequests()
	// var postgres = repository.GormDBWrapper{
	// 	DB: config.PostGresDB,
	// }
	// err := helpers.SeedSeats(&postgres)
	// if err != nil {
	// 	log.Fatal("Failed to seed seats:", err)
	// } else {
	// 	log.Println("Seats seeded successfully")
	// }

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

	//routes
	
	routes.EventRoutes(router)

	log.Infof("Server listening on http://localhost:8888/")
	if err := http.ListenAndServe("0.0.0.0:8888", router); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
