package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/middleware"
	"github.com/pick-cee/events-api/internal/models"
	"github.com/pick-cee/events-api/internal/utils"
	"gorm.io/gorm"
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
	var total int64

	params := utils.GetPaginationParams(c.Request)

	// Count total
	if err := database.DB.Model(&models.Event{}).Count(&total).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to count events")
		return
	}

	// preload creator information
	if err := database.DB.Scopes(utils.Paginate(params)).Preload("Creator").Find(&events).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch events")
		return
	}

	response := utils.NewPaginationResponse(events, total, params)

	utils.SuccessResponse(c, http.StatusOK, response)
}

func (h *EventHandler) GetEventById(c *gin.Context) {
	id := c.Param("id")

	var event models.Event
	if err := database.DB.Preload("Creator").Preload("Registrations.User").First(&event, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Event not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, event)
}

// create event 
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var request CreateEventRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "An error occured while trying to create events")
		return
	}

	// Load creator info
  database.DB.Preload("Creator").First(&event, event.ID)

	utils.SuccessResponse(c, http.StatusCreated, event)
}

func (h *EventHandler) UpdateEvent (c *gin.Context) {
	var request UpdateEventRequest
	id := c.Param("id")
  userID := middleware.GetUserId(c)

  var event models.Event
  if err := database.DB.First(&event, id).Error; err != nil {
    utils.ErrorResponse(c, http.StatusNotFound, "Event not found")
    return
  }

  // Check if user is the creator
  if event.CreatorID != userID {
    utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own events")
    return
  }


	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update event")
		return
	}

	database.DB.Preload("Creator").First(&event, event.ID)

	utils.SuccessResponse(c, http.StatusOK, event)
}

// Deletes event only by creator
func (h *EventHandler) DeleteEvent (c *gin.Context) {
	id := c.Param("id")
	userId := middleware.GetUserId(c)

	var event models.Event
	if err := database.DB.First(&event, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Event not found")
		return
	}

	// Check if user is the creator
  if event.CreatorID != userId {
    utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own events")
    return
  }

	if err := database.DB.Delete(&event, id).Error; err != nil{
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func Paginate(r *http.Request) func (db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		limit, _ := strconv.Atoi(q.Get("limit"))

		if page <= 0 {
			page = 1
		}

		if limit <= 0 {
			limit = 10
		}

		if limit > 100 {
			limit = 100
		}

		offset := (page - 1) * limit

		return db.Offset(offset).Limit(limit)
	}
}