// Package entities contains domain entities.
// Following DDD (Domain-Driven Design) principles.
package entities

import (
	"time"
)

// User represents a user entity.
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time // Soft delete
}

// IsActive returns true if the user is not deleted.
func (u *User) IsActive() bool {
	return u.DeletedAt == nil
}

// CanAccess returns true if the user can access the resource.
func (u *User) CanAccess(resource string) bool {
	return u.IsActive()
}

// Config represents a configuration entity.
type Config struct {
	ID        string
	Key       string
	Value     string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Workspace represents a workspace entity.
type Workspace struct {
	ID        string
	Name      string
	OwnerID   string
	Plan      Plan
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Plan represents a subscription plan.
type Plan string

const (
	PlanFree       Plan = "free"
	PlanPro        Plan = "pro"
	PlanEnterprise Plan = "enterprise"
)

// CanUseFeature returns true if the plan supports the feature.
func (p Plan) CanUseFeature(feature string) bool {
	switch p {
	case PlanFree:
		return freeFeatures[feature]
	case PlanPro:
		return proFeatures[feature]
	case PlanEnterprise:
		return true // All features
	default:
		return false
	}
}

var freeFeatures = map[string]bool{
	"basic_metrics": true,
	"3_workspaces":  true,
}

var proFeatures = map[string]bool{
	"basic_metrics":    true,
	"advanced_metrics": true,
	"unlimited_ws":     true,
}
