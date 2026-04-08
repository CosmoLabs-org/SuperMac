package module

import (
	"testing"
)

type testModule struct {
	name        string
	desc        string
	emoji       string
}

func (t *testModule) Name() string                { return t.name }
func (t *testModule) ShortDescription() string     { return t.desc }
func (t *testModule) Emoji() string                { return t.emoji }
func (t *testModule) Commands() []Command          { return nil }
func (t *testModule) Search(term string) []SearchResult { return nil }

func TestRegisterAndGet(t *testing.T) {
	// Clear registry for test isolation
	modulesMu.Lock()
	modules = make(map[string]Module)
	modulesMu.Unlock()

	mod := &testModule{name: "testmod", desc: "A test module", emoji: "🧪"}
	Register(mod)

	got, ok := Get("testmod")
	if !ok {
		t.Fatal("expected module to be registered")
	}
	if got.Name() != "testmod" {
		t.Errorf("expected name testmod, got %s", got.Name())
	}
}

func TestGetNonExistent(t *testing.T) {
	modulesMu.Lock()
	modules = make(map[string]Module)
	modulesMu.Unlock()

	_, ok := Get("nonexistent")
	if ok {
		t.Error("expected module not to be found")
	}
}

func TestAll(t *testing.T) {
	modulesMu.Lock()
	modules = make(map[string]Module)
	modulesMu.Unlock()

	Register(&testModule{name: "alpha", desc: "Alpha", emoji: "🅰️"})
	Register(&testModule{name: "beta", desc: "Beta", emoji: "🅱️"})

	all := All()
	if len(all) != 2 {
		t.Errorf("expected 2 modules, got %d", len(all))
	}
	if _, ok := all["alpha"]; !ok {
		t.Error("expected alpha in All()")
	}
	if _, ok := all["beta"]; !ok {
		t.Error("expected beta in All()")
	}
}

func TestNames(t *testing.T) {
	modulesMu.Lock()
	modules = make(map[string]Module)
	modulesMu.Unlock()

	Register(&testModule{name: "zebra", desc: "Z", emoji: "🦓"})
	Register(&testModule{name: "alpha", desc: "A", emoji: "🅰️"})

	names := Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
