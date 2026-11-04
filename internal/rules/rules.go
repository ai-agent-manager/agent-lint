// Package rules defines stable diagnostic IDs for agent-lint.
package rules

// Severity levels for diagnostics.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Rule describes a single lint rule.
type Rule struct {
	ID          string
	Severity    Severity
	Description string
}

// Stable rule IDs (namespace/name).
const (
	MissingSkillMD           = "skill/missing-skill-md"
	SkillMDCasing            = "skill/skill-md-casing"
	FrontmatterMissing       = "skill/frontmatter-missing"
	FrontmatterInvalidYAML   = "skill/frontmatter-invalid-yaml"
	FrontmatterNotObject     = "skill/frontmatter-not-object"
	NameMissing              = "skill/name-missing"
	NameFormat               = "skill/name-format"
	NameDirectoryMismatch    = "skill/name-directory-mismatch"
	NameDuplicate            = "skill/name-duplicate"
	DescriptionMissing       = "skill/description-missing"
	DescriptionLength        = "skill/description-length"
	CompatibilityLength      = "skill/compatibility-length"
	MetadataFormat           = "skill/metadata-format"
	AllowedToolsFormat       = "skill/allowed-tools-format"
	FrontmatterUnknownField  = "skill/frontmatter-unknown-field"
	BodyEmpty                = "skill/body-empty"
	BodyLineLimit            = "skill/body-line-limit"
	BodyTokenLimit           = "skill/body-token-limit"
	ReferenceDepth           = "skill/reference-depth"
	ReferenceMissing         = "skill/reference-missing"
	ReferenceOrphan          = "skill/reference-orphan"
	AssetTypeConflict        = "skill/asset-type-conflict"
	AssetTypeMissing         = "skill/asset-type-missing"
	ReadmeMissing            = "skill/readme-missing"
	ReadmeFrontmatterInvalid = "skill/readme-frontmatter-invalid"
	ReadmeNameMissing        = "skill/readme-name-missing"
	ReadmeDescriptionMissing = "skill/readme-description-missing"
	ReadmeFieldInvalid       = "skill/readme-field-invalid"
	ReadmeFieldEquals        = "skill/readme-field-equals"
	ReadmeFieldNotEquals     = "skill/readme-field-not-equals"
	ReadmeFieldMatches       = "skill/readme-field-matches"
	ReadmeFieldNotMatches    = "skill/readme-field-not-matches"
	SkillFieldInvalid        = "skill/skill-field-invalid"
	SkillFieldEquals         = "skill/skill-field-equals"
	SkillFieldNotEquals      = "skill/skill-field-not-equals"
	SkillFieldMatches        = "skill/skill-field-matches"
	SkillFieldNotMatches     = "skill/skill-field-not-matches"
	ClaudeReservedName       = "skill/claude-reserved-name"
	ClaudeAngleBrackets      = "skill/claude-angle-brackets"

	RovoYAMLInvalid   = "rovo/yaml-invalid"
	RovoYAMLNotObject = "rovo/yaml-not-object"
	RovoSchemaInvalid = "rovo/schema-invalid"
)

// All returns every known rule.
func All() []Rule {
	return []Rule{
		{MissingSkillMD, SeverityError, "Skill directory must contain SKILL.md"},
		{SkillMDCasing, SeverityWarning, "SKILL.md filename casing should be exact"},
		{FrontmatterMissing, SeverityError, "SKILL.md must start with YAML frontmatter delimited by ---"},
		{FrontmatterInvalidYAML, SeverityError, "SKILL.md frontmatter must be valid YAML"},
		{FrontmatterNotObject, SeverityError, "SKILL.md frontmatter must parse to a YAML object"},
		{NameMissing, SeverityError, "Frontmatter must include a non-empty name"},
		{NameFormat, SeverityError, "name must be 1-64 chars, lowercase alphanumeric and hyphens, no leading/trailing/consecutive hyphens"},
		{NameDirectoryMismatch, SeverityError, "name must match the parent directory name"},
		{NameDuplicate, SeverityError, "Duplicate skill name within the scan root"},
		{DescriptionMissing, SeverityError, "Frontmatter must include a non-empty description"},
		{DescriptionLength, SeverityError, "description must be 1-1024 characters"},
		{CompatibilityLength, SeverityError, "compatibility must be 1-500 characters when present"},
		{MetadataFormat, SeverityError, "metadata must be a mapping of string keys to string values"},
		{AllowedToolsFormat, SeverityError, "allowed-tools must be a space-separated string"},
		{FrontmatterUnknownField, SeverityError, "Unknown frontmatter field (agentskills.io allowlist only)"},
		{BodyEmpty, SeverityWarning, "SKILL.md has no instructions after the frontmatter"},
		{BodyLineLimit, SeverityError, "SKILL.md body exceeds agentskills.io max of 500 lines"},
		{BodyTokenLimit, SeverityError, "SKILL.md body exceeds agentskills.io max of 5000 tokens"},
		{ReferenceDepth, SeverityError, "File references must be at most one directory level deep from SKILL.md"},
		{ReferenceMissing, SeverityError, "Referenced file path does not exist on disk"},
		{ReferenceOrphan, SeverityWarning, "File under references/ is never mentioned from SKILL.md"},
		{AssetTypeConflict, SeverityError, "Directory contains both SKILL.md and rovo-agent.yaml"},
		{AssetTypeMissing, SeverityError, "Directory missing both SKILL.md and rovo-agent.yaml"},
		{ReadmeMissing, SeverityError, "README.md is required by config"},
		{ReadmeFrontmatterInvalid, SeverityError, "README.md frontmatter is invalid"},
		{ReadmeNameMissing, SeverityError, "README.md frontmatter must include a non-empty name"},
		{ReadmeDescriptionMissing, SeverityError, "README.md frontmatter must include a non-empty description"},
		{ReadmeFieldInvalid, SeverityError, "README.md constrained field path must resolve to a string"},
		{ReadmeFieldEquals, SeverityError, "README.md field does not equal the configured value"},
		{ReadmeFieldNotEquals, SeverityError, "README.md field equals a forbidden configured value"},
		{ReadmeFieldMatches, SeverityError, "README.md field does not match the configured regexp"},
		{ReadmeFieldNotMatches, SeverityError, "README.md field matches a forbidden configured regexp"},
		{SkillFieldInvalid, SeverityError, "SKILL.md constrained field path must resolve to a string"},
		{SkillFieldEquals, SeverityError, "SKILL.md field does not equal the configured value"},
		{SkillFieldNotEquals, SeverityError, "SKILL.md field equals a forbidden configured value"},
		{SkillFieldMatches, SeverityError, "SKILL.md field does not match the configured regexp"},
		{SkillFieldNotMatches, SeverityError, "SKILL.md field matches a forbidden configured regexp"},
		{ClaudeReservedName, SeverityError, "Claude.ai: name must not contain anthropic or claude"},
		{ClaudeAngleBrackets, SeverityError, "Claude.ai: description must not contain < or >"},
		{RovoYAMLInvalid, SeverityError, "rovo-agent.yaml must be valid YAML"},
		{RovoYAMLNotObject, SeverityError, "rovo-agent.yaml must parse to a YAML object"},
		{RovoSchemaInvalid, SeverityError, "rovo-agent.yaml does not conform to the vendored Rovo agent JSON Schema"},
	}
}

// IDs returns the set of all known rule IDs.
func IDs() map[string]struct{} {
	out := make(map[string]struct{}, len(All()))
	for _, r := range All() {
		out[r.ID] = struct{}{}
	}
	return out
}

// Lookup returns a rule by ID.
func Lookup(id string) (Rule, bool) {
	for _, r := range All() {
		if r.ID == id {
			return r, true
		}
	}
	return Rule{}, false
}
