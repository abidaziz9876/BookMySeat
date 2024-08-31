package config

import (
	"fmt"
	"log"
	"os"

	"github.com/abidaziz9876/booking-service/models"
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
	var booking models.Booking
	PostGresDB.Exec("CREATE SCHEMA IF NOT EXISTS ticketbooking")

	err := PostGresDB.Table("ticketbooking.booking").AutoMigrate(&booking)
	if err != nil {
		log.Fatal(err.Error())
	}

}

var RabbitMQ *amqp.Channel

func SetupRabbitMQ() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	fmt.Println("Connected to RabbitMQ")
	RabbitMQ = ch

	// Declare exchange
	err = ch.ExchangeDeclare(
		"booking-exchange", // name
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

	fmt.Println("Exchange, Queue, and Binding set up successfully.")

	// No need to defer `ch.Close()` here, since RabbitMQ is being used in the rest of the application.
	// Defer `conn.Close()` when the application shuts down.
}
