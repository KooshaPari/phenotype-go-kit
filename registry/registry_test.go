package registry_test

import (
	"testing"

	"github.com/KooshaPari/phenotype-go-kit/registry"
)

type testHook struct {
	registered   int
	unregistered int
}

func (h *testHook) OnRegister(ownerID string, key string, value int) { h.registered++ }
func (h *testHook) OnUnregister(ownerID string)                      { h.unregistered++ }

func TestRegisterAndGet(t *testing.T) {
	r := registry.New[string, int]()
	r.Register("owner1", "a", 1)
	v, ok := r.Get("a")
	if !ok || v != 1 {
		t.Fatalf("expected (1, true), got (%d, %v)", v, ok)
	}
}

func TestUnregister(t *testing.T) {
	r := registry.New[string, int]()
	r.Register("owner1", "a", 1)
	r.Unregister("owner1")
	_, ok := r.Get("a")
	if ok {
		t.Fatal("expected key to be removed after unregister")
	}
}

func TestCount(t *testing.T) {
	r := registry.New[string, int]()
	r.Register("o1", "a", 1)
	r.Register("o2", "a", 2)
	if r.Count("a") != 2 {
		t.Fatalf("expected count 2, got %d", r.Count("a"))
	}
	r.Unregister("o1")
	if r.Count("a") != 1 {
		t.Fatalf("expected count 1 after unregister, got %d", r.Count("a"))
	}
}

func TestList(t *testing.T) {
	r := registry.New[string, int]()
	r.Register("o1", "a", 1)
	r.Register("o1", "b", 2)
	m := r.List()
	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
}

func TestHooks(t *testing.T) {
	r := registry.New[string, int]()
	h := &testHook{}
	r.SetHook(h)
	r.Register("o1", "a", 1)
	r.Unregister("o1")
	if h.registered != 1 || h.unregistered != 1 {
		t.Fatalf("expected 1/1 hook calls, got %d/%d", h.registered, h.unregistered)
	}
}

func TestGetMissing(t *testing.T) {
	r := registry.New[string, int]()
	_, ok := r.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}
