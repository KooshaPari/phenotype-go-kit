# Metrics Taxonomy

This document defines the standardized metrics naming convention and structure for the Phenotype platform.

## Naming Convention

Format: `{namespace}_{subsystem}_{name}_{unit}`

- **namespace**: Application identifier (`phenotype`)
- **subsystem**: Functional area (`http`, `jobs`, `db`, `system`, `business`)
- **name**: Metric name (CamelCase to snake_case)
- **unit**: Unit of measurement (optional)

## HTTP Metrics

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `phenotype_http_requests_total` | Counter | Total HTTP requests | method, path, status |
| `phenotype_http_request_duration_seconds` | Histogram | Request duration | method, path |
| `phenotype_http_response_size_bytes` | Histogram | Response size | method, path |

### Usage
```promql
# Request rate by endpoint
sum(rate(phenotype_http_requests_total[5m])) by (path)

# P99 latency
histogram_quantile(0.99, sum(rate(phenotype_http_request_duration_seconds_bucket[5m])) by (le, path))

# Error rate
sum(rate(phenotype_http_requests_total{status=~"5.."}[5m])) / sum(rate(phenotype_http_requests_total[5m])) * 100
```

## Job Queue Metrics

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `phenotype_job_queue_depth` | Gauge | Current queue size | job_type |
| `phenotype_job_processing_duration_seconds` | Histogram | Job processing time | job_type, status |
| `phenotype_job_retries_total` | Counter | Total retries | job_type, attempt |

### Usage
```promql
# Queue backlog alert
phenotype_job_queue_depth > 1000

# Average processing time
rate(phenotype_job_processing_duration_seconds_sum[5m]) / rate(phenotype_job_processing_duration_seconds_count[5m])
```

## Database Metrics

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `phenotype_db_query_duration_seconds` | Histogram | Query execution time | query_type, table |
| `phenotype_db_query_errors_total` | Counter | Query errors | query_type, table, error_code |
| `phenotype_db_connections_active` | Gauge | Active connections | pool_name |
| `phenotype_db_connections_idle` | Gauge | Idle connections | pool_name |

### Usage
```promql
# Slow queries (>1s)
sum(rate(phenotype_db_query_duration_seconds_bucket{le="1"}[5m])) by (table) / ignoring(le) sum(rate(phenotype_db_query_duration_seconds_count[5m])) by (table)

# Connection pool usage
phenotype_db_connections_active / phenotype_db_connections_max * 100
```

## System Metrics

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `phenotype_system_cpu_usage_percent` | Gauge | CPU usage | instance |
| `phenotype_system_memory_usage_percent` | Gauge | Memory usage | instance |
| `phenotype_system_disk_usage_percent` | Gauge | Disk usage | instance |
| `phenotype_system_file_descriptors_used` | Gauge | Open file descriptors | instance |

## Business Metrics

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `phenotype_business_user_registrations_total` | Counter | New user registrations | source |
| `phenotype_business_api_calls_total` | Counter | API calls | endpoint, tier |
| `phenotype_business_revenue_total` | Counter | Revenue | currency, plan |

## Label Guidelines

- Use lowercase for label values
- Use `_` instead of `-` in label names
- Avoid high-cardinality labels (e.g., user IDs, session IDs)
- Include `instance`, `region`, `environment` for infrastructure metrics

## Alerting Thresholds

| Metric | Warning | Critical |
|--------|---------|----------|
| HTTP Error Rate | > 1% | > 5% |
| P99 Latency | > 500ms | > 1000ms |
| Job Queue Depth | > 500 | > 1000 |
| DB Connection Pool | > 70% | > 85% |
| CPU Usage | > 70% | > 90% |
| Memory Usage | > 80% | > 95% |
