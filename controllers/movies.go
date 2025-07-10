package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

func parseDate(dateMap map[string]interface{}) (time.Time, error) {
	day := int(dateMap["day"].(float64))
	month := time.Month(int(dateMap["month"].(float64)))
	year := int(dateMap["year"].(float64))
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}

func dateToMap(t time.Time) map[string]int {
	return map[string]int{
		"day":   t.Day(),
		"month": int(t.Month()),
		"year":  t.Year(),
	}
}

type MovieResponse struct {
	models.Movie
	StartDate map[string]int `json:"start_date"`
	EndDate   map[string]int `json:"end_date"`
}

func CreateMovie(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	languagesRaw, ok := input["languages"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "languages field is required"})
		return
	}
	languagesJSON, err := json.Marshal(languagesRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid languages format"})
		return
	}

	durationStr, ok := input["duration"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "duration must be a string"})
		return
	}
	durationInt, err := strconv.Atoi(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration value"})
		return
	}

	startDateRaw, ok := input["start_date"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be an object"})
		return
	}
	startDate, err := parseDate(startDateRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDateRaw, ok := input["end_date"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be an object"})
		return
	}
	endDate, err := parseDate(endDateRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}
	movieStatus, ok := input["movie_status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "movie_status is required and must be a string"})
		return
	}

	movie := models.Movie{
		MovieName:        input["movie_name"].(string),
		MovieDescription: input["movie_description"].(string),
		Duration:         durationInt,
		Languages:        languagesJSON,
		Genre:            input["genre"].(string),
		PosterURL:        input["poster_url"].(string),
		MovieStatus:      movieStatus,
		Rating:           0.0,
		StartDate:        startDate,
		EndDate:          endDate,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := MovieResponse{
		Movie:     movie,
		StartDate: dateToMap(movie.StartDate),
		EndDate:   dateToMap(movie.EndDate),
	}
	c.JSON(http.StatusCreated, response)
}

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []MovieResponse
	for _, movie := range movies {
		responses = append(responses, MovieResponse{
			Movie:     movie,
			StartDate: dateToMap(movie.StartDate),
			EndDate:   dateToMap(movie.EndDate),
		})
	}

	c.JSON(http.StatusOK, responses)
}

func GetMovieByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movie models.Movie
	db := c.MustGet("db").(*gorm.DB)

	if err := db.First(&movie, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	response := MovieResponse{
		Movie:     movie,
		StartDate: dateToMap(movie.StartDate),
		EndDate:   dateToMap(movie.EndDate),
	}
	c.JSON(http.StatusOK, response)
}

func UpdateMovie(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movie models.Movie
	db := c.MustGet("db").(*gorm.DB)

	if err := db.First(&movie, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	languagesRaw, ok := input["languages"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "languages field is required"})
		return
	}
	languagesJSON, err := json.Marshal(languagesRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid languages format"})
		return
	}

	durationStr, ok := input["duration"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "duration must be a string"})
		return
	}
	durationInt, err := strconv.Atoi(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration value"})
		return
	}

	startDateRaw, ok := input["start_date"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be an object"})
		return
	}
	startDate, err := parseDate(startDateRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDateRaw, ok := input["end_date"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be an object"})
		return
	}
	endDate, err := parseDate(endDateRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}
	movieStatus, ok := input["movie_status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "movie_status is required and must be a string"})
		return
	}

	movie.MovieName = input["movie_name"].(string)
	movie.MovieDescription = input["movie_description"].(string)
	movie.Duration = durationInt
	movie.Languages = languagesJSON
	movie.Genre = input["genre"].(string)
	movie.PosterURL = input["poster_url"].(string)
	movie.MovieStatus = movieStatus
	movie.StartDate = startDate
	movie.EndDate = endDate
	movie.UpdatedAt = time.Now()

	if err := db.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := MovieResponse{
		Movie:     movie,
		StartDate: dateToMap(movie.StartDate),
		EndDate:   dateToMap(movie.EndDate),
	}
	c.JSON(http.StatusOK, response)
}

func DeleteMovie(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Delete(&models.Movie{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
