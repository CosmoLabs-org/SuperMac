package module

import "sync"

// Package-level registry — NOT in init() — avoids race with module init() registrations.
var (
	modulesMu sync.RWMutex
	modules   = make(map[string]Module)
)

// Register adds a module to the global registry.
// Called from each module's init() function via blank import.
func Register(m Module) {
	modulesMu.Lock()
	defer modulesMu.Unlock()
	modules[m.Name()] = m
}

// All returns a copy of all registered modules.
func All() map[string]Module {
	modulesMu.RLock()
	defer modulesMu.RUnlock()
	out := make(map[string]Module, len(modules))
	for k, v := range modules {
		out[k] = v
	}
	return out
}

// Get returns a single module by name.
func Get(name string) (Module, bool) {
	modulesMu.RLock()
	defer modulesMu.RUnlock()
	m, ok := modules[name]
	return m, ok
}

// Names returns sorted module names.
func Names() []string {
	modulesMu.RLock()
	defer modulesMu.RUnlock()
	names := make([]string, 0, len(modules))
	for k := range modules {
		names = append(names, k)
	}
	return names
}
