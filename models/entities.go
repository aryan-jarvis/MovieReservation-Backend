package models

import "time"

type State struct {
	StateID   int    `gorm:"primaryKey;column:state_id" json:"state_id"`
	StateName string `gorm:"size:100;not null;column:state_name" json:"state_name"`
	Cities    []City `gorm:"constraint:OnDelete:CASCADE;foreignKey:StateID;references:StateID"`
}

type City struct {
	CityID   int    `gorm:"primaryKey;column:city_id" json:"city_id"`
	StateID  int    `gorm:"not null;column:state_id" json:"state_id"`
	CityName string `gorm:"size:100;not null;column:city_name" json:"city_name"`
}

type Theatre struct {
	TheatreID       int       `gorm:"primaryKey;column:theatre_id" json:"theatre_id"`
	TheatreName     string    `gorm:"size:255;not null;column:theatre_name" json:"theatre_name"`
	TheatreLocation string    `gorm:"size:255;column:theatre_location" json:"theatre_location"`
	CityID          int       `gorm:"not null;column:city_id" json:"city_id"`
	TotalSeats      int       `gorm:"not null;column:total_seats" json:"total_seats"`
	TheatreTiming   string    `gorm:"size:255;column:theatre_timing" json:"theatre_timing"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Movie struct {
	MovieID          int       `gorm:"primaryKey;column:movie_id" json:"movie_id"`
	MovieName        string    `gorm:"size:255;not null;column:movie_name" json:"movie_name"`
	MovieDescription string    `gorm:"type:text;column:movie_description" json:"movie_description"`
	Duration         int       `gorm:"column:duration" json:"duration"`
	Language         string    `gorm:"size:50;column:language" json:"language"`
	Genre            string    `gorm:"size:100;column:genre" json:"genre"`
	PosterURL        string    `gorm:"type:text;column:poster_url" json:"poster_url"`
	Rating           float64   `gorm:"type:numeric(2,1);column:rating" json:"rating"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Show struct {
	ShowID       int       `gorm:"primaryKey;column:show_id" json:"show_id"`
	TheatreID    int       `gorm:"not null;column:theatre_id" json:"theatre_id"`
	MovieID      int       `gorm:"not null;column:movie_id" json:"movie_id"`
	ShowTime     time.Time `gorm:"not null;column:show_time" json:"show_time"`
	TotalSeats   int       `gorm:"not null;column:total_seats" json:"total_seats"`
	Price        float64   `gorm:"type:numeric(10,2);not null;column:price" json:"price"`
	ShowDuration int       `gorm:"column:show_duration" json:"show_duration"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type User struct {
	UserID       int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name         string    `gorm:"size:255;not null;column:name" json:"name"`
	Email        string    `gorm:"size:255;unique;not null;column:email" json:"email"`
	Phone        string    `gorm:"size:20;column:phone" json:"phone"`
	PasswordHash string    `gorm:"type:text;not null;column:password_hash" json:"password_hash"`
	IsAdmin      bool      `gorm:"default:false;column:is_admin" json:"is_admin"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Booking struct {
	BookingID     int       `gorm:"primaryKey;column:booking_id" json:"booking_id"`
	BookingStatus string    `gorm:"size:50;not null;column:booking_status" json:"booking_status"`
	UserID        int       `gorm:"not null;column:user_id;constraint:OnDelete:CASCADE;" json:"user_id"`
	TotalTickets  int       `gorm:"not null;column:total_tickets" json:"total_tickets"`
	ShowID        int       `gorm:"not null;column:show_id;constraint:OnDelete:CASCADE;" json:"show_id"`
	BookingTime   time.Time `gorm:"column:booking_time" json:"booking_time"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Transaction struct {
	TransactionID   int       `gorm:"primaryKey;column:transaction_id" json:"transaction_id"`
	BookingID       int       `gorm:"not null;column:booking_id;constraint:OnDelete:CASCADE;" json:"booking_id"`
	TotalAmount     float64   `gorm:"type:numeric(10,2);not null;column:total_amount" json:"total_amount"`
	PaymentMethod   string    `gorm:"size:50;column:payment_method" json:"payment_method"`
	TransactionTime time.Time `gorm:"column:transaction_time" json:"transaction_time"`
	PaymentStatus   string    `gorm:"size:50;column:payment_status" json:"payment_status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Ticket struct {
	TicketID      int       `gorm:"primaryKey;column:ticket_id" json:"ticket_id"`
	Amount        float64   `gorm:"type:numeric(10,2);not null;column:amount" json:"amount"`
	TransactionID int       `gorm:"not null;column:transaction_id;constraint:OnDelete:CASCADE;" json:"transaction_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Review struct {
	ReviewID  int       `gorm:"primaryKey;column:review_id" json:"review_id"`
	UserID    int       `gorm:"not null;column:user_id;constraint:OnDelete:CASCADE;" json:"user_id"`
	ShowID    int       `gorm:"not null;column:show_id;constraint:OnDelete:CASCADE;" json:"show_id"`
	Rating    int       `gorm:"check:rating >= 1 AND rating <= 5;column:rating" json:"rating"`
	Comments  string    `gorm:"type:text;column:comments" json:"comments"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
