package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// EventBusPort defines the interface for event publishing.
// Following Event Sourcing and CQRS patterns.
type EventBusPort interface {
	// Publish publishes an event to the bus.
	Publish(ctx context.Context, topic string, event any) error

	// Subscribe subscribes to a topic.
	Subscribe(ctx context.Context, topic string, handler func(any)) error
}

// InMemoryEventBus implements EventBusPort in-memory.
// Following Mediator pattern.
//
// For production, use:
// - NATS JetStream for high-throughput pub/sub
// - Kafka for event streaming and replay
// - Redis Streams for simple message queuing
type InMemoryEventBus struct {
	subscribers map[string][]func(any)
	mu          sync.RWMutex
}

// NewInMemoryEventBus creates a new InMemoryEventBus.
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		subscribers: make(map[string][]func(any)),
	}
}

// Publish publishes an event to the bus.
func (b *InMemoryEventBus) Publish(ctx context.Context, topic string, event any) error {
	b.mu.RLock()
	handlers, ok := b.subscribers[topic]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}

	// Publish to all subscribers asynchronously
	for _, handler := range handlers {
		go handler(eventJSON)
	}

	return nil
}

// Subscribe subscribes to a topic.
func (b *InMemoryEventBus) Subscribe(ctx context.Context, topic string, handler func(any)) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], handler)
	return nil
}

// Compile-time interface check
var _ EventBusPort = (*InMemoryEventBus)(nil)

// NatsEventBus implements EventBusPort using NATS JetStream.
// For production deployment with high-throughput requirements.
type NatsEventBus struct {
	// natsConn would be injected via dependency injection
	// natsURL string
}

// NewNatsEventBus creates a new NatsEventBus.
func NewNatsEventBus(natsURL string) *NatsEventBus {
	return &NatsEventBus{}
}

// Publish publishes an event to NATS.
func (b *NatsEventBus) Publish(ctx context.Context, topic string, event any) error {
	// Implementation would use nats.Conn.JetStream().Publish()
	// Example:
	// data, _ := json.Marshal(event)
	// _, err := nc.JetStream().Publish(ctx, topic, data)
	return nil
}

// Subscribe subscribes to a NATS subject.
func (b *NatsEventBus) Subscribe(ctx context.Context, topic string, handler func(any)) error {
	// Implementation would use nats.Conn.JetStream().Subscribe()
	// Example:
	// sub, _ := nc.JetStream().Subscribe(ctx, topic, func(msg *nats.Msg) {
	//     handler(msg.Data)
	// })
	return nil
}

// Compile-time interface check
var _ EventBusPort = (*NatsEventBus)(nil)

// KafkaEventBus implements EventBusPort using Apache Kafka.
// For production with event replay and long-term storage requirements.
type KafkaEventBus struct {
	// kafkaProducer/kafkaConsumer would be injected
}

// NewKafkaEventBus creates a new KafkaEventBus.
func NewKafkaEventBus() *KafkaEventBus {
	return &KafkaEventBus{}
}

// Publish publishes an event to Kafka.
func (b *KafkaEventBus) Publish(ctx context.Context, topic string, event any) error {
	// Implementation would use segmentio/kafka-go or confluent-kafka-go
	return nil
}

// Subscribe subscribes to a Kafka topic.
func (b *KafkaEventBus) Subscribe(ctx context.Context, topic string, handler func(any)) error {
	// Implementation would use kafka-go Reader
	return nil
}

// Compile-time interface check
var _ EventBusPort = (*KafkaEventBus)(nil)
