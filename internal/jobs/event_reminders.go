package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pick-cee/events-api/internal/cache"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/models"
	"github.com/pick-cee/events-api/internal/services"
)

type EventReminderJob struct {
	emailService *services.EmailService
}

func NewEventReminderJob(emailService *services.EmailService) *EventReminderJob {
	return &EventReminderJob{
		emailService: emailService,
	}
}

// check 24hour remninders
func (j *EventReminderJob) Check24HourReminders() error {
	log.Println("⏰ Running 24-hour reminder job...")

	ctx := context.Background()
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	// Find events happening in approximately 24 hours (±30 minutes)
	var events []models.Event
	err := database.DB.
		Preload("Creator").
		Preload("Registrations.User").
		Where("date_time BETWEEN ? AND ?", tomorrow.Add(-30*time.Minute), tomorrow.Add(30*time.Minute)).
		Find(&events).Error

	if err != nil {
		return fmt.Errorf("failed to fetch events: %w", err)
	}

	log.Printf("📧 Found %d events for 24-hour reminders\n", len(events))

	// Send reminders
	for _, event := range events {
		// Check if we already sent 24h reminder
		cacheKey := fmt.Sprintf("reminder:24h:%d", event.ID)
		exists, _ := cache.Exists(ctx, cacheKey)

		if exists == true {
			log.Printf("⏭️  Already sent 24h reminder for event %d, skipping\n", event.ID)
			continue
		}

		// Send to all registered attendees
		for _, registration := range event.Registrations {
			err := j.emailService.Send24HourEventReminderEmail(registration.User.Email, registration.User.Name, &event)
			if err != nil {
				log.Printf("❌ Failed to send 24h reminder to %s: %v\n", registration.User.Email, err)
			} else {
				log.Printf("✅ Sent 24h reminder to %s for event: %s\n", registration.User.Email, event.Title)
			}
		}

		// Mark as sent (cache for 48 hours)
		cache.Set(ctx, cacheKey, "sent", 48*time.Hour)
	}
	return nil
}

func (j *EventReminderJob) Check1HourReminders() error {
	log.Println("⏰ Running 1-hour reminder job...")

	ctx := context.Background()
	now := time.Now()
	oneHourLater := now.Add(1 * time.Hour)
	log.Println(oneHourLater)

	// Find events happening in approximately 1 hour (±10 minutes)
	var events []models.Event
	err := database.DB.
		Preload("Creator").
		Preload("Registrations.User").
		Where("date_time BETWEEN ? AND ?", oneHourLater.Add(-10*time.Minute), oneHourLater.Add(10*time.Minute)).
		Find(&events).Error

	if err != nil {
		return fmt.Errorf("failed to fetch events: %w", err)
	}

	log.Printf("📧 Found %d events for 1-hour reminders\n", len(events))

	for _, event := range events {
		// Check if we already sent 1h reminder
		cacheKey := fmt.Sprintf("reminder:1h:%d", event.ID)
		exists, _ := cache.Exists(ctx, cacheKey)

		if exists == true {
			log.Printf("⏭️  Already sent 1h reminder for event %d, skipping\n", event.ID)
			continue
		}

		// Send to all registered attendees
		for _, registration := range event.Registrations {
			err := j.emailService.Send1HourEventReminderEmail(registration.User.Email, registration.User.Name, &event)
			if err != nil {
				log.Printf("❌ Failed to send 1h reminder to %s: %v\n", registration.User.Email, err)
			} else {
				log.Printf("✅ Sent 1h reminder to %s for event: %s\n", registration.User.Email, event.Title)
			}
		}

		// Mark as sent (cache for 24 hours)
		cache.Set(ctx, cacheKey, "sent", 24*time.Hour)
	}

	return nil
}
