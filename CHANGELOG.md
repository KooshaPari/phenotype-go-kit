# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
