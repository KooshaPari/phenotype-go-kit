package registry

import (
	"sync"
)

// Hook is called when registry entries change.
type Hook[K comparable, V any] interface {
	OnRegister(ownerID string, key K, value V)
	OnUnregister(ownerID string)
}

// entry tracks a value and its reference count.
type entry[V any] struct {
	value V
	count int
}

// Registry is a generic, thread-safe key-value store with owner tracking.
type Registry[K comparable, V any] struct {
	mu        sync.RWMutex
	entries   map[K]*entry[V]
	ownerKeys map[string][]K
	hook      Hook[K, V]
}

// New creates a new empty Registry.
func New[K comparable, V any]() *Registry[K, V] {
	return &Registry[K, V]{
		entries:   make(map[K]*entry[V]),
		ownerKeys: make(map[string][]K),
	}
}

// SetHook sets an optional hook for observing changes.
func (r *Registry[K, V]) SetHook(hook Hook[K, V]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hook = hook
}

// Register adds or increments a key under the given owner.
func (r *Registry[K, V]) Register(ownerID string, key K, value V) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if e, ok := r.entries[key]; ok {
		e.count++
		e.value = value
	} else {
		r.entries[key] = &entry[V]{value: value, count: 1}
	}
	r.ownerKeys[ownerID] = append(r.ownerKeys[ownerID], key)

	if r.hook != nil {
		r.hook.OnRegister(ownerID, key, value)
	}
}

// Unregister removes all keys owned by ownerID, decrementing counts.
func (r *Registry[K, V]) Unregister(ownerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	keys, ok := r.ownerKeys[ownerID]
	if !ok {
		return
	}
	for _, key := range keys {
		if e, exists := r.entries[key]; exists {
			e.count--
			if e.count <= 0 {
				delete(r.entries, key)
			}
		}
	}
	delete(r.ownerKeys, ownerID)

	if r.hook != nil {
		r.hook.OnUnregister(ownerID)
	}
}

// Get retrieves a value by key.
func (r *Registry[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if e, ok := r.entries[key]; ok {
		return e.value, true
	}
	var zero V
	return zero, false
}

// List returns a snapshot of all entries.
func (r *Registry[K, V]) List() map[K]V {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[K]V, len(r.entries))
	for k, e := range r.entries {
		result[k] = e.value
	}
	return result
}

// Count returns the reference count for a key.
func (r *Registry[K, V]) Count(key K) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if e, ok := r.entries[key]; ok {
		return e.count
	}
	return 0
}
