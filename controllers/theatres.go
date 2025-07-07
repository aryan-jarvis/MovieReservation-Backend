package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

func CreateTheatre(c *gin.Context) {
	var theatre models.Theatre

	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	theatre.CreatedAt = time.Now()
	theatre.UpdatedAt = time.Now()

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&theatre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, theatre)
}

func GetTheatres(c *gin.Context) {
	var theatres []models.Theatre
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&theatres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, theatres)
}

func GetTheatreByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre ID"})
		return
	}

	var theatre models.Theatre
	db := c.MustGet("db").(*gorm.DB)

	if err := db.First(&theatre, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Theatre not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, theatre)
}

func DeleteTheatre(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Delete(&models.Theatre{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theatre deleted successfully"})
}

func UpdateTheatre(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre ID"})
		return
	}

	var theatre models.Theatre
	db := c.MustGet("db").(*gorm.DB)

	if err := db.First(&theatre, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Theatre not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var input models.Theatre
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	theatre.TheatreName = input.TheatreName
	theatre.TheatreLocation = input.TheatreLocation
	theatre.CityID = input.CityID
	theatre.TotalSeats = input.TotalSeats
	theatre.TheatreTiming = input.TheatreTiming
	theatre.UpdatedAt = time.Now()

	if err := db.Save(&theatre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, theatre)
}
