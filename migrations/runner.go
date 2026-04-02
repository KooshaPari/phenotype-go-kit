package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
)

// Migration represents a database migration.
type Migration struct {
	Version string
	Name    string
	Up      func(*sql.Tx) error
	Down    func(*sql.Tx) error
}

// MigrationRunner manages database migrations.
type MigrationRunner struct {
	db         *sql.DB
	migrations []Migration
	table      string
	logger     *slog.Logger
}

// NewMigrationRunner creates a new migration runner.
func NewMigrationRunner(db *sql.DB, migrations []Migration, logger *slog.Logger) *MigrationRunner {
	return &MigrationRunner{
		db:         db,
		migrations: migrations,
		table:      "schema_migrations",
		logger:     logger,
	}
}

// Init creates the migrations table if it doesn't exist.
func (r *MigrationRunner) Init(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`, r.table)

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Up runs all pending migrations.
func (r *MigrationRunner) Up(ctx context.Context) error {
	applied, err := r.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	sort.Slice(r.migrations, func(i, j int) bool {
		return r.migrations[i].Version < r.migrations[j].Version
	})

	for _, m := range r.migrations {
		if _, ok := applied[m.Version]; ok {
			r.logger.Debug("migration already applied", "version", m.Version)
			continue
		}

		r.logger.Info("applying migration", "version", m.Version, "name", m.Name)

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if err := m.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", m.Version, err)
		}

		if err := r.recordMigration(ctx, tx, m.Version, m.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration: %w", err)
		}

		r.logger.Info("migration applied", "version", m.Version)
	}

	return nil
}

// Down rolls back the last migration.
func (r *MigrationRunner) Down(ctx context.Context) error {
	applied, err := r.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	sort.Slice(r.migrations, func(i, j int) bool {
		return r.migrations[i].Version > r.migrations[j].Version
	})

	for _, m := range r.migrations {
		if _, ok := applied[m.Version]; !ok {
			continue
		}

		r.logger.Info("rolling back migration", "version", m.Version, "name", m.Name)

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if err := m.Down(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("rollback %s failed: %w", m.Version, err)
		}

		if err := r.removeMigration(ctx, tx, m.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit rollback: %w", err)
		}

		r.logger.Info("migration rolled back", "version", m.Version)
		break // Only rollback one at a time
	}

	return nil
}

// Status returns the current migration status.
func (r *MigrationRunner) Status(ctx context.Context) (map[string]bool, error) {
	applied, err := r.getAppliedVersions(ctx)
	if err != nil {
		return nil, err
	}

	status := make(map[string]bool)
	for _, m := range r.migrations {
		status[m.Version] = applied[m.Version]
	}
	return status, nil
}

func (r *MigrationRunner) getAppliedVersions(ctx context.Context) (map[string]bool, error) {
	query := fmt.Sprintf("SELECT version FROM %s", r.table)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		versions[v] = true
	}
	return versions, rows.Err()
}

func (r *MigrationRunner) recordMigration(ctx context.Context, tx *sql.Tx, version, name string) error {
	query := fmt.Sprintf("INSERT INTO %s (version, name) VALUES (?, ?)", r.table)
	_, err := tx.ExecContext(ctx, query, version, name)
	return err
}

func (r *MigrationRunner) removeMigration(ctx context.Context, tx *sql.Tx, version string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE version = ?", r.table)
	_, err := tx.ExecContext(ctx, query, version)
	return err
}

// LoadMigrations loads migrations from a directory.
func LoadMigrations(dir string) ([]Migration, error) {
	// In production, scan directory for migration files
	// and dynamically load them
	return []Migration{
		{
			Version: "001",
			Name:    "initial_schema",
			Up: func(tx *sql.Tx) error {
				// Users table
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS users (
						id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
						email VARCHAR(255) UNIQUE NOT NULL,
						password_hash VARCHAR(255) NOT NULL,
						name VARCHAR(255),
						role VARCHAR(50) DEFAULT 'user',
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						deleted_at TIMESTAMP
					);
					CREATE INDEX idx_users_email ON users(email);
				`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec("DROP TABLE IF EXISTS users")
				return err
			},
		},
		{
			Version: "002",
			Name:    "webhooks_table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS webhooks (
						id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
						user_id UUID REFERENCES users(id) ON DELETE CASCADE,
						url VARCHAR(2048) NOT NULL,
						secret VARCHAR(255),
						events TEXT[],
						active BOOLEAN DEFAULT true,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
					);
					CREATE INDEX idx_webhooks_user_id ON webhooks(user_id);
				`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec("DROP TABLE IF EXISTS webhooks")
				return err
			},
		},
		{
			Version: "003",
			Name:    "jobs_table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS jobs (
						id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
						type VARCHAR(100) NOT NULL,
						payload JSONB,
						status VARCHAR(50) DEFAULT 'pending',
						retries INTEGER DEFAULT 0,
						max_retries INTEGER DEFAULT 3,
						error TEXT,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						started_at TIMESTAMP,
						completed_at TIMESTAMP,
						failed_at TIMESTAMP
					);
					CREATE INDEX idx_jobs_status ON jobs(status);
					CREATE INDEX idx_jobs_created_at ON jobs(created_at);
				`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec("DROP TABLE IF EXISTS jobs")
				return err
			},
		},
	}, nil
}

// SeedData holds seed data for the database.
type SeedData struct {
	Users    []UserSeed
	Webhooks []WebhookSeed
	Jobs     []JobSeed
}

// UserSeed represents a user to seed.
type UserSeed struct {
	Email    string
	Name     string
	Password string // Will be hashed
	Role     string
}

// WebhookSeed represents a webhook to seed.
type WebhookSeed struct {
	UserID string
	URL    string
	Events []string
}

// JobSeed represents a job to seed.
type JobSeed struct {
	Type    string
	Payload string
	Status  string
}

// Seeder seeds the database with initial data.
type Seeder struct {
	db     *sql.DB
	data   SeedData
	logger *slog.Logger
}

// NewSeeder creates a new database seeder.
func NewSeeder(db *sql.DB, data SeedData, logger *slog.Logger) *Seeder {
	return &Seeder{
		db:     db,
		data:   data,
		logger: logger,
	}
}

// Seed runs all seed operations.
func (s *Seeder) Seed(ctx context.Context) error {
	// Seed users
	for _, u := range s.data.Users {
		// In production, hash the password properly
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO users (email, name, password_hash, role)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (email) DO NOTHING
		`, u.Email, u.Name, u.Password, u.Role)
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.Email, err)
		}
	}
	s.logger.Info("seeded users")

	// Seed webhooks
	for _, w := range s.data.Webhooks {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO webhooks (user_id, url, events, active)
			VALUES ($1, $2, $3, true)
			ON CONFLICT DO NOTHING
		`, w.UserID, w.URL, w.Events)
		if err != nil {
			return fmt.Errorf("failed to seed webhook: %w", err)
		}
	}
	s.logger.Info("seeded webhooks")

	return nil
}

// DevelopmentSeedData returns seed data for development.
func DevelopmentSeedData() SeedData {
	return SeedData{
		Users: []UserSeed{
			{Email: "admin@phenotype.dev", Name: "Admin User", Password: "admin123", Role: "admin"},
			{Email: "dev@phenotype.dev", Name: "Developer", Password: "dev123", Role: "developer"},
			{Email: "user@phenotype.dev", Name: "Test User", Password: "user123", Role: "user"},
		},
		Webhooks: []WebhookSeed{
			{UserID: "00000000-0000-0000-0000-000000000001", URL: "https://example.com/webhook", Events: []string{"user.created", "user.updated"}},
		},
	}
}

var _ = filepath.Join // used in production for scanning migration files
