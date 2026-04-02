package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Message represents a message in the event bus.
type Message struct {
	ID        string
	EventType string
	Payload   interface{}
	Metadata  map[string]string
	Timestamp time.Time
}

// EventBus provides in-memory event bus functionality.
type EventBus struct {
	subscribers map[string][]chan Message
	mu          sync.RWMutex
	logger      *slog.Logger
}

// New creates a new event bus.
func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan Message),
		logger:      slog.Default(),
	}
}

// Subscribe adds a subscriber for an event type.
func (eb *EventBus) Subscribe(eventType string, handler func(ctx context.Context, msg Message) error) <-chan error {
	errChan := make(chan error, 1)

	ch := make(chan Message, 100)

	eb.mu.Lock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
	eb.mu.Unlock()

	go func() {
		for msg := range ch {
			if err := handler(context.Background(), msg); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}
		close(errChan)
	}()

	return errChan
}

// Publish sends a message to all subscribers.
func (eb *EventBus) Publish(eventType string, payload interface{}, metadata map[string]string) error {
	msg := Message{
		ID:        generateID(),
		EventType: eventType,
		Payload:   payload,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}

	eb.mu.RLock()
	subscribers, ok := eb.subscribers[eventType]
	eb.mu.RUnlock()

	if !ok || len(subscribers) == 0 {
		eb.logger.Debug("no subscribers for event", "type", eventType)
		return nil
	}

	for _, ch := range subscribers {
		select {
		case ch <- msg:
		default:
			eb.logger.Warn("subscriber channel full", "type", eventType)
		}
	}

	return nil
}

// PublishAsync publishes asynchronously.
func (eb *EventBus) PublishAsync(eventType string, payload interface{}, metadata map[string]string) {
	go func() {
		if err := eb.Publish(eventType, payload, metadata); err != nil {
			eb.logger.Error("async publish failed", "type", eventType, "error", err)
		}
	}()
}

// Unsubscribe removes all handlers for an event type.
func (eb *EventBus) Unsubscribe(eventType string) {
	eb.mu.Lock()
	delete(eb.subscribers, eventType)
	eb.mu.Unlock()
}

// Close closes all subscriber channels.
func (eb *Eb) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for _, chs := range eb.subscribers {
		for _, ch := range chs {
			close(ch)
		}
	}
	eb.subscribers = make(map[string][]chan Message)
}

type Eb = EventBus

func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randString(8))
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[int(time.Now().UnixNano())%len(letters)]
	}
	return string(b)
}

// JSONPayload represents a JSON-encodable message.
type JSONPayload struct {
	Type    string            `json:"type"`
	Data    json.RawMessage   `json:"data"`
	Meta    map[string]string `json:"meta,omitempty"`
	TraceID string            `json:"trace_id,omitempty"`
}

// EncodeJSON encodes a message as JSON.
func (m *Message) EncodeJSON() ([]byte, error) {
	payload := JSONPayload{
		Type:    m.EventType,
		Meta:    m.Metadata,
		TraceID: m.Metadata["trace_id"],
	}

	if data, ok := m.Payload.(json.RawMessage); ok {
		payload.Data = data
	} else {
		payload.Data, _ = json.Marshal(m.Payload)
	}

	return json.Marshal(payload)
}
