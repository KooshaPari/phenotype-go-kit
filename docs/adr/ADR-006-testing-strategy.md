# ADR-006: Testing Strategy

## Status
**Accepted** | 2024-01-15

## Context

We need a comprehensive testing strategy that balances speed, coverage, and maintenance burden.

## Decision

### Testing Pyramid

```
          /\
         /  \         E2E (5%)
        /    \        - Critical user paths
       /------\       - Smoke tests
      /        \      Integration (20%)
     /----------\     - Service boundaries
    /            \    - Adapter tests
   /--------------\  Unit (75%)
  /                \ - Business logic
 /                  \ - Edge cases
```

### Unit Tests (TDD)

```go
// Test naming: TestSubject_Method_Condition_Expected
func TestRingBuffer_Push_WhenFull_OverwritesOldest(t *testing.T) {
    rb := NewRingBuffer[string](3)
    rb.Push("a")
    rb.Push("b")
    rb.Push("c")
    rb.Push("d")  // Buffer full, should overwrite "a"
    
    items := rb.Items()
    require.Equal(t, []string{"b", "c", "d"}, items)
}
```

### Integration Tests

```go
func TestRedisCache_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    cache := NewRedisCache(client)
    
    err := cache.Set(ctx, "key", []byte("value"), time.Hour)
    require.NoError(t, err)
    
    val, err := cache.Get(ctx, "key")
    require.NoError(t, err)
    require.Equal(t, []byte("value"), val)
}
```

### Property-Based Tests (proptest)

```go
// Property: Ring buffer maintains FIFO order
func TestRingBuffer_Property_FIFO(t *testing.T) {
    proptest.PropTest(t, proptest.Config{
        MaxCases: 1000,
    }, func(t *proptest.T, size int, items []int) {
        size = (size % 100) + 1  // 1-100
        items = items[:len(items)%100]  // Limit items
        
        rb := NewRingBuffer[int](size)
        for _, item := range items {
            rb.Push(item)
        }
        
        // Get items in order - should be last 'size' items
        // in order of insertion
        result := rb.Items()
        require.Equal(t, min(len(items), size), len(result))
    })
}
```

### Mutation Testing

Use `mutation` package to verify test quality:

```go
// Run: go install github.com/zimnx/mutate@latest
// $ mutate ./...
// Goal: Kill >80% of mutations
```

## Coverage Requirements

| Layer | Minimum Coverage |
|-------|-----------------|
| Domain | 90% |
| Application | 85% |
| Infrastructure | 70% |
| Presentation | 60% |

## Consequences

### Positive
- Fast feedback (unit tests)
- Confidence in refactoring
- Living documentation
- Regression prevention

### Negative
- Test maintenance burden
- Slower CI (but worth it)
- Need discipline to maintain

## References
- [Go Testing Blog](https://go.dev/blog/testing)
- [proptest](github.com/proptest-rs/proptest)
- [testify](github.com/stretchr/testify)
