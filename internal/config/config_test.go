package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ai-agent-manager/agent-lint/internal/config"
)

func TestLoadConstraintsValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".agent-lint.yaml")
	content := `
skills:
  require_readme: true
  readme_constraints:
    - field: name
      equals: Valid Skill
    - field: tags
      must_include: [agent-skill]
    - field: phases
      allowed: [build_and_release]
  skill_constraints:
    - field: metadata.phases
      matches: "^(build_and_release)$"
    - field: metadata.author
      not_equals: ""
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Skills.ReadmeConstraints) != 3 {
		t.Fatalf("expected 3 readme constraints, got %d", len(cfg.Skills.ReadmeConstraints))
	}
	if len(cfg.Skills.SkillConstraints) != 2 {
		t.Fatalf("expected 2 skill constraints, got %d", len(cfg.Skills.SkillConstraints))
	}
	if cfg.Skills.SkillConstraints[0].Field != "metadata.phases" {
		t.Fatalf("unexpected skill field: %q", cfg.Skills.SkillConstraints[0].Field)
	}
}

func TestLoadConstraintOperatorCount(t *testing.T) {
	tests := []struct {
		name    string
		snippet string
		wantSub string
	}{
		{
			name: "zero operators",
			snippet: `
skills:
  readme_constraints:
    - field: name
`,
			wantSub: "exactly one",
		},
		{
			name: "two operators",
			snippet: `
skills:
  skill_constraints:
    - field: metadata.author
      equals: a
      not_equals: b
`,
			wantSub: "exactly one",
		},
		{
			name: "invalid regexp",
			snippet: `
skills:
  skill_constraints:
    - field: metadata.phases
      matches: "(unclosed"
`,
			wantSub: "invalid matches regexp",
		},
		{
			name: "missing field",
			snippet: `
skills:
  readme_constraints:
    - equals: x
`,
			wantSub: "field is required",
		},
		{
			name: "list ops rejected on skill_constraints",
			snippet: `
skills:
  skill_constraints:
    - field: metadata.tags
      must_include: [agent-skill]
`,
			wantSub: "only valid on readme_constraints",
		},
		{
			name: "mixed string and list ops",
			snippet: `
skills:
  readme_constraints:
    - field: tags
      equals: agent-skill
      must_include: [agent-skill]
`,
			wantSub: "do not mix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, ".agent-lint.yaml")
			if err := os.WriteFile(path, []byte(tt.snippet), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := config.Load(path)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("error %q does not contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}
