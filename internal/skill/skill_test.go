package skill_test

import (
	"path/filepath"
	"runtime"
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

func TestValidSkill(t *testing.T) {
	root := testdataRoot(t)
	cfg := config.Default()
	cfg.Skills.RequireReadme = true
	cfg.Skills.RequiredTags = []string{"agent-skill"}
	cfg.Skills.AllowedPhases = []string{"build_and_release"}

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
