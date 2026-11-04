// Package frontmatter parses YAML frontmatter from markdown files.
package frontmatter

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parsed is the result of splitting frontmatter from a markdown body.
type Parsed struct {
	Data map[string]any
	Body string
	Raw  string
}

// Parse extracts YAML frontmatter delimited by --- lines.
func Parse(content string) (*Parsed, error) {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return nil, fmt.Errorf("must start with YAML frontmatter delimited by ---")
	}

	rest := content[3:]
	if strings.HasPrefix(rest, "\r\n") {
		rest = rest[2:]
	} else if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	} else {
		return nil, fmt.Errorf("must start with YAML frontmatter delimited by ---")
	}

	closing := "\n---\n"
	closingCRLF := "\n---\r\n"
	closingEOF := "\n---"
	idx := strings.Index(rest, closing)
	endLen := len(closing)
	if idx == -1 {
		if i := strings.Index(rest, closingCRLF); i != -1 {
			idx = i
			endLen = len(closingCRLF)
		} else if strings.HasSuffix(rest, closingEOF) {
			idx = len(rest) - len(closingEOF)
			endLen = len(closingEOF)
		}
	}
	if idx == -1 {
		return nil, fmt.Errorf("frontmatter must end with a closing --- line")
	}

	raw := rest[:idx]
	body := rest[idx+endLen:]
	body = strings.TrimPrefix(body, "\r\n")
	body = strings.TrimPrefix(body, "\n")

	var data map[string]any
	if err := yaml.Unmarshal([]byte(raw), &data); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	if data == nil {
		data = map[string]any{}
	}

	return &Parsed{Data: data, Body: body, Raw: raw}, nil
}

// StringField returns a trimmed string field, or ok=false if missing/empty/non-string.
func StringField(data map[string]any, key string) (string, bool) {
	v, exists := data[key]
	if !exists || v == nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	return s, true
}

// StringSliceField returns a string slice field.
func StringSliceField(data map[string]any, key string) ([]string, bool) {
	v, exists := data[key]
	if !exists || v == nil {
		return nil, false
	}
	arr, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		if strings.TrimSpace(s) == "" {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}
