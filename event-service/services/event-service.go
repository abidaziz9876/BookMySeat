package services

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/abidaziz9876/event-service/config"
	"github.com/abidaziz9876/event-service/helpers"
	"github.com/abidaziz9876/event-service/models"
	"github.com/abidaziz9876/event-service/repository"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentConfirmationResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"` // Optional field, only present if the payment fails
}

type EventVenueAndTime struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	EventID  uint      `gorm:"not null;index" json:"eventId"`
	Venue    string    `gorm:"size:255;not null" json:"venue"`
	Date     time.Time `gorm:"not null" json:"date"`
	Duration int       `gorm:"not null" json:"duration"`
	Price    float64   `gorm:"not null" json:"price"`
}

type UserModel struct {
	ID        int64  `gorm:"column:id"`
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
	Phone     string `json:"phone" gorm:"column:phone"`
	Email     string `json:"email" gorm:"column:email"`
	Password  string `json:"password" gorm:"column:password"`
}
type Payment struct {
	UserID        uint   `json:"user_id"`
	Amount        int    `json:"amount"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	CorrelationID string `json:"correlation_id"` // Added field for tracking and correlating payment requests
}

func CreateEvent(postgres repository.PostgesDatabase, Event *models.Event) error {
	result := postgres.CreateEvent(Event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateOccurrence(postgres repository.PostgesDatabase, Occurrence []models.EventOccurrences) error {
	result := postgres.CreateEventOccurrences(Occurrence)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetEventList(postgres repository.PostgesDatabase) ([]models.Event, error) {

	var result []models.Event
	query := `SELECT * FROM ticketbooking.events`
	err := postgres.RawQueryWithFind(&result, query).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetEventVenueAndTime(postgres repository.PostgesDatabase, Id uint) ([]EventVenueAndTime, error) {
	var result []EventVenueAndTime
	query := `SELECT oc.id, oc.event_id, oc.venue, oc.date, oc.duration, oc.price
	FROM ticketbooking.occurrence oc 
	JOIN ticketbooking.events e ON oc.event_id = e.id 
	WHERE oc.event_id = ?`
	err := postgres.RawQueryWithFind(&result, query, Id).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteEvent(postgres repository.PostgesDatabase, id uint) error {
	type Event struct {
		ID uint `gorm:"primaryKey"`
	}

	// Create an instance of Event with the given ID
	event := Event{ID: id}

	// Perform the delete operation
	if err := postgres.Delete(&event).Error; err != nil {
		return err
	}

	return nil
}

func AllocateSeatsForOccurrence(postgres repository.PostgesDatabase, occurrence models.EventOccurrences) error {
	// Retrieve seats for the venue
	var seats []models.Seat
	err := postgres.RawQueryWithFind(&seats, "SELECT * FROM ticketbooking.seats WHERE venue = ?", occurrence.Venue).Error
	if err != nil {
		return err
	}

	// Allocate each seat to the occurrence
	for _, seat := range seats {
		seatAllocation := models.SeatAllocation{
			OccurrenceID: occurrence.ID,
			SeatID:       seat.ID,
			Status:       "Available",
		}
		err = postgres.CreateSeatAllocation(&seatAllocation).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func ConsumeBookingRequests() error {
	ch := config.RabbitMQ
	if ch == nil {
		return errors.New("RabbitMQ connection is not established")
	}
	msgs, err := ch.Consume(
		"booking-queue", // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Error("Can not consume the booking in the event: " + err.Error())
		return err
	}

	log.Info("Waiting for booking requests...")
	for {
		for msg := range msgs {
			log.Info("inside the messages of consume booking")
			var bookingRequest models.BookingRequest
			err := json.Unmarshal(msg.Body, &bookingRequest)
			if err != nil {
				log.Error("Failed to parse booking request: " + err.Error())
				continue // Continue to the next message
			}

			// Process the booking
			err = processBooking(bookingRequest)

			// Prepare the response
			var response string
			if err != nil {
				log.Error("Error occurred while processing booking: " + err.Error())
				response = "Booking Failed"
			} else {
				response = "Booking Success"
			}

			// Send a response back to the publisher
			err = ch.Publish(
				"",          // exchange
				msg.ReplyTo, // routing key (reply queue)
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: msg.CorrelationId,
					Body:          []byte(response),
				})
			if err != nil {
				log.Error("Failed to publish response: " + err.Error())
			}
		}
	}
}

// func processBooking(bookingRequest models.BookingRequest) error {
// 	db := config.PostGresDB
// 	var postgres = repository.GormDBWrapper{
// 		DB: db,
// 	}
// 	log.Info("inside processBooking")
// 	// 1. Validate the booking request
// 	if len(bookingRequest.SeatIDs) == 0 {
// 		return errors.New("no seats selected")
// 	}

// 	// 2. Check if the event and occurrence exist
// 	var occurrence models.EventOccurrences
// 	err := db.Where("id = ? AND event_id = ?", bookingRequest.OccurrenceID, bookingRequest.EventID).First(&occurrence).Error
// 	if err != nil {
// 		return errors.New("invalid event or occurrence")
// 	}

// 	// 3. Check seat availability
// 	var cnt int
// 	query := `SELECT count(*) FROM ticketbooking.seat_allocation WHERE occurrence_id = ? AND seat_id IN (?) AND status != 'Available'`
// 	err = postgres.RawQueryWithFind(&cnt, query, bookingRequest.OccurrenceID, bookingRequest.SeatIDs).Error
// 	if err != nil {
// 		return err
// 	}
// 	if cnt > 0 {
// 		return errors.New("some seats are already booked")
// 	}

// 	// 4. Calculate total amount
// 	var Price float64
// 	query = `SELECT price FROM ticketbooking.occurrence WHERE id = ?`
// 	err = postgres.RawQueryWithFind(&Price, query, bookingRequest.OccurrenceID).Error
// 	if err != nil {
// 		return err
// 	}
// 	price := int(math.Round(Price))
// 	totalAmount := len(bookingRequest.SeatIDs) * price
// 	var user UserModel
// 	query = `select * from ticketbooking.users u where u.id = ?`
// 	err = postgres.RawQueryWithFind(&user, query, bookingRequest.UserID).Error
// 	if err != nil {
// 		log.Error("can not find user details " + err.Error())
// 		return err
// 	}
// 	// 5. Create payment request and publish to the payment service
// 	corrID := uuid.New().String()
// 	payment := helpers.Payment{
// 		Amount:        totalAmount,
// 		UserID:        bookingRequest.UserID,
// 		Email:         user.Email,
// 		FirstName:     user.FirstName,
// 		LastName:      user.LastName,
// 		Phone:         user.Phone,
// 		CorrelationID: corrID, // Correlation ID
// 	}

// 	err=helpers.PublishPaymentRequest(payment)
// 	if err!=nil{
// 		log.Error("error occurred")
// 		return err
// 	}

// 	log.Info("return after published payment request")

// 	// 7. Proceed with seat allocation and booking creation
// 	query = `UPDATE ticketbooking.seat_allocation SET status=?, user_id =?, booked_at = ? WHERE seat_id IN (?) AND occurrence_id = ?`
// 	result := db.Exec(query, "booked", bookingRequest.UserID, time.Now(), bookingRequest.SeatIDs, bookingRequest.OccurrenceID)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	booking := models.Booking{
// 		UserID:       bookingRequest.UserID,
// 		EventID:      bookingRequest.EventID,
// 		OccurrenceID: bookingRequest.OccurrenceID,
// 		TotalAmount:  float64(totalAmount),
// 		Status:       "Confirmed",
// 	}
// 	err = postgres.CreateBooking(&booking).Error
// 	if err != nil {
// 		return err
// 	}
// 	log.Info("Booking processed successfully")
// 	return nil
// }

func processBooking(bookingRequest models.BookingRequest) error {
	db := config.PostGresDB
	// var postgres = repository.GormDBWrapper{
	// 	DB: db,
	// }
	// time.Sleep(10*time.Second)
	return db.Transaction(func(tx *gorm.DB) error {
		log.Info("inside processBooking")

		// 1. Validate the booking request

		if len(bookingRequest.SeatIDs) == 0 {
			return errors.New("no seats selected")
		}

		// 2. Check if the event and occurrence exist
		var occurrence models.EventOccurrences
		err := tx.Where("id = ? AND event_id = ?", bookingRequest.OccurrenceID, bookingRequest.EventID).First(&occurrence).Error
		if err != nil {
			return errors.New("invalid event or occurrence")
		}

		// 3. Lock the seats for this transaction
		var seats []models.SeatAllocation
		err = tx.Table("ticketbooking.seat_allocation").Clauses(clause.Locking{Strength: "UPDATE"}).Where("occurrence_id = ? AND seat_id IN (?) AND status = 'Available'", bookingRequest.OccurrenceID, bookingRequest.SeatIDs).Find(&seats).Error
		if err != nil {
			return err
		}
		if len(seats) != len(bookingRequest.SeatIDs) {
			return errors.New("some seats are already booked")
		}

		// 4. Calculate total amount
		var Price float64
		err = tx.Raw(`SELECT price FROM ticketbooking.occurrence WHERE id = ?`, bookingRequest.OccurrenceID).Scan(&Price).Error
		if err != nil {
			log.Error("error occuured in occurrence "+err.Error())
			return err
		}
		price := int(math.Round(Price))
		totalAmount := len(bookingRequest.SeatIDs) * price

		var user UserModel
		err = tx.Raw(`select * from ticketbooking.users u where u.id = ?`, bookingRequest.UserID).Scan(&user).Error
		if err != nil {
			log.Error("error occuured in users "+err.Error())
			return err
		}

		// 5. Create payment request and publish to the payment service
		corrID := uuid.New().String()
		payment := helpers.Payment{
			Amount:        totalAmount,
			UserID:        bookingRequest.UserID,
			Email:         user.Email,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Phone:         user.Phone,
			CorrelationID: corrID, // Correlation ID
		}

		
		// Publish the payment request
		err = helpers.PublishPaymentRequest(payment)
		if err != nil {
			return err
		}

		log.Info("waiting for payment confirmation")

		err = tx.Exec(`UPDATE ticketbooking.seat_allocation SET status = ?, user_id = ?, booked_at = ? WHERE seat_id IN (?) AND occurrence_id = ?`,
			"booked", bookingRequest.UserID, time.Now(), bookingRequest.SeatIDs, bookingRequest.OccurrenceID).Error
		if err != nil {
			return err
		}

		booking := models.Booking{
			UserID:       bookingRequest.UserID,
			EventID:      bookingRequest.EventID,
			OccurrenceID: bookingRequest.OccurrenceID,
			TotalAmount:  float64(totalAmount),
			Status:       "Confirmed",
		}
		err = tx.Table("ticketbooking.booking").Create(&booking).Error
		if err != nil {
			return err
		}
		log.Info("Booking processed successfully")

		return nil
	})
}
