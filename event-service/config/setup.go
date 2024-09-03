package config

import (
	"fmt"
	"log"
	"os"

	"github.com/abidaziz9876/event-service/models"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}
}

var PostGresDB *gorm.DB

func ConnectPostGresDB() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USERNAME")
	dbname := os.Getenv("POSTGRES_DATABASE")
	password := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, "disable")
	client, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgresDB")
	PostGresDB = client
}
func SetupPostgresSchemaAndTables() {
	var Events models.Event
	var Seat models.Seat
	var SeatAllocation models.SeatAllocation
	PostGresDB.Exec("CREATE SCHEMA IF NOT EXISTS ticketbooking")

	err := PostGresDB.Table("ticketbooking.events").AutoMigrate(&Events)
	if err != nil {
		log.Fatal("Error creating events table: ", err)
	} else {
		log.Println("Events table created successfully")
	}

	err = PostGresDB.Table("ticketbooking.seats").AutoMigrate(&Seat)
	if err != nil {
		log.Fatal("Error creating events table: ", err)
	} else {
		log.Println("seats table created successfully")
	}

	err = PostGresDB.Table("ticketbooking.seat_allocation").AutoMigrate(&SeatAllocation)
	if err != nil {
		log.Fatal("Error creating events table: ", err)
	} else {
		log.Println("seat allocation table created successfully")
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

	// Declare queue
	_, err = ch.QueueDeclare(
		"booking-queue", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		"booking-queue",    // queue name
		"booking.key",      // routing key
		"booking-exchange", // exchange name
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatal("Failed to bind queue:", err)
	}

	err = declareQueue(ch, "payment-confirmation-queue")
	if err != nil {
		log.Fatal(err.Error()) 
	}

	err = declareExchange(ch, "payment-confirmation-exchange", "direct")
	if err != nil {
		log.Fatal(err.Error()) 
	}

	err = bindQueue(ch, "payment-confirmation-queue", "payment-confirmation-exchange", "payment.confirmation.key")
	if err != nil {
		log.Fatal(err.Error()) 
	}

	fmt.Println("Queue and Binding set up successfully.")
}