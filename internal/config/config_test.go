package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ai-agent-manager/agent-lint/internal/config"
)

func TestLoadReadmeConstraintsValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".agent-lint.yaml")
	content := `
skills:
  require_readme: true
  readme_constraints:
    - field: tags
      must_include: [agent-skill]
    - field: phases
      allowed: [build_and_release]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Skills.ReadmeConstraints) != 2 {
		t.Fatalf("expected 2 constraints, got %d", len(cfg.Skills.ReadmeConstraints))
	}
	if cfg.Skills.ReadmeConstraints[0].Field != "tags" {
		t.Fatalf("unexpected first field: %q", cfg.Skills.ReadmeConstraints[0].Field)
	}
}

func TestLoadReadmeConstraintsEmptyEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".agent-lint.yaml")
	content := `
skills:
  readme_constraints:
    - field: tags
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for empty constraint")
	}
	if !strings.Contains(err.Error(), "must_include") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadReadmeConstraintsMissingField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".agent-lint.yaml")
	content := `
skills:
  readme_constraints:
    - must_include: [agent-skill]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing field")
	}
	if !strings.Contains(err.Error(), "field is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadSkillConstraintsValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".agent-lint.yaml")
	content := `
skills:
  skill_constraints:
    - field: phases
      allowed: [build_and_release]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Skills.SkillConstraints) != 1 || cfg.Skills.SkillConstraints[0].Field != "phases" {
		t.Fatalf("unexpected skill constraints: %#v", cfg.Skills.SkillConstraints)
	}
}

func TestLoadSkillConstraintsEmptyEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".agent-lint.yaml")
	content := `
skills:
  skill_constraints:
    - field: phases
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for empty skill constraint")
	}
	if !strings.Contains(err.Error(), "skill_constraints") {
		t.Fatalf("unexpected error: %v", err)
	}
}
