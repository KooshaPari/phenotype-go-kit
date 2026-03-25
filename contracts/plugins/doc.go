// Package plugins defines the plugin system interfaces following
// hexagonal architecture principles for extensibility.
//
// # Plugin Architecture
//
// Plugins are external adapters that can be loaded dynamically to extend
// the application's functionality. They follow the Dependency Inversion
// Principle - the core application depends on plugin interfaces (ports),
// not on concrete implementations.
//
// # Plugin Lifecycle
//
//	┌──────────┐    Register     ┌────────────┐
//	│ Registry  │ ──────────────► │  Plugin     │
//	└──────────┘                 │  Manifest   │
//	                              └────────────┘
//	                                     │
//	                                     │ Load
//	                                     ▼
//	                              ┌────────────┐
//	                              │  Plugin    │
//	                              │  Instance  │
//	                              └────────────┘
//	                                     │
//	                                     │ Init
//	                                     ▼
//	                              ┌────────────┐
//	                              │  Plugin    │
//	                              │  Running   │
//	                              └────────────┘
//	                                     │
//	                                     │ Shutdown
//	                                     ▼
//	                              ┌────────────┐
//	                              │  Plugin    │
//	                              │  Stopped   │
//	                              └────────────┘
package plugins
