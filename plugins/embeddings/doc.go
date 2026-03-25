// Package embeddings provides AI embeddings provider plugins.
//
// This package implements the plugin system for AI embeddings providers,
// following hexagonal architecture principles.
//
// # Supported Providers
//
//   - OpenAI (text-embedding-3-small, text-embedding-3-large)
//   - Anthropic (via OpenAI-compatible API)
//   - Ollama (local models)
//   - Azure OpenAI
//
// # Usage
//
//	provider := embeddings.NewOpenAIProvider(apiKey)
//	embeddings, err := provider.Embed(ctx, "Hello, world!")
//
// # Law of Demeter (LoD)
//
// This package follows the Law of Demeter:
//   - Only talk to immediate collaborators
//   - No "train-wreck" method calls (a.GetB().GetC().Do())
//   - Use dependency injection for dependencies
//
// See ADR-005 for architectural details.
package embeddings
