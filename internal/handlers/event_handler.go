package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/middleware"
	"github.com/pick-cee/events-api/internal/models"
)

type EventHandler struct {}

func NewEventHandler () *EventHandler {
	return &EventHandler{}
}

// Request/Response DTOs
type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	DateTime    time.Time `json:"date_time"`
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	var events []models.Event

	// preload creator information
	if err := database.DB.Preload("Creator").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) GetEventById(c *gin.Context) {
	id := c.Param("id")

	var event models.Event
	if err := database.DB.Preload("Creator").Preload("Registrations.User").First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

  c.JSON(http.StatusOK, event)
}

// create event 
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var request CreateEventRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get authenticated user
	userId := middleware.GetUserId(c)

	event := models.Event{
		Title: request.Title,
		Description: request.Description,
		Location: request.Location,
		CreatorID: userId,
		DateTime: request.DateTime,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured whilc trying to create events"})
		return
	}

	// Load creator info
  database.DB.Preload("Creator").First(&event, event.ID)

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) UpdateEvent (c *gin.Context) {
	var request UpdateEventRequest
	id := c.Param("id")
  userID := middleware.GetUserId(c)

  var event models.Event
  if err := database.DB.First(&event, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
    return
  }

  // Check if user is the creator
  if event.CreatorID != userID {
    c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own events"})
    return
  }


	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update field if required
	if request.Title != ""{
		event.Title = request.Title
	}

	if request.Description != "" {
		event.Description = request.Description
	}

	if request.Location != "" {
		event.Location = request.Location
	}

	if !request.DateTime.IsZero() {
		event.DateTime = request.DateTime
	}

	if err := database.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	database.DB.Preload("Creator").First(&event, event.ID)

	c.JSON(http.StatusOK, event)
}

// Deletes event only by creator
func (h *EventHandler) DeleteEvent (c *gin.Context) {
	id := c.Param("id")
	userId := middleware.GetUserId(c)

	var event models.Event
	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Check if user is the creator
  if event.CreatorID != userId {
    c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own events"})
    return
  }

	if err := database.DB.Delete(&event, id).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}