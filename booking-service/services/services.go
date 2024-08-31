package services

import (
	"encoding/json"

	"github.com/abidaziz9876/booking-service/config"
	"github.com/abidaziz9876/booking-service/models"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

// func PublishBookingRequest(bookingRequest models.BookingRequest) error {
// 	ch := config.RabbitMQ

// 	body, err := json.Marshal(bookingRequest)
// 	if err != nil {
// 		log.Error("Error occurred "+err.Error())
// 		return err
// 	}

// 	err = ch.Publish(
// 		"booking-exchange", // exchange
// 		"booking.key",      // routing key
// 		false,              // mandatory
// 		false,              // immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        body,
// 		})

// 	if err!=nil{
// 		log.Error("Error occurred while publishing request to event "+err.Error())
// 		return err
// 	}
// 	return nil
// }

func PublishBookingRequest(bookingRequest models.BookingRequest) (string, error) {
	ch := config.RabbitMQ

	// Declare a temporary queue to receive the response
	replyQueue, err := ch.QueueDeclare(
		"",    // Name (empty string means a generated name)
		false, // Durable
		false, // Delete when unused
		true,  // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Error("Failed to declare a reply queue: " + err.Error())
		return "", err
	}

	// Set up a channel to receive the response
	msgs, err := ch.Consume(
		replyQueue.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Error("Failed to set up consumer for reply queue: " + err.Error())
		return "", err
	}

	// Marshal the booking request
	body, err := json.Marshal(bookingRequest)
	if err != nil {
		log.Error("Error occurred " + err.Error())
		return "", err
	}

	// Generate a correlation ID to track the request
	corrID := uuid.New().String()

	// Publish the message with a reply-to and correlation ID
	err = ch.Publish(
		"booking-exchange", // exchange
		"booking.key",      // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			ReplyTo:       replyQueue.Name, // Reply to this queue
			CorrelationId: corrID,          // Track this message
		})

	if err != nil {
		log.Error("Error occurred while publishing request to event " + err.Error())
		return "", err
	}

	// Wait for a response
	for d := range msgs {
		if d.CorrelationId == corrID {
			// This is the response we're waiting for
			return string(d.Body), nil
		}
	}

	return "", nil // No response or timeout
}
