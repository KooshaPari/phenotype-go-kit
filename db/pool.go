package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var (
	DefaultMaxOpenConns    = 25
	DefaultMaxIdleConns    = 5
	DefaultConnMaxLifetime = 5 * time.Minute
	DefaultConnMaxIdleTime = 1 * time.Minute
)

// PoolConfig holds database connection pool configuration.
type PoolConfig struct {
	MaxOpenConns    int           `default:"25"`
	MaxIdleConns    int           `default:"5"`
	ConnMaxLifetime time.Duration `default:"5m"`
	ConnMaxIdleTime time.Duration `default:"1m"`
	DialTimeout     time.Duration `default:"5s"`
	QueryTimeout    time.Duration `default:"30s"`
}

// ConfigurePool sets up the connection pool with the given config.
func ConfigurePool(db *sql.DB, cfg PoolConfig) error {
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
	return nil
}

// PoolStats holds current pool statistics.
type PoolStats struct {
	OpenConnections  int
	IdleConnections  int
	InUseConnections int
	WaitCount        int64
	WaitDuration     time.Duration
}

// GetPoolStats returns current connection pool statistics.
func GetPoolStats(db *sql.DB) PoolStats {
	stats := db.Stats()
	return PoolStats{
		OpenConnections:  stats.OpenConnections,
		IdleConnections:  stats.Idle,
		InUseConnections: stats.InUse,
		WaitCount:        stats.WaitCount,
		WaitDuration:     stats.WaitDuration,
	}
}

// DBMetrics holds database performance metrics.
type DBMetrics struct {
	QueriesExecuted  int64
	QueriesFailed    int64
	AvgQueryDuration time.Duration
	SlowQueries      int64
}

// WithTimeout executes a query with the configured timeout.
func WithTimeout(ctx context.Context, db *sql.DB, query string, timeout time.Duration) (func() error, error) {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)

	done := make(chan struct{})
	var resultErr error

	go func() {
		defer close(done)
		_, resultErr = db.ExecContext(ctx, query)
	}()

	waitDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-done:
		}
		close(waitDone)
	}()

	<-waitDone
	return func() error { return resultErr }, resultErr
}

// ConnectionCheck verifies database connectivity.
func ConnectionCheck(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

// HealthCheck returns pool health status.
func HealthCheck(ctx context.Context, db *sql.DB) error {
	if err := ConnectionCheck(ctx, db); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	stats := GetPoolStats(db)
	if stats.OpenConnections >= DefaultMaxOpenConns {
		return fmt.Errorf("connection pool exhausted: %d/%d", stats.OpenConnections, DefaultMaxOpenConns)
	}

	return nil
}
