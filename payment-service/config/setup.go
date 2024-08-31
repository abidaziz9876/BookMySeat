package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}
}
func declareQueue(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	return nil
}
func declareExchange(ch *amqp.Channel, exchangeName, exchangeType string) error {
	err := ch.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // exchange type
		true,         // durable
		false,        // delete when unused
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func bindQueue(ch *amqp.Channel, queueName, exchangeName, routingKey string) error {
	err := ch.QueueBind(
		queueName,      // queue name
		routingKey,     // routing key
		exchangeName,   // exchange name
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}
	return nil
}
var RabbitMQ *amqp.Channel

func SetupRabbitMQ() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	fmt.Println("Connected to RabbitMQ")
	RabbitMQ = ch
	
	_, err = ch.QueueDeclare(
		"payment-queue", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}
	// Declare exchange
	err = ch.ExchangeDeclare(
		"payment-exchange", // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	// Declare queue
	

	// Bind queue to exchange
	err = ch.QueueBind(
		"payment-queue",    // queue name
		"payment.key",      // routing key
		"payment-exchange", // exchange name
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatal("Failed to bind queue:", err)
	}
	

	// Declare queue
	err = declareQueue(ch, "payment-confirmation-queue")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = declareExchange(ch, "payment-confirmation-exchange", "direct")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Bind queue to exchange
	err = bindQueue(ch, "payment-confirmation-queue", "payment-confirmation-exchange", "payment.confirmation.key")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Exchange, Queue, and Binding set up successfully.")
}
