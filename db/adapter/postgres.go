package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
)

// PostgresAdapter implements outbound.QueryExecutor for PostgreSQL.
type PostgresAdapter struct {
	db               *sql.DB
	poolConfig       outbound.PoolConfig
	slowQueryTimeout time.Duration
}

// NewPostgresAdapter creates a new PostgreSQL adapter.
func NewPostgresAdapter(dsn string, config outbound.PoolConfig) (*PostgresAdapter, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	adapter := &PostgresAdapter{
		db:               db,
		poolConfig:       config,
		slowQueryTimeout: time.Duration(config.QueryTimeout) * time.Millisecond,
	}

	if err := adapter.configurePool(config); err != nil {
		db.Close()
		return nil, err
	}

	return adapter, nil
}

// NewPostgresAdapterFromDB creates an adapter from an existing database connection.
func NewPostgresAdapterFromDB(db *sql.DB) *PostgresAdapter {
	return &PostgresAdapter{
		db: db,
		poolConfig: outbound.PoolConfig{
			MaxOpenConns:      25,
			MaxIdleConns:      5,
			ConnMaxLifetime:   300000,
			ConnMaxIdleTime:   60000,
			DialTimeout:       5000,
			QueryTimeout:      30000,
			EnableHealthCheck: true,
		},
		slowQueryTimeout: 30 * time.Second,
	}
}

// configurePool sets up the connection pool.
func (a *PostgresAdapter) configurePool(config outbound.PoolConfig) error {
	if config.MaxOpenConns > 0 {
		a.db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		a.db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		a.db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Millisecond)
	}
	if config.ConnMaxIdleTime > 0 {
		a.db.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Millisecond)
	}
	return nil
}

// Query implements outbound.QueryExecutor.
func (a *PostgresAdapter) Query(ctx context.Context, query string, args ...any) (outbound.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, a.slowQueryTimeout)
	defer cancel()

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return &postgresRows{rows: rows}, nil
}

// QueryRow implements outbound.QueryExecutor.
func (a *PostgresAdapter) QueryRow(ctx context.Context, query string, args ...any) outbound.Row {
	ctx, cancel := context.WithTimeout(ctx, a.slowQueryTimeout)
	defer cancel()

	return &postgresRow{row: a.db.QueryRowContext(ctx, query, args...)}
}

// Exec implements outbound.QueryExecutor.
func (a *PostgresAdapter) Exec(ctx context.Context, query string, args ...any) (outbound.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, a.slowQueryTimeout)
	defer cancel()

	result, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec failed: %w", err)
	}
	return &postgresResult{result: result}, nil
}

// BeginTx implements outbound.QueryExecutor.
func (a *PostgresAdapter) BeginTx(ctx context.Context, opts outbound.TxOptions) (outbound.Transaction, error) {
	sqlOpts := &sql.TxOptions{
		Isolation: sql.IsolationLevel(opts.Isolation),
		ReadOnly:  opts.ReadOnly,
	}

	tx, err := a.db.BeginTx(ctx, sqlOpts)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}

	return &postgresTransaction{tx: tx, queryTimeout: a.slowQueryTimeout}, nil
}

// Ping implements outbound.QueryExecutor.
func (a *PostgresAdapter) Ping(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

// Stats implements outbound.QueryExecutor.
func (a *PostgresAdapter) Stats() outbound.PoolStats {
	stats := a.db.Stats()
	return outbound.PoolStats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUseConnections:   stats.InUse,
		IdleConnections:    stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration.Nanoseconds(),
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
	}
}

// Close closes the database connection.
func (a *PostgresAdapter) Close() error {
	return a.db.Close()
}

// postgresRow wraps sql.Row.
type postgresRow struct {
	row *sql.Row
}

func (r *postgresRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func (r *postgresRow) Err() error {
	return nil // sql.Row doesn't have an Err method
}

// postgresRows wraps sql.Rows.
type postgresRows struct {
	rows *sql.Rows
}

func (r *postgresRows) Next() bool {
	return r.rows.Next()
}

func (r *postgresRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *postgresRows) Err() error {
	return r.rows.Err()
}

func (r *postgresRows) Close() error {
	return r.rows.Close()
}

func (r *postgresRows) Columns() ([]string, error) {
	return r.rows.Columns()
}

func (r *postgresRows) ScanSlice(dest any) error {
	return r.rows.Scan(dest)
}

func (r *postgresRows) ScanMap(dest map[string]any) error {
	columns, err := r.rows.Columns()
	if err != nil {
		return err
	}

	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(any)
	}

	if err := r.rows.Scan(values...); err != nil {
		return err
	}

	for i, col := range columns {
		dest[col] = *values[i].(*any)
	}
	return nil
}

// postgresResult wraps sql.Result.
type postgresResult struct {
	result sql.Result
}

func (r *postgresResult) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

func (r *postgresResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}

// postgresTransaction wraps sql.Tx.
type postgresTransaction struct {
	tx           *sql.Tx
	queryTimeout time.Duration
}

func (t *postgresTransaction) Query(ctx context.Context, query string, args ...any) (outbound.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, t.queryTimeout)
	defer cancel()

	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &postgresRows{rows: rows}, nil
}

func (t *postgresTransaction) QueryRow(ctx context.Context, query string, args ...any) outbound.Row {
	ctx, cancel := context.WithTimeout(ctx, t.queryTimeout)
	defer cancel()

	return &postgresRow{row: t.tx.QueryRowContext(ctx, query, args...)}
}

func (t *postgresTransaction) Exec(ctx context.Context, query string, args ...any) (outbound.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, t.queryTimeout)
	defer cancel()

	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &postgresResult{result: result}, nil
}

func (t *postgresTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *postgresTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

func (t *postgresTransaction) CommitAsync() error {
	return t.tx.Commit()
}

func (t *postgresTransaction) Ping(ctx context.Context) error {
	var exists bool
	return t.tx.QueryRowContext(ctx, "SELECT true").Scan(&exists)
}

func (t *postgresTransaction) Stats() outbound.PoolStats {
	return outbound.PoolStats{}
}

func (t *postgresTransaction) BeginTx(ctx context.Context, opts outbound.TxOptions) (outbound.Transaction, error) {
	return nil, fmt.Errorf("nested transactions not supported in postgres")
}
