package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/middleware"
	"github.com/pick-cee/events-api/internal/models"
)

type RegistrationHandler struct {}

func NewRegistrationHandler() *RegistrationHandler {
	return &RegistrationHandler{}
}

func (h *RegistrationHandler) RegisterForEvent(c *gin.Context) {
	eventId := c.Param("id")
	userId := middleware.GetUserId(c)

	// check if event exists
	var event models.Event
	if err := database.DB.First(&event, eventId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// check if already registered
	var existingReg models.Registration
	if err := database.DB.Where("user_id = ? AND event_id = ?", userId, eventId).First(&existingReg).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already registered for this event"})
		return
	}
	
	// create registration
	registration := models.Registration{
		UserID: userId,
		EventID: event.ID,
	}

	if err := database.DB.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register, try again"})
		return
	}

	database.DB.Preload("Event").Preload("User").Preload("Event.Creator").First(&registration, registration.ID)

	c.JSON(http.StatusCreated, registration)
}

func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	eventId := c.Param("id")
	userId := middleware.GetUserId(c)

	var registration models.Registration
	if err := database.DB.Where("user_id = ? AND event_id = ?", userId, eventId).First(&registration).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registration not found"})
		return
	}

	// delete registration
	if err := database.DB.Delete(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel registration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration canceled successfully"})
}

func (h *RegistrationHandler) GetEventAttendees(c *gin.Context) {
	eventId := c.Param("id")
	
	// check if event exists
	var event models.Event
	if err := database.DB.First(&event, eventId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event does not exists"})
		return
	}

	// get all registrations with user info
	var registrations []models.Registration
	if err := database.DB.Where("event_id = ?", event.ID).Preload("User").Preload("Event.Creator").Find(&registrations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendees"})
		return
	}

	c.JSON(http.StatusOK, registrations)
}

func (h *RegistrationHandler) GetMyRegistrations(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var registrations []models.Registration
	if err := database.DB.Where("user_id", userId).Preload("Event.Creator").Find(&registrations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, registrations)
}