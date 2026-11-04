package rovo_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ai-agent-manager/agent-lint/internal/report"
	"github.com/ai-agent-manager/agent-lint/internal/rovo"
	"github.com/ai-agent-manager/agent-lint/internal/rules"
)

func hasRule(diags []report.Diagnostic, id string) bool {
	for _, d := range diags {
		if d.RuleID == id {
			return true
		}
	}
	return false
}

func writeTemp(t *testing.T, content string) (file, root string) {
	t.Helper()
	root = t.TempDir()
	file = filepath.Join(root, "rovo-agent.yaml")
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return file, root
}

func TestValidV2Beta(t *testing.T) {
	file, root := writeTemp(t, `
apiVersion: "rovo.atlassian.com/v2-beta"
kind: "StudioAgent"
name: "Ok Agent"
description: "A valid agent"
instructions: "Do the thing."
`)
	diags := rovo.Check(file, root)
	for _, d := range diags {
		if d.Severity == rules.SeverityError {
			t.Fatalf("unexpected error: %s: %s", d.RuleID, d.Message)
		}
	}
}

func TestMissingInstructions(t *testing.T) {
	file, root := writeTemp(t, `
apiVersion: "rovo.atlassian.com/v2-beta"
kind: "StudioAgent"
name: "Ok Agent"
description: "A valid agent"
`)
	diags := rovo.Check(file, root)
	if !hasRule(diags, rules.RovoSchemaInvalid) {
		t.Fatalf("expected %s; got %#v", rules.RovoSchemaInvalid, diags)
	}
}

func TestArchivedV1(t *testing.T) {
	file, root := writeTemp(t, `
apiVersion: "rovo.atlassian.com/v1"
kind: "StudioAgent"
identity:
  name: "Ok"
  description: "Ok"
`)
	diags := rovo.Check(file, root)
	if !hasRule(diags, rules.RovoSchemaInvalid) {
		t.Fatalf("expected %s for archived v1; got %#v", rules.RovoSchemaInvalid, diags)
	}
}

func TestNameTooLong(t *testing.T) {
	file, root := writeTemp(t, `
apiVersion: "rovo.atlassian.com/v2-beta"
kind: "StudioAgent"
name: "This name is way too long for Studio"
description: "A valid agent"
instructions: "Do the thing."
`)
	diags := rovo.Check(file, root)
	if !hasRule(diags, rules.RovoSchemaInvalid) {
		t.Fatalf("expected %s; got %#v", rules.RovoSchemaInvalid, diags)
	}
	found := false
	for _, d := range diags {
		if d.RuleID == rules.RovoSchemaInvalid && strings.Contains(d.Message, "/name") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected /name in schema error; got %#v", diags)
	}
}

func TestInvalidYAML(t *testing.T) {
	file, root := writeTemp(t, ":\n  - bad\n")
	diags := rovo.Check(file, root)
	if !hasRule(diags, rules.RovoYAMLInvalid) {
		t.Fatalf("expected %s; got %#v", rules.RovoYAMLInvalid, diags)
	}
}

func TestYAMLNotObject(t *testing.T) {
	file, root := writeTemp(t, "- just\n- a\n- list\n")
	diags := rovo.Check(file, root)
	if !hasRule(diags, rules.RovoYAMLNotObject) {
		t.Fatalf("expected %s; got %#v", rules.RovoYAMLNotObject, diags)
	}
}

func TestTooManyConversationStarters(t *testing.T) {
	file, root := writeTemp(t, `
apiVersion: "rovo.atlassian.com/v2-beta"
kind: "StudioAgent"
name: "Ok Agent"
description: "A valid agent"
instructions: "Do the thing."
conversationStarters:
  - "one"
  - "two"
  - "three"
  - "four"
`)
	diags := rovo.Check(file, root)
	if !hasRule(diags, rules.RovoSchemaInvalid) {
		t.Fatalf("expected %s; got %#v", rules.RovoSchemaInvalid, diags)
	}
}

func TestValidSubagent(t *testing.T) {
	file, root := writeTemp(t, `
apiVersion: "rovo.atlassian.com/v2-beta"
kind: "StudioAgent"
name: "Ok Agent"
description: "A valid agent"
instructions: "Do the thing."
subagents:
  helper:
    name: "Helper"
    enabled: true
    trigger: "When asked"
    instructions: "Help"
    conversationStarters: ["Hi"]
    skills: ["a"]
    knowledge: "all"
    webSearch: false
    deepResearch: true
`)
	diags := rovo.Check(file, root)
	for _, d := range diags {
		if d.Severity == rules.SeverityError {
			t.Fatalf("unexpected error: %s: %s", d.RuleID, d.Message)
		}
	}
}

func TestInvalidSubagent(t *testing.T) {
	file, root := writeTemp(t, `
apiVersion: "rovo.atlassian.com/v2-beta"
kind: "StudioAgent"
name: "Ok Agent"
description: "A valid agent"
instructions: "Do the thing."
subagents:
  helper:
    name: "Helper"
    enabled: "yes"
`)
	diags := rovo.Check(file, root)
	if !hasRule(diags, rules.RovoSchemaInvalid) {
		t.Fatalf("expected %s; got %#v", rules.RovoSchemaInvalid, diags)
	}
}
