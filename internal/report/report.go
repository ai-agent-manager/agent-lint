// Package report defines diagnostics and output formatting.
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ai-agent-manager/agent-lint/internal/rules"
)

// Diagnostic is a single lint finding.
type Diagnostic struct {
	RuleID   string         `json:"ruleId"`
	Severity rules.Severity `json:"severity"`
	Path     string         `json:"path"`
	Message  string         `json:"message"`
	Ignored  bool           `json:"ignored,omitempty"`
}

// Result is the aggregate lint result.
type Result struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
	Ignored     []Diagnostic `json:"ignored"`
	ErrorCount  int          `json:"errorCount"`
	WarningCount int         `json:"warningCount"`
	IgnoredCount int         `json:"ignoredCount"`
}

// NewResult partitions diagnostics into active vs ignored and counts them.
func NewResult(all []Diagnostic) Result {
	var active, ignored []Diagnostic
	var errors, warnings int
	for _, d := range all {
		if d.Ignored {
			ignored = append(ignored, d)
			continue
		}
		active = append(active, d)
		switch d.Severity {
		case rules.SeverityError:
			errors++
		case rules.SeverityWarning:
			warnings++
		}
	}
	return Result{
		Diagnostics:  active,
		Ignored:      ignored,
		ErrorCount:   errors,
		WarningCount: warnings,
		IgnoredCount: len(ignored),
	}
}

// ExitCode returns 0 on success, 1 on lint failures, respecting strict mode.
func (r Result) ExitCode(strict bool) int {
	if r.ErrorCount > 0 {
		return 1
	}
	if strict && r.WarningCount > 0 {
		return 1
	}
	return 0
}

const (
	reset  = "\x1b[0m"
	bold   = "\x1b[1m"
	dim    = "\x1b[2m"
	red    = "\x1b[31m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	cyan   = "\x1b[36m"
)

func useColour(w io.Writer) bool {
	f, ok := w.(*os.File)
	return ok && isTerminal(f)
}

func colour(enabled bool, text string, codes ...string) string {
	if !enabled {
		return text
	}
	return strings.Join(codes, "") + text + reset
}

// WriteHuman prints a coloured human-readable report.
func WriteHuman(w io.Writer, r Result, quiet bool) {
	c := useColour(w)
	if !quiet {
		fmt.Fprintln(w, colour(c, fmt.Sprintf("Found %d error(s), %d warning(s), %d ignored.", r.ErrorCount, r.WarningCount, r.IgnoredCount), bold))
		fmt.Fprintln(w)
	}

	for _, d := range r.Diagnostics {
		sevColour := red
		if d.Severity == rules.SeverityWarning {
			sevColour = yellow
		}
		label := colour(c, strings.ToUpper(string(d.Severity)), bold, sevColour)
		fmt.Fprintf(w, "%s %s: %s\n", label, d.Path, d.Message)
		fmt.Fprintf(w, "  %s\n", colour(c, d.RuleID, dim, cyan))
	}

	if r.ErrorCount == 0 && r.WarningCount == 0 {
		fmt.Fprintln(w, colour(c, "All checks passed.", bold, green))
	}
}

// WriteJSON prints a machine-readable JSON report.
func WriteJSON(w io.Writer, r Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
