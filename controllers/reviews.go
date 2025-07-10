package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

type CreateReviewInput struct {
	MovieID  uint   `json:"movie_id" binding:"required"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Comments string `json:"comments"`
}

type UpdateReviewInput struct {
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Comments string `json:"comments"`
}

func CreateReview(c *gin.Context) {
	var input CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	db := c.MustGet("db").(*gorm.DB)

	var existing models.Review
	if err := db.Where("user_id = ? AND movie_id = ?", userID, input.MovieID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this movie"})
		return
	}

	review := models.Review{
		UserID:    userID,
		MovieID:   input.MovieID,
		Rating:    input.Rating,
		Comments:  input.Comments,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func GetReviews(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var reviews []models.Review
	if err := db.Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func GetReviewByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var review models.Review
	if err := db.First(&review, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, review)
}

func UpdateReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var review models.Review
	if err := db.First(&review, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var input UpdateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.Rating = input.Rating
	review.Comments = input.Comments
	review.UpdatedAt = time.Now()

	if err := db.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

func DeleteReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := db.Delete(&models.Review{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted"})
}
