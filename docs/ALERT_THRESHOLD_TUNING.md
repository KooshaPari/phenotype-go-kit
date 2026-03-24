# Alert Threshold Tuning Guide

This guide provides best practices for tuning alert thresholds to reduce alert fatigue while maintaining system reliability.

## Principles

1. **SLA-driven thresholds**: Base thresholds on customer-facing SLAs
2. **Graduated escalation**: Use warning/critical levels with different response actions
3. **Historical analysis**: Review P95/P99 metrics to set realistic thresholds
4. **Zero false positives**: Every alert should be actionable

## Initial Threshold Calculation

### Error Rate

| SLA | Calculation | Initial Threshold |
|-----|-------------|-------------------|
| 99.9% | 100% - 99.9% = 0.1% | Warning: 0.5%, Critical: 1% |
| 99% | 100% - 99% = 1% | Warning: 2%, Critical: 5% |

```promql
-- Error rate calculation
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) * 100
```

### Latency

| Percentile | Use Case | Initial Threshold |
|------------|----------|-------------------|
| P50 | Median user experience | 200ms |
| P95 |大多数用户体验 | 500ms |
| P99 | Critical operations | 1000ms |

```promql
-- P99 latency
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```

### Queue Depth

Calculate based on processing rate:
- If processing 100 jobs/minute, queue of 1000 = 10 minute backlog
- Warning at 5 min, Critical at 10 min backlog

```promql
phenotype_job_queue_depth > 500  -- Warning
phenotype_job_queue_depth > 1000 -- Critical
```

## Tuning Process

### Week 1: Baseline

1. Deploy initial thresholds
2. Document all alerts received
3. Categorize: actionable vs. noise

### Week 2-4: Adjustment

1. For each alert type, analyze:
   - Is it actionable?
   - Was the threshold too sensitive?
   - Is there a pattern (time of day, deployment)?

2. Adjust thresholds based on data:
   - Too many alerts → increase threshold
   - Missing incidents → decrease threshold

### Ongoing: Maintenance

- Review quarterly
- Adjust for traffic growth
- Account for new services

## Common Adjustments

| Alert Type | Common Issue | Fix |
|------------|--------------|-----|
| Error rate | Normal retries trigger alerts | Exclude retried requests |
| Latency | Cold starts cause spikes | Add warm-up period |
| Queue depth | Large batch jobs queue | Exclude batch jobs |
| Disk usage | Log rotation delays | Adjust rotation config |

## Alert Suppression Rules

Use suppression rules to reduce noise:

```yaml
# AlertManager route configuration
route:
  routes:
    # Suppress during maintenance window
    - match:
        severity: warning
      continue: true
      routes:
        - match:
            maintenance: "true"
          receiver: null  # Do not alert

    # Group related alerts
    - match:
        component: database
      group_by: [alertname, database]
      receiver: database-team
```

## Escalation Path

Every alert should have a clear path:

1. **Warning** → On-call dashboard
2. **Critical** → PagerDuty → On-call engineer
3. **SLO breach** → Incident channel

## Measurement

Track alert quality metrics:

- **Precision**: % of alerts that are actionable
- **Recall**: % of incidents detected
- **MTTR**: Mean time to resolution

```sql
-- Query to analyze alert quality
SELECT 
  alert_name,
  count(*) as total_alerts,
  count(case when action_taken = true then 1 end) as actionable,
  count(case when action_taken = false then 1 end) as noise
FROM alerts
WHERE created_at > now() - 30 days
GROUP BY alert_name
ORDER BY actionable DESC
```

## Template Alert Rules

```yaml
groups:
- name: phenotype-alerts
  rules:
    - alert: HighErrorRate
      expr: sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) * 100 > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Error rate above 1%"
        runbook_url: "https://wiki.phenotype.dev/runbooks/high-error-rate"

    - alert: HighLatencyP99
      expr: histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "P99 latency above 1s"
```
