package ci

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Pipeline represents a CI/CD pipeline.
type Pipeline struct {
	name   string
	stages []*Stage
	logger *slog.Logger
}

// Stage represents a pipeline stage.
type Stage struct {
	Name     string
	Commands []string
	Env      map[string]string
	Timeout  time.Duration
}

// New creates a new pipeline.
func New(name string) *Pipeline {
	return &Pipeline{
		name:   name,
		stages: make([]*Stage, 0),
		logger: slog.Default(),
	}
}

// AddStage adds a stage to the pipeline.
func (p *Pipeline) AddStage(stage *Stage) {
	p.stages = append(p.stages, stage)
}

// Run executes the pipeline.
func (p *Pipeline) Run(ctx context.Context) error {
	p.logger.Info("starting pipeline", "name", p.name)

	for i, stage := range p.stages {
		p.logger.Info("running stage", "stage", stage.Name, "progress", fmt.Sprintf("%d/%d", i+1, len(p.stages)))

		if err := p.runStage(ctx, stage); err != nil {
			p.logger.Error("stage failed", "stage", stage.Name, "error", err)
			return fmt.Errorf("stage %s failed: %w", stage.Name, err)
		}
	}

	p.logger.Info("pipeline completed", "name", p.name)
	return nil
}

func (p *Pipeline) runStage(ctx context.Context, stage *Stage) error {
	// Set environment variables
	for k, v := range stage.Env {
		_ = os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	for _, cmd := range stage.Commands {
		if err := p.runCommand(ctx, cmd, stage.Timeout); err != nil {
			return err
		}
	}

	return nil
}

func (p *Pipeline) runCommand(ctx context.Context, cmd string, timeout time.Duration) error {
	parts := strings.Split(cmd, " ")
	name := parts[0]
	args := parts[1:]

	command := exec.CommandContext(ctx, name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	return command.Run()
}

// BuildStage creates a build stage.
func BuildStage(commands ...string) *Stage {
	return &Stage{
		Name:     "build",
		Commands: commands,
		Timeout:  5 * time.Minute,
	}
}

// TestStage creates a test stage.
func TestStage(commands ...string) *Stage {
	return &Stage{
		Name:     "test",
		Commands: commands,
		Timeout:  10 * time.Minute,
	}
}

// DeployStage creates a deploy stage.
func DeployStage(commands ...string) *Stage {
	return &Stage{
		Name:     "deploy",
		Commands: commands,
		Timeout:  15 * time.Minute,
	}
}

// LintStage creates a lint stage.
func LintStage(commands ...string) *Stage {
	return &Stage{
		Name:     "lint",
		Commands: commands,
		Timeout:  5 * time.Minute,
	}
}

// SecurityStage creates a security scan stage.
func SecurityStage(commands ...string) *Stage {
	return &Stage{
		Name:     "security",
		Commands: commands,
		Timeout:  10 * time.Minute,
	}
}
