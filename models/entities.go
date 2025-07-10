package models

import (
	"time"

	"gorm.io/datatypes"
)

type State struct {
	StateID   int    `gorm:"primaryKey;column:state_id" json:"state_id"`
	StateName string `gorm:"size:100;not null;column:state_name" json:"state_name"`
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
	TheatreImage    string    `gorm:"type:text;column:theatre_image" json:"theatre_image"`
	TheatreStatus   string    `gorm:"size:255;not null;column:theatre_status" json:"theatre_status" binding:"required"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Movie struct {
	MovieID          int            `gorm:"primaryKey;column:movie_id" json:"movie_id"`
	MovieName        string         `gorm:"size:255;not null;column:movie_name" json:"movie_name" binding:"required"`
	MovieDescription string         `gorm:"type:text;column:movie_description" json:"movie_description" binding:"required"`
	Duration         int            `gorm:"column:duration" json:"duration" binding:"required"`
	Languages        datatypes.JSON `gorm:"type:json;column:languages" json:"languages" binding:"required"`
	Genre            string         `gorm:"size:100;column:genre" json:"genre" binding:"required"`
	PosterURL        string         `gorm:"type:text;column:poster_url" json:"poster_url"`
	Rating           float64        `gorm:"type:numeric(2,1);column:rating" json:"rating"`
	StartDate        time.Time      `gorm:"column:start_date" json:"-"`
	EndDate          time.Time      `gorm:"column:end_date" json:"-"`
	MovieStatus      string         `gorm:"size:255;not null;column:movie_status" json:"movie_status" binding:"required"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

type Show struct {
	ShowID    uint           `gorm:"primaryKey" json:"show_id"`
	MovieID   int            `gorm:"column:movie_id" json:"movie_id"`
	Movie     Movie          `gorm:"foreignKey:MovieID;references:MovieID"`
	TheatreID int            `gorm:"column:theatre_id" json:"theatre_id"`
	Theatre   Theatre        `gorm:"foreignKey:TheatreID;references:TheatreID"`
	Date      time.Time      `json:"date"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Languages datatypes.JSON `json:"languages"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type User struct {
	UserID       int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name         string    `gorm:"size:255;not null;unique;column:name" json:"name"`
	Email        string    `gorm:"size:255;unique;not null;column:email" json:"email"`
	Phone        string    `gorm:"size:20;column:phone" json:"phone"`
	PasswordHash string    `gorm:"type:text;not null;column:password_hash" json:"-"`
	IsAdmin      bool      `gorm:"default:false;column:is_admin" json:"is_admin"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Booking struct {
	BookingID   int       `gorm:"primaryKey;column:booking_id" json:"booking_id"`
	UserID      int       `gorm:"not null;column:user_id" json:"user_id"`
	ShowID      uint      `gorm:"not null;column:show_id" json:"show_id"`
	TxnID       string    `gorm:"size:100;column:txn_id" json:"txn_id"`
	Amount      int       `gorm:"not null;column:amount" json:"amount"`
	Status      string    `gorm:"size:50;not null;column:status" json:"status"`
	Seats       string    `gorm:"size:255;column:seats" json:"seats"`
	BookingTime time.Time `gorm:"column:booking_time" json:"booking_time"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Transaction struct {
	TransactionID   int       `gorm:"primaryKey;column:transaction_id" json:"transaction_id"`
	BookingID       int       `gorm:"not null;column:booking_id" json:"booking_id"`
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
	TransactionID int       `gorm:"not null;column:transaction_id" json:"transaction_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Review struct {
	ReviewID  uint      `gorm:"primaryKey;column:review_id" json:"review_id"`
	UserID    uint      `gorm:"not null;index:idx_user_movie,unique;" json:"user_id"`
	MovieID   uint      `gorm:"not null;index:idx_user_movie,unique;" json:"movie_id"`
	Rating    int       `gorm:"check:rating >= 1 AND rating <= 5;column:rating" json:"rating"`
	Comments  string    `gorm:"type:text;column:comments" json:"comments"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type SeatBooking struct {
	ID        uint   `gorm:"primaryKey" json:"seat_id"`
	ShowID    uint   `json:"show_id"`
	Seat      string `json:"seat"`
	UserID    uint   `json:"user_id"`
	BarcodeID string `json:"barcode_id"`
}

type Payment struct {
	PaymentID       uint           `gorm:"primaryKey;column:payment_id" json:"payment_id"`
	BookingID       int            `gorm:"not null;column:booking_id" json:"booking_id"`
	GatewayTxnID    string         `gorm:"size:100;column:gateway_txn_id" json:"gateway_txn_id"`
	Amount          float64        `gorm:"type:numeric(10,2);not null;column:amount" json:"amount"`
	Status          string         `gorm:"size:50;not null;column:status" json:"status"`
	PaymentMethod   string         `gorm:"size:50;column:payment_method" json:"payment_method"`
	PaymentResponse datatypes.JSON `gorm:"type:json;column:payment_response" json:"payment_response"`
	TransactionTime time.Time      `gorm:"column:transaction_time" json:"transaction_time"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
}
type BookingDetailsResponse struct {
	TxnID    string   `json:"txnid"`
	Amount   float64  `json:"amount"`
	Status   string   `json:"status"`
	Seats    []string `json:"selectedSeats"`
	Movie    string   `json:"movie"`
	Theatre  string   `json:"theatre"`
	Date     string   `json:"date"`
	ShowTime string   `json:"time"`
	ShowID   int      `json:"show_id"`
}
