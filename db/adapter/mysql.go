package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
)

// MySQLAdapter implements outbound.QueryExecutor for MySQL.
type MySQLAdapter struct {
	db           *sql.DB
	queryTimeout time.Duration
}

// NewMySQLAdapter creates a new MySQL adapter.
func NewMySQLAdapter(dsn string, config outbound.PoolConfig) (*MySQLAdapter, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	adapter := &MySQLAdapter{
		db:           db,
		queryTimeout: time.Duration(config.QueryTimeout) * time.Millisecond,
	}

	if err := adapter.configurePool(config); err != nil {
		db.Close()
		return nil, err
	}

	return adapter, nil
}

// NewMySQLAdapterFromDB creates an adapter from an existing database connection.
func NewMySQLAdapterFromDB(db *sql.DB) *MySQLAdapter {
	return &MySQLAdapter{
		db:           db,
		queryTimeout: 30 * time.Second,
	}
}

// configurePool sets up the connection pool.
func (a *MySQLAdapter) configurePool(config outbound.PoolConfig) error {
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
func (a *MySQLAdapter) Query(ctx context.Context, query string, args ...any) (outbound.Rows, error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &mysqlRows{rows: rows}, nil
}

// QueryRow implements outbound.QueryExecutor.
func (a *MySQLAdapter) QueryRow(ctx context.Context, query string, args ...any) outbound.Row {
	return &mysqlRow{row: a.db.QueryRowContext(ctx, query, args...)}
}

// Exec implements outbound.QueryExecutor.
func (a *MySQLAdapter) Exec(ctx context.Context, query string, args ...any) (outbound.Result, error) {
	result, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &mysqlResult{result: result}, nil
}

// BeginTx implements outbound.QueryExecutor.
func (a *MySQLAdapter) BeginTx(ctx context.Context, opts outbound.TxOptions) (outbound.Transaction, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.IsolationLevel(opts.Isolation),
		ReadOnly:  opts.ReadOnly,
	})
	if err != nil {
		return nil, err
	}
	return &mysqlTransaction{tx: tx}, nil
}

// Ping implements outbound.QueryExecutor.
func (a *MySQLAdapter) Ping(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

// Stats implements outbound.QueryExecutor.
func (a *MySQLAdapter) Stats() outbound.PoolStats {
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
func (a *MySQLAdapter) Close() error {
	return a.db.Close()
}

// mysqlRow wraps sql.Row.
type mysqlRow struct {
	row *sql.Row
}

func (r *mysqlRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func (r *mysqlRow) Err() error {
	return nil
}

// mysqlRows wraps sql.Rows.
type mysqlRows struct {
	rows *sql.Rows
}

func (r *mysqlRows) Next() bool {
	return r.rows.Next()
}

func (r *mysqlRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *mysqlRows) Err() error {
	return r.rows.Err()
}

func (r *mysqlRows) Close() error {
	return r.rows.Close()
}

func (r *mysqlRows) Columns() ([]string, error) {
	return r.rows.Columns()
}

func (r *mysqlRows) ScanSlice(dest any) error {
	return r.rows.Scan(dest)
}

func (r *mysqlRows) ScanMap(dest map[string]any) error {
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

// mysqlResult wraps sql.Result.
type mysqlResult struct {
	result sql.Result
}

func (r *mysqlResult) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

func (r *mysqlResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}

// mysqlTransaction wraps sql.Tx.
type mysqlTransaction struct {
	tx *sql.Tx
}

func (t *mysqlTransaction) Query(ctx context.Context, query string, args ...any) (outbound.Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &mysqlRows{rows: rows}, nil
}

func (t *mysqlTransaction) QueryRow(ctx context.Context, query string, args ...any) outbound.Row {
	return &mysqlRow{row: t.tx.QueryRowContext(ctx, query, args...)}
}

func (t *mysqlTransaction) Exec(ctx context.Context, query string, args ...any) (outbound.Result, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &mysqlResult{result: result}, nil
}

func (t *mysqlTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *mysqlTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

func (t *mysqlTransaction) CommitAsync() error {
	return t.tx.Commit()
}

func (t *mysqlTransaction) Ping(ctx context.Context) error {
	var exists bool
	return t.tx.QueryRowContext(ctx, "SELECT 1").Scan(&exists)
}

func (t *mysqlTransaction) Stats() outbound.PoolStats {
	return outbound.PoolStats{}
}

func (t *mysqlTransaction) BeginTx(ctx context.Context, opts outbound.TxOptions) (outbound.Transaction, error) {
	return nil, fmt.Errorf("nested transactions not supported in MySQL")
}
