package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gorankrgovic/dai/internal/config"
	"github.com/gorankrgovic/dai/internal/project"
	"github.com/gorankrgovic/dai/internal/triage"
)

var (
	flagTriageExt    string
	flagTriageDryRun bool
	flagModel        string
	flagMaxKB        int
	flagIgnorePath   string
	flagAlwaysOpen   bool
	flagDiffContext  int
)

func init() {
	rootCmd.AddCommand(triageCmd)

	triageCmd.Flags().StringVar(&flagTriageExt, "ext", ".js,.jsx,.ts,.tsx,.vue,.php,.py,.go", "Comma-separated file extensions to analyze")
	triageCmd.Flags().BoolVar(&flagTriageDryRun, "dry-run", false, "Print the would-be GitHub issue without creating it")
	triageCmd.Flags().StringVar(&flagModel, "model", "", "Override OpenAI model from config (optional)")
	triageCmd.Flags().IntVar(&flagMaxKB, "max-file-kb", 80, "Max file size per analyzed file (KB)")
	triageCmd.Flags().StringVar(&flagIgnorePath, "ignore", ".daiignore", "Path to ignore file (gitignore syntax), relative to project root")
	triageCmd.Flags().BoolVar(&flagAlwaysOpen, "always-open", false, "Always create a GitHub issue even when no findings")
	triageCmd.Flags().IntVar(&flagDiffContext, "diff-context", 3, "Number of context lines per diff hunk")
}

var triageCmd = &cobra.Command{
	Use:   "triage [commit]",
	Short: "Analyze a commit and open a single GitHub issue with findings",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := ensureProjectRoot()
		if err != nil {
			return err
		}

		// project.yaml
		prj, err := project.Load(wd)
		if err != nil {
			return fmt.Errorf("project config not found — run 'dai init' first: %w", err)
		}
		if prj.Owner == "" || prj.Repo == "" {
			return fmt.Errorf("project config missing owner/repo")
		}

		// GitHub token
		token, err := config.LoadGitHubToken()
		if err != nil || strings.TrimSpace(token) == "" {
			return fmt.Errorf("GitHub token not found — run 'dai auth' first: %w", err)
		}
		token = strings.TrimSpace(token)

		// OpenAI config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("global config not found — run 'dai config' first: %w", err)
		}
		if flagModel != "" {
			cfg.Model = flagModel
		}
		if strings.TrimSpace(cfg.OpenAIKey) == "" {
			return fmt.Errorf("OpenAI key missing in global config — run 'dai config'")
		}

		// Commit
		var commit string
		if len(args) == 1 {
			commit = strings.TrimSpace(args[0])
		}

		// Include extensions
		exts := splitCSV(flagTriageExt)

		opts := triage.Options{
			Root:         wd,
			Owner:        prj.Owner,
			Repo:         prj.Repo,
			GitHubToken:  token,
			OpenAIKey:    cfg.OpenAIKey,
			Model:        cfg.Model,
			Commit:       commit, // empty == HEAD
			IncludeExts:  exts,
			MaxFileBytes: int64(flagMaxKB) * 1024,
			IgnoreFile:   filepath.Join(wd, flagIgnorePath),
			DryRun:       flagTriageDryRun,
			AlwaysOpen:   flagAlwaysOpen,
			DiffContext:  flagDiffContext, // NEW
		}

		result, err := triage.Run(cmd.Context(), opts)
		if err != nil {
			return err
		}
		if opts.DryRun {
			fmt.Println("— DRY RUN —")
			fmt.Println(result.Body) // već formatiran body
			return nil
		}
		if result.Skipped {
			fmt.Println("No findings from diff hunks. Skipped creating a GitHub issue. (use --always-open to force)")
			return nil
		}
		fmt.Printf("✓ Issue #%d created: %s\n", result.Number, result.URL)
		return nil
	},
}

func ensureProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := filepath.Join(wd, ".dai", "project.yaml")
	if _, statErr := os.Stat(p); statErr != nil {
		if os.IsNotExist(statErr) {
			return "", fmt.Errorf("'.dai/project.yaml' not found in current directory (%s)\nRun 'dai triage' from the project root (where you ran 'dai init')", wd)
		}
		return "", statErr
	}
	return wd, nil
}

func splitCSV(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p == "" {
			continue
		}
		if !strings.HasPrefix(p, ".") {
			p = "." + p
		}
		out = append(out, p)
	}
	return out
}
