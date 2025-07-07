package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.AutoMigrate(&State{}); err != nil {
		log.Fatal("Failed to migrate State:", err)
	}

	if err := database.AutoMigrate(&City{}); err != nil {
		log.Fatal("Failed to migrate City:", err)
	}

	err = database.Exec(`
    ALTER TABLE cities DROP CONSTRAINT IF EXISTS fk_states_cities;
    ALTER TABLE cities
    ADD CONSTRAINT fk_states_cities
    FOREIGN KEY (state_id)
    REFERENCES states(state_id)
    ON DELETE CASCADE
`).Error
	if err != nil {
		log.Fatal("Failed to add foreign key constraint:", err)
	}

	if err := database.AutoMigrate(
		&Theatre{},
		&Movie{},
		&Show{},
		&User{},
		&Booking{},
		&Transaction{},
		&Ticket{},
		&Review{},
	); err != nil {
		log.Fatal("Failed to migrate other models:", err)
	}

	DB = database
}
