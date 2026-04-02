package outbound

import (
	"context"
)

// QueryResult represents the result of a database query with pagination.
type QueryResult[T any] struct {
	Data       []T
	TotalCount int64
	Page       int
	PageSize   int
	HasNext    bool
	HasPrev    bool
}

// PaginationParams holds pagination parameters.
type PaginationParams struct {
	Page     int
	PageSize int
}

// NewPaginationParams creates pagination params with defaults.
func NewPaginationParams(page, pageSize int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PaginationParams{Page: page, PageSize: pageSize}
}

// Offset returns the SQL offset.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// QueryOptions holds optional query parameters.
type QueryOptions struct {
	Timeout       int // milliseconds, 0 = default
	SlowThreshold int // milliseconds, 0 = disabled
	MaxRows       int // 0 = unlimited
	ReadOnly      bool
	NoWait        bool
}

// QueryExecutor defines the interface for executing queries.
// This is the primary outbound port for database operations.
type QueryExecutor interface {
	// Query executes a SELECT query and returns rows.
	Query(ctx context.Context, query string, args ...any) (Rows, error)

	// QueryRow executes a query that returns a single row.
	QueryRow(ctx context.Context, query string, args ...any) Row

	// Exec executes a query that doesn't return rows (INSERT, UPDATE, DELETE).
	Exec(ctx context.Context, query string, args ...any) (Result, error)

	// BeginTx starts a new transaction with options.
	BeginTx(ctx context.Context, opts TxOptions) (Transaction, error)

	// Ping checks database connectivity.
	Ping(ctx context.Context) error

	// Stats returns connection pool statistics.
	Stats() PoolStats
}

// TxOptions holds transaction options.
type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
	NoWait    bool
	Immediate bool
}

// IsolationLevel represents SQL isolation levels.
type IsolationLevel int

const (
	IsolationDefault         IsolationLevel = 0
	IsolationReadUncommitted IsolationLevel = 1
	IsolationReadCommitted   IsolationLevel = 2
	IsolationWriteCommitted  IsolationLevel = 3
	IsolationRepeatableRead  IsolationLevel = 4
	IsolationSnapshot        IsolationLevel = 5
	IsolationSerializable    IsolationLevel = 6
	IsolationLinearizable    IsolationLevel = 7
)

// Transaction represents a database transaction.
type Transaction interface {
	QueryExecutor

	// Commit commits the transaction.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction.
	Rollback(ctx context.Context) error

	// CommitAsync commits asynchronously.
	CommitAsync() error
}

// Row represents a single row returned by QueryRow.
type Row interface {
	Scan(dest ...any) error
	Err() error
}

// Rows represents a set of rows returned by Query.
type Rows interface {
	Row

	// Next advances to the next row.
	Next() bool

	// Close closes the rows.
	Close() error

	// Columns returns the column names.
	Columns() ([]string, error)

	// ScanSlice scans the current row into a slice.
	ScanSlice(dest any) error

	// ScanMap scans the current row into a map.
	ScanMap(dest map[string]any) error
}

// Result represents the result of an Exec operation.
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// PoolStats holds connection pool statistics.
type PoolStats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUseConnections   int
	IdleConnections    int
	WaitCount          int64
	WaitDuration       int64 // nanoseconds
	MaxIdleClosed      int64
	MaxLifetimeClosed  int64
}

// ConnectionPool provides connection pool management.
type ConnectionPool interface {
	// ConfigurePool sets up the connection pool.
	ConfigurePool(config PoolConfig) error

	// GetPoolStats returns current pool statistics.
	GetPoolStats(ctx context.Context) (PoolStats, error)

	// HealthCheck verifies database health.
	HealthCheck(ctx context.Context) error

	// Close closes all connections.
	Close() error
}

// PoolConfig holds database connection pool configuration.
type PoolConfig struct {
	MaxOpenConns      int  // default: 25
	MaxIdleConns      int  // default: 5
	ConnMaxLifetime   int  // milliseconds, default: 5 minutes
	ConnMaxIdleTime   int  // milliseconds, default: 1 minute
	DialTimeout       int  // milliseconds, default: 5 seconds
	QueryTimeout      int  // milliseconds, default: 30 seconds
	EnableHealthCheck bool // default: true
}

// DefaultPoolConfig returns the default pool configuration.
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:      25,
		MaxIdleConns:      5,
		ConnMaxLifetime:   300000, // 5 minutes
		ConnMaxIdleTime:   60000,  // 1 minute
		DialTimeout:       5000,   // 5 seconds
		QueryTimeout:      30000,  // 30 seconds
		EnableHealthCheck: true,
	}
}

// IndexManager defines the interface for managing database indexes.
type IndexManager interface {
	// CreateIndex creates an index.
	CreateIndex(ctx context.Context, def IndexDefinition) error

	// DropIndex drops an index.
	DropIndex(ctx context.Context, name string) error

	// ListIndexes returns all indexes for a table.
	ListIndexes(ctx context.Context, table string) ([]IndexDefinition, error)

	// CreateIndexes creates multiple indexes efficiently.
	CreateIndexes(ctx context.Context, defs []IndexDefinition) error
}

// IndexDefinition represents a database index definition.
type IndexDefinition struct {
	Name         string
	Table        string
	Columns      []string
	Unique       bool
	Partial      string // SQL WHERE clause for partial index
	Concurrently bool   // CREATE INDEX CONCURRENTLY
}

// MigrationExecutor defines the interface for running migrations.
type MigrationExecutor interface {
	// Up runs pending migrations.
	Up(ctx context.Context) error

	// Down rolls back the last migration.
	Down(ctx context.Context) error

	// Version returns the current migration version.
	Version(ctx context.Context) (int, error)

	// Pending returns pending migrations.
	Pending(ctx context.Context) ([]Migration, error)
}

// Migration represents a database migration.
type Migration struct {
	Version   int
	Name      string
	AppliedAt int64
}
