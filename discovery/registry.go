package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Service represents a service instance.
type Service struct {
	ID       string
	Name     string
	Address  string
	Port     int
	Metadata map[string]string
	Healthy  bool
	LastSeen time.Time
}

// Registry manages service registration and discovery.
type Registry struct {
	services map[string]map[string]*Service // name -> id -> service
	mu       sync.RWMutex
	logger   *slog.Logger
}

// New creates a new service registry.
func New() *Registry {
	return &Registry{
		services: make(map[string]map[string]*Service),
		logger:   slog.Default(),
	}
}

// Register adds a service to the registry.
func (r *Registry) Register(service *Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.services[service.Name] == nil {
		r.services[service.Name] = make(map[string]*Service)
	}

	service.LastSeen = time.Now()
	r.services[service.Name][service.ID] = service

	r.logger.Info("service registered", "name", service.Name, "id", service.ID, "address", service.Address)
	return nil
}

// Deregister removes a service from the registry.
func (r *Registry) Deregister(name, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.services[name] == nil {
		return fmt.Errorf("service not found: %s", name)
	}

	delete(r.services[name], id)
	r.logger.Info("service deregistered", "name", name, "id", id)
	return nil
}

// Discover returns all healthy instances of a service.
func (r *Registry) Discover(name string) []*Service {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services, ok := r.services[name]
	if !ok {
		return nil
	}

	result := make([]*Service, 0)
	for _, s := range services {
		if s.Healthy {
			result = append(result, s)
		}
	}

	return result
}

// GetService returns a specific service instance.
func (r *Registry) GetService(name, id string) (*Service, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services, ok := r.services[name]
	if !ok {
		return nil, false
	}

	svc, ok := services[id]
	return svc, ok
}

// Heartbeat updates the last seen time for a service.
func (r *Registry) Heartbeat(name, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	services, ok := r.services[name]
	if !ok {
		return fmt.Errorf("service not found: %s", name)
	}

	svc, ok := services[id]
	if !ok {
		return fmt.Errorf("service instance not found: %s", id)
	}

	svc.LastSeen = time.Now()
	return nil
}

// SetHealthy marks a service as healthy or not.
func (r *Registry) SetHealthy(name, id string, healthy bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	services, ok := r.services[name]
	if !ok {
		return fmt.Errorf("service not found: %s", name)
	}

	svc, ok := services[id]
	if !ok {
		return fmt.Errorf("service instance not found: %s", id)
	}

	svc.Healthy = healthy
	return nil
}

// ListServices returns all registered services.
func (r *Registry) ListServices() map[string][]*Service {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string][]*Service)
	for name, services := range r.services {
		result[name] = make([]*Service, 0, len(services))
		for _, s := range services {
			result[name] = append(result[name], s)
		}
	}

	return result
}

// HealthChecker periodically checks service health.
type HealthChecker struct {
	registry *Registry
	client   *http.Client
	interval time.Duration
	timeout  time.Duration
}

// NewHealthChecker creates a new health checker.
func NewHealthChecker(registry *Registry, interval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		registry: registry,
		client:   &http.Client{Timeout: timeout},
		interval: interval,
		timeout:  timeout,
	}
}

// Start begins the health check loop.
func (hc *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hc.checkServices()
		}
	}
}

func (hc *HealthChecker) checkServices() {
	services := hc.registry.ListServices()

	for _, instances := range services {
		for _, svc := range instances {
			healthURL := fmt.Sprintf("http://%s:%d/health", svc.Address, svc.Port)

			ctx, cancel := context.WithTimeout(context.Background(), hc.timeout)
			defer cancel()

			req, _ := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
			resp, err := hc.client.Do(req)

			if err != nil || resp.StatusCode != http.StatusOK {
				hc.registry.SetHealthy(svc.Name, svc.ID, false)
				continue
			}
			resp.Body.Close()

			hc.registry.SetHealthy(svc.Name, svc.ID, true)
		}
	}
}

// LoadBalancer provides service instance selection.
type LoadBalancer interface {
	Next(serviceName string) (*Service, error)
}

// RoundRobinLB is a round-robin load balancer.
type RoundRobinLB struct {
	registry *Registry
	counters map[string]int
	mu       sync.Mutex
}

// NewRoundRobinLB creates a round-robin load balancer.
func NewRoundRobinLB(registry *Registry) *RoundRobinLB {
	return &RoundRobinLB{
		registry: registry,
		counters: make(map[string]int),
	}
}

// Next returns the next service instance.
func (lb *RoundRobinLB) Next(serviceName string) (*Service, error) {
	services := lb.registry.Discover(serviceName)
	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances for service: %s", serviceName)
	}

	lb.mu.Lock()
	count := lb.counters[serviceName]
	lb.counters[serviceName] = (count + 1) % len(services)
	lb.mu.Unlock()

	return services[count], nil
}
