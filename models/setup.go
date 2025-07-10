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

	if err := database.AutoMigrate(&Theatre{}); err != nil {
		log.Fatal("Failed to migrate Theatre:", err)
	}

	if err := database.AutoMigrate(&Movie{}); err != nil {
		log.Fatal("Failed to migrate Movie:", err)
	}

	if err := database.AutoMigrate(&Show{}); err != nil {
		log.Fatal("Failed to migrate Show:", err)
	}

	if err := database.AutoMigrate(&User{}); err != nil {
		log.Fatal("Failed to migrate User:", err)
	}

	if err := database.AutoMigrate(&Booking{}); err != nil {
		log.Fatal("Failed to migrate Booking:", err)
	}

	if err := database.AutoMigrate(&Transaction{}); err != nil {
		log.Fatal("Failed to migrate Transaction:", err)
	}

	if err := database.AutoMigrate(&Ticket{}); err != nil {
		log.Fatal("Failed to migrate Ticket:", err)
	}

	if err := database.AutoMigrate(&Review{}); err != nil {
		log.Fatal("Failed to migrate Review:", err)
	}

	if err := database.AutoMigrate(&SeatBooking{}); err != nil {
		log.Fatal("Failed to migrate SeatBooking:", err)
	}

	if err := database.AutoMigrate(&Payment{}); err != nil {
		log.Fatal("Failed to migrate Payment:", err)
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
		log.Fatal("Failed to add foreign key constraint for cities:", err)
	}

	err = database.Exec(`
    ALTER TABLE shows DROP CONSTRAINT IF EXISTS fk_theatres_shows;
    ALTER TABLE shows
    ADD CONSTRAINT fk_theatres_shows
    FOREIGN KEY (theatre_id)
    REFERENCES theatres(theatre_id)
    ON DELETE CASCADE;

    ALTER TABLE shows DROP CONSTRAINT IF EXISTS fk_movies_shows;
    ALTER TABLE shows
    ADD CONSTRAINT fk_movies_shows
    FOREIGN KEY (movie_id)
    REFERENCES movies(movie_id)
    ON DELETE CASCADE;
`).Error
	if err != nil {
		log.Fatal("Failed to add foreign key constraints for shows:", err)
	}

	err = database.Exec(`
	ALTER TABLE seat_bookings DROP CONSTRAINT IF EXISTS fk_shows_seat_bookings;
	ALTER TABLE seat_bookings
	ADD CONSTRAINT fk_shows_seat_bookings
	FOREIGN KEY (show_id)
	REFERENCES shows(show_id)
	ON DELETE CASCADE;

	ALTER TABLE seat_bookings DROP CONSTRAINT IF EXISTS fk_users_seat_bookings;
	ALTER TABLE seat_bookings
	ADD CONSTRAINT fk_users_seat_bookings
	FOREIGN KEY (user_id)
	REFERENCES users(user_id)
	ON DELETE CASCADE;

	ALTER TABLE seat_bookings DROP CONSTRAINT IF EXISTS uq_show_seat;
	ALTER TABLE seat_bookings
	ADD CONSTRAINT uq_show_seat
	UNIQUE (show_id, seat);
`).Error
	if err != nil {
		log.Fatal("Failed to add constraints for seat_bookings:", err)
	}

	DB = database
}
