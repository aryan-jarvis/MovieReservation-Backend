package controllers

import (
	"backend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBookingDetails(c *gin.Context) {
	txnID := c.Param("txnid")

	type result struct {
		TxnID    string
		Amount   float64
		Status   string
		Seats    string
		ShowID   int
		Movie    string
		Theatre  string
		Date     string
		ShowTime string
	}

	var r result
	if err := models.DB.Table("bookings").
		Select(`
		bookings.txn_id,
		bookings.amount,
		bookings.status,
		bookings.seats,
		bookings.show_id,
		movies.movie_name AS movie,
		theatres.theatre_name AS theatre,
		shows.date,
		shows.start_time
	`).
		Joins(`JOIN shows ON bookings.show_id = shows.show_id`).
		Joins(`JOIN movies ON shows.movie_id = movies.movie_id`).
		Joins(`JOIN theatres ON shows.theatre_id = theatres.theatre_id`).
		Where("bookings.txn_id = ?", txnID).
		Scan(&r).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	var seatList []string
	if r.Seats != "" {
		seatList = strings.Split(r.Seats, ",")
	}

	response := models.BookingDetailsResponse{
		TxnID:    r.TxnID,
		Amount:   r.Amount,
		Status:   r.Status,
		Seats:    seatList,
		Movie:    r.Movie,
		Theatre:  r.Theatre,
		Date:     r.Date,
		ShowTime: r.ShowTime,
		ShowID:   r.ShowID,
	}

	c.JSON(http.StatusOK, response)
}
