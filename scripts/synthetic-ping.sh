#!/bin/bash
# Synthetic Ping Testing Script
# Tests public endpoints periodically and reports availability

set -euo pipefail

# Configuration
ENDPOINTS="${ENDPOINTS:-http://localhost:8080/health,http://localhost:8080/ready}"
INTERVAL="${INTERVAL:-60}"  # seconds
SLACK_WEBHOOK="${SLACK_WEBHOOK:-}"
PAGERDUTY_KEY="${PAGERDUTY_KEY:-}"

LOG_FILE="/var/log/phenotype/synthetic-tests.log"

log() {
    echo "[$(date -u +"%Y-%m-%dT%H:%M:%SZ")] $*" >> "${LOG_FILE}"
}

check_endpoint() {
    local url=$1
    local timeout=10
    
    local http_code
    http_code=$(curl -s -o /dev/null -w "%{http_code}" --max-time "$timeout" "$url" 2>/dev/null || echo "000")
    
    echo "$http_code"
}

run_checks() {
    IFS=',' read -ra ENDPOINT_ARRAY <<< "$ENDPOINTS"
    
    local failed=0
    local total=0
    
    for endpoint in "${ENDPOINT_ARRAY[@]}"; do
        total=$((total + 1))
        local code
        code=$(check_endpoint "$endpoint")
        
        if [ "$code" = "200" ]; then
            log "OK: $endpoint (HTTP $code)"
        else
            log "FAIL: $endpoint (HTTP $code)"
            failed=$((failed + 1))
        fi
    done
    
    if [ $failed -gt 0 ]; then
        log "WARN: $failed/$total endpoints failed"
        send_alert "$failed" "$total"
    else
        log "OK: All $total endpoints healthy"
    fi
}

send_alert() {
    local failed=$1
    local total=$2
    
    if [ -n "${SLACK_WEBHOOK}" ]; then
        curl -s -X POST "${SLACK_WEBHOOK}" -H "Content-Type: application/json" -d "{
            \"text\": \"🚨 Synthetic Test Alert: $failed/$total endpoints failed\",
            \"attachments\": [{
                \"color\": \"danger\",
                \"fields\": [
                    {\"title\": \"Failed\", \"value\": \"$failed\", \"short\": true},
                    {\"title\": \"Total\", \"value\": \"$total\", \"short\": true}
                ]
            }]
        }" 2>/dev/null || true
    fi
    
    if [ -n "${PAGERDUTY_KEY}" ]; then
        curl -s -X POST "https://events.pagerduty.com/v2/enqueue" -H "Content-Type: application/json" -d "{
            \"routing_key\": \"${PAGERDUTY_KEY}\",
            \"event_action\": \"trigger\",
            \"payload\": {
                \"summary\": \"Synthetic Test Failure: $failed/$total endpoints down\",
                \"severity\": \"critical\",
                \"source\": \"synthetic-ping-test\"
            }
        }" 2>/dev/null || true
    fi
}

main() {
    log "Starting synthetic ping tests for endpoints: $ENDPOINTS"
    
    while true; do
        run_checks
        sleep "$INTERVAL"
    done
}

main "$@"
