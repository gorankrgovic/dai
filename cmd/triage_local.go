package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/gorankrgovic/dai/internal/config"
	"github.com/gorankrgovic/dai/internal/triage"
)

var (
	flagLocalModel    string
	flagLocalMaxKB    int
	flagLocalLogPath  string
	flagLocalFormat   string // md|json
	flagLocalNoStdout bool
)

func init() {
	rootCmd.AddCommand(triageLocalCmd)

	triageLocalCmd.Flags().StringVar(&flagLocalModel, "model", "", "Override OpenAI model from config (optional)")
	triageLocalCmd.Flags().IntVar(&flagLocalMaxKB, "max-file-kb", 200, "Max bytes per analyzed file (KB)")
	triageLocalCmd.Flags().StringVar(&flagLocalLogPath, "log", ".dai/local.log", "Path to local log file (relative to project root)")
	triageLocalCmd.Flags().StringVar(&flagLocalFormat, "format", "md", "Log format: md | json")
	triageLocalCmd.Flags().BoolVar(&flagLocalNoStdout, "no-stdout", false, "Do not print findings to stdout (log only)")
}

var triageLocalCmd = &cobra.Command{
	Use:   "triage-local <file>",
	Short: "Analyze a single local file and append findings to .dai/local.log",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// moramo biti u project root-u zbog .dai/
		root, err := ensureProjectRoot()
		if err != nil {
			return err
		}

		// globalni config (OpenAI key + model)
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("global config not found — run 'dai config' first: %w", err)
		}
		model := cfg.Model
		if flagLocalModel != "" {
			model = flagLocalModel
		}
		if strings.TrimSpace(cfg.OpenAIKey) == "" {
			return errors.New("OpenAI key missing — run 'dai config'")
		}

		// rezolvoj putanju fajla
		p := args[0]
		if !filepath.IsAbs(p) {
			p = filepath.Join(root, p)
		}
		info, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("cannot stat file: %w", err)
		}
		if info.IsDir() {
			return fmt.Errorf("path is a directory, expected a file: %s", p)
		}

		// analiza
		finding, truncated, err := triage.AnalyzeLocal(cmd.Context(), cfg.OpenAIKey, model, p, int64(flagLocalMaxKB)*1024)
		if err != nil {
			return err
		}

		// pripremi log path
		logPath := flagLocalLogPath
		if !filepath.IsAbs(logPath) {
			logPath = filepath.Join(root, logPath)
		}
		if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
			return err
		}

		// upiši log u traženom formatu
		now := time.Now().Format(time.RFC3339)
		switch strings.ToLower(flagLocalFormat) {
		case "json":
			entry := map[string]any{
				"time":      now,
				"file":      relOrSame(root, p),
				"model":     model,
				"truncated": truncated,
				"type":      finding.Type,
				"title":     finding.Title,
				"severity":  finding.Severity,
				"lineHints": finding.LineHints,
				"details":   finding.Details,
			}
			b, _ := json.Marshal(entry)
			if err := appendLine(logPath, string(b)); err != nil {
				return err
			}
		default: // md
			var sb strings.Builder
			fmt.Fprintf(&sb, "### %s — %s (model: %s, truncated: %v)\n", now, relOrSame(root, p), model, truncated)
			fmt.Fprintf(&sb, "- Type: **%s**\n", strings.ToUpper(finding.Type))
			if finding.Title != "" {
				fmt.Fprintf(&sb, "- Title: %s\n", finding.Title)
			}
			if finding.Severity != "" {
				fmt.Fprintf(&sb, "- Severity: %s\n", strings.ToUpper(finding.Severity))
			}
			if finding.LineHints != "" {
				fmt.Fprintf(&sb, "- Lines: %s\n", finding.LineHints)
			}
			if finding.Details != "" {
				fmt.Fprintf(&sb, "\n%s\n", finding.Details)
			}
			fmt.Fprintf(&sb, "\n---\n")
			if err := appendLine(logPath, sb.String()); err != nil {
				return err
			}
		}

		if !flagLocalNoStdout {
			fmt.Printf("File: %s\n", relOrSame(root, p))
			fmt.Printf("Type: %s | Severity: %s\n", finding.Type, strings.ToUpper(finding.Severity))
			if finding.Title != "" {
				fmt.Println("Title:", finding.Title)
			}
			if finding.LineHints != "" {
				fmt.Println("Lines:", finding.LineHints)
			}
			if finding.Details != "" {
				fmt.Println("Details:", finding.Details)
			}
			fmt.Printf("→ Logged to %s\n", relOrSame(root, logPath))
		}
		return nil
	},
}

func appendLine(path, s string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(s)
	return err
}

func relOrSame(base, full string) string {
	if r, err := filepath.Rel(base, full); err == nil && !strings.HasPrefix(r, "..") {
		return r
	}
	return full
}
