package controllers

import (
	"net/http"
	"time"

	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateShow(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var show models.Show
	if err := c.ShouldBindJSON(&show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var movie models.Movie
	if err := db.First(&movie, "movie_id = ?", show.MovieID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie_id: movie not found"})
		return
	}

	var theatre models.Theatre
	if err := db.First(&theatre, "theatre_id = ?", show.TheatreID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre_id: theatre not found"})
		return
	}

	show.CreatedAt = time.Now()
	show.UpdatedAt = time.Now()

	if err := db.Create(&show).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, show)
}

func GetShowByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var show models.Show
	if err := db.Preload("Movie").Preload("Theatre").First(&show, "show_id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, show)
}

func GetShows(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var shows []models.Show
	if err := db.Preload("Movie").Preload("Theatre").Find(&shows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}

func UpdateShow(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var show models.Show
	if err := db.First(&show, "show_id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input models.Show
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var movie models.Movie
	if err := db.First(&movie, "movie_id = ?", input.MovieID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie_id: movie not found"})
		return
	}

	var theatre models.Theatre
	if err := db.First(&theatre, "theatre_id = ?", input.TheatreID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre_id: theatre not found"})
		return
	}

	show.MovieID = input.MovieID
	show.TheatreID = input.TheatreID
	show.ShowTime = input.ShowTime
	show.TotalSeats = input.TotalSeats
	show.Price = input.Price
	show.ShowDuration = input.ShowDuration
	show.UpdatedAt = time.Now()

	if err := db.Save(&show).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, show)
}

func DeleteShow(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var show models.Show
	if err := db.First(&show, "show_id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Delete(&show).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Show deleted successfully"})
}
