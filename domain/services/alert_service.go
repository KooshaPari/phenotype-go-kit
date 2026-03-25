// Package services contains the domain services.
// Domain services orchestrate operations across entities and value objects.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/domain/entities"
	"github.com/KooshaPari/phenotype-go-kit/domain/ports"
	"github.com/KooshaPari/phenotype-go-kit/domain/valueobjects"
)

// AlertService handles alert-related business logic.
type AlertService struct {
	repo     ports.Repository
	cache    ports.Cache
	eventBus ports.EventBus
}

// NewAlertService creates a new alert service.
func NewAlertService(repo ports.Repository, cache ports.Cache, eventBus ports.EventBus) *AlertService {
	return &AlertService{
		repo:     repo,
		cache:    cache,
		eventBus: eventBus,
	}
}

// EvaluateRule evaluates an alert rule against current metrics.
func (s *AlertService) EvaluateRule(ctx context.Context, rule *valueobjects.AlertRule, metrics map[string]float64) error {
	metricValue, ok := metrics[rule.Condition.Metric]
	if !ok {
		return fmt.Errorf("metric %s not found", rule.Condition.Metric)
	}

	if s.evaluateCondition(metricValue, rule.Condition) {
		alert := &entities.Alert{
			ID:       generateID(),
			RuleID:   rule.ID,
			Severity: entities.Severity(rule.Severity),
			Message:  fmt.Sprintf("Alert %s is firing", rule.Name),
			Status:   entities.AlertStatusFiring,
		}

		if err := s.repo.Create(ctx, alert); err != nil {
			return fmt.Errorf("failed to create alert: %w", err)
		}

		return s.eventBus.Publish(ctx, AlertFiredEvent{Alert: alert})
	}

	return nil
}

func (s *AlertService) evaluateCondition(value float64, cond valueobjects.Condition) bool {
	switch cond.Operator {
	case ">":
		return value > cond.Threshold
	case ">=":
		return value >= cond.Threshold
	case "<":
		return value < cond.Threshold
	case "<=":
		return value <= cond.Threshold
	case "==":
		return value == cond.Threshold
	case "!=":
		return value != cond.Threshold
	default:
		return false
	}
}

// AlertFiredEvent represents a fired alert event.
type AlertFiredEvent struct {
	Alert *entities.Alert
}

func generateID() string {
	return fmt.Sprintf("alert-%d", time.Now().UnixNano())
}
