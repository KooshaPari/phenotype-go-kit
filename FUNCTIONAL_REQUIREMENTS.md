# Functional Requirements - phenotype-go-kit

## FR-LOG-001: Logger Injection
The system SHALL store a `*slog.Logger` in a `context.Context` via `WithLogger` and retrieve it via `From`.

## FR-LOG-002: Missing Logger Panic
`From` SHALL panic if no logger has been injected into the context.

## FR-BUF-001: Fixed Capacity
`ringbuffer.New[T](cap)` SHALL create a buffer that holds at most `cap` elements.

## FR-BUF-002: Overwrite on Full
`Push` SHALL overwrite the oldest element when the buffer is at capacity.

## FR-BUF-003: Oldest-First Retrieval
`GetAll` SHALL return elements in insertion order (oldest first).

## FR-WAIT-001: Exponential Backoff Polling
`WaitFor` SHALL poll a condition function with exponential backoff between `MinInterval` and `MaxInterval`.

## FR-WAIT-002: Timeout
`WaitFor` SHALL return `ErrTimedOut` if the condition is not met within `Timeout`.

## FR-WAIT-003: Testable Clock
`WaitFor` SHALL accept a `quartz.Clock` for deterministic testing.

## FR-REG-001: Owner-Scoped Registration
`Register` SHALL associate an owner ID with a key-value entry, supporting multiple owners per key.

## FR-REG-002: Ref-Counted Removal
`Unregister` SHALL remove the owner; the entry SHALL be deleted only when the last owner unregisters.

## FR-REG-003: Change Hooks
The registry SHALL support a `Hook` interface with `OnRegister` and `OnUnregister` callbacks.

## FR-REG-004: Thread Safety
All registry operations SHALL be safe for concurrent use via `sync.RWMutex`.
