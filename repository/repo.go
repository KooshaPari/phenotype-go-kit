package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

// Repository defines the interface for data access.
type Repository interface {
	Create(ctx context.Context, model interface{}) error
	Read(ctx context.Context, id string) (interface{}, error)
	Update(ctx context.Context, model interface{}) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter interface{}, pagination interface{}) ([]interface{}, error)
}

// SQLRepository provides SQL database operations.
type SQLRepository struct {
	db     *sql.DB
	table  string
	logger *slog.Logger
}

// NewSQLRepository creates a new SQL repository.
func NewSQLRepository(db *sql.DB, table string) *SQLRepository {
	return &SQLRepository{
		db:     db,
		table:  table,
		logger: slog.Default(),
	}
}

// Create inserts a new record.
func (r *SQLRepository) Create(ctx context.Context, model interface{}) error {
	query, args, err := r.buildInsert(model)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("create failed", "query", query, "error", err)
		return err
	}

	return nil
}

// Read retrieves a record by ID.
func (r *SQLRepository) Read(ctx context.Context, id string) (interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", r.table)

	row := r.db.QueryRowContext(ctx, query, id)

	var result map[string]interface{}
	if err := row.Scan(result); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not found: %s", id)
		}
		return nil, err
	}

	return result, nil
}

// Update updates an existing record.
func (r *SQLRepository) Update(ctx context.Context, model interface{}) error {
	query, args, err := r.buildUpdate(model)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

// Delete removes a record.
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.table)

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

// List retrieves multiple records with filtering and pagination.
func (r *SQLRepository) List(ctx context.Context, filter interface{}, pagination interface{}) ([]interface{}, error) {
	query, args := r.buildListQuery(filter, pagination)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []interface{}
	for rows.Next() {
		var result map[string]interface{}
		// Scan would go here
		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *SQLRepository) buildInsert(model interface{}) (string, []interface{}, error) {
	// Simplified - would need reflection to build actual query
	return "", nil, nil
}

func (r *SQLRepository) buildUpdate(model interface{}) (string, []interface{}, error) {
	return "", nil, nil
}

func (r *SQLRepository) buildListQuery(filter, pagination interface{}) (string, []interface{}) {
	return fmt.Sprintf("SELECT * FROM %s", r.table), nil
}

// Filter holds query filter parameters.
type Filter struct {
	Field    string
	Value    interface{}
	Operator string // eq, ne, gt, lt, like, in
}

// Pagination holds pagination parameters.
type Pagination struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string // asc, desc
}

// NewPagination creates default pagination.
func NewPagination(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset calculates the offset for pagination.
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// BuildWhereClause builds WHERE clause from filters.
func BuildWhereClause(filters []Filter) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	clauses := make([]string, 0, len(filters))
	args := make([]interface{}, 0, len(filters))

	for _, f := range filters {
		switch f.Operator {
		case "eq":
			clauses = append(clauses, fmt.Sprintf("%s = ?", f.Field))
			args = append(args, f.Value)
		case "ne":
			clauses = append(clauses, fmt.Sprintf("%s != ?", f.Field))
			args = append(args, f.Value)
		case "gt":
			clauses = append(clauses, fmt.Sprintf("%s > ?", f.Field))
			args = append(args, f.Value)
		case "lt":
			clauses = append(clauses, fmt.Sprintf("%s < ?", f.Field))
			args = append(args, f.Value)
		case "like":
			clauses = append(clauses, fmt.Sprintf("%s LIKE ?", f.Field))
			args = append(args, f.Value)
		case "in":
			clauses = append(clauses, fmt.Sprintf("%s IN (?)", f.Field))
			args = append(args, f.Value)
		}
	}

	where := " WHERE " + clauses[0]
	for i := 1; i < len(clauses); i++ {
		where += " AND " + clauses[i]
	}

	return where, args
}

// BuildPaginationClause builds ORDER BY and LIMIT clauses.
func BuildPaginationClause(p *Pagination) string {
	clause := fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d",
		p.SortBy, p.SortDir, p.PageSize, p.Offset())
	return clause
}
