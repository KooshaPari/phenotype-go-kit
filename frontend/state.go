package frontend

import (
	"context"
	"sync"
	"time"
)

// State represents application state.
type State[T any] struct {
	value     T
	mu        sync.RWMutex
	listeners []chan T
}

// NewState creates a new state container.
func NewState[T any](initial T) *State[T] {
	return &State[T]{
		value:     initial,
		listeners: make([]chan T, 0),
	}
}

// Get returns the current value.
func (s *State[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set updates the value and notifies listeners.
func (s *State[T]) Set(value T) {
	s.mu.Lock()
	s.value = value
	s.mu.Unlock()

	s.notify()
}

// Subscribe returns a channel that receives value updates.
func (s *State[T]) Subscribe() <-chan T {
	ch := make(chan T, 1)
	s.mu.Lock()
	s.listeners = append(s.listeners, ch)
	s.mu.Unlock()
	return ch
}

// Unsubscribe removes a listener.
func (s *State[T]) Unsubscribe(ch <-chan T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, listener := range s.listeners {
		if listener == ch {
			s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
			return
		}
	}
}

func (s *State[T]) notify() {
	s.mu.RLock()
	listeners := make([]chan T, len(s.listeners))
	copy(listeners, s.listeners)
	s.mu.RUnlock()

	for _, ch := range listeners {
		select {
		case ch <- s.value:
		default:
		}
	}
}

// Store provides a state store with actions.
type Store[T any] struct {
	state   *State[T]
	actions map[string]func(context.Context, T, ...interface{}) (T, error)
}

// NewStore creates a new store.
func NewStore[T any](initial T) *Store[T] {
	return &Store[T]{
		state:   NewState(initial),
		actions: make(map[string]func(context.Context, T, ...interface{}) (T, error)),
	}
}

// RegisterAction registers an action handler.
func (s *Store[T]) RegisterAction(name string, handler func(context.Context, T, ...interface{}) (T, error)) {
	s.actions[name] = handler
}

// Dispatch dispatches an action.
func (s *Store[T]) Dispatch(ctx context.Context, action string, params ...interface{}) error {
	handler, ok := s.actions[action]
	if !ok {
		return nil
	}

	newState, err := handler(ctx, s.state.Get(), params...)
	if err != nil {
		return err
	}

	s.state.Set(newState)
	return nil
}

// GetState returns the state container.
func (s *Store[T]) GetState() *State[T] {
	return s.state
}

// Reducer is a state reducer function.
type Reducer[T any] func(T, interface{}) T

// ReducersStore provides a store with reducers.
type ReducersStore[T any] struct {
	state    *State[T]
	reducers map[string]Reducer[T]
}

// NewReducersStore creates a new reducers store.
func NewReducersStore[T any](initial T) *ReducersStore[T] {
	return &ReducersStore[T]{
		state:    NewState(initial),
		reducers: make(map[string]Reducer[T]),
	}
}

// RegisterReducer registers a reducer.
func (s *ReducersStore[T]) RegisterReducer(name string, reducer Reducer[T]) {
	s.reducers[name] = reducer
}

// Dispatch dispatches an action to a reducer.
func (s *ReducersStore[T]) Dispatch(action string, payload interface{}) {
	reducer, ok := s.reducers[action]
	if !ok {
		return
	}

	newState := reducer(s.state.Get(), payload)
	s.state.Set(newState)
}

// GetState returns the current state.
func (s *ReducersStore[T]) GetState() T {
	return s.state.Get()
}

// Subscribe returns a channel for state updates.
func (s *ReducersStore[T]) Subscribe() <-chan T {
	return s.state.Subscribe()
}

// Pagination represents pagination state.
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
}

// NewPagination creates default pagination.
func NewPagination(page, pageSize int) Pagination {
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// LoadingState represents loading state.
type LoadingState struct {
	Loading bool   `json:"loading"`
	Error   string `json:"error,omitempty"`
}

// ErrorState represents error state.
type ErrorState struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Timestamp wraps time.Time for JSON serialization.
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var ts time.Time
	if err := ts.UnmarshalJSON(data); err != nil {
		return err
	}
	*t = Timestamp(ts)
	return nil
}
