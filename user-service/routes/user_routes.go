package routes

import (
	"github.com/abidaziz9876/user-service/config"
	"github.com/abidaziz9876/user-service/controllers"
	"github.com/abidaziz9876/user-service/repository"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	var postgres = repository.GormDBWrapper{
		DB: config.PostGresDB,
	}
	
	
	router.POST("/users/signup",controllers.SignUp(&postgres))
	router.POST("/users/signin",controllers.LogIn(&postgres))
	router.PUT("users/updateuser", controllers.UpdateUser(&postgres))
	router.DELETE("/users/deleteuser", controllers.DeleteUser(&postgres))
	router.GET("/",controllers.Check())
}
