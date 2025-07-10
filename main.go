package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/controllers"
	"backend/middlewares"
	"backend/models"
)

func main() {
	models.ConnectDatabase()

	router := gin.Default()

	router.POST("/api/payment/success", controllers.PaymentSuccessHandler)
	router.POST("/api/payment/failure", controllers.PaymentFailureHandler)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(func(c *gin.Context) {
		c.Set("db", models.DB)
		c.Next()
	})

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.GET("/me", middlewares.AuthMiddleware(), controllers.Me)
	router.GET("/api/booking/:txnid", controllers.GetBookingDetails)

	protected := router.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(200, gin.H{"user": user})
		})

		// Payment initiation (requires authentication)
		protected.POST("/payment/initiate", controllers.InitiatePayment)
	}

	movieRoutes := router.Group("/movies")
	{
		movieRoutes.POST("", controllers.CreateMovie)
		movieRoutes.GET("", controllers.GetMovies)
		movieRoutes.GET("/:id", controllers.GetMovieByID)
		movieRoutes.PUT("/:id", controllers.UpdateMovie)
		movieRoutes.DELETE("/:id", controllers.DeleteMovie)
	}

	theatreRoutes := router.Group("/theatres")
	{
		theatreRoutes.POST("", controllers.CreateTheatre)
		theatreRoutes.GET("", controllers.GetTheatres)
		theatreRoutes.GET("/:id", controllers.GetTheatreByID)
		theatreRoutes.PUT("/:id", controllers.UpdateTheatre)
		theatreRoutes.DELETE("/:id", controllers.DeleteTheatre)
	}

	showRoutes := router.Group("/shows")
	{
		showRoutes.POST("", controllers.CreateShow)
		showRoutes.GET("", controllers.GetShows)
		showRoutes.GET("/:id", controllers.GetShowByID)
		showRoutes.PUT("/:id", controllers.UpdateShow)
		showRoutes.DELETE("/:id", controllers.DeleteShow)
	}

	reviewRoutes := router.Group("/reviews")
	{
		reviewRoutes.POST("", controllers.CreateReview)
		reviewRoutes.GET("", controllers.GetReviews)
		reviewRoutes.GET("/:id", controllers.GetReviewByID)
		reviewRoutes.PUT("/:id", controllers.UpdateReview)
		reviewRoutes.DELETE("/:id", controllers.DeleteReview)
	}

	stateRoutes := router.Group("/states")
	{
		stateRoutes.POST("", controllers.CreateState)
		stateRoutes.GET("", controllers.GetStates)
		stateRoutes.GET("/:id", controllers.GetStateByID)
		stateRoutes.PUT("/:id", controllers.UpdateState)
		stateRoutes.DELETE("/:id", controllers.DeleteState)
	}

	cityRoutes := router.Group("/cities")
	{
		cityRoutes.POST("", controllers.CreateCity)
		cityRoutes.GET("", controllers.GetCities)
		cityRoutes.GET("/:id", controllers.GetCityByID)
		cityRoutes.PUT("/:id", controllers.UpdateCity)
		cityRoutes.DELETE("/:id", controllers.DeleteCity)
	}

	seatRoutes := router.Group("/seats")
	{
		seatRoutes.GET("/show/:id", controllers.GetBookedSeats)
		seatRoutes.POST("/book", controllers.BookSeat)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
