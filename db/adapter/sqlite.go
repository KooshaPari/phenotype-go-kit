package adapter

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
)

// SQLiteAdapter implements outbound.QueryExecutor for SQLite.
type SQLiteAdapter struct {
	db           *sql.DB
	queryTimeout int // milliseconds
}

// NewSQLiteAdapter creates a new SQLite adapter.
func NewSQLiteAdapter(dsn string, config outbound.PoolConfig) (*SQLiteAdapter, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	adapter := &SQLiteAdapter{
		db:           db,
		queryTimeout: config.QueryTimeout,
	}

	// SQLite-specific settings
	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")

	return adapter, nil
}

// NewSQLiteAdapterFromDB creates an adapter from an existing database connection.
func NewSQLiteAdapterFromDB(db *sql.DB) *SQLiteAdapter {
	return &SQLiteAdapter{
		db:           db,
		queryTimeout: 30000,
	}
}

// Query implements outbound.QueryExecutor.
func (a *SQLiteAdapter) Query(ctx context.Context, query string, args ...any) (outbound.Rows, error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqliteRows{rows: rows}, nil
}

// QueryRow implements outbound.QueryExecutor.
func (a *SQLiteAdapter) QueryRow(ctx context.Context, query string, args ...any) outbound.Row {
	return &sqliteRow{row: a.db.QueryRowContext(ctx, query, args...)}
}

// Exec implements outbound.QueryExecutor.
func (a *SQLiteAdapter) Exec(ctx context.Context, query string, args ...any) (outbound.Result, error) {
	result, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqliteResult{result: result}, nil
}

// BeginTx implements outbound.QueryExecutor.
func (a *SQLiteAdapter) BeginTx(ctx context.Context, opts outbound.TxOptions) (outbound.Transaction, error) {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqliteTransaction{tx: tx}, nil
}

// Ping implements outbound.QueryExecutor.
func (a *SQLiteAdapter) Ping(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

// Stats implements outbound.QueryExecutor.
func (a *SQLiteAdapter) Stats() outbound.PoolStats {
	stats := a.db.Stats()
	return outbound.PoolStats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUseConnections:   stats.InUse,
		IdleConnections:    stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration.Nanoseconds(),
	}
}

// Close closes the database connection.
func (a *SQLiteAdapter) Close() error {
	return a.db.Close()
}

// sqliteRow wraps sql.Row.
type sqliteRow struct {
	row *sql.Row
}

func (r *sqliteRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func (r *sqliteRow) Err() error {
	return nil
}

// sqliteRows wraps sql.Rows.
type sqliteRows struct {
	rows *sql.Rows
}

func (r *sqliteRows) Next() bool {
	return r.rows.Next()
}

func (r *sqliteRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *sqliteRows) Err() error {
	return r.rows.Err()
}

func (r *sqliteRows) Close() error {
	return r.rows.Close()
}

func (r *sqliteRows) Columns() ([]string, error) {
	return r.rows.Columns()
}

func (r *sqliteRows) ScanSlice(dest any) error {
	return r.rows.Scan(dest)
}

func (r *sqliteRows) ScanMap(dest map[string]any) error {
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

// sqliteResult wraps sql.Result.
type sqliteResult struct {
	result sql.Result
}

func (r *sqliteResult) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

func (r *sqliteResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}

// sqliteTransaction wraps sql.Tx.
type sqliteTransaction struct {
	tx *sql.Tx
}

func (t *sqliteTransaction) Query(ctx context.Context, query string, args ...any) (outbound.Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqliteRows{rows: rows}, nil
}

func (t *sqliteTransaction) QueryRow(ctx context.Context, query string, args ...any) outbound.Row {
	return &sqliteRow{row: t.tx.QueryRowContext(ctx, query, args...)}
}

func (t *sqliteTransaction) Exec(ctx context.Context, query string, args ...any) (outbound.Result, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqliteResult{result: result}, nil
}

func (t *sqliteTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *sqliteTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

func (t *sqliteTransaction) CommitAsync() error {
	return t.tx.Commit()
}

func (t *sqliteTransaction) Ping(ctx context.Context) error {
	var exists bool
	return t.tx.QueryRowContext(ctx, "SELECT 1").Scan(&exists)
}

func (t *sqliteTransaction) Stats() outbound.PoolStats {
	return outbound.PoolStats{}
}

func (t *sqliteTransaction) BeginTx(ctx context.Context, opts outbound.TxOptions) (outbound.Transaction, error) {
	return nil, fmt.Errorf("nested transactions not supported in SQLite")
}
