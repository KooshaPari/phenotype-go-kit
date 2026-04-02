package entities

import (
	"time"
)

// FeatureStatus represents the state of a feature.
type FeatureStatus string

const (
	StatusDraft        FeatureStatus = "draft"
	StatusSpecifying   FeatureStatus = "specifying"
	StatusResearched   FeatureStatus = "researched"
	StatusPlanned      FeatureStatus = "planned"
	StatusImplementing FeatureStatus = "implementing"
	StatusValidated    FeatureStatus = "validated"
	StatusShipped      FeatureStatus = "shipped"
	StatusArchived     FeatureStatus = "archived"
)

// Feature represents a unit of work from idea to shipment.
// Following DDD Aggregate pattern.
type Feature struct {
	ID          string
	Name        string
	Description string
	Status      FeatureStatus
	Mission     string
	KittySpec   string
	SpecMD      string
	ResearchMD  string
	PlanMD      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ShippedAt   *time.Time
}

// IsActive returns true if the feature is not shipped or archived.
func (f *Feature) IsActive() bool {
	return f.Status != StatusShipped && f.Status != StatusArchived
}

// CanTransitionTo checks if the feature can transition to the given status.
// Following State Machine pattern.
func (f *Feature) CanTransitionTo(target FeatureStatus) bool {
	validTransitions := map[FeatureStatus][]FeatureStatus{
		StatusDraft:        {StatusSpecifying},
		StatusSpecifying:   {StatusResearched, StatusDraft},
		StatusResearched:   {StatusPlanned, StatusSpecifying},
		StatusPlanned:      {StatusImplementing, StatusResearched},
		StatusImplementing: {StatusValidated, StatusPlanned},
		StatusValidated:    {StatusShipped, StatusImplementing},
		StatusShipped:      {StatusArchived},
	}

	allowed, ok := validTransitions[f.Status]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == target {
			return true
		}
	}
	return false
}

// WorkPackage represents a decomposed unit of implementation.
// Following DDD Aggregate pattern.
type WorkPackage struct {
	ID                 string
	FeatureID          string
	Name               string
	Description        string
	AcceptanceCriteria []string
	Dependencies       []string
	AssignedAgent      string
	PRURL              string
	Status             WPStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// WPStatus represents the state of a work package.
type WPStatus string

const (
	WPStatusPlanned WPStatus = "planned"
	WPStatusDoing   WPStatus = "doing"
	WPStatusReview  WPStatus = "review"
	WPStatusDone    WPStatus = "done"
)

// GovernanceContract defines required evidence for state transitions.
// Following Policy pattern.
type GovernanceContract struct {
	ID        string
	FeatureID string
	Version   int
	FRs       []FunctionalRequirement
	Policies  []PolicyRule
	CreatedAt time.Time
}

// FunctionalRequirement represents a functional requirement.
// Following BDD/SpecDD patterns.
type FunctionalRequirement struct {
	ID          string
	Code        string // e.g., "FR-001"
	Description string
	Evidence    []Evidence
	Satisfied   bool
}

// PolicyRule defines a quality/security/reliability check.
type PolicyRule struct {
	ID        string
	Domain    string // quality, security, reliability
	Name      string
	Command   string
	Threshold float64
}

// Evidence represents test results, CI output, or review approvals.
type Evidence struct {
	ID         string
	FRCode     string
	Type       string // test, ci, review
	URL        string
	VerifiedAt time.Time
	VerifiedBy string
}

// AuditEntry represents a hash-chained audit record.
// Following Event Sourcing pattern.
type AuditEntry struct {
	ID             string
	FeatureID      string
	TransitionType string
	FromStatus     string
	ToStatus       string
	EvidenceRefs   []string
	PreviousHash   string
	Hash           string
	Timestamp      time.Time
	Actor          string
}
