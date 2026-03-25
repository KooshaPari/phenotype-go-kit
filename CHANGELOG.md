# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Documentation & Delivery** (WBS #271-300):
  - **API Documentation** (#271-273): OpenAPI 3.0 and Markdown generation
  - **Architecture Docs** (#274-276): System architecture overview and components
  - **User Guides** (#277-279): Usage patterns and examples
  - **Deployment Guides** (#280-282): Docker, Kubernetes, Helm deployment
  - **Contributing Guidelines** (#283-285): Development workflow and standards
  - **Release Notes** (#286-288): Release notes generation

- **CI/CD Pipeline** (WBS #241-270):
  - **Pipeline Definition** (#241-243): Pipeline stages, build/test/deploy/lint/security
  - **Docker Configuration** (#244-246): Docker image parsing, docker-compose config
  - **Kubernetes Deployment** (#247-249): K8s deployer, Helm deployer, rollout/rollback
  - **Secrets Management** (#250-252): Vault/AWS secrets integration, env generation

- **Frontend Infrastructure** (WBS #211-240):
  - **HTTP Client** (#211-213): API client with auth, retry, error handling
  - **State Management** (#214-216): State containers, store with actions, reducers
  - **Form Validation** (#217-219): Form fields, validation rules, required/email/min/max/pattern
  - **UI Utilities** (#220-222): HTML escaping, date formatting, string helpers, slugify

- **Data Layer** (WBS #181-210):
  - **Data Validation** (#181-183): Validator with required, email, min/max, pattern rules
  - **Data Transformation** (#184-186): Struct to map conversion, flattening, merging
  - **Repository Pattern** (#187-189): SQL repository with CRUD, filtering, pagination
  - **Data Transformation Utilities** (#190-192): JSON conversion, field mapping

- **Microservices Communication** (WBS #151-180):
  - **Event Bus** (#151-153): Local event bus with pub/sub, async publishing
  - **Service Discovery** (#154-156): Registry with registration, health checks, load balancing
  - **Circuit Breaker** (#157-159): Circuit breaker pattern with configurable thresholds
  - **Retry with Backoff** (#160-162): Exponential backoff with jitter, permanent error handling

- **API Gateway & Auth** (WBS #121-150):
  - **JWT Authentication** (#121-123): JWT generation/validation, token refresh, middleware
  - **API Key Management** (#124): API key generation, validation, revocation
  - **Rate Limiting** (#125-127): Token bucket rate limiter, middleware, distributed support
  - **CORS Configuration** (#128): Allowed origins, methods, headers, credentials
  - **API Versioning** (#129): Header/query parameter versioning, Accept header parsing
  - **OAuth2 Provider** (#130): Google/GitHub OAuth2 integration, PKCE support, session management

- **Database Scaling Infrastructure** (WBS #91-120):
  - **Connection Pool Configuration** (#96): Configurable pool with max open/idle connections, lifetime settings
  - **Database Indexes** (#91-94): Composite and partial indexes for users, webhooks, jobs tables
  - **Redis Caching Layer** (#101): Get/Set/Delete/JSON operations, sorted sets, TTL management
  - **Cache Invalidation** (#102): Tag-based invalidation strategies, pattern-based deletion
  - **Cache Warming** (#103): Preload frequently accessed data, fill missing cache entries
  - **Query Optimization** (#97-100): Query builder, slow query logging, pagination, explain plans

- **Observability Infrastructure** (WBS #61-90):
  - **Structured Logging** (#61): JSON log schema with timestamp, level, message, attributes, trace/span IDs
  - **Global Logging Interceptor** (#62): HTTP middleware with automatic trace ID propagation
  - **Log Rotation** (#63): Log rotation and retention policies with configurable max size, age, backups
  - **Centralized Aggregation Configs** (#64): Datadog, Logstash, Filebeat configurations for log aggregation
  - **OpenTelemetry Distributed Tracing** (#65-68): Tracing with OTLP export, HTTP instrumentation, baggage propagation
  - **Prometheus Metrics Collection** (#69-72): HTTP, job queue, database, business metrics with Prometheus client
  - **Grafana Dashboards** (#73-75): Operational, database, health dashboard definitions
  - **Alerting Rules** (#76-79): Prometheus alerting rules with PagerDuty/OpsGenie integration
  - **Health Check Endpoints** (#80-82): Liveness/readiness endpoints with component health checking
  - **SLI/SLO Reporting** (#83): Daily SLO reporting script for availability, latency, error rate
  - **Synthetic Ping Testing** (#85): Endpoint availability testing with Slack/PagerDuty alerts
  - **Sentry Frontend Integration** (#86-87): Browser error tracking, source map upload, issue routing
  - **Metrics Taxonomy Documentation** (#88): Standardized metrics naming convention
  - **Chaos Engineering Setup** (#89): LitmusChaos experiments for pod failure, network latency, CPU/memory stress
  - **Alert Threshold Tuning Guide** (#90): SLA-driven threshold calculation methodology

- **Background Job Queue** (WBS #49): In-memory job queue with worker pool, concurrent job processing, retry logic with exponential backoff
- **Email Notification Job** (WBS #50): Email job handler with SMTP integration stub (ready for SendGrid, AWS SES, Mailgun)
- **SMS Notification Job** (WBS #51): SMS job handler with Twilio and AWS SNS provider implementations
- **Webhook Delivery System** (WBS #52): HTTP webhook delivery with signature generation (HMAC-SHA256), retry logic with exponential backoff
- **Webhook Signature Verification** (WBS #53): Signature validation utilities for incoming webhooks
- **OpenAPI Documentation** (WBS #54): Full OpenAPI 3.0.3 specification with auth, users, webhooks, jobs endpoints
- **Database Migration Runner** (WBS #56): Migration framework with up/down migrations, version tracking
- **Initial Schema Migrations** (WBS #57): Users, webhooks, and jobs tables with indexes
- **Database Seed Scripts** (WBS #58): Development seed data for users, webhooks
- **Object Storage Connection** (WBS #59): S3 and GCS storage implementations with unified interface
- **File Upload/Download** (WBS #60): File service with multipart form support, signed URLs
- **Microservices Communication** (WBS #151-180):
  - **Event Bus** (#151-153): Local event bus with pub/sub, async publishing
  - **Service Discovery** (#154-156): Registry with registration, health checks, load balancing
  - **Circuit Breaker** (#157-159): Circuit breaker pattern with configurable thresholds
  - **Retry with Backoff** (#160-162): Exponential backoff with jitter, permanent error handling

- **API Gateway & Auth** (WBS #121-150):
  - **JWT Authentication** (#121-123): JWT generation/validation, token refresh, middleware
  - **API Key Management** (#124): API key generation, validation, revocation
  - **Rate Limiting** (#125-127): Token bucket rate limiter, middleware, distributed support
  - **CORS Configuration** (#128): Allowed origins, methods, headers, credentials
  - **API Versioning** (#129): Header/query parameter versioning, Accept header parsing
  - **OAuth2 Provider** (#130): Google/GitHub OAuth2 integration, PKCE support, session management

- **Database Scaling Infrastructure** (WBS #91-120):
  - **Connection Pool Configuration** (#96): Configurable pool with max open/idle connections, lifetime settings
  - **Database Indexes** (#91-94): Composite and partial indexes for users, webhooks, jobs tables
  - **Redis Caching Layer** (#101): Get/Set/Delete/JSON operations, sorted sets, TTL management
  - **Cache Invalidation** (#102): Tag-based invalidation strategies, pattern-based deletion
  - **Cache Warming** (#103): Preload frequently accessed data, fill missing cache entries
  - **Query Optimization** (#97-100): Query builder, slow query logging, pagination, explain plans

- **Observability Infrastructure** (WBS #61-90):
  - **Structured Logging** (#61): JSON log schema with timestamp, level, message, attributes, trace/span IDs
  - **Global Logging Interceptor** (#62): HTTP middleware with automatic trace ID propagation
  - **Log Rotation** (#63): Log rotation and retention policies with configurable max size, age, backups
  - **Centralized Aggregation Configs** (#64): Datadog, Logstash, Filebeat configurations for log aggregation
  - **OpenTelemetry Distributed Tracing** (#65-68): Tracing with OTLP export, HTTP instrumentation, baggage propagation
  - **Prometheus Metrics Collection** (#69-72): HTTP, job queue, database, business metrics with Prometheus client
  - **Grafana Dashboards** (#73-75): Operational, database, health dashboard definitions
  - **Alerting Rules** (#76-79): Prometheus alerting rules with PagerDuty/OpsGenie integration
  - **Health Check Endpoints** (#80-82): Liveness/readiness endpoints with component health checking
  - **SLI/SLO Reporting** (#83): Daily SLO reporting script for availability, latency, error rate
  - **Synthetic Ping Testing** (#85): Endpoint availability testing with Slack/PagerDuty alerts
  - **Sentry Frontend Integration** (#86-87): Browser error tracking, source map upload, issue routing
  - **Metrics Taxonomy Documentation** (#88): Standardized metrics naming convention
  - **Chaos Engineering Setup** (#89): LitmusChaos experiments for pod failure, network latency, CPU/memory stress
  - **Alert Threshold Tuning Guide** (#90): SLA-driven threshold calculation methodology

- **Background Job Queue** (WBS #49): In-memory job queue with worker pool, concurrent job processing, retry logic with exponential backoff
- **Email Notification Job** (WBS #50): Email job handler with SMTP integration stub (ready for SendGrid, AWS SES, Mailgun)
- **SMS Notification Job** (WBS #51): SMS job handler with Twilio and AWS SNS provider implementations
- **Webhook Delivery System** (WBS #52): HTTP webhook delivery with signature generation (HMAC-SHA256), retry logic with exponential backoff
- **Webhook Signature Verification** (WBS #53): Signature validation utilities for incoming webhooks
- **OpenAPI Documentation** (WBS #54): Full OpenAPI 3.0.3 specification with auth, users, webhooks, jobs endpoints
- **Database Migration Runner** (WBS #56): Migration framework with up/down migrations, version tracking
- **Initial Schema Migrations** (WBS #57): Users, webhooks, and jobs tables with indexes
- **Database Seed Scripts** (WBS #58): Development seed data for users, webhooks
- **Object Storage Connection** (WBS #59): S3 and GCS storage implementations with unified interface
- **File Upload/Download** (WBS #60): File service with multipart form support, signed URLs

### Testing
- Unit tests for job queue (enqueue, retry, job status tracking)

### Dependencies
- Added AWS SDK v2 for S3 integration
- Added Google Cloud Storage client

## [0.1.0] - 2026-03-24

### Initial Release
- phenotye-go-kit: Core Go utility packages
- logctx: Structured logging context utilities
- registry: Service registry pattern
- ringbuffer: Circular buffer implementation
- waitfor: Wait-for-condition utilities
- Governance documentation and worktree policies
