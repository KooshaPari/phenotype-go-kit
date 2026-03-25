// Package queries contains read operations (CQRS).
package queries

import (
	"context"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/domain/entities"
)

// GetAlertQuery represents a query to get a single alert.
type GetAlertQuery struct {
	AlertID string
}

// GetAlertsQuery represents a query to list alerts.
type GetAlertsQuery struct {
	Limit   int
	Offset  int
	Status  *entities.AlertStatus
	Severity *entities.Severity
}

// GetAlertsHandler handles GetAlertsQuery.
type GetAlertsHandler struct{}

// NewGetAlertsHandler creates a new handler.
func NewGetAlertsHandler() *GetAlertsHandler {
	return &GetAlertsHandler{}
}

// Handle executes the query.
func (h *GetAlertsHandler) Handle(ctx context.Context, query *GetAlertsQuery) ([]*entities.Alert, error) {
	// Query logic would go here
	return nil, nil
}

// GetAlertStatsQuery represents a query to get alert statistics.
type GetAlertStatsQuery struct {
	Since time.Time
}

// AlertStats represents aggregated alert statistics.
type AlertStats struct {
	Total      int
	Critical   int
	Warning    int
	Resolved   int
	Firing     int
}
