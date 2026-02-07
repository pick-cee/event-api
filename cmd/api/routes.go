package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/handlers"
	"github.com/pick-cee/events-api/internal/middleware"
	"github.com/pick-cee/events-api/internal/services"
)

func setupRoutes(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// initialize services needed by handlers
	emailService := services.NewEmailService(cfg)

	// initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, emailService)
	eventHandler := handlers.NewEventHandler()
	registrationHandler := handlers.NewRegistrationHandler(emailService)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API V1 routes
	v1 := r.Group("/api/v1")
	{
		// public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
		}

		// public event routes
		events := v1.Group("/events")
		{
			events.GET("", eventHandler.ListEvents)
			events.GET("/:id", eventHandler.GetEventById)
			events.GET("/:id/attendees", registrationHandler.GetEventAttendees)
		}

		// protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMidleware(cfg))
		{
			// Event management (authenticated users)
			protected.POST("/events", eventHandler.CreateEvent)       // POST /api/v1/events
			protected.PUT("/events/:id", eventHandler.UpdateEvent)    // PUT /api/v1/events/:id
			protected.DELETE("/events/:id", eventHandler.DeleteEvent) // DELETE /api/v1/events/:id

			// Event registration (authenticated users)
			protected.POST("/events/:id/register", registrationHandler.RegisterForEvent)   // POST /api/v1/events/:id/register
			protected.DELETE("/events/:id/cancel", registrationHandler.CancelRegistration) // DELETE /api/v1/events/:id/register
			protected.GET("/my-registrations", registrationHandler.GetMyRegistrations)     // GET /api/v1/my-registrations
		}
	}
	return r
}
