// Package ignore filters diagnostics by path/glob and rule ID.
package ignore

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ai-agent-manager/agent-lint/internal/config"
	"github.com/ai-agent-manager/agent-lint/internal/report"
	"github.com/ai-agent-manager/agent-lint/internal/rules"
	"github.com/bmatcuk/doublestar/v4"
)

// Filter applies ignore entries to diagnostics.
type Filter struct {
	entries []config.IgnoreEntry
}

// New builds a filter and validates rule IDs.
func New(entries []config.IgnoreEntry) (*Filter, error) {
	known := rules.IDs()
	for i, e := range entries {
		if strings.TrimSpace(e.Path) == "" {
			return nil, fmt.Errorf("ignore entry %d: path is required", i)
		}
		for _, id := range e.Rules {
			if _, ok := known[id]; !ok {
				return nil, fmt.Errorf("ignore entry %d: unknown rule id %q", i, id)
			}
		}
	}
	return &Filter{entries: entries}, nil
}

// Apply marks matching diagnostics as Ignored.
func (f *Filter) Apply(diags []report.Diagnostic) []report.Diagnostic {
	if f == nil || len(f.entries) == 0 {
		return diags
	}
	out := make([]report.Diagnostic, len(diags))
	copy(out, diags)
	for i := range out {
		if f.shouldIgnore(out[i]) {
			out[i].Ignored = true
		}
	}
	return out
}

func (f *Filter) shouldIgnore(d report.Diagnostic) bool {
	path := filepath.ToSlash(d.Path)
	for _, e := range f.entries {
		pattern := filepath.ToSlash(e.Path)
		if !matchPath(pattern, path) {
			continue
		}
		// Empty rules list = ignore all findings on matching paths.
		if len(e.Rules) == 0 {
			return true
		}
		for _, id := range e.Rules {
			if id == d.RuleID {
				return true
			}
		}
	}
	return false
}

func matchPath(pattern, path string) bool {
	pattern = strings.TrimPrefix(pattern, "./")
	path = strings.TrimPrefix(path, "./")

	// Exact match
	if pattern == path {
		return true
	}

	// Directory prefix without glob: "skills/foo" matches "skills/foo" and children
	if !hasGlob(pattern) {
		if path == pattern || strings.HasPrefix(path, strings.TrimSuffix(pattern, "/")+"/") {
			return true
		}
		return false
	}

	ok, err := doublestar.Match(pattern, path)
	if err == nil && ok {
		return true
	}
	// Also try matching if pattern is a directory glob without trailing /**
	if !strings.HasSuffix(pattern, "/**") && !strings.Contains(filepath.Base(pattern), "*") {
		ok, err = doublestar.Match(pattern+"/**", path)
		if err == nil && ok {
			return true
		}
	}
	return false
}

func hasGlob(s string) bool {
	return strings.ContainsAny(s, "*?[")
}
