package plugins

// Manifest represents a plugin manifest file.
// Used for declarative plugin configuration.
type Manifest struct {
	// Version is the manifest version.
	Version string `json:"manifest_version"`

	// Plugins contains the list of plugins to load.
	Plugins []ManifestPlugin `json:"plugins"`
}

// ManifestPlugin represents a plugin in the manifest.
type ManifestPlugin struct {
	// ID is the unique identifier of the plugin.
	ID string `json:"id"`

	// Type is the type of the plugin.
	Type PluginType `json:"type"`

	// Source is the source of the plugin (file path, URL, or registry reference).
	Source string `json:"source"`

	// Version is the version constraint.
	Version string `json:"version,omitempty"`

	// Enabled indicates if the plugin is enabled.
	Enabled bool `json:"enabled,omitempty"`

	// Config contains the plugin configuration.
	Config map[string]any `json:"config,omitempty"`

	// Dependencies are the plugin dependencies.
	Dependencies []Dependency `json:"dependencies,omitempty"`

	// Priority is the plugin priority for loading order.
	Priority int `json:"priority,omitempty"`
}

// DefaultManifest returns a default manifest structure.
func DefaultManifest() *Manifest {
	return &Manifest{
		Version: "1.0.0",
		Plugins: []ManifestPlugin{},
	}
}

// Validate validates the manifest structure.
func (m *Manifest) Validate() error {
	if m.Version == "" {
		m.Version = "1.0.0"
	}

	for i, p := range m.Plugins {
		if p.ID == "" {
			return &ManifestError{
				Field:   "plugins",
				Index:   i,
				Message: "plugin ID is required",
			}
		}
		if p.Type == "" {
			return &ManifestError{
				Field:   "type",
				Index:   i,
				Message: "plugin type is required",
			}
		}
	}

	return nil
}

// ManifestError represents a manifest validation error.
type ManifestError struct {
	Field   string
	Index   int
	Message string
}

func (e *ManifestError) Error() string {
	if e.Index >= 0 {
		return "manifest: " + e.Field + "[" + string(rune(e.Index)) + "]: " + e.Message
	}
	return "manifest: " + e.Field + ": " + e.Message
}
