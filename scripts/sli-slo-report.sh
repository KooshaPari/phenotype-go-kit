#!/bin/bash
# SLI/SLO Daily Reporting Script
# Usage: ./scripts/sli-slo-report.sh

set -euo pipefail

# Configuration
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
SLACK_WEBHOOK="${SLACK_WEBHOOK:-}"
EMAIL_TO="${EMAIL_TO:-}"

DATE=$(date -u +"%Y-%m-%d")
REPORT_FILE="/tmp/slo-report-${DATE}.json"

# SLO Definitions
declare -A SLO_TARGETS
SLO_TARGETS["availability"]=99.9
SLO_TARGETS["latency_p99"]=1000  # ms
SLO_TARGETS["latency_p95"]=500   # ms
SLO_TARGETS["error_rate"]=1.0    # percent
SLO_TARGETS["db_query_success"]=99.5

get_metric() {
    local query=$1
    local result
    result=$(curl -s -G --data-urlencode "query=${query}" "${PROMETHEUS_URL}/api/v1/query" | \
        jq -r '.data.result[0].value[1] // "0"')
    echo "$result"
}

calculate_availability() {
    local total_requests
    local error_requests
    
    total_requests=$(get_metric "sum(rate(phenotype_http_requests_total[24h]))")
    error_requests=$(get_metric "sum(rate(phenotype_http_requests_total{status=~'5..'}[24h]))")
    
    if [ "$total_requests" = "0" ] || [ "$total_requests" = "null" ]; then
        echo "100.0"
        return
    fi
    
    local availability
    availability=$(echo "scale=4; (1 - $error_requests / $total_requests) * 100" | bc)
    printf "%.2f" "$availability"
}

calculate_latency_p99() {
    local latency
    latency=$(get_metric "histogram_quantile(0.99, sum(rate(phenotype_http_request_duration_seconds_bucket[24h])) by (le))")
    if [ "$latency" = "null" ] || [ "$latency" = "" ]; then
        echo "0"
        return
    fi
    # Convert to milliseconds
    echo "$latency * 1000" | bc | cut -d. -f1
}

calculate_latency_p95() {
    local latency
    latency=$(get_metric "histogram_quantile(0.95, sum(rate(phenotype_http_request_duration_seconds_bucket[24h])) by (le))")
    if [ "$latency" = "null" ] || [ "$latency" = "" ]; then
        echo "0"
        return
    fi
    echo "$latency * 1000" | bc | cut -d. -f1
}

calculate_error_rate() {
    local total_requests
    local error_requests
    
    total_requests=$(get_metric "sum(rate(phenotype_http_requests_total[24h]))")
    error_requests=$(get_metric "sum(rate(phenotype_http_requests_total{status=~'5..'}[24h]))")
    
    if [ "$total_requests" = "0" ] || [ "$total_requests" = "null" ]; then
        echo "0.0"
        return
    fi
    
    local error_rate
    error_rate=$(echo "scale=4; ($error_requests / $total_requests) * 100" | bc)
    printf "%.2f" "$error_rate"
}

generate_report() {
    echo "Generating SLO Report for ${DATE}..."
    
    local availability latency_p99 latency_p95 error_rate
    
    availability=$(calculate_availability)
    latency_p99=$(calculate_latency_p99)
    latency_p95=$(calculate_latency_p95)
    error_rate=$(calculate_error_rate)
    
    cat > "${REPORT_FILE}" << EOF
{
  "date": "${DATE}",
  "slo_targets": {
    "availability": ${SLO_TARGETS[availability]},
    "latency_p99": ${SLO_TARGETS[latency_p99]},
    "latency_p95": ${SLO_TARGETS[latency_p95]},
    "error_rate": ${SLO_TARGETS[error_rate]},
    "db_query_success": ${SLO_TARGETS[db_query_success]}
  },
  "actual_values": {
    "availability": ${availability},
    "latency_p99_ms": ${latency_p99},
    "latency_p95_ms": ${latency_p95},
    "error_rate_percent": ${error_rate}
  },
  "status": {
    "availability": $([ "$(echo "$availability >= ${SLO_TARGETS[availability]}" | bc)" = "1" ] && echo "met" || echo "breached"),
    "latency_p99": $([ "$(echo "$latency_p99 <= ${SLO_TARGETS[latency_p99]}" | bc)" = "1" ] && echo "met" || echo "breached"),
    "latency_p95": $([ "$(echo "$latency_p95 <= ${SLO_TARGETS[latency_p95]}" | bc)" = "1" ] && echo "met" || echo "breached"),
    "error_rate": $([ "$(echo "$error_rate <= ${SLO_TARGETS[error_rate]}" | bc)" = "1" ] && echo "met" || echo "breached")
  },
  "generated_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
}
EOF
    
    echo "Report saved to ${REPORT_FILE}"
}

send_slack_notification() {
    if [ -z "${SLACK_WEBHOOK}" ]; then
        echo "No Slack webhook configured, skipping notification"
        return
    fi
    
    local status_emoji
    local availability_error
    availability_error=$(echo "${SLO_TARGETS[availability]} - $(calculate_availability)" | bc | cut -d. -f1)
    
    if [ "$availability_error" -gt 0 ]; then
        status_emoji="🔴"
    else
        status_emoji="🟢"
    fi
    
    curl -s -X POST "${SLACK_WEBHOOK}" -H "Content-Type: application/json" -d "{
        \"text\": \"${status_emoji} SLO Report - ${DATE}\",
        \"blocks\": [
            {
                \"type\": \"header",
                \"text\": {\"type\": \"plain_text\", \"text\": \"SLO Daily Report - ${DATE}\"}
            },
            {
                \"type\": \"section\",
                \"fields\": [
                    {\"type\": \"mrkdwn\", \"text\": \"*Availability:*\\n$(calculate_availability)% / ${SLO_TARGETS[availability]}%\"},
                    {\"type\": \"mrkdwn\", \"text\": \"*P99 Latency:*\\n$(calculate_latency_p99)ms / ${SLO_TARGETS[latency_p99]}ms\"}
                ]
            }
        ]
    }"
}

send_email() {
    if [ -z "${EMAIL_TO}" ]; then
        echo "No email configured, skipping"
        return
    fi
    
    # Email sending logic would go here
    echo "Email would be sent to: ${EMAIL_TO}"
}

main() {
    generate_report
    send_slack_notification
    send_email
    echo "SLO reporting complete"
}

main "$@"
