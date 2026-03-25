// Package config provides configuration management utilities.
package config

import (
	"github.com/spf13/viper"
)

// ConfigLoader represents a configuration loader using viper.
type ConfigLoader struct {
	viper *viper.Viper
}

// NewConfigLoader creates a new ConfigLoader instance.
func NewConfigLoader(filePath string) *ConfigLoader {
	v := viper.New()
	if filePath != "" {
		v.SetConfigFile(filePath)
	}
	return &ConfigLoader{
		viper: v,
	}
}

// Load loads configuration from the configured file path.
func (cl *ConfigLoader) Load() error {
	if err := cl.viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// LoadDefaults sets default values for configuration keys.
func (cl *ConfigLoader) LoadDefaults(defaults map[string]any) {
	for key, value := range defaults {
		cl.viper.SetDefault(key, value)
	}
}

// Unmarshal unmarshals the configuration into a struct.
func (cl *ConfigLoader) Unmarshal(rawVal interface{}) error {
	return cl.viper.Unmarshal(rawVal)
}

// Get retrieves a value from the configuration.
func (cl *ConfigLoader) Get(key string) interface{} {
	return cl.viper.Get(key)
}

// GetString retrieves a string value from the configuration.
func (cl *ConfigLoader) GetString(key string) string {
	return cl.viper.GetString(key)
}

// GetInt retrieves an int value from the configuration.
func (cl *ConfigLoader) GetInt(key string) int {
	return cl.viper.GetInt(key)
}

// GetBool retrieves a bool value from the configuration.
func (cl *ConfigLoader) GetBool(key string) bool {
	return cl.viper.GetBool(key)
}

// Config represents the application configuration.
type Config struct {
	viper *viper.Viper
}

// New creates a new Config instance.
func New() *Config {
	return &Config{
		viper: viper.New(),
	}
}

// Load loads configuration from a file.
func (c *Config) Load(configPath string) error {
	c.viper.SetConfigFile(configPath)
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// Get retrieves a value from the configuration.
func (c *Config) Get(key string) interface{} {
	return c.viper.Get(key)
}

// GetString retrieves a string value from the configuration.
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}
