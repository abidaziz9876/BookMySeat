package routes

import (
	"github.com/abidaziz9876/event-service/config"
	"github.com/abidaziz9876/event-service/controllers"
	"github.com/abidaziz9876/event-service/repository"
	"github.com/gin-gonic/gin"
)

func EventRoutes(router *gin.Engine) {
	var postgres = repository.GormDBWrapper{
		DB: config.PostGresDB,
	}
	router.GET("/getevents", controllers.GetEventList(&postgres))
	router.POST("/events/create-events", controllers.CreateEvents(&postgres))
	router.DELETE("/delete-events", controllers.DeleteEvents(&postgres))
	router.GET("/events/getevent", controllers.GetEvent(&postgres))
}
