# ADR-004: Plugin System & Extensibility Architecture

**Status**: Accepted
**Date**: 2026-03-25
**Deciders**: Phenotype Architecture Team

---

## Context

Phenotype needs a plugin system to support:
- Dynamic feature loading without recompilation
- Third-party extensions
- Modular architecture deployment
- Hot reloading capabilities
- Clear boundaries between core and extensions

---

## Decision

### Plugin Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Plugin Host                          │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Plugin Registry                     │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐         │   │
│  │  │Plugin A │ │Plugin B │ │Plugin C │         │   │
│  │  │Manifest │ │Manifest │ │Manifest │         │   │
│  │  └────┬────┘ └────┬────┘ └────┬────┘         │   │
│  └───────┼────────────┼───────────┼─────────────────┘   │
│          │            │           │                      │
│          └────────────┼───────────┘                      │
│                       ▼                                   │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Plugin Manager                      │   │
│  │  • Load/Unload lifecycle                         │   │
│  │  • Dependency resolution                         │   │
│  │  • Health monitoring                             │   │
│  │  • Version compatibility                         │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Plugin Interface (Contract)         │   │
│  │  • Init()                                        │   │
│  │  • Name()                                        │   │
│  │  │  Version()                                    │   │
│  │  │  Execute(ctx)                                 │   │
│  │  │  Shutdown()                                   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Plugin Interface

```go
// contracts/plugins/plugin.go

// Plugin is the base interface all plugins must implement
type Plugin interface {
    // Metadata
    Manifest() (*Manifest, error)

    // Lifecycle
    Init(ctx context.Context) error
    Execute(ctx context.Context, req *Request) (*Response, error)
    Shutdown(ctx context.Context) error

    // Health
    Health(ctx context.Context) (*HealthStatus, error)
}

// Manifest contains plugin metadata
type Manifest struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Author      string            `json:"author"`
    License     string            `json:"license"`
    Tags        []string          `json:"tags"`
    Requires    map[string]string `json:"requires"` // service:version
    Provides    []string          `json:"provides"` // capabilities
}

// Plugin Registry
type Registry interface {
    Register(plugin Plugin) error
    Unregister(name string) error
    Get(name string) (Plugin, error)
    List() ([]*Manifest, error)
}
```

### Plugin Types

```
┌─────────────────────────────────────────────────────┐
│                   Plugin Taxonomy                   │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐│
│  │   Input     │  │   Output    │  │   Transform ││
│  │   Plugin    │  │   Plugin    │  │   Plugin    ││
│  ├─────────────┤  ├─────────────┤  ├─────────────┤│
│  │ Parse files │  │ Write JSON  │  │ Enrich data ││
│  │ Fetch API   │  │ Send email  │  │ Filter      ││
│  │ Read DB     │  │ Log to sys  │  │ Aggregate   ││
│  └─────────────┘  └─────────────┘  └─────────────┘│
│                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐│
│  │  Middleware │  │   Storage   │  │    AI/LLM   ││
│  │   Plugin    │  │   Plugin    │  │   Plugin    ││
│  ├─────────────┤  ├─────────────┤  ├─────────────┤│
│  │ Auth        │  │ PostgreSQL  │  │ OpenAI      ││
│  │ Rate limit  │  │ MongoDB     │  │ Anthropic   ││
│  │ Logging     │  │ Redis       │  │ Ollama      ││
│  └─────────────┘  └─────────────┘  └─────────────┘│
│                                                     │
└─────────────────────────────────────────────────────┘
```

### Loading Mechanism

```go
// Plugin loader with security and lifecycle management
type Loader struct {
    pluginDir  string
    registry   Registry
    logger     Logger
    maxMemory  int64
    allowedAPI map[string]bool
}

func (l *Loader) Load(ctx context.Context, path string) error {
    // 1. Verify plugin signature
    if err := l.verifySignature(path); err != nil {
        return fmt.Errorf("plugin signature invalid: %w", err)
    }

    // 2. Parse manifest
    manifest, err := l.parseManifest(path)
    if err != nil {
        return fmt.Errorf("invalid manifest: %w", err)
    }

    // 3. Check dependencies
    if err := l.checkDependencies(manifest); err != nil {
        return fmt.Errorf("dependency check failed: %w", err)
    }

    // 4. Load plugin binary
    plugin, err := plugin.Open(path)
    if err != nil {
        return fmt.Errorf("failed to load plugin: %w", err)
    }

    // 5. Initialize plugin
    if err := plugin.Init(ctx); err != nil {
        return fmt.Errorf("plugin init failed: %w", err)
    }

    // 6. Register
    return l.registry.Register(plugin)
}
```

### Plugin Directory Structure

```
plugins/
├── manifest.yaml              # Plugin registry manifest
│
├── builtin/
│   └── health/
│       ├── health.go          # Built-in health plugin
│       └── manifest.yaml
│
└── external/
    ├── embeddings-openrouter/
    │   ├── plugin.so          # Compiled plugin
    │   ├── manifest.yaml     # Plugin manifest
    │   └── README.md
    │
    └── custom-plugin/
        ├── plugin.so
        └── manifest.yaml
```

### Plugin Manifest Example

```yaml
# plugins/external/embeddings-openrouter/manifest.yaml
name: embeddings-openrouter
version: 1.0.0
description: OpenRouter embeddings provider
author: Phenotype Team
license: Apache-2.0

tags:
  - embeddings
  - ai
  - openrouter

requires:
  phenotype-core: ">=1.0.0"
  phenotype-http: ">=1.0.0"

provides:
  - embeddings-provider
  - ai-completion

config:
  api_endpoint:
    type: string
    required: true
  api_key:
    type: secret
    required: true
  model:
    type: string
    default: "openai/text-embedding-3-small"
```

---

## Consequences

### Positive
- Extensibility without core changes
- Third-party contributions
- Modular deployment
- Isolated failures
- Technology flexibility

### Negative
- Complexity overhead
- Security considerations (untrusted code)
- Debugging challenges
- Version compatibility

### Risks
- Plugin stability → Version constraints in manifest
- Security vulnerabilities → Sandboxing and signing
- Dependency hell → Registry and resolution logic

---

## Alternatives Considered

### 1. Compile-time Plugins (Build Tags)
- Simple
- Not extensible post-build
- Not chosen

### 2. Scripting Engine (Lua/JavaScript)
- Sandboxed
- Limited ecosystem
- Not chosen

### 3. External Service (Microservices)
- Maximum isolation
- Network overhead
- Overkill for plugin level

---

## Security Considerations

```
┌─────────────────────────────────────────────────────┐
│              Plugin Security Model                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│  1. Signing & Verification                         │
│     • All plugins must be signed                    │
│     • Registry validates signatures                 │
│                                                     │
│  2. Sandboxing                                      │
│     • Memory limits                                │
│     • CPU limits                                   │
│     • Network restrictions                         │
│     • File system access control                   │
│                                                     │
│  3. API Restrictions                               │
│     • Plugins only access declared capabilities    │
│     • Capability-based permissions                 │
│                                                     │
│  4. Audit Logging                                  │
│     • All plugin operations logged                 │
│     • Tamper-evident logs                          │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## References

- [Go Plugin System](pkg.go.dev/plugin)
- [HashiCorp Plugin System](github.com/hashicorp/go-plugin)
- [Kubernetes Plugin Architecture](kubernetes.io/docs/concepts/extend-kubernetes/)
