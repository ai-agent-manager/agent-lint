package skill_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ai-agent-manager/agent-lint/internal/config"
	"github.com/ai-agent-manager/agent-lint/internal/ignore"
	"github.com/ai-agent-manager/agent-lint/internal/report"
	"github.com/ai-agent-manager/agent-lint/internal/rules"
	"github.com/ai-agent-manager/agent-lint/internal/skill"
)

func testdataRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "testdata"))
}

func hasRule(diags []report.Diagnostic, id string) bool {
	for _, d := range diags {
		if d.RuleID == id && !d.Ignored {
			return true
		}
	}
	return false
}

func strPtr(s string) *string { return &s }

func TestValidSkill(t *testing.T) {
	root := testdataRoot(t)
	cfg := config.Default()
	cfg.Skills.RequireReadme = true
	cfg.Skills.ReadmeConstraints = []config.FrontmatterConstraint{
		{Field: "tags", Equals: strPtr("agent-skill")},
		{Field: "phases", Matches: strPtr("^build_and_release$")},
	}
	cfg.Skills.SkillConstraints = []config.FrontmatterConstraint{
		{Field: "metadata.phases", Equals: strPtr("build_and_release")},
		{Field: "metadata.author", NotEquals: strPtr("")},
	}

	diags, err := skill.Check(skill.Options{
		Config:     cfg,
		ScanRoot:   root,
		SkillsRoot: root,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range diags {
		if d.Path == "valid-skill/SKILL.md" || d.Path == "valid-skill/README.md" || d.Path == "valid-skill" {
			if d.Severity == rules.SeverityError {
				t.Errorf("unexpected error on valid skill: %s: %s", d.RuleID, d.Message)
			}
		}
	}
}

func TestReadmeConstraints(t *testing.T) {
	tests := []struct {
		name          string
		readme        string
		requireReadme bool
		constraints   []config.FrontmatterConstraint
		wantIDs       []string
	}{
		{
			name: "equals failure",
			readme: `---
name: demo-skill
description: Demo
tags: other
---
`,
			requireReadme: true,
			constraints: []config.FrontmatterConstraint{
				{Field: "tags", Equals: strPtr("agent-skill")},
			},
			wantIDs: []string{rules.ReadmeFieldEquals},
		},
		{
			name: "matches failure",
			readme: `---
name: demo-skill
description: Demo
phases: unknown_phase
---
`,
			requireReadme: true,
			constraints: []config.FrontmatterConstraint{
				{Field: "phases", Matches: strPtr("^build_and_release$")},
			},
			wantIDs: []string{rules.ReadmeFieldMatches},
		},
		{
			name: "not_equals failure",
			readme: `---
name: demo-skill
description: Demo
tags: forbidden
---
`,
			requireReadme: true,
			constraints: []config.FrontmatterConstraint{
				{Field: "tags", NotEquals: strPtr("forbidden")},
			},
			wantIDs: []string{rules.ReadmeFieldNotEquals},
		},
		{
			name: "not_matches failure",
			readme: `---
name: demo-skill
description: Demo
tags: temp-123
---
`,
			requireReadme: true,
			constraints: []config.FrontmatterConstraint{
				{Field: "tags", NotMatches: strPtr("^temp-")},
			},
			wantIDs: []string{rules.ReadmeFieldNotMatches},
		},
		{
			name: "field not a string",
			readme: `---
name: demo-skill
description: Demo
tags:
  - agent-skill
---
`,
			requireReadme: true,
			constraints: []config.FrontmatterConstraint{
				{Field: "tags", Equals: strPtr("agent-skill")},
			},
			wantIDs: []string{rules.ReadmeFieldInvalid},
		},
		{
			name: "constraints ignored when require_readme is false",
			readme: `---
name: demo-skill
description: Demo
tags: other
---
`,
			requireReadme: false,
			constraints: []config.FrontmatterConstraint{
				{Field: "tags", Equals: strPtr("agent-skill")},
			},
			wantIDs: nil,
		},
		{
			name: "valid constraints",
			readme: `---
name: demo-skill
description: Demo
tags: agent-skill
phases: build_and_release
---
`,
			requireReadme: true,
			constraints: []config.FrontmatterConstraint{
				{Field: "tags", Equals: strPtr("agent-skill")},
				{Field: "phases", Matches: strPtr("^build_and_release$")},
			},
			wantIDs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			skillDir := filepath.Join(root, "demo-skill")
			if err := os.MkdirAll(skillDir, 0o755); err != nil {
				t.Fatal(err)
			}
			skillMD := `---
name: demo-skill
description: Demo skill for README constraint tests.
---

# Demo
`
			if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillMD), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(skillDir, "README.md"), []byte(tt.readme), 0o644); err != nil {
				t.Fatal(err)
			}

			cfg := config.Default()
			cfg.Skills.RequireReadme = tt.requireReadme
			cfg.Skills.ReadmeConstraints = tt.constraints

			diags, err := skill.Check(skill.Options{
				Config:     cfg,
				ScanRoot:   root,
				SkillsRoot: root,
			})
			if err != nil {
				t.Fatal(err)
			}

			for _, id := range tt.wantIDs {
				if !hasRule(diags, id) {
					t.Fatalf("expected rule %s; got %#v", id, diags)
				}
			}
			if tt.wantIDs == nil {
				for _, d := range diags {
					if strings.HasPrefix(d.RuleID, "skill/readme-field-") && !d.Ignored {
						t.Fatalf("unexpected constraint error: %s: %s", d.RuleID, d.Message)
					}
				}
			}
		})
	}
}

func TestSkillConstraints(t *testing.T) {
	tests := []struct {
		name        string
		frontmatter string
		constraints []config.FrontmatterConstraint
		wantIDs     []string
		forbidIDs   []string
	}{
		{
			name: "metadata path passes without allowlist carve-out",
			frontmatter: `---
name: demo-skill
description: Demo skill for skill constraint tests.
metadata:
  phases: build_and_release
---
`,
			constraints: []config.FrontmatterConstraint{
				{Field: "metadata.phases", Equals: strPtr("build_and_release")},
			},
			forbidIDs: []string{rules.FrontmatterUnknownField, rules.SkillFieldInvalid, rules.SkillFieldEquals},
		},
		{
			name: "top-level constrained field is carved out",
			frontmatter: `---
name: demo-skill
description: Demo skill for skill constraint tests.
phases: build_and_release
---
`,
			constraints: []config.FrontmatterConstraint{
				{Field: "phases", Equals: strPtr("build_and_release")},
			},
			forbidIDs: []string{rules.FrontmatterUnknownField, rules.SkillFieldInvalid, rules.SkillFieldEquals},
		},
		{
			name: "equals failure on metadata",
			frontmatter: `---
name: demo-skill
description: Demo skill for skill constraint tests.
metadata:
  author: other
---
`,
			constraints: []config.FrontmatterConstraint{
				{Field: "metadata.author", Equals: strPtr("agent-lint")},
			},
			wantIDs: []string{rules.SkillFieldEquals},
		},
		{
			name: "matches failure",
			frontmatter: `---
name: demo-skill
description: Demo skill for skill constraint tests.
metadata:
  phases: unknown_phase
---
`,
			constraints: []config.FrontmatterConstraint{
				{Field: "metadata.phases", Matches: strPtr("^build_and_release$")},
			},
			wantIDs: []string{rules.SkillFieldMatches},
		},
		{
			name: "field missing",
			frontmatter: `---
name: demo-skill
description: Demo skill for skill constraint tests.
---
`,
			constraints: []config.FrontmatterConstraint{
				{Field: "metadata.phases", Equals: strPtr("build_and_release")},
			},
			wantIDs: []string{rules.SkillFieldInvalid},
		},
		{
			name: "unconstrained unknown field still rejected",
			frontmatter: `---
name: demo-skill
description: Demo skill for skill constraint tests.
owner: platform
metadata:
  phases: build_and_release
---
`,
			constraints: []config.FrontmatterConstraint{
				{Field: "metadata.phases", Equals: strPtr("build_and_release")},
			},
			wantIDs: []string{rules.FrontmatterUnknownField},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			skillDir := filepath.Join(root, "demo-skill")
			if err := os.MkdirAll(skillDir, 0o755); err != nil {
				t.Fatal(err)
			}
			content := tt.frontmatter + "\n# Demo\n"
			if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}

			cfg := config.Default()
			cfg.Skills.SkillConstraints = tt.constraints

			diags, err := skill.Check(skill.Options{
				Config:     cfg,
				ScanRoot:   root,
				SkillsRoot: root,
			})
			if err != nil {
				t.Fatal(err)
			}

			for _, id := range tt.wantIDs {
				if !hasRule(diags, id) {
					t.Fatalf("expected rule %s; got %#v", id, diags)
				}
			}
			for _, id := range tt.forbidIDs {
				if hasRule(diags, id) {
					t.Fatalf("did not expect rule %s; got %#v", id, diags)
				}
			}
		})
	}
}

func TestBadName(t *testing.T) {
	root := testdataRoot(t)
	diags, err := skill.Check(skill.Options{
		Config:     config.Default(),
		ScanRoot:   root,
		SkillsRoot: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasRule(diags, rules.NameFormat) && !hasRule(diags, rules.NameDirectoryMismatch) {
		t.Fatalf("expected name format or directory mismatch; got %#v", diags)
	}
}

func TestBodyLineLimit(t *testing.T) {
	root := testdataRoot(t)
	diags, err := skill.Check(skill.Options{
		Config:     config.Default(),
		ScanRoot:   root,
		SkillsRoot: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasRule(diags, rules.BodyLineLimit) {
		t.Fatalf("expected body line limit error; got %#v", diags)
	}
}

func TestIgnoreBodyLineLimit(t *testing.T) {
	root := testdataRoot(t)
	diags, err := skill.Check(skill.Options{
		Config:     config.Default(),
		ScanRoot:   root,
		SkillsRoot: root,
	})
	if err != nil {
		t.Fatal(err)
	}

	filter, err := ignore.New([]config.IgnoreEntry{{
		Path:  "oversized/SKILL.md",
		Rules: []string{rules.BodyLineLimit, rules.BodyTokenLimit},
	}})
	if err != nil {
		t.Fatal(err)
	}
	diags = filter.Apply(diags)
	if hasRule(diags, rules.BodyLineLimit) {
		t.Fatal("body line limit should be ignored")
	}
}

func TestIgnoreUnknownRule(t *testing.T) {
	_, err := ignore.New([]config.IgnoreEntry{{
		Path:  "oversized/**",
		Rules: []string{"skill/does-not-exist"},
	}})
	if err == nil {
		t.Fatal("expected error for unknown rule id")
	}
}

func TestEstimateAndName(t *testing.T) {
	// Frontmatter name matching via full check on bad-name
	root := testdataRoot(t)
	diags, err := skill.Check(skill.Options{
		Config:     config.Default(),
		ScanRoot:   root,
		SkillsRoot: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasRule(diags, rules.NameDirectoryMismatch) {
		t.Fatal("expected name/directory mismatch for bad-name")
	}
}

func TestRovoOnlyValidated(t *testing.T) {
	root := testdataRoot(t)
	diags, err := skill.Check(skill.Options{
		Config:     config.Default(),
		ScanRoot:   root,
		SkillsRoot: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasRule(diags, rules.RovoSchemaInvalid) {
		t.Fatalf("expected rovo instructions missing from fixture; got %#v", diags)
	}
	if !hasRule(diags, rules.RovoSchemaInvalid) {
		t.Fatalf("expected unsupported apiVersion from fixture; got %#v", diags)
	}
	for _, d := range diags {
		if d.Path == "valid-rovo-v2/rovo-agent.yaml" && d.Severity == rules.SeverityError {
			t.Errorf("unexpected error on valid-rovo-v2: %s: %s", d.RuleID, d.Message)
		}
	}
}

func TestRovoRequireReadme(t *testing.T) {
	root := testdataRoot(t)
	cfg := config.Default()
	cfg.Skills.RequireReadme = true
	diags, err := skill.Check(skill.Options{
		Config:     cfg,
		ScanRoot:   root,
		SkillsRoot: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range diags {
		if d.RuleID == rules.ReadmeMissing && d.Path == "valid-rovo-v2" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected skill/readme-missing for valid-rovo-v2; got %#v", diags)
	}
}
