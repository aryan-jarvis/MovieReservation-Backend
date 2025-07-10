package controllers

import (
	"net/http"
	"strconv"

	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BookSeat(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var booking models.SeatBooking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if booking.Seat == "" || booking.ShowID == 0 || booking.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seat, ShowID, and UserID are required"})
		return
	}

	var existing models.SeatBooking
	if err := db.Where("show_id = ? AND seat = ?", booking.ShowID, booking.Seat).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Seat already booked"})
		return
	}

	if err := db.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book seat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Seat booked successfully"})
}

func GetBookedSeats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	showIDStr := c.Param("id")
	showID, err := strconv.Atoi(showIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid show ID"})
		return
	}

	var bookings []models.SeatBooking
	if err := db.Where("show_id = ?", showID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}

	seats := make([]string, len(bookings))
	for i, b := range bookings {
		seats[i] = b.Seat
	}

	c.JSON(http.StatusOK, gin.H{"bookedSeats": seats})
}
