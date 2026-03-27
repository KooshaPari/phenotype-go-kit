package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/KooshaPari/phenotype-go-kit/domain/entities"
	"github.com/KooshaPari/phenotype-go-kit/domain/ports"
)

// PostgresAuditRepository implements AuditPort using PostgreSQL.
// Following Event Sourcing pattern with hash chains.
type PostgresAuditRepository struct {
	db *sql.DB
}

// NewPostgresAuditRepository creates a new PostgresAuditRepository.
func NewPostgresAuditRepository(db *sql.DB) *PostgresAuditRepository {
	return &PostgresAuditRepository{db: db}
}

// Record records an audit entry.
func (r *PostgresAuditRepository) Record(ctx context.Context, entry *entities.AuditEntry) error {
	evidenceJSON, _ := json.Marshal(entry.EvidenceRefs)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO audit_entries (id, feature_id, transition_type, from_status, to_status, evidence_refs, previous_hash, hash, timestamp, actor)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, entry.ID, entry.FeatureID, entry.TransitionType, entry.FromStatus, entry.ToStatus,
		evidenceJSON, entry.PreviousHash, entry.Hash, entry.Timestamp, entry.Actor)

	return err
}

// GetLatestHash returns the latest hash for a feature's audit chain.
func (r *PostgresAuditRepository) GetLatestHash(ctx context.Context, featureID string) (string, error) {
	var hash sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT hash FROM audit_entries WHERE feature_id = $1 ORDER BY timestamp DESC LIMIT 1
	`, featureID).Scan(&hash)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("getting latest hash: %w", err)
	}

	return hash.String, nil
}

// GetChain returns the full audit chain for a feature.
func (r *PostgresAuditRepository) GetChain(ctx context.Context, featureID string) ([]entities.AuditEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, feature_id, transition_type, from_status, to_status, evidence_refs, previous_hash, hash, timestamp, actor
		FROM audit_entries WHERE feature_id = $1 ORDER BY timestamp ASC
	`, featureID)
	if err != nil {
		return nil, fmt.Errorf("getting audit chain: %w", err)
	}
	defer rows.Close()

	var entries []entities.AuditEntry
	for rows.Next() {
		var e entities.AuditEntry
		var evidenceJSON []byte
		var fromStatus, toStatus sql.NullString

		err := rows.Scan(&e.ID, &e.FeatureID, &e.TransitionType, &fromStatus, &toStatus,
			&evidenceJSON, &e.PreviousHash, &e.Hash, &e.Timestamp, &e.Actor)
		if err != nil {
			return nil, fmt.Errorf("scanning audit entry: %w", err)
		}

		e.FromStatus = fromStatus.String
		e.ToStatus = toStatus.String
		json.Unmarshal(evidenceJSON, &e.EvidenceRefs)

		entries = append(entries, e)
	}

	return entries, nil
}

// VerifyChain verifies the integrity of an audit chain.
// Following Event Sourcing verification pattern.
func (r *PostgresAuditRepository) VerifyChain(ctx context.Context, featureID string) (bool, error) {
	entries, err := r.GetChain(ctx, featureID)
	if err != nil {
		return false, err
	}

	for i, entry := range entries {
		// Recompute hash
		computedHash := entities.ComputeHash(&entry)

		// Check stored hash matches
		if computedHash != entry.Hash {
			return false, fmt.Errorf("hash mismatch at entry %d", i)
		}

		// Check hash chain continuity
		if i > 0 {
			if entry.PreviousHash != entries[i-1].Hash {
				return false, fmt.Errorf("chain broken at entry %d", i)
			}
		}
	}

	return true, nil
}

// Compile-time interface check
var _ ports.AuditPort = (*PostgresAuditRepository)(nil)
