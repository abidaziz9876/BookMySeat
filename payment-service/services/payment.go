package services

import (
	"encoding/json"

	"os"

	"github.com/abidaziz9876/payment-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type PaymentConfirmationResponse struct {
	CorrelationID string `json:"correlation_id"`
	Success       bool   `json:"success"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

type Payment struct {
	UserID        uint   `json:"user_id"`
	Amount        int    `json:"amount"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	CorrelationID string `json:"correlation_id"`
}

func CreatePayment(payment Payment) error {
	stripe.Key = os.Getenv("STRIPE_KEY")
	FullName := payment.FirstName + " " + payment.LastName

	_, err := charge.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(int64(payment.Amount)),
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Description:  stripe.String("payment for ticket booking"),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		ReceiptEmail: stripe.String(payment.Email),
		Shipping: &stripe.ShippingDetailsParams{
			Phone: &payment.Phone,
			Name:  &FullName,
			Address: &stripe.AddressParams{
				Line1:      stripe.String("1234 Main street"),
				City:       stripe.String("San Francisco"),
				State:      stripe.String("CA"),
				PostalCode: stripe.String("94111"),
				Country:    stripe.String("US"),
			},
		},
	})

	if err != nil {

		return err
	}

	log.Println("Payment successful")

	return nil
}

func publishPaymentConfirmation(corrID string, success bool, errorMsg string) error {
	response := PaymentConfirmationResponse{
		Success:      success,
		ErrorMessage: errorMsg,
	}
	ch := config.RabbitMQ
	body, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// Publish without re-declaring the queue
	err = ch.Publish(
		"payment-confirmation-exchange", // exchange
		"payment.confirmation.key",      // routing key
		false,                           // mandatory
		false,                           // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: corrID, // Include the CorrelationID here
		})
	if err != nil {
		return err
	}
	log.Infof("Publishing payment confirmation with CorrelationID: %s", corrID)
	log.Info("Payment confirmation published to RabbitMQ")
	return nil
}

func ConsumePaymentRequests() {
	ch := config.RabbitMQ
	log.Info("Waiting for payment requests...")

	msgs, err := ch.Consume(
		"payment-queue", // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}
	for{
		for msg := range msgs {
			var paymentRequest Payment
			err := json.Unmarshal(msg.Body, &paymentRequest)
			if err != nil {
				log.Error("Failed to unmarshal payment request:", err)
				return
			}
	
			// Process the payment
			err = CreatePayment(paymentRequest)
			if err != nil {
				errPublish := publishPaymentConfirmation(paymentRequest.CorrelationID, false, err.Error())
				if errPublish != nil {
					log.Printf("Failed to publish payment confirmation: %v", errPublish)
				}
			}else{
				errPublish:=publishPaymentConfirmation(paymentRequest.CorrelationID, true, "")
				if errPublish != nil {
					log.Printf("Failed to publish payment confirmation: %v", errPublish)
				}
			}
		}
	}
}
