// Package service provides application services for cache domain.
//
// This layer orchestrates CQRS commands and queries, applying domain logic
// and coordinating between ports (inbound handlers and outbound adapters).
//
// Application services follow GRASP Controller pattern, receiving input
// and delegating to domain services and repositories.
//
// See ADR-001 for hexagonal architecture context.
package service
