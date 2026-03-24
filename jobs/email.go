package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// EmailPayload holds the data for an email notification job.
type EmailPayload struct {
	To      string   `json:"to"`
	From    string   `json:"from"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	HTML    string   `json:"html,omitempty"`
	Cc      []string `json:"cc,omitempty"`
	Bcc     []string `json:"bcc,omitempty"`
}

// EmailJobHandler creates a handler for processing email jobs.
func EmailJobHandler(smtpHost string, smtpPort int, logger *slog.Logger) JobHandler {
	return func(ctx context.Context, job *Job) error {
		var payload EmailPayload
		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal email payload: %w", err)
		}

		logger.Info("sending email",
			"to", payload.To,
			"subject", payload.Subject,
			"job_id", job.ID)

		// In production, integrate with actual SMTP or email service
		// e.g., sendgrid, aws-ses, mailgun
		email := EmailMessage{
			From:    payload.From,
			To:      payload.To,
			Subject: payload.Subject,
			Body:    payload.Body,
			HTML:    payload.HTML,
			Cc:      payload.Cc,
			Bcc:     payload.Bcc,
		}

		return sendEmail(ctx, smtpHost, smtpPort, email, logger)
	}
}

// EmailMessage represents an email to be sent.
type EmailMessage struct {
	From    string
	To      string
	Subject string
	Body    string
	HTML    string
	Cc      []string
	Bcc     []string
}

// sendEmail sends an email via SMTP or email service.
// This is a stub - integrate with actual provider in production.
func sendEmail(ctx context.Context, host string, port int, msg EmailMessage, logger *slog.Logger) error {
	// Stub: In production, implement actual email sending
	// - SMTP via net/smtp
	// - Or third-party API (SendGrid, AWS SES, Mailgun)
	logger.Debug("email queued for delivery",
		"from", msg.From,
		"to", msg.To,
		"subject", msg.Subject)

	// Simulate successful send
	time.Sleep(10 * time.Millisecond)
	return nil
}

// NewEmailJob creates a new email notification job.
func NewEmailJob(payload EmailPayload) (*Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email payload: %w", err)
	}

	return &Job{
		Type:    "email",
		Payload: payloadBytes,
	}, nil
}
