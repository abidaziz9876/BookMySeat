package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/abidaziz9876/event-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

var ErrorForPayment error

type Payment struct {
	UserID        uint   `json:"user_id"`
	Amount        int    `json:"amount"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	CorrelationID string `json:"correlation_id"` // Added field for tracking and correlating payment requests
}
type PaymentConfirmationResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"` // Optional field, only present if the payment fails
}


func PublishPaymentRequest(paymentRequest Payment) error {
	ch := config.RabbitMQ
	body, err := json.Marshal(paymentRequest)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"payment-exchange", // exchange
		"payment.key",      // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: paymentRequest.CorrelationID, // Include Correlation ID
		})
	if err != nil {
		return err
	}
	log.Info("inside publishPaymentRequest function")
	msgs, err := ch.Consume(
		"payment-confirmation-queue", // queue
		"",                           // consumer
		true,                         // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		return err
	}
	errorChan := make(chan error)
	go func() {
		for msg := range msgs {
			log.Infof("Received message with CorrelationID: %s", msg.CorrelationId)
			if msg.CorrelationId == paymentRequest.CorrelationID {
				var response PaymentConfirmationResponse
				err := json.Unmarshal(msg.Body, &response)
				if err != nil {
					errorChan <- fmt.Errorf("error unmarshaling response: %v", err)
				}
				if response.ErrorMessage != "" {
					errorChan <- fmt.Errorf("could not complete payment: %s", response.ErrorMessage)
				}
				if response.Success {
					log.Info("Payment Done Successfully")
				} else {
					errorChan <- fmt.Errorf("payment failed: %s", response.ErrorMessage)
				}
			}
		}
		close(errorChan)
	}()
	go func() error {
		for err := range errorChan {
			log.Error(err)
		}
		return err
	}()

	return nil
}
