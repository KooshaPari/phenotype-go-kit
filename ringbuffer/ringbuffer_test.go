package ringbuffer_test

import (
	"math/rand"
	"testing"

	"github.com/KooshaPari/phenotype-go-kit/ringbuffer"
)

// Test naming convention following BDD:
// Test<Subject>_<Method>_<Condition>_<Expected>

func TestPushAndGetAll(t *testing.T) {
	rb := ringbuffer.New[int](3)
	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	got := rb.GetAll()
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("expected [1 2 3], got %v", got)
	}
}

func TestOverflow(t *testing.T) {
	rb := ringbuffer.New[int](2)
	rb.Push(1)
	rb.Push(2)
	rb.Push(3) // Should overwrite 1
	got := rb.GetAll()
	if len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("expected [2 3], got %v", got)
	}
}

func TestEmpty(t *testing.T) {
	rb := ringbuffer.New[string](5)
	if rb.Len() != 0 {
		t.Fatal("expected len 0")
	}
	got := rb.GetAll()
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestLen(t *testing.T) {
	rb := ringbuffer.New[int](3)
	rb.Push(1)
	if rb.Len() != 1 {
		t.Fatalf("expected len 1, got %d", rb.Len())
	}
	rb.Push(2)
	rb.Push(3)
	rb.Push(4) // Overflow, should stay at 3
	if rb.Len() != 3 {
		t.Fatalf("expected len 3 (capped), got %d", rb.Len())
	}
}

// Property-based tests (following TDD + Property Testing patterns)

func TestProperty_FIFOOrder(t *testing.T) {
	// Property: Items are always returned in FIFO order
	for capacity := 1; capacity <= 10; capacity++ {
		for nPush := 1; nPush <= 20; nPush++ {
			rb := ringbuffer.New[int](capacity)
			var expected []int

			for i := 0; i < nPush; i++ {
				rb.Push(i)
				expected = append(expected, i)
				if len(expected) > capacity {
					expected = expected[len(expected)-capacity:]
				}
			}

			result := rb.GetAll()
			if len(result) != len(expected) {
				t.Errorf("capacity=%d, nPush=%d: len=%d, want %d",
					capacity, nPush, len(result), len(expected))
			}

			for i := range result {
				if result[i] != expected[i] {
					t.Errorf("capacity=%d, nPush=%d: result[%d]=%d, want %d",
						capacity, nPush, i, result[i], expected[i])
				}
			}
		}
	}
}

func TestProperty_CapacityNeverExceeds(t *testing.T) {
	// Property: Len() never exceeds Cap()
	rb := ringbuffer.New[int](3)

	for i := 0; i < 100; i++ {
		rb.Push(i)
		if len := rb.Len(); len > rb.Cap() {
			t.Errorf("Len()=%d exceeds Cap()=%d", len, rb.Cap())
		}
	}
}

func TestProperty_MaxCapacityAfterOverflow(t *testing.T) {
	// Property: After capacity+1 pushes, Len() == Cap()
	for cap := 1; cap <= 10; cap++ {
		rb := ringbuffer.New[int](cap)
		for i := 0; i < cap+5; i++ {
			rb.Push(i)
		}
		if len := rb.Len(); len != cap {
			t.Errorf("cap=%d: after overflow, Len()=%d, want %d", cap, len, cap)
		}
	}
}

func TestProperty_EmptyBufferReturnsEmptySlice(t *testing.T) {
	// Property: New buffer has empty GetAll()
	for cap := 1; cap <= 10; cap++ {
		rb := ringbuffer.New[int](cap)
		got := rb.GetAll()
		if len(got) != 0 {
			t.Errorf("cap=%d: new buffer GetAll() len=%d, want 0", cap, len(got))
		}
	}
}

func TestProperty_AllTypes(t *testing.T) {
	// Property: Works with any comparable type
	rb := ringbuffer.New[string](3)
	rb.Push("a")
	rb.Push("b")
	rb.Push("c")
	got := rb.GetAll()
	want := []string{"a", "b", "c"}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("got[%d]=%q, want %q", i, got[i], v)
		}
	}
}

// Fuzz test (following mutation/property testing)
func TestFuzz_RandomOperations(t *testing.T) {
	// Fuzz: Random operations maintain invariants
	seed := int64(12345)
	rng := rand.New(rand.NewSource(seed))

	capacity := 5
	rb := ringbuffer.New[int](capacity)

	for i := 0; i < 1000; i++ {
		op := rng.Intn(3)
		switch op {
		case 0: // Push
			rb.Push(rng.Intn(100))
		case 1: // Check length invariant
			if len := rb.Len(); len < 0 || len > capacity {
				t.Errorf("Len()=%d out of bounds [0,%d]", len, capacity)
			}
		case 2: // Check no nil panic
			_ = rb.GetAll()
		}
	}
}

func TestProperty_PushPreservesNewestItems(t *testing.T) {
	// Property: After overflow, newest items are preserved
	capacity := 3
	rb := ringbuffer.New[int](capacity)

	// Push 5 items into capacity 3
	pushes := []int{1, 2, 3, 4, 5}
	for _, v := range pushes {
		rb.Push(v)
	}

	// Should have [3, 4, 5]
	got := rb.GetAll()
	want := []int{3, 4, 5}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("got[%d]=%d, want %d", i, got[i], v)
		}
	}
}

func TestProperty_ZeroCapacity(t *testing.T) {
	// Edge case: capacity of 1
	rb := ringbuffer.New[int](1)
	rb.Push(10)
	rb.Push(20)
	got := rb.GetAll()
	if len(got) != 1 || got[0] != 20 {
		t.Errorf("got %v, want [20]", got)
	}
}

func TestProperty_ManyItems(t *testing.T) {
	// Stress test: large number of pushes
	rb := ringbuffer.New[int](100)
	for i := 0; i < 10000; i++ {
		rb.Push(i)
	}

	if len := rb.Len(); len != 100 {
		t.Errorf("Len()=%d, want 100", len)
	}

	// Last item should be 9999
	got := rb.GetAll()
	if got[99] != 9999 {
		t.Errorf("got[99]=%d, want 9999", got[99])
	}
}
