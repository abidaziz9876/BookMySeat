package routes

import (
	"github.com/abidaziz9876/payment-service/controllers"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine) {
	router.POST("/create-payment", controllers.CreatePayment())
}
