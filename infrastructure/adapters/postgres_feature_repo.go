package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/KooshaPari/phenotype-go-kit/domain/entities"
	"github.com/KooshaPari/phenotype-go-kit/domain/ports"
)

// PostgresFeatureRepository implements FeatureRepositoryPort using PostgreSQL.
// Following Adapter pattern from Hexagonal Architecture.
//
// This adapter translates domain ports to infrastructure concerns,
// keeping the domain layer pure and independent of persistence details.
type PostgresFeatureRepository struct {
	db *sql.DB
}

// NewPostgresFeatureRepository creates a new PostgresFeatureRepository.
func NewPostgresFeatureRepository(db *sql.DB) *PostgresFeatureRepository {
	return &PostgresFeatureRepository{db: db}
}

// Save creates or updates a feature.
// Following upsert pattern.
func (r *PostgresFeatureRepository) Save(ctx context.Context, feature entities.Feature) (entities.Feature, error) {
	// Check if exists
	var existing string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM features WHERE id = $1", feature.ID).Scan(&existing)

	if err == sql.ErrNoRows {
		// Insert new
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO features (id, name, description, status, mission, kitty_spec, spec_md, research_md, plan_md, created_at, updated_at, shipped_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, feature.ID, feature.Name, feature.Description, feature.Status, feature.Mission, feature.KittySpec,
			feature.SpecMD, feature.ResearchMD, feature.PlanMD, feature.CreatedAt, feature.UpdatedAt, feature.ShippedAt)
	} else if err != nil {
		return feature, fmt.Errorf("checking existence: %w", err)
	} else {
		// Update existing
		_, err = r.db.ExecContext(ctx, `
			UPDATE features SET
				name = $2, description = $3, status = $4, mission = $5, kitty_spec = $6,
				spec_md = $7, research_md = $8, plan_md = $9, updated_at = $10, shipped_at = $11
			WHERE id = $1
		`, feature.ID, feature.Name, feature.Description, feature.Status, feature.Mission, feature.KittySpec,
			feature.SpecMD, feature.ResearchMD, feature.PlanMD, feature.UpdatedAt, feature.ShippedAt)
	}

	if err != nil {
		return feature, fmt.Errorf("saving feature: %w", err)
	}

	return feature, nil
}

// FindByID retrieves a feature by ID.
func (r *PostgresFeatureRepository) FindByID(ctx context.Context, id string) (*entities.Feature, error) {
	var f entities.Feature
	var specMD, researchMD, planMD sql.NullString
	var shippedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, status, mission, kitty_spec, spec_md, research_md, plan_md, created_at, updated_at, shipped_at
		FROM features WHERE id = $1
	`, id).Scan(&f.ID, &f.Name, &f.Description, &f.Status, &f.Mission, &f.KittySpec,
		&specMD, &researchMD, &planMD, &f.CreatedAt, &f.UpdatedAt, &shippedAt)

	if err == sql.ErrNoRows {
		return nil, entities.ErrFeatureNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding feature: %w", err)
	}

	f.SpecMD = specMD.String
	f.ResearchMD = researchMD.String
	f.PlanMD = planMD.String
	if shippedAt.Valid {
		f.ShippedAt = &shippedAt.Time
	}

	return &f, nil
}

// List returns all features matching the filter.
func (r *PostgresFeatureRepository) List(ctx context.Context, filter map[string]any) ([]entities.Feature, error) {
	query := "SELECT id, name, description, status, mission, kitty_spec, spec_md, research_md, plan_md, created_at, updated_at, shipped_at FROM features WHERE 1=1"
	args := []any{}
	argNum := 1

	if status, ok := filter["status"]; ok {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, status)
		argNum++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing features: %w", err)
	}
	defer rows.Close()

	var features []entities.Feature
	for rows.Next() {
		var f entities.Feature
		var specMD, researchMD, planMD sql.NullString
		var shippedAt sql.NullTime

		err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.Status, &f.Mission, &f.KittySpec,
			&specMD, &researchMD, &planMD, &f.CreatedAt, &f.UpdatedAt, &shippedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning feature: %w", err)
		}

		f.SpecMD = specMD.String
		f.ResearchMD = researchMD.String
		f.PlanMD = planMD.String
		if shippedAt.Valid {
			f.ShippedAt = &shippedAt.Time
		}

		features = append(features, f)
	}

	return features, nil
}

// Delete removes a feature by ID.
func (r *PostgresFeatureRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM features WHERE id = $1", id)
	return err
}

// GetGovernanceContract retrieves the governance contract for a feature.
func (r *PostgresFeatureRepository) GetGovernanceContract(ctx context.Context, featureID string) (*entities.GovernanceContract, error) {
	var contractJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT contract FROM governance_contracts WHERE feature_id = $1 AND active = true
	`, featureID).Scan(&contractJSON)

	if err == sql.ErrNoRows {
		return nil, &entities.DomainError{Code: "CONTRACT_NOT_FOUND", Message: "No active governance contract"}
	}
	if err != nil {
		return nil, fmt.Errorf("getting contract: %w", err)
	}

	var contract entities.GovernanceContract
	if err := json.Unmarshal(contractJSON, &contract); err != nil {
		return nil, fmt.Errorf("unmarshaling contract: %w", err)
	}

	return &contract, nil
}

// SaveGovernanceContract saves a governance contract.
func (r *PostgresFeatureRepository) SaveGovernanceContract(ctx context.Context, contract entities.GovernanceContract) error {
	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("marshaling contract: %w", err)
	}

	// Deactivate old contracts
	_, err = r.db.ExecContext(ctx, `
		UPDATE governance_contracts SET active = false WHERE feature_id = $1
	`, contract.FeatureID)
	if err != nil {
		return fmt.Errorf("deactivating old contracts: %w", err)
	}

	// Insert new contract
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO governance_contracts (feature_id, contract, active)
		VALUES ($1, $2, true)
	`, contract.FeatureID, contractJSON)

	return err
}

// Compile-time interface check - ensures adapter implements port
var _ ports.FeatureRepositoryPort = (*PostgresFeatureRepository)(nil)
