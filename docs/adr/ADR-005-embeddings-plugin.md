# ADR-005: AI Embeddings Plugin System

**Status**: Accepted
**Date**: 2026-03-25
**Deciders**: Phenotype Architecture Team

---

## Context

The Phenotype ecosystem needs a pluggable AI embeddings system that supports:
- Multiple embeddings providers (OpenAI, Anthropic, Ollama, Azure)
- Runtime provider selection
- Consistent API across providers
- Testability with mock providers
- Local and cloud deployment options

---

## Decision

### Plugin Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Embeddings Plugin System                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    Registry                             │   │
│   │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │   │
│   │  │ OpenAI  │  │ Anthrop │  │ Ollama  │  │  Azure  │  │   │
│   │  │Provider │  │ Provider│  │Provider │  │Provider │  │   │
│   │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  │   │
│   └───────┼─────────────┼─────────────┼─────────────┼───────┘   │
│           └─────────────┴─────────────┴─────────────┘            │
│                            │                                      │
│                   Provider Interface                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Provider Interface

```go
type Provider interface {
    Embed(ctx context.Context, texts []string) (*EmbedResult, error)
    EmbedSingle(ctx context.Context, text string) (*Embedding, error)
    Name() string
    Model() string
    Dimensions() int
}
```

### Law of Demeter (LoD) Compliance

This plugin system follows the Law of Demeter (Principle of Least Knowledge):

```
┌────────────────────────────────────────────────────────────────┐
│                Law of Demeter (LoD) Rules                      │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ✓ DO: Accept dependencies via constructor                     │
│  ✓ DO: Return value objects (not internal state)               │
│  ✓ DO: Single responsibility per method                        │
│  ✓ DO: Use dependency injection                                │
│                                                                 │
│  ✗ DON'T: Train-wreck calls (a.GetB().GetC().Do())            │
│  ✗ DON'T: Global state or singletons                           │
│  ✗ DON'T: Service locator anti-pattern                         │
│  ✗ DON'T: Expose internal components via getters               │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

### Directory Structure

```
plugins/embeddings/
├── doc.go           # Package documentation
├── provider.go      # Provider interface, Config, Options
├── registry.go      # Provider registry
├── openai.go        # OpenAI implementation
├── anthropic.go     # Anthropic implementation (future)
├── ollama.go        # Ollama implementation
└── azure.go         # Azure OpenAI implementation (future)
```

---

## Consequences

### Positive
- Easy to add new providers
- Consistent API across providers
- Testable with mock providers
- Supports local (Ollama) and cloud (OpenAI) deployment

### Negative
- Interface may need changes for provider-specific features
- Additional abstraction overhead

### Risks
- API compatibility across versions → Semantic versioning

---

## Alternatives Considered

### 1. Direct API Calls (No Plugin)
- Simple
- Not extensible
- Not chosen

### 2. External Service
- Maximum isolation
- Network overhead
- Overkill for embeddings

---

## References

- [Law of Demeter - Wikipedia](https://en.wikipedia.org/wiki/Law_of_Demeter)
- [OpenAI Embeddings API](https://platform.openai.com/docs/guides/embeddings)
- [Ollama API](https://github.com/ollama/ollama/blob/main/docs/api.md)
