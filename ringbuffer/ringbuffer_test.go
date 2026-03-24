package ringbuffer_test

import (
	"testing"

	"github.com/KooshaPari/phenotype-go-kit/ringbuffer"
)

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
	rb.Push(3)
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
	rb.Push(4)
	if rb.Len() != 3 {
		t.Fatalf("expected len 3 (capped), got %d", rb.Len())
	}
}
