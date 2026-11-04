// Package config loads agent-lint configuration.
package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the top-level agent-lint configuration.
type Config struct {
	Skills SkillsConfig `yaml:"skills"`
	Ignore IgnoreConfig `yaml:"ignore"`
}

// SkillsConfig holds skill-checker settings.
type SkillsConfig struct {
	MaxBodyLines          int                     `yaml:"max_body_lines"`
	MaxBodyTokens         int                     `yaml:"max_body_tokens"`
	RequireReadme         bool                    `yaml:"require_readme"`
	ReadmeConstraints     []FrontmatterConstraint `yaml:"readme_constraints"`
	SkillConstraints      []FrontmatterConstraint `yaml:"skill_constraints"`
	AllowClaudeCodeFields bool                    `yaml:"allow_claude_code_fields"`
	ClaudeCodeFields      []string                `yaml:"claude_code_fields"`
}

// FrontmatterConstraint validates a string-list field in README.md or SKILL.md frontmatter.
type FrontmatterConstraint struct {
	Field       string   `yaml:"field"`
	MustInclude []string `yaml:"must_include"`
	Allowed     []string `yaml:"allowed"`
}

// IgnoreConfig embeds ignore entries in the main config file.
type IgnoreConfig struct {
	Ignores []IgnoreEntry `yaml:"ignores"`
}

// IgnoreEntry suppresses findings for a path/glob.
type IgnoreEntry struct {
	Path  string   `yaml:"path"`
	Rules []string `yaml:"rules"`
}

// IgnoreFile is the standalone .agent-lint-ignore.yaml schema.
type IgnoreFile struct {
	Ignores []IgnoreEntry `yaml:"ignores"`
}

// Default returns configuration with agentskills.io progressive-disclosure defaults.
func Default() Config {
	return Config{
		Skills: SkillsConfig{
			MaxBodyLines:  500,
			MaxBodyTokens: 5000,
			ClaudeCodeFields: []string{
				"when_to_use",
				"argument-hint",
				"arguments",
				"disable-model-invocation",
				"user-invocable",
				"disallowed-tools",
				"model",
				"effort",
				"context",
				"agent",
				"background",
				"hooks",
				"paths",
				"shell",
			},
		},
	}
}

// Load reads a config YAML file. Missing file returns defaults.
func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	if cfg.Skills.MaxBodyLines <= 0 {
		cfg.Skills.MaxBodyLines = 500
	}
	if cfg.Skills.MaxBodyTokens <= 0 {
		cfg.Skills.MaxBodyTokens = 5000
	}
	if err := validateConstraints("skills.readme_constraints", cfg.Skills.ReadmeConstraints); err != nil {
		return Config{}, fmt.Errorf("config %s: %w", path, err)
	}
	if err := validateConstraints("skills.skill_constraints", cfg.Skills.SkillConstraints); err != nil {
		return Config{}, fmt.Errorf("config %s: %w", path, err)
	}
	return cfg, nil
}

func validateConstraints(path string, constraints []FrontmatterConstraint) error {
	for i, c := range constraints {
		field := strings.TrimSpace(c.Field)
		if field == "" {
			return fmt.Errorf("%s[%d]: field is required", path, i)
		}
		if len(c.MustInclude) == 0 && len(c.Allowed) == 0 {
			return fmt.Errorf("%s[%d] (field %q): must set must_include and/or allowed", path, i, field)
		}
	}
	return nil
}

// LoadIgnoreFile reads .agent-lint-ignore.yaml. Missing file is OK.
func LoadIgnoreFile(path string) ([]IgnoreEntry, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var f IgnoreFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse ignore file %s: %w", path, err)
	}
	return f.Ignores, nil
}
