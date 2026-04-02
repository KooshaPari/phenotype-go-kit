package deploy

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// Config holds deployment configuration.
type Config struct {
	Environment string
	Namespace   string
	Replicas    int
	Image       string
}

// KubernetesDeployer handles Kubernetes deployments.
type KubernetesDeployer struct {
	config *Config
	logger *slog.Logger
}

// NewKubernetesDeployer creates a new K8s deployer.
func NewKubernetesDeployer(cfg Config) *KubernetesDeployer {
	return &KubernetesDeployer{
		config: &cfg,
		logger: slog.Default(),
	}
}

// Apply applies a Kubernetes manifest.
func (k *KubernetesDeployer) Apply(ctx context.Context, manifest string) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(manifest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Deploy deploys an application.
func (k *KubernetesDeployer) Deploy(ctx context.Context) error {
	manifest := k.generateManifest()
	return k.Apply(ctx, manifest)
}

func (k *KubernetesDeployer) generateManifest() string {
	return fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: phenotype-app
  namespace: %s
spec:
  replicas: %d
  selector:
    matchLabels:
      app: phenotype-app
  template:
    metadata:
      labels:
        app: phenotype-app
    spec:
      containers:
      - name: app
        image: %s
        ports:
        - containerPort: 8080
`, k.config.Namespace, k.config.Replicas, k.config.Image)
}

// Rollback rolls back a deployment.
func (k *KubernetesDeployer) Rollback(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "rollout", "undo", "deployment/phenotype-app", "-n", k.config.Namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Status returns deployment status.
func (k *KubernetesDeployer) Status(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "rollout", "status", "deployment/phenotype-app", "-n", k.config.Namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return "", cmd.Run()
}

// HelmDeployer handles Helm deployments.
type HelmDeployer struct {
	chartPath string
	namespace string
	values    map[string]string
	logger    *slog.Logger
}

// NewHelmDeployer creates a new Helm deployer.
func NewHelmDeployer(chartPath, namespace string) *HelmDeployer {
	return &HelmDeployer{
		chartPath: chartPath,
		namespace: namespace,
		values:    make(map[string]string),
		logger:    slog.Default(),
	}
}

// SetValue sets a Helm value.
func (h *HelmDeployer) SetValue(key, value string) {
	h.values[key] = value
}

// Install installs a Helm chart.
func (h *HelmDeployer) Install(ctx context.Context, releaseName string) error {
	args := []string{"install", releaseName, h.chartPath, "-n", h.namespace}

	for k, v := range h.values {
		args = append(args, "--set", fmt.Sprintf("%s=%s", k, v))
	}

	cmd := exec.CommandContext(ctx, "helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Upgrade upgrades a Helm release.
func (h *HelmDeployer) Upgrade(ctx context.Context, releaseName string) error {
	args := []string{"upgrade", releaseName, h.chartPath, "-n", h.namespace}

	for k, v := range h.values {
		args = append(args, "--set", fmt.Sprintf("%s=%s", k, v))
	}

	cmd := exec.CommandContext(ctx, "helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Rollback rolls back a Helm release.
func (h *HelmDeployer) Rollback(ctx context.Context, releaseName string) error {
	cmd := exec.CommandContext(ctx, "helm", "rollback", releaseName, "-n", h.namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ReleaseInfo holds release information.
type ReleaseInfo struct {
	Name      string
	Namespace string
	Revision  string
	Status    string
}

// List lists Helm releases.
func (h *HelmDeployer) List(ctx context.Context) ([]ReleaseInfo, error) {
	cmd := exec.CommandContext(ctx, "helm", "list", "-n", h.namespace, "-o", "json")
	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Simplified - would need JSON parsing
	var releases []ReleaseInfo
	return releases, nil
}
