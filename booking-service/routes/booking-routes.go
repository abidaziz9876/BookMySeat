package routes

import (
	"github.com/abidaziz9876/booking-service/config"
	"github.com/abidaziz9876/booking-service/controllers"
	"github.com/abidaziz9876/booking-service/repository"
	"github.com/gin-gonic/gin"
)

func BookingRoutes(router *gin.Engine) {
	var postgres = repository.GormDBWrapper{
		DB: config.PostGresDB,
	}
	router.GET("/",controllers.CheckBooking())
	router.POST("/bookticket", controllers.BookingRequest(&postgres))
}
