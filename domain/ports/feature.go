package ports

import (
	"context"

	"github.com/KooshaPari/phenotype-go-kit/domain/entities"
)

// FeatureRepositoryPort defines the interface for feature persistence.
// Following Repository pattern from DDD and Hexagonal Architecture.
//
// This port isolates the domain from persistence concerns,
// allowing the domain to remain pure and testable.
type FeatureRepositoryPort interface {
	// Save creates or updates a feature.
	Save(ctx context.Context, feature entities.Feature) (entities.Feature, error)

	// FindByID retrieves a feature by ID.
	// Returns (nil, ErrFeatureNotFound) if not found.
	FindByID(ctx context.Context, id string) (*entities.Feature, error)

	// List returns all features matching the filter.
	List(ctx context.Context, filter map[string]any) ([]entities.Feature, error)

	// Delete removes a feature by ID.
	Delete(ctx context.Context, id string) error

	// GetGovernanceContract retrieves the governance contract for a feature.
	GetGovernanceContract(ctx context.Context, featureID string) (*entities.GovernanceContract, error)

	// SaveGovernanceContract saves a governance contract.
	SaveGovernanceContract(ctx context.Context, contract entities.GovernanceContract) error
}

// WorkPackageRepositoryPort defines the interface for work package persistence.
// Following DDD Repository pattern.
type WorkPackageRepositoryPort interface {
	// Save creates or updates a work package.
	Save(ctx context.Context, wp entities.WorkPackage) (entities.WorkPackage, error)

	// FindByID retrieves a work package by ID.
	FindByID(ctx context.Context, id string) (*entities.WorkPackage, error)

	// FindByFeatureID retrieves all work packages for a feature.
	FindByFeatureID(ctx context.Context, featureID string) ([]entities.WorkPackage, error)

	// Delete removes a work package by ID.
	Delete(ctx context.Context, id string) error
}

// AuditPort defines the interface for audit logging.
// Following Event Sourcing pattern.
type AuditPort interface {
	// Record records an audit entry.
	Record(ctx context.Context, entry *entities.AuditEntry) error

	// GetLatestHash returns the latest hash for a feature's audit chain.
	GetLatestHash(ctx context.Context, featureID string) (string, error)

	// GetChain returns the full audit chain for a feature.
	GetChain(ctx context.Context, featureID string) ([]entities.AuditEntry, error)

	// VerifyChain verifies the integrity of an audit chain.
	VerifyChain(ctx context.Context, featureID string) (bool, error)
}
