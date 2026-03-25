// Package adapter provides hexagonal architecture adapters for the cache package.
//
// Primary Adapters (Driving):
//   - Handler adapters for HTTP/gRPC endpoints
//
// Secondary Adapters (Driven):
//   - RedisCacheAdapter implements CachePort
//
// See ADR-001 for full architectural context.
package adapter
