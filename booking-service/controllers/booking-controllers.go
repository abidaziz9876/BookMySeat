package controllers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/abidaziz9876/booking-service/models"
	"github.com/abidaziz9876/booking-service/repository"
	"github.com/abidaziz9876/booking-service/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func BookingRequest(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bookingRequest models.BookingRequest
		// Bind the JSON payload to the BookingRequest struct
		if err := ctx.BindJSON(&bookingRequest); err != nil {
			// Handle error
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Publish the booking request to the event service
		log.Info("booking process started")
		var msg string
		var err error
		if msg, err = services.PublishBookingRequest(bookingRequest); err != nil {
			log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			log.Info(msg)
			var seatDetails []models.Seat
			var message string
			query := `select * from ticketbooking.seats where id IN (?)`
			res := postgres.RawQueryWithFind(&seatDetails, query, bookingRequest.SeatIDs)
			if res.Error != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": res.Error.Error()})
				return
			} else {
				for i := 0; i < len(seatDetails); i++ {
					if i == 0 {
						message += fmt.Sprintf("Venue: %s\n", seatDetails[i].Venue)
					}
					message += fmt.Sprintf("person %d Row %s seat no: %d \n", i+1, seatDetails[i].Row, seatDetails[i].Number)
				}
				password:=os.Getenv("APP_PASSWORD")
				auth := smtp.PlainAuth("", "abidaziz9876@gmail.com",password , "smtp.gmail.com")

				// Here we do it all: connect to our server, set up a message and send it

				to := []string{"azizabid98766@gmail.com"}

				msg := []byte("To: azizabid98766@gmail.com\r\n" +

					"Subject: Seat Details of Movie Ticket\r\n" +

					"\r\n" +

					message+"\r\n")

				err := smtp.SendMail("smtp.gmail.com:587", auth, "abidaziz9876@gmail.com", to, msg)

				if err != nil {
					log.Error(err)
				}
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"message": msg})
	}
}


func CheckBooking() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,"booking ok")
	}
}