// Package inbound contains driving (inbound) ports - interfaces for use cases
// that are called by driving adapters (REST handlers, gRPC services, CLI commands).
//
// Following GRASP patterns:
//   - Controller: Handles input and delegates to domain services
//   - Creator: Responsibility for creating domain objects
//   - Expert: Information Expert pattern for operations
package inbound
