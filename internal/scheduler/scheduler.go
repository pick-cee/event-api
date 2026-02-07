package scheduler

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/pick-cee/events-api/internal/jobs"
	"github.com/pick-cee/events-api/internal/services"
)

func StartScheduler(emailService *services.EmailService) (gocron.Scheduler, error) {
	// create a new scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	reminderJob := jobs.NewEventReminderJob(emailService)

	// run 24-hour reminder every hour
	_, err = scheduler.NewJob(
		gocron.DurationJob(1*time.Hour),
		gocron.NewTask(func() {
			if err := reminderJob.Check24HourReminders(); err != nil {
				log.Printf("❌ 24h reminder job failed: %v\n", err)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	// run 1-hour reminder every 10 minutes
	_, err = scheduler.NewJob(
		gocron.DurationJob(10*time.Minute),
		gocron.NewTask(func() {
			if err := reminderJob.Check1HourReminders(); err != nil {
				log.Printf("❌ 1h reminder job failed: %v\n", err)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Scheduler started")
	log.Println("  - 24h reminders: Every 1 hour")
	log.Println("  - 1h reminders: Every 10 minutes")

	// Start scheduler
	scheduler.Start()

	return scheduler, nil
}
