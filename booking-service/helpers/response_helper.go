package helpers

import (
	"github.com/abidaziz9876/booking-service/response"
	"github.com/gin-gonic/gin"
)

func ReturnResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	if data == nil {
		ctx.JSON(statusCode, response.ApiResponse{
			Status:  statusCode,
			Message: message,
		})
	} else {
		ctx.JSON(statusCode, response.ApiResponse{
			Status:  statusCode,
			Message: message,
			Data:    data,
		})
	}
}
