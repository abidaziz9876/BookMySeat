package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abidaziz9876/event-service/helpers"
	"github.com/abidaziz9876/event-service/models"
	"github.com/abidaziz9876/event-service/repository"
	"github.com/abidaziz9876/event-service/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetEventList(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var Events []models.Event
		result, err := services.GetEventList(postgres)
		if err != nil {
			log.Error("something went wrong " + err.Error())
			helpers.ReturnResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		helpers.ReturnResponse(ctx, http.StatusOK, "success", result)
	}
}

func CreateEvents(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var Event models.Event
		err := ctx.BindJSON(&Event)
		if err != nil {
			log.Error("could not bind json")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		// newEvent := Event{
		// 	Title:       "Avengers: Endgame",
		// 	Description: "The epic conclusion to the Infinity Saga.",
		// 	Language:    "English",
		// 	Genre:       "Action",
		// 	Rating:      9.5,
		// }

		// Save the event to the database
		err = services.CreateEvent(postgres, &Event)
		if err != nil {
			log.Error("could not create event into database")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		// Create event occurrences for the new event
		occurrences := []models.EventOccurrences{
			{
				EventID:  Event.ID,
				Venue:    "Cinema Hall 1",
				Date:     time.Date(2024, 8, 25, 19, 0, 0, 0, time.UTC),
				TimeSlot: "Evening",
				Duration: 120,
				Price:    140,
			},
			{
				EventID:  Event.ID,
				Venue:    "Cinema Hall 2",
				Date:     time.Date(2024, 8, 26, 14, 0, 0, 0, time.UTC),
				TimeSlot: "Afternoon",
				Duration: 120,
				Price:    100,
			},
		}

		// Save the event occurrences to the database
		err = services.CreateOccurrence(postgres, occurrences)
		if err != nil {
			log.Error("could not insert event into database")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		for _, occurrence := range occurrences {
			err = services.AllocateSeatsForOccurrence(postgres, occurrence)
			if err != nil {
				log.Error("could not allocate seats for occurrence")
				helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
				return
			}
		}

		helpers.ReturnResponse(ctx, http.StatusOK, "events inserted", nil)
	}
}

func GetEvent(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			log.Error("id not found in the query")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "id not found in the query", nil)
			return
		}
		var id uint
		if i, err := strconv.ParseUint(ID, 10, 64); err == nil {
			id = uint(i)
		} else {
			log.Error("invalid id format")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "invalid id format", nil)
			return
		}
		result, err := services.GetEventVenueAndTime(postgres, id)
		if err != nil {
			log.Error("something went wrong " + err.Error())
			helpers.ReturnResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		helpers.ReturnResponse(ctx, http.StatusOK, "success", result)
	}
}

func DeleteEvents(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			log.Error("id not found in the query")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "id not found in the query", nil)
			return
		}
		var id uint
		if i, err := strconv.ParseUint(ID, 10, 64); err == nil {
			id = uint(i)
		} else {
			log.Error("invalid id format")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "invalid id format", nil)
			return
		}
		err := services.DeleteEvent(postgres, id)
		if err != nil {
			log.Error("Error deleting event " + err.Error())
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "error occurred "+err.Error(), nil)
			return
		}
		helpers.ReturnResponse(ctx, http.StatusOK, "event deleted", nil)
	}
}
