package frontmatter_test

import (
	"testing"

	"github.com/ai-agent-manager/agent-lint/internal/frontmatter"
)

func TestParseOK(t *testing.T) {
	content := "---\nname: foo\ndescription: bar\n---\n\n# Body\n"
	p, err := frontmatter.Parse(content)
	if err != nil {
		t.Fatal(err)
	}
	name, ok := frontmatter.StringField(p.Data, "name")
	if !ok || name != "foo" {
		t.Fatalf("name=%q ok=%v", name, ok)
	}
	if p.Body != "# Body\n" && p.Body != "# Body" {
		t.Fatalf("body=%q", p.Body)
	}
}

func TestParseMissing(t *testing.T) {
	_, err := frontmatter.Parse("# no frontmatter\n")
	if err == nil {
		t.Fatal("expected error")
	}
}
