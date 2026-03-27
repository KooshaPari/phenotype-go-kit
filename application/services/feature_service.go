package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	outbound "github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
	"github.com/KooshaPari/phenotype-go-kit/domain/entities"
	"github.com/KooshaPari/phenotype-go-kit/domain/ports"
)

// FeatureService handles feature-related use cases.
// Following Clean Architecture + DDD + CQRS principles.
//
// Application Layer contains:
// - Use cases (orchestration logic)
// - Input/Output DTOs
// - Transaction boundaries
// - Cross-cutting concerns (logging, observability)
type FeatureService struct {
	featureRepo   ports.FeatureRepositoryPort
	wpRepo        ports.WorkPackageRepositoryPort
	cache         outbound.CacheInvalidationPort
	eventBus      outbound.EventBusPort
	audit         ports.AuditPort
	observability outbound.ObservabilityPort
}

// NewFeatureService creates a new FeatureService.
// Following Dependency Injection pattern.
func NewFeatureService(
	featureRepo ports.FeatureRepositoryPort,
	wpRepo ports.WorkPackageRepositoryPort,
	cache outbound.CacheInvalidationPort,
	eventBus outbound.EventBusPort,
	audit ports.AuditPort,
	observability outbound.ObservabilityPort,
) *FeatureService {
	return &FeatureService{
		featureRepo:   featureRepo,
		wpRepo:        wpRepo,
		cache:         cache,
		eventBus:      eventBus,
		audit:         audit,
		observability: observability,
	}
}

// CreateFeatureInput represents the input for creating a feature.
type CreateFeatureInput struct {
	Name        string
	Description string
	Mission     string
	KittySpec   string
	CreatedBy   string
}

// CreateFeature creates a new feature.
// Following TDD and BDD patterns - use case returns Result type.
func (s *FeatureService) CreateFeature(ctx context.Context, input CreateFeatureInput) (*entities.Feature, error) {
	// Validation (following DDD validation patterns)
	if input.Name == "" {
		return nil, entities.NewValidationError("Name", "cannot be empty")
	}

	// Domain logic - create aggregate
	feature := &entities.Feature{
		ID:          string(entities.NewAggregateID("feat")),
		Name:        input.Name,
		Description: input.Description,
		Mission:     input.Mission,
		KittySpec:   input.KittySpec,
		Status:      entities.StatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Persistence (following Repository pattern)
	saved, err := s.featureRepo.Save(ctx, *feature)
	if err != nil {
		s.observability.RecordError(ctx, "CreateFeature", err)
		return nil, fmt.Errorf("saving feature: %w", err)
	}

	// Audit (following Event Sourcing)
	if err := s.audit.Record(ctx, &entities.AuditEntry{
		ID:             string(entities.NewAggregateID("audit")),
		FeatureID:      saved.ID,
		TransitionType: "created",
		FromStatus:     "",
		ToStatus:       string(saved.Status),
		EvidenceRefs:   []string{},
		PreviousHash:   "",
		Hash:           "",
		Timestamp:      time.Now(),
		Actor:          input.CreatedBy,
	}); err != nil {
		s.observability.RecordError(ctx, "CreateFeature.audit", err)
	}

	// Publish event (following CQRS/Event Sourcing)
	if err := s.eventBus.Publish(ctx, "feature.created", saved); err != nil {
		s.observability.RecordError(ctx, "CreateFeature.event", err)
	}

	// Invalidate cache
	if _, err := s.cache.DeletePattern(ctx, "features:list*"); err != nil {
		s.observability.RecordError(ctx, "CreateFeature.cache", err)
	}

	return &saved, nil
}

// TransitionFeature transitions a feature to a new status.
// Following State Machine pattern.
func (s *FeatureService) TransitionFeature(ctx context.Context, featureID string, targetStatus entities.FeatureStatus, actor string) (*entities.Feature, error) {
	// Load aggregate
	feature, err := s.featureRepo.FindByID(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("finding feature: %w", err)
	}

	// Validate transition (state machine)
	if !feature.CanTransitionTo(targetStatus) {
		return nil, entities.ErrInvalidTransition
	}

	// Business validation (governance contract check)
	if err := s.validateTransition(ctx, feature, targetStatus); err != nil {
		return nil, fmt.Errorf("validating transition: %w", err)
	}

	// State change
	fromStatus := feature.Status
	feature.Status = targetStatus
	feature.UpdatedAt = time.Now()

	if targetStatus == entities.StatusShipped {
		now := time.Now()
		feature.ShippedAt = &now
	}

	// Persist
	saved, err := s.featureRepo.Save(ctx, *feature)
	if err != nil {
		return nil, fmt.Errorf("saving feature: %w", err)
	}

	// Audit with hash chain
	previousHash, _ := s.audit.GetLatestHash(ctx, featureID)
	auditEntry := &entities.AuditEntry{
		ID:             string(entities.NewAggregateID("audit")),
		FeatureID:      featureID,
		TransitionType: "transitioned",
		FromStatus:     string(fromStatus),
		ToStatus:       string(targetStatus),
		Timestamp:      time.Now(),
		Actor:          actor,
		PreviousHash:   previousHash,
	}
	auditEntry.Hash = entities.ComputeHash(auditEntry)

	if err := s.audit.Record(ctx, auditEntry); err != nil {
		s.observability.RecordError(ctx, "TransitionFeature.audit", err)
	}

	// Publish event
	if err := s.eventBus.Publish(ctx, "feature.transitioned", saved); err != nil {
		s.observability.RecordError(ctx, "TransitionFeature.event", err)
	}

	return &saved, nil
}

// validateTransition performs business validation for transitions.
// Following Policy pattern.
func (s *FeatureService) validateTransition(ctx context.Context, feature *entities.Feature, target entities.FeatureStatus) error {
	// Check governance contract compliance for validation gate
	if target == entities.StatusValidated {
		contract, err := s.featureRepo.GetGovernanceContract(ctx, feature.ID)
		if err != nil {
			return err
		}

		for _, fr := range contract.FRs {
			if !fr.Satisfied {
				return &entities.DomainError{
					Code:    "FR_NOT_SATISFIED",
					Message: fmt.Sprintf("FR %s is not satisfied", fr.Code),
				}
			}
		}
	}
	return nil
}

// ListFeatures returns all features with optional filtering.
// Following CQRS - read model query.
func (s *FeatureService) ListFeatures(ctx context.Context, filter map[string]any) ([]entities.Feature, error) {
	// Try cache first (following Cache-Aside pattern)
	cacheKey := fmt.Sprintf("features:list:%v", filter)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		return s.unmarshalFeatures(cached)
	}

	// Load from repository
	features, err := s.featureRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("listing features: %w", err)
	}

	// Cache result (5 minute TTL)
	_ = s.cache.Set(ctx, cacheKey, s.marshalFeatures(features), 5*time.Minute)

	return features, nil
}

// GetFeature returns a single feature by ID.
func (s *FeatureService) GetFeature(ctx context.Context, id string) (*entities.Feature, error) {
	return s.featureRepo.FindByID(ctx, id)
}

// marshalFeatures serializes features for caching.
func (s *FeatureService) marshalFeatures(features []entities.Feature) string {
	data, err := json.Marshal(features)
	if err != nil {
		return "[]"
	}
	return string(data)
}

// unmarshalFeatures deserializes features from cache.
func (s *FeatureService) unmarshalFeatures(data string) ([]entities.Feature, error) {
	var features []entities.Feature
	if err := json.Unmarshal([]byte(data), &features); err != nil {
		return nil, err
	}
	return features, nil
}
