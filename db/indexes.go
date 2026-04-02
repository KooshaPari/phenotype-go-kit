package db

// IndexDefinition represents a database index.
type IndexDefinition struct {
	Name    string   `json:"name"`
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
	Partial string   `json:"partial,omitempty"` // SQL WHERE clause for partial index
}

// Common indexes for the application.
var Indexes = []IndexDefinition{
	// Users table indexes
	{
		Name:    "idx_users_email",
		Table:   "users",
		Columns: []string{"email"},
		Unique:  true,
	},
	{
		Name:    "idx_users_status",
		Table:   "users",
		Columns: []string{"status"},
	},
	{
		Name:    "idx_users_created_at",
		Table:   "users",
		Columns: []string{"created_at"},
	},
	{
		Name:    "idx_users_updated_at",
		Table:   "users",
		Columns: []string{"updated_at"},
	},

	// Webhooks table indexes
	{
		Name:    "idx_webhooks_user_id",
		Table:   "webhooks",
		Columns: []string{"user_id"},
	},
	{
		Name:    "idx_webhooks_url",
		Table:   "webhooks",
		Columns: []string{"url"},
	},
	{
		Name:    "idx_webhooks_event_type",
		Table:   "webhooks",
		Columns: []string{"event_type"},
	},
	{
		Name:    "idx_webhooks_active",
		Table:   "webhooks",
		Columns: []string{"active", "user_id"},
	},
	{
		Name:    "idx_webhooks_created_at",
		Table:   "webhooks",
		Columns: []string{"created_at"},
	},

	// Jobs table indexes
	{
		Name:    "idx_jobs_user_id",
		Table:   "jobs",
		Columns: []string{"user_id"},
	},
	{
		Name:    "idx_jobs_status",
		Table:   "jobs",
		Columns: []string{"status"},
	},
	{
		Name:    "idx_jobs_type",
		Table:   "jobs",
		Columns: []string{"type"},
	},
	{
		Name:    "idx_jobs_scheduled_at",
		Table:   "jobs",
		Columns: []string{"scheduled_at"},
	},
	{
		Name:    "idx_jobs_priority",
		Table:   "jobs",
		Columns: []string{"priority", "scheduled_at"},
	},
	{
		Name:    "idx_jobs_created_at",
		Table:   "jobs",
		Columns: []string{"created_at"},
	},
	{
		Name:    "idx_jobs_updated_at",
		Table:   "jobs",
		Columns: []string{"updated_at"},
	},

	// Composite indexes for common queries
	{
		Name:    "idx_jobs_status_scheduled",
		Table:   "jobs",
		Columns: []string{"status", "scheduled_at"},
	},
	{
		Name:    "idx_jobs_type_status",
		Table:   "jobs",
		Columns: []string{"type", "status"},
	},
	{
		Name:    "idx_webhooks_user_active",
		Table:   "webhooks",
		Columns: []string{"user_id", "active"},
	},

	// Partial indexes
	{
		Name:    "idx_jobs_pending",
		Table:   "jobs",
		Columns: []string{"scheduled_at"},
		Partial: "status = 'pending'",
	},
	{
		Name:    "idx_webhooks_active_only",
		Table:   "webhooks",
		Columns: []string{"user_id"},
		Partial: "active = true",
	},
}

// GenerateCreateIndexSQL generates SQL for creating an index.
func (i *IndexDefinition) GenerateCreateIndexSQL() string {
	columns := ""
	for idx, col := range i.Columns {
		if idx > 0 {
			columns += ", "
		}
		columns += col
	}

	sql := "CREATE"
	if i.Unique {
		sql += " UNIQUE"
	}
	sql += " INDEX IF NOT EXISTS " + i.Name + " ON " + i.Table + " (" + columns + ")"

	if i.Partial != "" {
		sql += " WHERE " + i.Partial
	}

	return sql
}

// GenerateDropIndexSQL generates SQL for dropping an index.
func (i *IndexDefinition) GenerateDropIndexSQL() string {
	return "DROP INDEX IF EXISTS " + i.Name
}
