// Package rovo validates Atlassian Rovo / Studio agent definitions (rovo-agent.yaml)
// against the vendored agent-manager JSON Schema.
package rovo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ai-agent-manager/agent-lint/internal/report"
	"github.com/ai-agent-manager/agent-lint/internal/rules"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

var (
	schemaOnce sync.Once
	schemaObj  *jsonschema.Schema
	schemaErr  error
)

// Check validates a rovo-agent.yaml file. scanRoot is used to relativise paths in diagnostics.
func Check(rovoFile, scanRoot string) []report.Diagnostic {
	relPath := rel(scanRoot, rovoFile)
	content, err := os.ReadFile(rovoFile)
	if err != nil {
		return []report.Diagnostic{diag(rules.RovoYAMLInvalid, relPath, "cannot read rovo-agent.yaml: "+err.Error())}
	}

	var root any
	if err := yaml.Unmarshal(content, &root); err != nil {
		return []report.Diagnostic{diag(rules.RovoYAMLInvalid, relPath, "invalid YAML: "+err.Error())}
	}
	if root == nil {
		return []report.Diagnostic{diag(rules.RovoYAMLNotObject, relPath, "rovo-agent.yaml must parse to an object")}
	}
	data, ok := asStringKeyedMap(root)
	if !ok {
		return []report.Diagnostic{diag(rules.RovoYAMLNotObject, relPath, "rovo-agent.yaml must parse to an object")}
	}

	sch, err := compiledSchema()
	if err != nil {
		return []report.Diagnostic{diag(rules.RovoSchemaInvalid, relPath, "failed to load Rovo JSON Schema: "+err.Error())}
	}

	instance, err := jsonRoundTrip(data)
	if err != nil {
		return []report.Diagnostic{diag(rules.RovoYAMLInvalid, relPath, "cannot normalise YAML for schema validation: "+err.Error())}
	}

	if err := sch.Validate(instance); err != nil {
		return schemaDiagnostics(relPath, err)
	}
	return nil
}

func compiledSchema() (*jsonschema.Schema, error) {
	schemaOnce.Do(func() {
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
		if err != nil {
			schemaErr = fmt.Errorf("parse schema JSON: %w", err)
			return
		}
		c := jsonschema.NewCompiler()
		c.UseRegexpEngine(compileECMARegexp)
		const schemaURL = "rovo-agent.schema.json"
		if err := c.AddResource(schemaURL, doc); err != nil {
			schemaErr = fmt.Errorf("add schema resource: %w", err)
			return
		}
		sch, err := c.Compile(schemaURL)
		if err != nil {
			schemaErr = fmt.Errorf("compile schema: %w", err)
			return
		}
		schemaObj = sch
	})
	return schemaObj, schemaErr
}

func jsonRoundTrip(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return jsonschema.UnmarshalJSON(bytes.NewReader(b))
}

func schemaDiagnostics(relPath string, err error) []report.Diagnostic {
	ve, ok := err.(*jsonschema.ValidationError)
	if !ok {
		return []report.Diagnostic{diag(rules.RovoSchemaInvalid, relPath, err.Error())}
	}
	basic := ve.BasicOutput()
	var diags []report.Diagnostic
	for _, e := range basic.Errors {
		msg := ""
		if e.Error != nil {
			msg = e.Error.String()
		}
		if msg == "" || msg == "validation failed" {
			continue
		}
		loc := e.InstanceLocation
		if loc == "" {
			loc = "/"
		}
		diags = append(diags, diag(rules.RovoSchemaInvalid, relPath,
			fmt.Sprintf("%s: %s", loc, msg)))
	}
	if len(diags) == 0 {
		diags = append(diags, diag(rules.RovoSchemaInvalid, relPath, ve.Error()))
	}
	return diags
}

func asStringKeyedMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return normalizeMap(m), true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, val := range m {
			ks, ok := k.(string)
			if !ok {
				return nil, false
			}
			out[ks] = val
		}
		return normalizeMap(out), true
	default:
		return nil, false
	}
}

func normalizeMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = normalizeValue(v)
	}
	return out
}

func normalizeValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return normalizeMap(t)
	case map[any]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			ks, ok := k.(string)
			if !ok {
				return v
			}
			out[ks] = normalizeValue(val)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = normalizeValue(item)
		}
		return out
	default:
		return v
	}
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(r)
}

func diag(ruleID, path, message string) report.Diagnostic {
	r, ok := rules.Lookup(ruleID)
	sev := rules.SeverityError
	if ok {
		sev = r.Severity
	}
	return report.Diagnostic{
		RuleID:   ruleID,
		Severity: sev,
		Path:     path,
		Message:  message,
	}
}
