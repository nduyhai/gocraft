package embed_registry

import (
	"fmt"
	"sort"

	"github.com/nduyhai/gocraft/internal/core/ports"
)

// Registry is an in-memory module registry with a embed_registry planner.
// It keeps deterministic registration order for List(), while Apply() resolves
// transitive Requires(), detects Conflicts() and cycles, and applies modules
// in a valid topological order.
type Registry struct {
	order  []ports.Module
	byName map[string]ports.Module
}

func New() *Registry {
	return &Registry{byName: make(map[string]ports.Module)}
}

func (r *Registry) Register(m ports.Module) {
	name := m.Name()
	if _, exists := r.byName[name]; exists {
		// overwrite allowed; but keep first registration order stable
		r.byName[name] = m
		return
	}
	r.byName[name] = m
	r.order = append(r.order, m)
}

func (r *Registry) List() []ports.Module {
	// return a copy to avoid external modification
	out := make([]ports.Module, len(r.order))
	copy(out, r.order)
	return out
}

func (r *Registry) Get(name string) (ports.Module, bool) { m, ok := r.byName[name]; return m, ok }

func (r *Registry) Apply(ctx ports.Ctx, names ...string) error {
	if len(names) == 0 {
		return nil
	}
	// Expand requires transitively
	expanded, err := r.expandRequires(names)
	if err != nil {
		return err
	}
	// Conflicts detection
	if err := r.checkConflicts(expanded); err != nil {
		return err
	}
	// Toposort using DFS with cycle detection
	ordered, err := r.toposort(expanded)
	if err != nil {
		return err
	}
	// Apply in order
	for _, name := range ordered {
		m := r.byName[name]
		if m == nil {
			return fmt.Errorf("module not found during apply: %s", name)
		}
		if !m.Applies(ctx) {
			continue
		}
		if err := m.Apply(ctx); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
	}
	return nil
}

func (r *Registry) expandRequires(names []string) (map[string]struct{}, error) {
	seen := make(map[string]struct{})
	var visit func(string) error
	visit = func(name string) error {
		if _, ok := seen[name]; ok {
			return nil
		}
		m, ok := r.byName[name]
		if !ok {
			return fmt.Errorf("unknown module: %s", name)
		}
		seen[name] = struct{}{}
		for _, req := range m.Requires() {
			if err := visit(req); err != nil {
				return err
			}
		}
		return nil
	}
	for _, n := range names {
		if err := visit(n); err != nil {
			return nil, err
		}
	}
	return seen, nil
}

func (r *Registry) checkConflicts(set map[string]struct{}) error {
	for name := range set {
		m := r.byName[name]
		if m == nil {
			continue
		}
		for _, c := range m.Conflicts() {
			if _, ok := set[c]; ok {
				return fmt.Errorf("conflict: %s conflicts with %s", name, c)
			}
		}
	}
	return nil
}

func (r *Registry) toposort(set map[string]struct{}) ([]string, error) {
	// Build adjacency: name -> requires
	adj := make(map[string][]string, len(set))
	for name := range set {
		m := r.byName[name]
		if m == nil {
			return nil, fmt.Errorf("unknown module in set: %s", name)
		}
		adj[name] = append([]string(nil), m.Requires()...)
	}
	// DFS states: 0 = unvisited; 1 = visiting; 2 = visited
	state := make(map[string]int, len(set))
	var order []string
	var dfs func(string) error
	dfs = func(u string) error {
		s := state[u]
		if s == 1 {
			return fmt.Errorf("cycle detected at %s", u)
		}
		if s == 2 {
			return nil
		}
		state[u] = 1
		for _, v := range adj[u] {
			if _, ok := set[v]; !ok {
				// If a required module isn't requested explicitly but is part of dependencies, we still included it in set via expandRequires
				// However, if a module requires something unregistered, we error early here
				if _, exists := r.byName[v]; !exists {
					return fmt.Errorf("required module not registered: %s (needed by %s)", v, u)
				}
				// It is registered but not in set: add and continue DFS
				set[v] = struct{}{}
				// also extend adjacency
				adj[v] = append([]string(nil), r.byName[v].Requires()...)
			}
			if err := dfs(v); err != nil {
				return err
			}
		}
		state[u] = 2
		order = append(order, u)
		return nil
	}
	// Use deterministic traversal order: sort keys for stable output
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := dfs(k); err != nil {
			return nil, err
		}
	}
	// reverse postorder gives topological order with requires before dependents
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, nil
}
