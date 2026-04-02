package dashboards

import (
	"encoding/json"
	"time"
)

// GrafanaDashboard represents a Grafana dashboard.
type GrafanaDashboard struct {
	ID         int                 `json:"id,omitempty"`
	UID        string              `json:"uid"`
	Title      string              `json:"title"`
	Tags       []string            `json:"tags"`
	Time       GrafanaTimeSettings `json:"time"`
	Timepicker GrafanaTimepicker   `json:"timepicker"`
	Panels     []GrafanaPanel      `json:"panels"`
	Schema     int                 `json:"schemaVersion"`
	Version    int                 `json:"version"`
	Refresh    string              `json:"refresh"`
}

// GrafanaTimeSettings holds time settings.
type GrafanaTimeSettings struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// GrafanaTimepicker holds timepicker settings.
type GrafanaTimepicker struct {
	RefreshIntervals []string `json:"refresh_intervals"`
	TimeOptions      []string `json:"time_options"`
}

// GrafanaPanel represents a dashboard panel.
type GrafanaPanel struct {
	ID          int                    `json:"id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	GridPos     GrafanaGridPos         `json:"gridPos"`
	Targets     []GrafanaTarget        `json:"targets"`
	Options     map[string]interface{} `json:"options,omitempty"`
	FieldConfig GrafanaFieldConfig     `json:"fieldConfig,omitempty"`
}

// GrafanaGridPos holds grid position.
type GrafanaGridPos struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// GrafanaTarget represents a data source target.
type GrafanaTarget struct {
	Expr         string `json:"expr"`
	LegendFormat string `json:"legendFormat"`
	RefID        string `json:"refId"`
	Format       string `json:"format"`
	Instant      bool   `json:"instant,omitempty"`
}

// GrafanaFieldConfig holds field configuration.
type GrafanaFieldConfig struct {
	Defaults GrafanaDefaults `json:"defaults"`
}

// GrafanaDefaults holds default field settings.
type GrafanaDefaults struct {
	Unit       string            `json:"unit"`
	Min        int               `json:"min,omitempty"`
	Max        int               `json:"max,omitempty"`
	Decimals   int               `json:"decimals,omitempty"`
	Thresholds GrafanaThresholds `json:"thresholds,omitempty"`
}

// GrafanaThresholds holds threshold configuration.
type GrafanaThresholds struct {
	Mode  string          `json:"mode"`
	Steps []ThresholdStep `json:"steps"`
}

// ThresholdStep represents a threshold step.
type ThresholdStep struct {
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

// OperationalDashboard creates the primary operational dashboard.
func OperationalDashboard() GrafanaDashboard {
	return GrafanaDashboard{
		UID:     "phenotype-operational",
		Title:   "Phenotype Operational Dashboard",
		Tags:    []string{"operational", "phenotype"},
		Schema:  30,
		Version: 1,
		Refresh: "30s",
		Time: GrafanaTimeSettings{
			From: "now-1h",
			To:   "now",
		},
		Timepicker: GrafanaTimepicker{
			RefreshIntervals: []string{"10s", "30s", "1m", "5m", "15m"},
			TimeOptions:      []string{"1h", "6h", "12h", "24h", "7d"},
		},
		Panels: []GrafanaPanel{
			// HTTP Request Rate
			{
				ID:      1,
				Type:    "graph",
				Title:   "HTTP Request Rate",
				GridPos: GrafanaGridPos{X: 0, Y: 0, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "sum(rate(phenotype_http_requests_total[1m])) by (method, path)",
						LegendFormat: "{{method}} {{path}}",
						RefID:        "A",
					},
				},
			},
			// Error Rate
			{
				ID:      2,
				Type:    "graph",
				Title:   "Error Rate (%)",
				GridPos: GrafanaGridPos{X: 12, Y: 0, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "sum(rate(phenotype_http_requests_total{status=~'5..'}[1m])) / sum(rate(phenotype_http_requests_total[1m])) * 100",
						LegendFormat: "Error Rate",
						RefID:        "A",
					},
				},
				FieldConfig: GrafanaFieldConfig{
					Defaults: GrafanaDefaults{
						Unit: "percent",
						Thresholds: GrafanaThresholds{
							Mode: "absolute",
							Steps: []ThresholdStep{
								{Value: 0, Color: "green"},
								{Value: 5, Color: "yellow"},
								{Value: 10, Color: "red"},
							},
						},
					},
				},
			},
			// P99 Latency
			{
				ID:      3,
				Type:    "graph",
				Title:   "P99 Latency",
				GridPos: GrafanaGridPos{X: 0, Y: 8, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "histogram_quantile(0.99, sum(rate(phenotype_http_request_duration_seconds_bucket[1m])) by (le, method, path))",
						LegendFormat: "P99 {{method}} {{path}}",
						RefID:        "A",
					},
				},
				FieldConfig: GrafanaFieldConfig{
					Defaults: GrafanaDefaults{
						Unit: "s",
						Thresholds: GrafanaThresholds{
							Mode: "absolute",
							Steps: []ThresholdStep{
								{Value: 0, Color: "green"},
								{Value: 1, Color: "yellow"},
								{Value: 5, Color: "red"},
							},
						},
					},
				},
			},
			// P95 Latency
			{
				ID:      4,
				Type:    "graph",
				Title:   "P95 Latency",
				GridPos: GrafanaGridPos{X: 12, Y: 8, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "histogram_quantile(0.95, sum(rate(phenotype_http_request_duration_seconds_bucket[1m])) by (le, method, path))",
						LegendFormat: "P95 {{method}} {{path}}",
						RefID:        "A",
					},
				},
				FieldConfig: GrafanaFieldConfig{
					Defaults: GrafanaDefaults{
						Unit: "s",
					},
				},
			},
			// CPU Usage
			{
				ID:      5,
				Type:    "graph",
				Title:   "CPU Usage",
				GridPos: GrafanaGridPos{X: 0, Y: 16, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "rate(process_cpu_seconds_total[1m]) * 100",
						LegendFormat: "CPU %",
						RefID:        "A",
					},
				},
				FieldConfig: GrafanaFieldConfig{
					Defaults: GrafanaDefaults{
						Unit: "percent",
						Max:  100,
					},
				},
			},
			// Memory Usage
			{
				ID:      6,
				Type:    "graph",
				Title:   "Memory Usage",
				GridPos: GrafanaGridPos{X: 12, Y: 16, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "process_resident_memory_bytes / 1024 / 1024",
						LegendFormat: "Memory MB",
						RefID:        "A",
					},
				},
				FieldConfig: GrafanaFieldConfig{
					Defaults: GrafanaDefaults{
						Unit: "MB",
					},
				},
			},
			// Active Connections
			{
				ID:      7,
				Type:    "stat",
				Title:   "Active Connections",
				GridPos: GrafanaGridPos{X: 0, Y: 24, W: 6, H: 4},
				Targets: []GrafanaTarget{
					{
						Expr:         "sum(http_connections_active)",
						LegendFormat: "Active",
						RefID:        "A",
					},
				},
			},
			// Job Queue Depth
			{
				ID:      8,
				Type:    "stat",
				Title:   "Job Queue Depth",
				GridPos: GrafanaGridPos{X: 6, Y: 24, W: 6, H: 4},
				Targets: []GrafanaTarget{
					{
						Expr:         "sum(phenotype_job_queue_depth)",
						LegendFormat: "Queue Depth",
						RefID:        "A",
					},
				},
			},
			// Database Pool
			{
				ID:      9,
				Type:    "gauge",
				Title:   "DB Connection Pool",
				GridPos: GrafanaGridPos{X: 12, Y: 24, W: 6, H: 4},
				Targets: []GrafanaTarget{
					{
						Expr:         "db_pool_connections_active / db_pool_connections_max * 100",
						LegendFormat: "Pool %",
						RefID:        "A",
					},
				},
				FieldConfig: GrafanaFieldConfig{
					Defaults: GrafanaDefaults{
						Unit: "percent",
						Max:  100,
						Thresholds: GrafanaThresholds{
							Mode: "absolute",
							Steps: []ThresholdStep{
								{Value: 70, Color: "green"},
								{Value: 85, Color: "yellow"},
								{Value: 95, Color: "red"},
							},
						},
					},
				},
			},
		},
	}
}

// DatabaseDashboard creates the database performance dashboard.
func DatabaseDashboard() GrafanaDashboard {
	return GrafanaDashboard{
		UID:     "phenotype-database",
		Title:   "Phenotype Database Performance",
		Tags:    []string{"database", "phenotype"},
		Schema:  30,
		Version: 1,
		Refresh: "30s",
		Time: GrafanaTimeSettings{
			From: "now-1h",
			To:   "now",
		},
		Panels: []GrafanaPanel{
			{
				ID:      1,
				Type:    "graph",
				Title:   "Query Duration (P95)",
				GridPos: GrafanaGridPos{X: 0, Y: 0, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "histogram_quantile(0.95, sum(rate(phenotype_db_query_duration_seconds_bucket[1m])) by (le, query_type))",
						LegendFormat: "P95 {{query_type}}",
						RefID:        "A",
					},
				},
			},
			{
				ID:      2,
				Type:    "graph",
				Title:   "Query Error Rate",
				GridPos: GrafanaGridPos{X: 12, Y: 0, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "sum(rate(phenotype_db_query_errors_total[1m])) by (query_type)",
						LegendFormat: "{{query_type}}",
						RefID:        "A",
					},
				},
			},
			{
				ID:      3,
				Type:    "graph",
				Title:   "Active Connections",
				GridPos: GrafanaGridPos{X: 0, Y: 8, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "db_connections_active",
						LegendFormat: "Active",
						RefID:        "A",
					},
					{
						Expr:         "db_connections_idle",
						LegendFormat: "Idle",
						RefID:        "B",
					},
					{
						Expr:         "db_connections_max",
						LegendFormat: "Max",
						RefID:        "C",
					},
				},
			},
			{
				ID:      4,
				Type:    "graph",
				Title:   "Table Index Usage",
				GridPos: GrafanaGridPos{X: 12, Y: 8, W: 12, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "db_table_index_hit_ratio",
						LegendFormat: "{{table}}",
						RefID:        "A",
					},
				},
				FieldConfig: GrafanaFieldConfig{
					Defaults: GrafanaDefaults{
						Unit: "percentunit",
						Min:  0,
						Max:  1,
						Thresholds: GrafanaThresholds{
							Mode: "absolute",
							Steps: []ThresholdStep{
								{Value: 0.9, Color: "red"},
								{Value: 0.95, Color: "yellow"},
								{Value: 1, Color: "green"},
							},
						},
					},
				},
			},
		},
	}
}

// HealthDashboard creates the application health dashboard.
func HealthDashboard() GrafanaDashboard {
	return GrafanaDashboard{
		UID:     "phenotype-health",
		Title:   "Phenotype Health Dashboard",
		Tags:    []string{"health", "phenotype"},
		Schema:  30,
		Version: 1,
		Refresh: "10s",
		Time: GrafanaTimeSettings{
			From: "now-5m",
			To:   "now",
		},
		Panels: []GrafanaPanel{
			{
				ID:      1,
				Type:    "stat",
				Title:   "Service Status",
				GridPos: GrafanaGridPos{X: 0, Y: 0, W: 4, H: 4},
				Targets: []GrafanaTarget{
					{
						Expr:         "up{job='phenotype-go-kit'}",
						LegendFormat: "Instance {{instance}}",
						RefID:        "A",
						Format:       "table",
						Instant:      true,
					},
				},
			},
			{
				ID:      2,
				Type:    "stat",
				Title:   "Liveness Checks",
				GridPos: GrafanaGridPos{X: 4, Y: 0, W: 4, H: 4},
				Targets: []GrafanaTarget{
					{
						Expr:  "probe_success{target='liveness'}",
						RefID: "A",
					},
				},
			},
			{
				ID:      3,
				Type:    "stat",
				Title:   "Readiness Checks",
				GridPos: GrafanaGridPos{X: 8, Y: 0, W: 4, H: 4},
				Targets: []GrafanaTarget{
					{
						Expr:  "probe_success{target='readiness'}",
						RefID: "A",
					},
				},
			},
			{
				ID:      4,
				Type:    "table",
				Title:   "Component Health",
				GridPos: GrafanaGridPos{X: 0, Y: 4, W: 24, H: 8},
				Targets: []GrafanaTarget{
					{
						Expr:         "health_check_status",
						LegendFormat: "{{component}}",
						RefID:        "A",
						Format:       "table",
					},
				},
			},
		},
	}
}

// ToJSON converts dashboard to JSON bytes.
func (d GrafanaDashboard) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func init() {
	_ = time.Second // ensure time is imported
}
