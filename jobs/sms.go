package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// SMSPayload holds the data for an SMS notification job.
type SMSPayload struct {
	To       string `json:"to"`
	From     string `json:"from,omitempty"`
	Message  string `json:"message"`
	Schedule string `json:"schedule,omitempty"` // RFC3339 for scheduled delivery
}

// SMSJobHandler creates a handler for processing SMS jobs.
func SMSJobHandler(provider SMSProvider, logger *slog.Logger) JobHandler {
	return func(ctx context.Context, job *Job) error {
		var payload SMSPayload
		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal SMS payload: %w", err)
		}

		logger.Info("sending SMS",
			"to", payload.To,
			"message_length", len(payload.Message),
			"job_id", job.ID)

		// Handle scheduled delivery
		if payload.Schedule != "" {
			scheduledTime, err := time.Parse(time.RFC3339, payload.Schedule)
			if err != nil {
				return fmt.Errorf("invalid schedule time: %w", err)
			}
			// In production, store in scheduled queue and process at time
			logger.Info("SMS scheduled", "at", scheduledTime)
		}

		sms := SMSMessage{
			To:      payload.To,
			From:    payload.From,
			Message: payload.Message,
		}

		return provider.Send(ctx, sms, logger)
	}
}

// SMSMessage represents an SMS to be sent.
type SMSMessage struct {
	To      string
	From    string
	Message string
}

// SMSProvider interface for SMS delivery services.
type SMSProvider interface {
	Send(ctx context.Context, msg SMSMessage, logger *slog.Logger) error
}

// TwilioProvider implements SMSProvider using Twilio.
type TwilioProvider struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// NewTwilioProvider creates a new Twilio SMS provider.
func NewTwilioProvider(accountSID, authToken, fromNumber string) *TwilioProvider {
	return &TwilioProvider{
		AccountSID: accountSID,
		AuthToken:  authToken,
		FromNumber: fromNumber,
	}
}

// Send sends an SMS via Twilio API.
func (p *TwilioProvider) Send(ctx context.Context, msg SMSMessage, logger *slog.Logger) error {
	// Stub: In production, use Twilio SDK
	// client := twilio.NewClient(p.AccountSID, p.AuthToken)
	// _, err := client.Messages.Send(&twilio.Message{
	//     From: p.FromNumber,
	//     To:   msg.To,
	//     Body: msg.Message,
	// })
	logger.Debug("SMS sent via Twilio",
		"to", msg.To,
		"from", p.FromNumber)
	return nil
}

// AWS SNS Provider
type SNSProvider struct {
	Region string
	Topic  string
}

// NewSNSProvider creates a new AWS SNS provider.
func NewSNSProvider(region, topic string) *SNSProvider {
	return &SNSProvider{
		Region: region,
		Topic:  topic,
	}
}

// Send sends an SMS via AWS SNS.
func (p *SNSProvider) Send(ctx context.Context, msg SMSMessage, logger *slog.Logger) error {
	// Stub: In production, use AWS SDK
	// svc := sns.New(session.Must(session.NewSession(&aws.Config{Region: &p.Region}))
	// _, err := svc.PublishWithContext(ctx, &sns.PublishInput{
	//     PhoneNumber: &msg.To,
	//     Message:     &msg.Message,
	// })
	logger.Debug("SMS sent via AWS SNS",
		"to", msg.To,
		"region", p.Region)
	return nil
}

// NewSMSJob creates a new SMS notification job.
func NewSMSJob(payload SMSPayload) (*Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SMS payload: %w", err)
	}

	return &Job{
		Type:    "sms",
		Payload: payloadBytes,
	}, nil
}
