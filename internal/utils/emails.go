package utils

import (
	"context"
	"log"

	novugo "github.com/novuhq/novu-go"
	"github.com/novuhq/novu-go/models/components"
	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/models"
)

var _ = config.Load()
var secretKey = config.GetEnv("NOVU_SECRET_KEY", "")
var s = novugo.New(novugo.WithSecurity(secretKey))

func SendWelcomeEmail(toEmail, name string) {
	// Run in a goroutine to avoid blocking response
	go func() {
		ctx := context.Background()
		_, err := s.Trigger(ctx, components.TriggerEventRequestDto{
			WorkflowID: "golang-welcome-email",
			Payload: map[string]any{
				"name": name,
			},
			To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
				Email:        &toEmail,
				SubscriberID: toEmail,
			}),
		}, nil)

		if err != nil {
			log.Println("Error sending welcome email", err)
		}
	}()
}

func SendEventRegistrarionSuccessEmail(toEmail, name string, event *models.Event) {
	go func() {
		ctx := context.Background()
		_, err := s.Trigger(ctx, components.TriggerEventRequestDto{
			WorkflowID: "event-registration-success-email",
			Payload: map[string]any{
				"name":          name,
				"eventTitle":    &event.Title,
				"eventTime":     &event.DateTime,
				"eventLocation": &event.Location,
			},
			To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
				Email:        &toEmail,
				SubscriberID: toEmail,
			}),
		}, nil)

		if err != nil {
			log.Println("Error sending event registration success email", err)
		}
	}()
}

func SendEventCancellationSuccessEmail(toEmail, name string, event *models.Event) {
	go func() {
		ctx := context.Background()
		_, err := s.Trigger(ctx, components.TriggerEventRequestDto{
			WorkflowID: "golang-event-registration-cancellation-email",
			Payload: map[string]any{
				"name":       name,
				"eventTitle": &event.Title,
			},
			To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
				Email:        &toEmail,
				SubscriberID: toEmail,
			}),
		}, nil)

		if err != nil {
			log.Println("Error sending event cancellation success email", err)
		}
	}()
}
