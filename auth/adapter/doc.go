// Package adapter provides hexagonal architecture adapters for auth.
//
// Primary Adapters (Driving):
//   - HTTP middleware for authentication
//
// Secondary Adapters (Driven):
//   - JWTValidator implements outbound.AuthPort
//   - APIKeyManager implements outbound.APIKeyPort
//
// See ADR-001 for hexagonal architecture context.
package adapter
