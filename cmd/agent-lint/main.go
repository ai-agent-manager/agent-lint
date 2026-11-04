package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ai-agent-manager/agent-lint/internal/config"
	"github.com/ai-agent-manager/agent-lint/internal/ignore"
	"github.com/ai-agent-manager/agent-lint/internal/report"
	"github.com/ai-agent-manager/agent-lint/internal/rules"
	"github.com/ai-agent-manager/agent-lint/internal/skill"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	root := &cobra.Command{
		Use:           "agent-lint",
		Short:         "Lint Agent Skills against the agentskills.io specification",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.Version = fmt.Sprintf("%s (commit=%s date=%s)", version, commit, date)

	root.AddCommand(newCheckCmd())
	root.AddCommand(newListRulesCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func newCheckCmd() *cobra.Command {
	var (
		rootDir     string
		configPath  string
		ignorePath  string
		jsonOut     bool
		quiet       bool
		strict      bool
		claude      bool
		allowClaude bool
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Validate skill directories against agentskills.io",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if configPath == "" {
				if _, err := os.Stat(filepath.Join(cwd, ".agent-lint.yaml")); err == nil {
					configPath = filepath.Join(cwd, ".agent-lint.yaml")
				}
			}
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}

			if ignorePath == "" {
				candidate := filepath.Join(cwd, ".agent-lint-ignore.yaml")
				if _, err := os.Stat(candidate); err == nil {
					ignorePath = candidate
				}
			}
			ignoreEntries, err := config.LoadIgnoreFile(ignorePath)
			if err != nil {
				return err
			}
			ignoreEntries = append(cfg.Ignore.Ignores, ignoreEntries...)

			filter, err := ignore.New(ignoreEntries)
			if err != nil {
				return err
			}

			skillsRoot := rootDir
			if !filepath.IsAbs(skillsRoot) {
				skillsRoot = filepath.Join(cwd, skillsRoot)
			}
			if _, err := os.Stat(skillsRoot); err != nil {
				return fmt.Errorf("skills root %q: %w", skillsRoot, err)
			}

			diags, err := skill.Check(skill.Options{
				Config:          cfg,
				Claude:          claude,
				AllowClaudeCode: allowClaude || cfg.Skills.AllowClaudeCodeFields,
				ScanRoot:        cwd,
				SkillsRoot:      skillsRoot,
			})
			if err != nil {
				return err
			}

			diags = filter.Apply(diags)
			result := report.NewResult(diags)

			if jsonOut {
				if err := report.WriteJSON(os.Stdout, result); err != nil {
					return err
				}
			} else {
				report.WriteHuman(os.Stdout, result, quiet)
			}

			code := result.ExitCode(strict)
			if code != 0 {
				os.Exit(code)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&rootDir, "root", "skills", "Directory containing skill folders")
	cmd.Flags().StringVar(&configPath, "config", "", "Path to .agent-lint.yaml")
	cmd.Flags().StringVar(&ignorePath, "ignore-file", "", "Path to .agent-lint-ignore.yaml")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output results as JSON")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only show diagnostics")
	cmd.Flags().BoolVar(&strict, "strict", false, "Treat warnings as errors")
	cmd.Flags().BoolVar(&claude, "claude", false, "Enable Claude.ai-specific checks")
	cmd.Flags().BoolVar(&allowClaude, "allow-claude-code", false, "Allow Claude Code frontmatter extensions")

	return cmd
}

func newListRulesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-rules",
		Short: "List stable rule IDs",
		Run: func(cmd *cobra.Command, args []string) {
			for _, r := range rules.All() {
				fmt.Printf("%-36s %-8s %s\n", r.ID, r.Severity, r.Description)
			}
		},
	}
}
