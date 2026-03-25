// Package adapter provides hexagonal architecture adapters for secrets management.
// Adapters implement the outbound ports defined in contracts/ports/outbound.
//
// Supported Providers:
//   - Vault: HashiCorp Vault integration
//   - AWS: AWS Secrets Manager integration
//   - Env: Environment variable based secrets (development)
package adapter
