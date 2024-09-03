package config

import (
	"fmt"
	"log"
	"os"

	"github.com/abidaziz9876/user-service/models"
	"github.com/joho/godotenv"
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
	var Users models.UserModel

	PostGresDB.Exec("CREATE SCHEMA IF NOT EXISTS ticketbooking")

	err := PostGresDB.Table("ticketbooking.users").AutoMigrate(&Users)
	if err != nil {
		log.Fatal(err.Error())
	}

}
