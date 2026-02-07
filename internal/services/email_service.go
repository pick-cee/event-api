package services

import (
	"context"

	novugo "github.com/novuhq/novu-go"
	"github.com/novuhq/novu-go/models/components"
	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/models"
)

type EmailService struct {
	novuClient *novugo.Novu
}

func NewEmailService(cfg *config.Config) *EmailService {
	secretKey := config.GetEnv("NOVU_SECRET_KEY", "")
	client := novugo.New(novugo.WithSecurity(secretKey))

	return &EmailService{
		novuClient: client,
	}
}

func (s *EmailService) SendWelcomeEmail(email, name string) error {
	ctx := context.Background()
	_, err := s.novuClient.Trigger(ctx, components.TriggerEventRequestDto{
		WorkflowID: "golang-welcome-email",
		Payload: map[string]any{
			"name": name,
		},
		To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
			Email:        &email,
			SubscriberID: email,
		}),
	}, nil)

	return err
}

func (s *EmailService) SendEventRegistrarionSuccessEmail(email, name string, event *models.Event) error {
	ctx := context.Background()
	_, err := s.novuClient.Trigger(ctx, components.TriggerEventRequestDto{
		WorkflowID: "event-registration-success-email",
		Payload: map[string]any{
			"name":          name,
			"eventTitle":    event.Title,
			"eventTime":     event.DateTime.Format("Monday, January 2, 2006 at 3:04 PM"),
			"eventLocation": event.Location,
		},
		To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
			Email:        &email,
			SubscriberID: email,
		}),
	}, nil)

	return err
}

func (s *EmailService) SendEventCancellationSuccessEmail(email, name string, event *models.Event) error {
	ctx := context.Background()
	_, err := s.novuClient.Trigger(ctx, components.TriggerEventRequestDto{
		WorkflowID: "golang-event-registration-cancellation-email",
		Payload: map[string]any{
			"name":       name,
			"eventTitle": event.Title,
		},
		To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
			Email:        &email,
			SubscriberID: email,
		}),
	}, nil)

	return err
}

func (s *EmailService) Send24HourEventReminderEmail(email, name string, event *models.Event) error {
	ctx := context.Background()
	_, err := s.novuClient.Trigger(ctx, components.TriggerEventRequestDto{
		WorkflowID: "golang-event-24h-reminder",
		Payload: map[string]any{
			"name":             name,
			"eventTitle":       event.Title,
			"eventLocation":    event.Location,
			"eventTime":        event.DateTime.Format("Monday, January 2, 2006 at 3:04 PM"),
			"eventDescription": event.Description,
		},
		To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
			Email:        &email,
			SubscriberID: email,
		}),
	}, nil)

	return err
}

func (s *EmailService) Send1HourEventReminderEmail(email, name string, event *models.Event) error {
	ctx := context.Background()
	_, err := s.novuClient.Trigger(ctx, components.TriggerEventRequestDto{
		WorkflowID: "golang-event-1h-reminder",
		Payload: map[string]any{
			"name":             name,
			"eventTitle":       event.Title,
			"eventLocation":    event.Location,
			"eventTime":        event.DateTime.Format("Monday, January 2, 2006 at 3:04 PM"),
			"eventDescription": event.Description,
		},
		To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
			Email:        &email,
			SubscriberID: email,
		}),
	}, nil)

	return err
}
