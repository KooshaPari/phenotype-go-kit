package plugins

import (
	"context"
)

// PluginType defines the type/category of a plugin.
type PluginType string

const (
	// PluginTypeBus represents a message bus plugin.
	PluginTypeBus PluginType = "bus"
	// PluginTypeCache represents a cache plugin.
	PluginTypeCache PluginType = "cache"
	// PluginTypeDatabase represents a database plugin.
	PluginTypeDatabase PluginType = "database"
	// PluginTypeAuth represents an authentication plugin.
	PluginTypeAuth PluginType = "auth"
	// PluginTypeStorage represents a storage plugin.
	PluginTypeStorage PluginType = "storage"
	// PluginTypeLogger represents a logging plugin.
	PluginTypeLogger PluginType = "logger"
	// PluginTypeMetrics represents a metrics plugin.
	PluginTypeMetrics PluginType = "metrics"
	// PluginTypeTracer represents a tracing plugin.
	PluginTypeTracer PluginType = "tracer"
	// PluginTypeCustom represents a custom plugin type.
	PluginTypeCustom PluginType = "custom"
)

// PluginState represents the current state of a plugin.
type PluginState string

const (
	// PluginStateRegistered indicates the plugin is registered but not loaded.
	PluginStateRegistered PluginState = "registered"
	// PluginStateLoaded indicates the plugin is loaded.
	PluginStateLoaded PluginState = "loaded"
	// PluginStateInitialized indicates the plugin is initialized and ready.
	PluginStateInitialized PluginState = "initialized"
	// PluginStateRunning indicates the plugin is running.
	PluginStateRunning PluginState = "running"
	// PluginStateStopping indicates the plugin is stopping.
	PluginStateStopping PluginState = "stopping"
	// PluginStateStopped indicates the plugin has stopped.
	PluginStateStopped PluginState = "stopped"
	// PluginStateError indicates the plugin encountered an error.
	PluginStateError PluginState = "error"
)

// Plugin is the base interface that all plugins must implement.
// Following SRP: Single responsibility for plugin lifecycle.
type Plugin interface {
	// Metadata returns the plugin metadata.
	Metadata() *Metadata

	// Init initializes the plugin with the given configuration.
	Init(ctx context.Context, config map[string]any) error

	// Start starts the plugin.
	Start(ctx context.Context) error

	// Stop stops the plugin gracefully.
	Stop(ctx context.Context) error

	// Health returns the plugin health status.
	Health(ctx context.Context) (*HealthStatus, error)
}

// Metadata contains plugin metadata information.
type Metadata struct {
	// ID is the unique identifier of the plugin.
	ID string `json:"id"`

	// Name is the human-readable name of the plugin.
	Name string `json:"name"`

	// Version is the version of the plugin (semver).
	Version string `json:"version"`

	// Description describes what the plugin does.
	Description string `json:"description"`

	// Type is the type of the plugin.
	Type PluginType `json:"type"`

	// Author is the plugin author.
	Author string `json:"author,omitempty"`

	// License is the plugin license.
	License string `json:"license,omitempty"`

	// Homepage is the plugin homepage URL.
	Homepage string `json:"homepage,omitempty"`

	// Repository is the plugin repository URL.
	Repository string `json:"repository,omitempty"`

	// Dependencies are the plugin dependencies.
	Dependencies []Dependency `json:"dependencies,omitempty"`

	// Tags are the plugin tags for categorization.
	Tags []string `json:"tags,omitempty"`
}

// Dependency represents a plugin dependency.
type Dependency struct {
	// ID is the dependency plugin ID.
	ID string `json:"id"`

	// Version specifies the version constraint.
	Version string `json:"version"`

	// Optional indicates if the dependency is optional.
	Optional bool `json:"optional,omitempty"`
}

// HealthStatus represents the health status of a plugin.
type HealthStatus struct {
	// State is the current state of the plugin.
	State PluginState `json:"state"`

	// Healthy indicates if the plugin is healthy.
	Healthy bool `json:"healthy"`

	// Message contains a status message.
	Message string `json:"message,omitempty"`

	// Error contains error information if unhealthy.
	Error string `json:"error,omitempty"`

	// LastCheck is when the health was last checked.
	LastCheck string `json:"last_check,omitempty"`
}

// PluginFactory creates plugin instances.
// Following Factory pattern for plugin instantiation.
type PluginFactory interface {
	// Create creates a new plugin instance.
	Create(ctx context.Context, config map[string]any) (Plugin, error)
}

// PluginRegistry manages plugin registration and discovery.
// Following Registry pattern for plugin lifecycle management.
type PluginRegistry interface {
	// Register registers a plugin factory.
	Register(pluginType PluginType, factory PluginFactory) error

	// Unregister unregisters a plugin factory.
	Unregister(pluginType PluginType) error

	// GetFactory returns the factory for a plugin type.
	GetFactory(pluginType PluginType) (PluginFactory, error)

	// List returns all registered plugin types.
	List() []PluginType

	// Create creates a new plugin instance by type.
	Create(ctx context.Context, pluginType PluginType, config map[string]any) (Plugin, error)
}

// PluginLoader loads plugins from various sources.
type PluginLoader interface {
	// Load loads plugins from the given source.
	Load(ctx context.Context, source string) ([]Plugin, error)

	// LoadFromFile loads plugins from a plugin manifest file.
	LoadFromFile(ctx context.Context, path string) ([]Plugin, error)

	// LoadFromDir loads plugins from a directory.
	LoadFromDir(ctx context.Context, dir string) ([]Plugin, error)

	// LoadFromRegistry loads plugins from a registry.
	LoadFromRegistry(ctx context.Context, registryURL string, names []string) ([]Plugin, error)
}

// PluginValidator validates plugin configurations and dependencies.
type PluginValidator interface {
	// Validate validates a plugin configuration.
	Validate(config map[string]any) error

	// ValidateDependencies checks if dependencies are satisfied.
	ValidateDependencies(deps []Dependency, available map[string]bool) error
}
