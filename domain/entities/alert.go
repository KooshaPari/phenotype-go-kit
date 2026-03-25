// Package entities contains the core domain entities.
// These are the main business objects with identity.
package entities

import "time"

// Alert represents a monitoring alert.
type Alert struct {
	ID        string
	RuleID    string
	Severity  Severity
	Message   string
	Status    AlertStatus
	CreatedAt time.Time
	ResolvedAt *time.Time
}

// Severity represents alert severity levels.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

// AlertStatus represents the current status of an alert.
type AlertStatus string

const (
	AlertStatusFiring   AlertStatus = "firing"
	AlertStatusResolved AlertStatus = "resolved"
	AlertStatusSilenced AlertStatus = "silenced"
)

// HealthStatus represents a component's health.
type HealthStatus struct {
	Component string
	Status    HealthState
	Message   string
	CheckedAt time.Time
}

// HealthState represents the health state.
type HealthState string

const (
	HealthStateHealthy   HealthState = "healthy"
	HealthStateDegraded  HealthState = "degraded"
	HealthStateUnhealthy HealthState = "unhealthy"
)
