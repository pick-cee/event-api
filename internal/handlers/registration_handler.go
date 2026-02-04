package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/middleware"
	"github.com/pick-cee/events-api/internal/models"
	"github.com/pick-cee/events-api/internal/utils"
)

type RegistrationHandler struct{}

func NewRegistrationHandler() *RegistrationHandler {
	return &RegistrationHandler{}
}

func (h *RegistrationHandler) RegisterForEvent(c *gin.Context) {
	eventId := c.Param("id")
	userId := middleware.GetUserId(c)

	// check if event exists
	var event models.Event
	if err := database.DB.First(&event, eventId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Event not found")
		return
	}

	// check if already registered
	var existingReg models.Registration
	if err := database.DB.Where("user_id = ? AND event_id = ?", userId, eventId).First(&existingReg).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Already registered for this event")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// create registration
	registration := models.Registration{
		UserID:  userId,
		EventID: event.ID,
	}

	if err := database.DB.Create(&registration).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register, try again")
		return
	}

	database.DB.Preload("Event").Preload("User").Preload("Event.Creator").First(&registration, registration.ID)

	utils.SendEventRegistrarionSuccessEmail(user.Email, user.Name, &event)

	utils.SuccessResponse(c, http.StatusCreated, registration)
}

func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	eventId := c.Param("id")
	userId := middleware.GetUserId(c)

	var registration models.Registration
	if err := database.DB.Where("user_id = ? AND event_id = ?", userId, eventId).First(&registration).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Registration not found")
		return
	}

	var event models.Event
	if err := database.DB.First(&event, eventId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Event not found")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// delete registration
	if err := database.DB.Delete(&registration).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel registration")
		return
	}

	utils.SendEventCancellationSuccessEmail(user.Email, user.Name, &event)
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Registration canceled successfully"})
}

func (h *RegistrationHandler) GetEventAttendees(c *gin.Context) {
	eventId := c.Param("id")

	// check if event exists
	var event models.Event
	if err := database.DB.First(&event, eventId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Event does not exists")
		return
	}

	// get all registrations with user info
	var registrations []models.Registration
	if err := database.DB.Where("event_id = ?", event.ID).Preload("User").Preload("Event.Creator").Find(&registrations).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch attendees")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, registrations)
}

func (h *RegistrationHandler) GetMyRegistrations(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var registrations []models.Registration
	if err := database.DB.Where("user_id", userId).Preload("Event.Creator").Find(&registrations).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch events")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, registrations)
}
