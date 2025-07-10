package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateShow(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input struct {
		MovieID   int      `json:"movie_id" binding:"required"`
		TheatreID int      `json:"theatre_id" binding:"required"`
		Date      string   `json:"date" binding:"required"`
		Languages []string `json:"languages" binding:"required"`
		Times     []struct {
			StartTime string `json:"start_time" binding:"required"`
			EndTime   string `json:"end_time" binding:"required"`
		} `json:"times" binding:"required"`
	}

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

	dateParsed, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	languagesJSON, err := json.Marshal(input.Languages)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid languages"})
		return
	}

	var createdShows []models.Show

	for _, t := range input.Times {
		startTimeParsed, err := time.Parse("15:04", t.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format. Use HH:MM"})
			return
		}
		endTimeParsed, err := time.Parse("15:04", t.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format. Use HH:MM"})
			return
		}

		show := models.Show{
			MovieID:   input.MovieID,
			TheatreID: input.TheatreID,
			Date:      dateParsed,
			StartTime: startTimeParsed,
			EndTime:   endTimeParsed,
			Languages: languagesJSON,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&show).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		createdShows = append(createdShows, show)
	}

	c.JSON(http.StatusCreated, createdShows)
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

	var input struct {
		MovieID   int      `json:"movie_id" binding:"required"`
		TheatreID int      `json:"theatre_id" binding:"required"`
		Date      string   `json:"date" binding:"required"`
		StartTime string   `json:"start_time" binding:"required"`
		EndTime   string   `json:"end_time" binding:"required"`
		Languages []string `json:"languages" binding:"required"`
	}

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

	dateParsed, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	startTimeParsed, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format. Use HH:MM"})
		return
	}
	endTimeParsed, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format. Use HH:MM"})
		return
	}
	languagesJSON, err := json.Marshal(input.Languages)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid languages"})
		return
	}

	show.MovieID = input.MovieID
	show.TheatreID = input.TheatreID
	show.Date = dateParsed
	show.StartTime = startTimeParsed
	show.EndTime = endTimeParsed
	show.Languages = languagesJSON
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
