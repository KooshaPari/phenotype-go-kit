package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// QueryConfig holds query optimization configuration.
type QueryConfig struct {
	Timeout       time.Duration
	SlowThreshold time.Duration
	MaxRows       int
	EnableAnalyze bool
	EnableExplain bool
}

// DefaultQueryConfig returns default configuration.
func DefaultQueryConfig() QueryConfig {
	return QueryConfig{
		Timeout:       30 * time.Second,
		SlowThreshold: 1 * time.Second,
		MaxRows:       10000,
		EnableAnalyze: false,
		EnableExplain: false,
	}
}

// QueryMetrics holds query performance metrics.
type QueryMetrics struct {
	Query        string
	Duration     time.Duration
	RowsAffected int64
	Timestamp    time.Time
	Slow         bool
	Error        error
}

// SlowQueryLogger logs slow queries.
type SlowQueryLogger struct {
	logger    *slog.Logger
	threshold time.Duration
}

// NewSlowQueryLogger creates a slow query logger.
func NewSlowQueryLogger(logger *slog.Logger, threshold time.Duration) *SlowQueryLogger {
	return &SlowQueryLogger{
		logger:    logger,
		threshold: threshold,
	}
}

// Log records a query execution.
func (l *SlowQueryLogger) Log(ctx context.Context, metrics QueryMetrics) {
	if metrics.Duration >= l.threshold {
		l.logger.Warn("slow query detected",
			"query", metrics.Query,
			"duration_ms", metrics.Duration.Milliseconds(),
			"rows", metrics.RowsAffected,
		)
	}
}

// QueryBuilder provides a fluent interface for building queries.
type QueryBuilder struct {
	table       string
	columns     []string
	where       []string
	orderBy     []string
	limitVal    int
	offsetVal   int
	groupBy     []string
	joinClauses []string
	params      []interface{}
}

// NewQueryBuilder creates a new query builder.
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table: table,
	}
}

// Select specifies columns to select.
func (qb *QueryBuilder) Select(cols ...string) *QueryBuilder {
	if len(cols) == 0 {
		qb.columns = []string{"*"}
	} else {
		qb.columns = cols
	}
	return qb
}

// Where adds a WHERE condition.
func (qb *QueryBuilder) Where(condition string, params ...interface{}) *QueryBuilder {
	qb.where = append(qb.where, condition)
	qb.params = append(qb.params, params...)
	return qb
}

// OrderBy adds ORDER BY clause.
func (qb *QueryBuilder) OrderBy(col string, desc bool) *QueryBuilder {
	if desc {
		qb.orderBy = append(qb.orderBy, col+" DESC")
	} else {
		qb.orderBy = append(qb.orderBy, col+" ASC")
	}
	return qb
}

// Limit adds LIMIT clause.
func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
	qb.limitVal = n
	return qb
}

// Offset adds OFFSET clause.
func (qb *QueryBuilder) Offset(n int) *QueryBuilder {
	qb.offsetVal = n
	return qb
}

// GroupBy adds GROUP BY clause.
func (qb *QueryBuilder) GroupBy(cols ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, cols...)
	return qb
}

// Join adds a JOIN clause.
func (qb *QueryBuilder) Join(table, condition string) *QueryBuilder {
	qb.joinClauses = append(qb.joinClauses, "JOIN "+table+" ON "+condition)
	return qb
}

// LeftJoin adds a LEFT JOIN clause.
func (qb *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	qb.joinClauses = append(qb.joinClauses, "LEFT JOIN "+table+" ON "+condition)
	return qb
}

// Build generates the SQL query.
func (qb *QueryBuilder) Build() (string, []interface{}) {
	sql := "SELECT "
	sql += joinStrings(qb.columns, ", ")
	sql += " FROM " + qb.table

	for _, join := range qb.joinClauses {
		sql += " " + join
	}

	if len(qb.where) > 0 {
		sql += " WHERE " + joinStrings(qb.where, " AND ")
	}

	if len(qb.groupBy) > 0 {
		sql += " GROUP BY " + joinStrings(qb.groupBy, ", ")
	}

	if len(qb.orderBy) > 0 {
		sql += " ORDER BY " + joinStrings(qb.orderBy, ", ")
	}

	if qb.limitVal > 0 {
		sql += fmt.Sprintf(" LIMIT %d", qb.limitVal)
	}

	if qb.offsetVal > 0 {
		sql += fmt.Sprintf(" OFFSET %d", qb.offsetVal)
	}

	return sql, qb.params
}

// Count builds a count query.
func (qb *QueryBuilder) Count() (string, []interface{}) {
	sql := "SELECT COUNT(*) FROM " + qb.table

	if len(qb.where) > 0 {
		sql += " WHERE " + joinStrings(qb.where, " AND ")
	}

	return sql, qb.params
}

// Exists builds an EXISTS query.
func (qb *QueryBuilder) Exists() (string, []interface{}) {
	sql := "SELECT EXISTS(SELECT 1 FROM " + qb.table

	if len(qb.where) > 0 {
		sql += " WHERE " + joinStrings(qb.where, " AND ")
	}

	sql += ")"

	return sql, qb.params
}

// Paginate adds pagination using limit and offset.
func (qb *QueryBuilder) Paginate(page, pageSize int) *QueryBuilder {
	qb.limitVal = pageSize
	qb.offsetVal = (page - 1) * pageSize
	return qb
}

// PaginationResult holds pagination metadata.
type PaginationResult struct {
	Page       int
	PageSize   int
	TotalCount int64
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

// CalculatePagination calculates pagination metadata.
func CalculatePagination(page, pageSize int, totalCount int64) PaginationResult {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

func joinStrings(items []string, sep string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += sep
		}
		result += item
	}
	return result
}

// ExplainQuery explains a query execution plan.
func ExplainQuery(ctx context.Context, db *sql.DB, query string, params ...interface{}) (string, error) {
	rows, err := db.QueryContext(ctx, "EXPLAIN "+query, params...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return "", err
		}
		result += line + "\n"
	}

	return result, rows.Err()
}
