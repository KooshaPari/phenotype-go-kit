package inbound

import (
	"context"

	"github.com/KooshaPari/phenotype-go-kit/contracts/models"
)

// UseCase defines the interface for application use cases.
// Following SRP: Each use case has a single responsibility.
type UseCase interface {
	// Execute runs the use case with the given context and input.
	// Returns the output or error.
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}

// CommandHandler handles commands (CQRS pattern).
// Commands are operations that change state.
type CommandHandler[C any] interface {
	// Handle processes a command and returns a result.
	Handle(ctx context.Context, cmd C) (*models.CommandResult, error)
}

// QueryHandler handles queries (CQRS pattern).
// Queries are operations that read state without side effects.
type QueryHandler[Q any, R any] interface {
	// Handle processes a query and returns results.
	Handle(ctx context.Context, query Q) (R, error)
}

// EventHandler handles domain events.
type EventHandler[E any] interface {
	// Handle processes an event.
	Handle(ctx context.Context, event E) error
}

// InputPort is a generic interface for input ports.
// Used by driving adapters to interact with the application.
type InputPort[In any, Out any] interface {
	// Execute runs the port with input and returns output.
	Execute(ctx context.Context, in In) (Out, error)
}

// Validator defines interface for input validation.
// Following ISP: Separate interfaces for different validation concerns.
type Validator[In any] interface {
	// Validate checks if the input is valid.
	Validate(in In) error
}

// PreProcessor defines interface for preprocessing input.
// Applies cross-cutting concerns before use case execution.
type PreProcessor[In any] interface {
	// PreProcess runs preprocessing on input.
	PreProcess(ctx context.Context, in In) (In, error)
}

// PostProcessor defines interface for postprocessing output.
// Applies cross-cutting concerns after use case execution.
type PostProcessor[Out any] interface {
	// PostProcess runs postprocessing on output.
	PostProcess(ctx context.Context, out Out) (Out, error)
}

// Interceptor defines middleware interface for use cases.
// Combines PreProcessor and PostProcessor for full request/response interception.
type Interceptor[In any, Out any] interface {
	PreProcessor[In]
	PostProcessor[Out]
}
