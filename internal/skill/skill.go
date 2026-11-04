// Package skill validates Agent Skills against agentskills.io.
package skill

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ai-agent-manager/agent-lint/internal/config"
	"github.com/ai-agent-manager/agent-lint/internal/frontmatter"
	"github.com/ai-agent-manager/agent-lint/internal/report"
	"github.com/ai-agent-manager/agent-lint/internal/rovo"
	"github.com/ai-agent-manager/agent-lint/internal/rules"
)

var (
	namePattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	// Markdown links: [text](path) — skip http(s) and anchors.
	mdLinkPattern = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)
	// Plain relative path mentions like references/foo.md or scripts/bar.py
	plainPathPattern = regexp.MustCompile(`(?m)(?:^|[\s` + "`" + `("'])((?:scripts|references|assets)/[A-Za-z0-9._/-]+)`)
)

// Options controls skill checking behaviour.
type Options struct {
	Config          config.Config
	Claude          bool
	AllowClaudeCode bool
	ScanRoot        string // absolute path used as relative root for reports
	SkillsRoot      string // absolute path to skills directory
}

// Check walks SkillsRoot and validates each child directory.
func Check(opts Options) ([]report.Diagnostic, error) {
	entries, err := os.ReadDir(opts.SkillsRoot)
	if err != nil {
		return nil, fmt.Errorf("read skills root: %w", err)
	}

	var diags []report.Diagnostic
	seenNames := map[string]string{} // name -> first path

	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		dirPath := filepath.Join(opts.SkillsRoot, ent.Name())
		diags = append(diags, checkDirectory(dirPath, ent.Name(), opts, seenNames)...)
	}

	return diags, nil
}

func checkDirectory(dirPath, dirName string, opts Options, seenNames map[string]string) []report.Diagnostic {
	var diags []report.Diagnostic
	relDir := rel(opts.ScanRoot, dirPath)

	skillFile := filepath.Join(dirPath, "SKILL.md")
	rovoFile := filepath.Join(dirPath, "rovo-agent.yaml")
	readmeFile := filepath.Join(dirPath, "README.md")

	hasSkill := fileExists(skillFile)
	hasRovo := fileExists(rovoFile)

	if hasSkill && hasRovo {
		diags = append(diags, diag(rules.AssetTypeConflict, relDir,
			"contains both SKILL.md and rovo-agent.yaml (must have only one)"))
		return diags
	}
	if !hasSkill && !hasRovo {
		// Check for wrong casing of skill.md
		if wrong := findWrongCasing(dirPath); wrong != "" {
			diags = append(diags, diag(rules.SkillMDCasing, rel(opts.ScanRoot, wrong),
				"found "+filepath.Base(wrong)+" but clients require exact SKILL.md casing"))
		}
		diags = append(diags, diag(rules.AssetTypeMissing, relDir,
			"missing both SKILL.md and rovo-agent.yaml (must have exactly one)"))
		return diags
	}
	if hasRovo {
		diags = append(diags, rovo.Check(rovoFile, opts.ScanRoot)...)
		if opts.Config.Skills.RequireReadme {
			if !fileExists(readmeFile) {
				diags = append(diags, diag(rules.ReadmeMissing, relDir, "missing README.md"))
			} else {
				diags = append(diags, validateReadme(readmeFile, opts)...)
			}
		}
		return diags
	}

	diags = append(diags, validateSkillFile(skillFile, dirName, opts, seenNames)...)

	if opts.Config.Skills.RequireReadme {
		if !fileExists(readmeFile) {
			diags = append(diags, diag(rules.ReadmeMissing, relDir, "missing README.md"))
		} else {
			diags = append(diags, validateReadme(readmeFile, opts)...)
		}
	}

	return diags
}

func validateSkillFile(skillFile, dirName string, opts Options, seenNames map[string]string) []report.Diagnostic {
	var diags []report.Diagnostic
	relPath := rel(opts.ScanRoot, skillFile)

	content, err := os.ReadFile(skillFile)
	if err != nil {
		diags = append(diags, diag(rules.MissingSkillMD, relPath, "cannot read SKILL.md: "+err.Error()))
		return diags
	}

	parsed, err := frontmatter.Parse(string(content))
	if err != nil {
		msg := err.Error()
		id := rules.FrontmatterMissing
		if strings.Contains(msg, "invalid YAML") {
			id = rules.FrontmatterInvalidYAML
		} else if strings.Contains(msg, "object") {
			id = rules.FrontmatterNotObject
		}
		diags = append(diags, diag(id, relPath, "SKILL.md "+msg))
		return diags
	}

	diags = append(diags, validateFrontmatter(parsed, relPath, dirName, opts, seenNames)...)
	diags = append(diags, validateBody(parsed.Body, skillFile, relPath, opts)...)
	diags = append(diags, validateReferences(parsed.Body, filepath.Dir(skillFile), relPath, opts)...)

	return diags
}

func validateFrontmatter(parsed *frontmatter.Parsed, relPath, dirName string, opts Options, seenNames map[string]string) []report.Diagnostic {
	var diags []report.Diagnostic
	data := parsed.Data

	allowed := map[string]struct{}{
		"name": {}, "description": {}, "license": {},
		"compatibility": {}, "metadata": {}, "allowed-tools": {},
	}
	allowClaude := opts.AllowClaudeCode || opts.Config.Skills.AllowClaudeCodeFields
	claudeFields := opts.Config.Skills.ClaudeCodeFields
	if allowClaude {
		for _, f := range claudeFields {
			allowed[f] = struct{}{}
		}
	}
	for _, c := range opts.Config.Skills.SkillConstraints {
		if f := strings.TrimSpace(c.Field); f != "" {
			allowed[f] = struct{}{}
		}
	}

	for key := range data {
		if _, ok := allowed[key]; !ok {
			diags = append(diags, diag(rules.FrontmatterUnknownField, relPath,
				fmt.Sprintf("unknown frontmatter field %q (agentskills.io allowlist only)", key)))
		}
	}

	name, ok := frontmatter.StringField(data, "name")
	if !ok {
		diags = append(diags, diag(rules.NameMissing, relPath, "frontmatter must include a non-empty name"))
	} else {
		if utf8.RuneCountInString(name) > 64 || !isValidName(name) {
			diags = append(diags, diag(rules.NameFormat, relPath,
				"name must be 1-64 chars, lowercase alphanumeric and hyphens, no leading/trailing/consecutive hyphens"))
		}
		if name != dirName {
			diags = append(diags, diag(rules.NameDirectoryMismatch, relPath,
				fmt.Sprintf("name %q must match parent directory name %q", name, dirName)))
		}
		if prev, exists := seenNames[name]; exists {
			diags = append(diags, diag(rules.NameDuplicate, relPath,
				fmt.Sprintf("duplicate name %q (also declared in %s)", name, prev)))
		} else {
			seenNames[name] = relPath
		}
		if opts.Claude {
			lower := strings.ToLower(name)
			if strings.Contains(lower, "anthropic") || strings.Contains(lower, "claude") {
				diags = append(diags, diag(rules.ClaudeReservedName, relPath,
					"Claude.ai: name must not contain anthropic or claude"))
			}
		}
	}

	desc, ok := frontmatter.StringField(data, "description")
	if !ok {
		diags = append(diags, diag(rules.DescriptionMissing, relPath, "frontmatter must include a non-empty description"))
	} else {
		if utf8.RuneCountInString(desc) > 1024 {
			diags = append(diags, diag(rules.DescriptionLength, relPath,
				fmt.Sprintf("description exceeds agentskills.io max of 1024 characters (%d)", utf8.RuneCountInString(desc))))
		}
		if opts.Claude && (strings.Contains(desc, "<") || strings.Contains(desc, ">")) {
			diags = append(diags, diag(rules.ClaudeAngleBrackets, relPath,
				"Claude.ai: description must not contain < or >"))
		}
	}

	if v, exists := data["compatibility"]; exists && v != nil {
		s, ok := v.(string)
		if !ok || strings.TrimSpace(s) == "" || utf8.RuneCountInString(s) > 500 {
			diags = append(diags, diag(rules.CompatibilityLength, relPath,
				"compatibility must be a non-empty string of at most 500 characters"))
		}
	}

	if v, exists := data["metadata"]; exists && v != nil {
		m, ok := v.(map[string]any)
		if !ok {
			diags = append(diags, diag(rules.MetadataFormat, relPath,
				"metadata must be a mapping of string keys to string values"))
		} else {
			for k, val := range m {
				if _, ok := val.(string); !ok {
					diags = append(diags, diag(rules.MetadataFormat, relPath,
						fmt.Sprintf("metadata.%s must be a string", k)))
				}
			}
		}
	}

	if v, exists := data["allowed-tools"]; exists && v != nil {
		s, ok := v.(string)
		if !ok {
			diags = append(diags, diag(rules.AllowedToolsFormat, relPath,
				"allowed-tools must be a space-separated string"))
		} else if strings.Contains(s, ",") {
			diags = append(diags, diag(rules.AllowedToolsFormat, relPath,
				"allowed-tools must be space-separated (comma-separated lists are not allowed)"))
		}
	}

	if v, exists := data["license"]; exists && v != nil {
		if _, ok := v.(string); !ok {
			diags = append(diags, diag(rules.FrontmatterUnknownField, relPath, "license must be a string"))
		}
	}

	for _, c := range opts.Config.Skills.SkillConstraints {
		diags = append(diags, validateFrontmatterConstraint(data, relPath, "SKILL.md", c, fieldRuleIDs{
			invalid:       rules.SkillFieldInvalid,
			requiredValue: rules.SkillFieldRequiredValue,
			allowedValue:  rules.SkillFieldAllowedValue,
		})...)
	}

	return diags
}

func validateBody(body, skillFile, relPath string, opts Options) []report.Diagnostic {
	var diags []report.Diagnostic
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		diags = append(diags, diag(rules.BodyEmpty, relPath, "SKILL.md has no instructions after the frontmatter"))
		return diags
	}

	lines := countLines(body)
	maxLines := opts.Config.Skills.MaxBodyLines
	if lines > maxLines {
		diags = append(diags, diag(rules.BodyLineLimit, relPath,
			fmt.Sprintf("SKILL.md body has %d lines; exceeds agentskills.io recommended max of %d lines — move detail to references/", lines, maxLines)))
	}

	tokens := estimateTokens(body)
	maxTokens := opts.Config.Skills.MaxBodyTokens
	if tokens >= maxTokens {
		diags = append(diags, diag(rules.BodyTokenLimit, relPath,
			fmt.Sprintf("SKILL.md body is ~%d tokens; exceeds agentskills.io recommended max of %d tokens — move detail to references/", tokens, maxTokens)))
	}

	_ = skillFile
	return diags
}

func validateReferences(body, skillDir, relPath string, opts Options) []report.Diagnostic {
	var diags []report.Diagnostic
	linkRefs, plainRefs := extractReferences(body)

	mentioned := map[string]struct{}{}
	checkRef := func(ref string, requireExists bool) {
		clean := strings.SplitN(ref, "#", 2)[0]
		clean = strings.TrimSpace(clean)
		if clean == "" || isExternal(clean) || isPlaceholderPath(clean) {
			return
		}
		clean = filepath.ToSlash(clean)
		clean = strings.TrimSuffix(clean, "/")
		if strings.HasPrefix(clean, "/") || strings.Contains(clean, "..") {
			diags = append(diags, diag(rules.ReferenceDepth, relPath,
				fmt.Sprintf("reference %q must be a relative path from the skill root (no .. or absolute paths)", ref)))
			return
		}
		abs := filepath.Join(skillDir, filepath.FromSlash(clean))
		mentioned[filepath.Clean(abs)] = struct{}{}
		exists := pathExists(abs)
		// Depth limits apply to real skill-local references (links, or plain paths that exist).
		if (requireExists || exists) && referenceDepth(clean) > 1 {
			diags = append(diags, diag(rules.ReferenceDepth, relPath,
				fmt.Sprintf("reference %q is %d levels deep; agentskills.io requires at most one directory level from SKILL.md", clean, referenceDepth(clean))))
		}
		if requireExists && !exists {
			diags = append(diags, diag(rules.ReferenceMissing, relPath,
				fmt.Sprintf("referenced path %q does not exist", clean)))
		}
	}

	// Markdown links are authoritative references and must exist.
	for _, ref := range linkRefs {
		checkRef(ref, true)
	}
	// Plain path mentions only contribute to orphan detection (examples often use fictional paths).
	for _, ref := range plainRefs {
		checkRef(ref, false)
	}

	refsDir := filepath.Join(skillDir, "references")
	if info, err := os.Stat(refsDir); err == nil && info.IsDir() {
		_ = filepath.WalkDir(refsDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if _, ok := mentioned[filepath.Clean(path)]; !ok {
				// Also check if mentioned via relative path string search in body
				relRef, _ := filepath.Rel(skillDir, path)
				relRef = filepath.ToSlash(relRef)
				if !strings.Contains(body, relRef) && !strings.Contains(body, filepath.Base(path)) {
					diags = append(diags, diag(rules.ReferenceOrphan, rel(opts.ScanRoot, path),
						fmt.Sprintf("references file %q is never mentioned from SKILL.md", relRef)))
				}
			}
			return nil
		})
	}

	return diags
}

func validateReadme(readmeFile string, opts Options) []report.Diagnostic {
	var diags []report.Diagnostic
	relPath := rel(opts.ScanRoot, readmeFile)
	content, err := os.ReadFile(readmeFile)
	if err != nil {
		diags = append(diags, diag(rules.ReadmeFrontmatterInvalid, relPath, "cannot read README.md: "+err.Error()))
		return diags
	}
	parsed, err := frontmatter.Parse(string(content))
	if err != nil {
		diags = append(diags, diag(rules.ReadmeFrontmatterInvalid, relPath, "README.md "+err.Error()))
		return diags
	}
	if _, ok := frontmatter.StringField(parsed.Data, "name"); !ok {
		diags = append(diags, diag(rules.ReadmeNameMissing, relPath, "README.md frontmatter must include a non-empty name"))
	}
	if _, ok := frontmatter.StringField(parsed.Data, "description"); !ok {
		diags = append(diags, diag(rules.ReadmeDescriptionMissing, relPath, "README.md frontmatter must include a non-empty description"))
	}

	for _, c := range opts.Config.Skills.ReadmeConstraints {
		diags = append(diags, validateFrontmatterConstraint(parsed.Data, relPath, "README.md", c, fieldRuleIDs{
			invalid:       rules.ReadmeFieldInvalid,
			requiredValue: rules.ReadmeFieldRequiredValue,
			allowedValue:  rules.ReadmeFieldAllowedValue,
		})...)
	}

	return diags
}

type fieldRuleIDs struct {
	invalid       string
	requiredValue string
	allowedValue  string
}

func validateFrontmatterConstraint(data map[string]any, relPath, fileLabel string, c config.FrontmatterConstraint, ids fieldRuleIDs) []report.Diagnostic {
	var diags []report.Diagnostic
	values, ok := frontmatter.StringSliceField(data, c.Field)
	if !ok || len(values) == 0 {
		return []report.Diagnostic{diag(ids.invalid, relPath,
			fmt.Sprintf("%s %q must be a non-empty array of strings", fileLabel, c.Field))}
	}

	for _, req := range c.MustInclude {
		if !contains(values, req) {
			diags = append(diags, diag(ids.requiredValue, relPath,
				fmt.Sprintf("%s %q must include %q", fileLabel, c.Field, req)))
		}
	}

	if len(c.Allowed) > 0 {
		allowed := map[string]struct{}{}
		for _, a := range c.Allowed {
			allowed[a] = struct{}{}
		}
		for _, v := range values {
			if _, ok := allowed[v]; !ok {
				diags = append(diags, diag(ids.allowedValue, relPath,
					fmt.Sprintf("%s %q value %q is not in the configured allowlist", fileLabel, c.Field, v)))
			}
		}
	}

	return diags
}

func isValidName(name string) bool {
	if name == "" || utf8.RuneCountInString(name) > 64 {
		return false
	}
	for _, r := range name {
		if r > unicode.MaxASCII {
			// Spec allows unicode lowercase alphanumeric; keep ASCII subset tight + allow unicode letters/digits via pattern alternative
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
				return false
			}
			if unicode.IsLetter(r) && !unicode.IsLower(r) {
				return false
			}
		}
	}
	// Prefer the ASCII agentskills.io pattern; for unicode, fall back to hyphen rules
	if isASCII(name) {
		return namePattern.MatchString(name)
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") || strings.Contains(name, "--") {
		return false
	}
	return true
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func estimateTokens(s string) int {
	// Documented estimator: ceil(rune_count / 4), matching common rough token heuristics.
	n := utf8.RuneCountInString(s)
	if n == 0 {
		return 0
	}
	return (n + 3) / 4
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n") + 1
	if strings.HasSuffix(s, "\n") {
		n--
	}
	return n
}

func extractReferences(body string) (links []string, plain []string) {
	seenLink := map[string]struct{}{}
	seenPlain := map[string]struct{}{}
	add := func(dst *[]string, seen map[string]struct{}, p string) {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"'")
		if p == "" {
			return
		}
		if _, ok := seen[p]; ok {
			return
		}
		seen[p] = struct{}{}
		*dst = append(*dst, p)
	}
	for _, m := range mdLinkPattern.FindAllStringSubmatch(body, -1) {
		add(&links, seenLink, m[1])
	}
	for _, m := range plainPathPattern.FindAllStringSubmatch(body, -1) {
		add(&plain, seenPlain, m[1])
	}
	return links, plain
}

func isPlaceholderPath(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	switch lower {
	case "url", "uri", "path", "file", "filename", "example", "link":
		return true
	}
	if strings.Contains(path, "*") {
		return true
	}
	return false
}

func referenceDepth(path string) int {
	path = filepath.ToSlash(path)
	parts := strings.Split(path, "/")
	// "references/foo.md" => depth 1; "references/sub/foo.md" => depth 2
	if len(parts) <= 1 {
		return 0
	}
	return len(parts) - 1
}

func isExternal(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "mailto:") ||
		strings.HasPrefix(path, "#")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func findWrongCasing(dirPath string) string {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(e.Name(), "SKILL.md") && e.Name() != "SKILL.md" {
			return filepath.Join(dirPath, e.Name())
		}
	}
	return ""
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(r)
}

func contains(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
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
