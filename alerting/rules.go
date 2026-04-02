// Copyright (c) 2026 Phenotype Enterprise. All rights reserved.
// Licensed under the Phenotype Standard License.

package alerting

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Alert represents an alert definition.
type Alert struct {
	Name        string            `yaml:"name"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for"`
	Source      string            `yaml:"source"`
}

// AlertRuleSet holds a collection of alert rules.
type AlertRuleSet struct {
	Groups []AlertGroup `yaml:"groups"`
}

// AlertGroup groups alerts together.
type AlertGroup struct {
	Name  string  `yaml:"name"`
	Rules []Alert `yaml:"rules"`
}

// NewAlertRuleSet creates a new alert rule set.
func NewAlertRuleSet() *AlertRuleSet {
	return &AlertRuleSet{
		Groups: make([]AlertGroup, 0),
	}
}

// AddGroup adds an alert group.
func (a *AlertRuleSet) AddGroup(name string, rules []Alert) {
	a.Groups = append(a.Groups, AlertGroup{Name: name, Rules: rules})
}

// HighErrorRate creates an alert for high error rates.
func HighErrorRate(threshold float64, duration time.Duration) Alert {
	return Alert{
		Name: "HighErrorRate",
		Labels: map[string]string{
			"severity": "critical",
			"team":     "backend",
			"category": "reliability",
		},
		Annotations: map[string]string{
			"summary":     "High error rate detected",
			"description": "Error rate is above {{ $value }}% for the last 5 minutes",
		},
		Expr:   fmt.Sprintf("sum(rate(phenotype_http_requests_total{status=~'5..'}[5m])) / sum(rate(phenotype_http_requests_total[5m])) * 100 > %f", threshold),
		For:    duration.String(),
		Source: "metrics",
	}
}

// HighLatency creates an alert for high latency.
func HighLatency(threshold float64, duration time.Duration) Alert {
	return Alert{
		Name: "HighLatency",
		Labels: map[string]string{
			"severity": "warning",
			"team":     "backend",
			"category": "performance",
		},
		Annotations: map[string]string{
			"summary":     "High latency detected",
			"description": "P99 latency is above %dms for the last 5 minutes",
		},
		Expr:   fmt.Sprintf("histogram_quantile(0.99, sum(rate(phenotype_http_request_duration_seconds_bucket[5m])) by (le)) > %f", threshold),
		For:    duration.String(),
		Source: "metrics",
	}
}

// ResourceExhaustion creates an alert for system resource exhaustion.
func ResourceExhaustion(resourceType, threshold string) Alert {
	exprMap := map[string]string{
		"memory": "container_memory_usage_bytes / container_spec_memory_limit_bytes > " + threshold,
		"cpu":    "container_cpu_usage_seconds_total / container_spec_cpu_quota / container_spec_cpu_period > " + threshold,
		"disk":   "disk_usage_bytes / disk_capacity_bytes > " + threshold,
		"fd":     "process_open_fds / process_max_fds > " + threshold,
	}

	return Alert{
		Name: fmt.Sprintf("ResourceExhaustion_%s", resourceType),
		Labels: map[string]string{
			"severity": "critical",
			"team":     "platform",
			"category": "capacity",
		},
		Annotations: map[string]string{
			"summary":     fmt.Sprintf("System resource %s exhaustion", resourceType),
			"description": fmt.Sprintf("%s usage is above %s threshold", resourceType, threshold),
		},
		Expr:   exprMap[resourceType],
		For:    "2m",
		Source: "metrics",
	}
}

// JobQueueBacklog creates an alert for job queue backlog.
func JobQueueBacklog(threshold int, duration time.Duration) Alert {
	return Alert{
		Name: "JobQueueBacklog",
		Labels: map[string]string{
			"severity": "warning",
			"team":     "backend",
			"category": "reliability",
		},
		Annotations: map[string]string{
			"summary":     "Job queue backlog detected",
			"description": "Job queue depth exceeds %d for the last 10 minutes",
		},
		Expr:   fmt.Sprintf("phenotype_job_queue_depth > %d", threshold),
		For:    duration.String(),
		Source: "metrics",
	}
}

// DatabaseConnections creates an alert for database connection exhaustion.
func DatabaseConnections(threshold float64) Alert {
	return Alert{
		Name: "DatabaseConnectionExhaustion",
		Labels: map[string]string{
			"severity": "critical",
			"team":     "backend",
			"category": "database",
		},
		Annotations: map[string]string{
			"summary":     "Database connection exhaustion",
			"description": "Database connection pool is above %d%% utilization",
		},
		Expr:   fmt.Sprintf("db_pool_connections_active / db_pool_connections_max * 100 > %f", threshold),
		For:    "2m",
		Source: "metrics",
	}
}

// PagerDutyConfig holds PagerDuty integration configuration.
type PagerDutyConfig struct {
	APIKey           string            `yaml:"api_key"`
	ServiceID        string            `yaml:"service_id"`
	IntegrationKey   string            `yaml:"integration_key"`
	EscalationPolicy string            `yaml:"escalation_policy"`
	PriorityMapping  map[string]string `yaml:"priority_mapping"`
}

// OpsGenieConfig holds OpsGenie integration configuration.
type OpsGenieConfig struct {
	APIKey          string            `yaml:"api_key"`
	TeamName        string            `yaml:"team_name"`
	PriorityMapping map[string]string `yaml:"priority_mapping"`
}

// PrometheusAlertManagerConfig generates AlertManager config.
func PrometheusAlertManagerConfig(pdConfig PagerDutyConfig, ogConfig OpsGenieConfig) string {
	config := `global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default'
  routes:
    - match:
        severity: critical
      receiver: 'critical'
      continue: true
    - match:
        severity: warning
      receiver: 'warning'

receivers:
  - name: 'default'
    webhook_configs:
      - url: 'http://notification-service/alerts'
        send_resolved: true
  - name: 'critical'
    pagerduty_configs:
      - service_key: '` + pdConfig.IntegrationKey + `'
        severity: critical
        details:
          - key: 'team'
            value: 'backend'
  - name: 'warning'
    webhook_configs:
      - url: 'http://notification-service/alerts/warning'
        send_resolved: true
`

	return config
}

// AlertThresholds holds threshold configurations.
type AlertThresholds struct {
	ErrorRate          float64 `yaml:"error_rate"`
	P99LatencyMs       float64 `yaml:"p99_latency_ms"`
	P95LatencyMs       float64 `yaml:"p95_latency_ms"`
	MemoryUsagePercent float64 `yaml:"memory_usage_percent"`
	CPUUsagePercent    float64 `yaml:"cpu_usage_percent"`
	DiskUsagePercent   float64 `yaml:"disk_usage_percent"`
	QueueDepth         int     `yaml:"queue_depth"`
	DBPoolUsagePercent float64 `yaml:"db_pool_usage_percent"`
}

// DefaultThresholds returns default alert thresholds.
func DefaultThresholds() AlertThresholds {
	return AlertThresholds{
		ErrorRate:          5.0,
		P99LatencyMs:       1000.0,
		P95LatencyMs:       500.0,
		MemoryUsagePercent: 85.0,
		CPUUsagePercent:    80.0,
		DiskUsagePercent:   90.0,
		QueueDepth:         1000,
		DBPoolUsagePercent: 80.0,
	}
}

// ValidateThresholds validates alert thresholds.
func ValidateThresholds(t AlertThresholds) []string {
	var errors []string

	if t.ErrorRate <= 0 || t.ErrorRate > 100 {
		errors = append(errors, "error_rate must be between 0 and 100")
	}
	if t.P99LatencyMs <= 0 {
		errors = append(errors, "p99_latency_ms must be positive")
	}
	if t.MemoryUsagePercent <= 0 || t.MemoryUsagePercent > 100 {
		errors = append(errors, "memory_usage_percent must be between 0 and 100")
	}
	if t.CPUUsagePercent <= 0 || t.CPUUsagePercent > 100 {
		errors = append(errors, "cpu_usage_percent must be between 0 and 100")
	}

	return errors
}

// ParseDuration parses duration string like "5m".
func ParseDuration(s string) (time.Duration, error) {
	re := regexp.MustCompile(`^(\d+)([smhd])$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}

	var value int
	_, _ = fmt.Sscanf(matches[1], "%d", &value)

	switch matches[2] {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", matches[2])
	}
}

// SilenceRule represents an alert silence rule.
type SilenceRule struct {
	ID        string    `json:"id"`
	Matchers  []string  `json:"matchers"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	CreatedBy string    `json:"created_by"`
	Comment   string    `json:"comment"`
}

// Match checks if an alert matches the silence rule.
func (s SilenceRule) Match(labels map[string]string) bool {
	for _, matcher := range s.Matchers {
		parts := strings.Split(matcher, "=")
		if len(parts) != 2 {
			continue
		}
		key, expected := parts[0], parts[1]
		if actual, ok := labels[key]; !ok || actual != expected {
			return false
		}
	}
	return true
}
