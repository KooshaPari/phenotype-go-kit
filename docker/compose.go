package docker

import (
	"fmt"
	"strings"
)

// Image represents a Docker image.
type Image struct {
	Registry string
	Name     string
	Tag      string
}

// ParseImage parses an image string.
func ParseImage(s string) (*Image, error) {
	parts := strings.Split(s, "/")

	nameParts := strings.Split(parts[len(parts)-1], ":")

	registry := ""
	if len(parts) > 1 && !strings.Contains(parts[0], ":") {
		registry = parts[0]
	}

	name := nameParts[0]
	tag := "latest"
	if len(nameParts) > 1 {
		tag = nameParts[1]
	}

	return &Image{
		Registry: registry,
		Name:     name,
		Tag:      tag,
	}, nil
}

// String returns the full image string.
func (i *Image) String() string {
	var result string
	if i.Registry != "" {
		result = i.Registry + "/"
	}
	result += i.Name + ":" + i.Tag
	return result
}

// BuildConfig holds Docker build configuration.
type BuildConfig struct {
	Context    string
	Dockerfile string
	ImageName  string
	ImageTag   string
	BuildArgs  map[string]string
	Labels     map[string]string
	NoCache    bool
	Target     string
}

// GetBuildArgs converts build args to CLI args.
func (c *BuildConfig) GetBuildArgs() []string {
	var args []string

	for k, v := range c.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}

	return args
}

// DockerComposeConfig represents docker-compose configuration.
type DockerComposeConfig struct {
	Version  string
	Services map[string]ServiceConfig
}

// ServiceConfig holds service configuration.
type ServiceConfig struct {
	Image       string
	Build       string
	Ports       []string
	Environment map[string]string
	Volumes     []string
	DependsOn   []string
	Command     string
}

// NewCompose creates a new docker-compose config.
func NewCompose(version string) *DockerComposeConfig {
	return &DockerComposeConfig{
		Version:  version,
		Services: make(map[string]ServiceConfig),
	}
}

// AddService adds a service.
func (c *DockerComposeConfig) AddService(name string, config ServiceConfig) {
	c.Services[name] = config
}

// Generate generates docker-compose.yml content.
func (c *DockerComposeConfig) Generate() string {
	var lines []string
	lines = append(lines, "version: '"+c.Version+"'")
	lines = append(lines, "services:")

	for name, svc := range c.Services {
		lines = append(lines, "  "+name+":")

		if svc.Image != "" {
			lines = append(lines, "    image: "+svc.Image)
		}
		if svc.Build != "" {
			lines = append(lines, "    build: "+svc.Build)
		}

		if len(svc.Ports) > 0 {
			lines = append(lines, "    ports:")
			for _, p := range svc.Ports {
				lines = append(lines, "      - \""+p+"\"")
			}
		}

		if len(svc.Environment) > 0 {
			lines = append(lines, "    environment:")
			for k, v := range svc.Environment {
				lines = append(lines, "      - "+k+"="+v)
			}
		}

		if len(svc.Volumes) > 0 {
			lines = append(lines, "    volumes:")
			for _, v := range svc.Volumes {
				lines = append(lines, "      - "+v)
			}
		}
	}

	return strings.Join(lines, "\n")
}
