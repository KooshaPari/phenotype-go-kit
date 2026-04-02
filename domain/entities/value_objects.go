package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// ValidationError represents a domain validation error.
// Following DDD Value Object pattern.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// DomainError represents a domain-level error.
// Following Error wrapping best practices.
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Message + " - " + e.Err.Error()
	}
	return e.Code + ": " + e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Common domain errors following Error handling best practices.
var (
	ErrFeatureNotFound     = &DomainError{Code: "FEATURE_NOT_FOUND", Message: "Feature not found"}
	ErrWPNotFound          = &DomainError{Code: "WP_NOT_FOUND", Message: "Work package not found"}
	ErrInvalidTransition   = &DomainError{Code: "INVALID_TRANSITION", Message: "Invalid state transition"}
	ErrMissingEvidence     = &DomainError{Code: "MISSING_EVIDENCE", Message: "Missing required evidence"}
	ErrGovernanceViolation = &DomainError{Code: "GOVERNANCE_VIOLATION", Message: "Governance contract violated"}
)

// AggregateID represents a unique identifier for aggregates.
// Following DDD Identifier pattern.
type AggregateID string

// NewAggregateID creates a new aggregate ID with prefix.
// Following Factory pattern.
func NewAggregateID(prefix string) AggregateID {
	return AggregateID(fmt.Sprintf("%s-%s", prefix, generateUUID()))
}

// ComputeHash computes a SHA-256 hash of an audit entry.
// Following Event Sourcing pattern.
func ComputeHash(entry *AuditEntry) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%v|%s",
		entry.ID,
		entry.FeatureID,
		entry.TransitionType,
		entry.FromStatus,
		entry.ToStatus,
		entry.EvidenceRefs,
		entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// generateUUID generates a simple UUID-like string.
// In production, use google/uuid or similar.
func generateUUID() string {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}
