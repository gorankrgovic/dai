package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/gorankrgovic/dai/internal/ignore"
)

var (
	flagIgnorePathCLI string
	flagIgnoreForce   bool
	flagIgnoreYes     bool
)

func init() {
	rootCmd.AddCommand(ignoreCmd)

	ignoreCmd.Flags().StringVar(&flagIgnorePathCLI, "path", ".daiignore", "Path to create the .daiignore file (relative to project root)")
	ignoreCmd.Flags().BoolVar(&flagIgnoreForce, "force", false, "Overwrite existing file without prompting")
	ignoreCmd.Flags().BoolVar(&flagIgnoreYes, "yes", false, "Non-interactive: assume 'yes' to prompts")
}

var ignoreCmd = &cobra.Command{
	Use:   "ignore",
	Short: "Create a default .daiignore file (gitignore syntax) in the project root",
	Long:  "Create a default .daiignore file to exclude files from DAI triage (uses gitignore syntax).",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := ensureProjectRoot()
		if err != nil {
			return err
		}

		path := flagIgnorePathCLI
		if !filepath.IsAbs(path) {
			path = filepath.Join(wd, path)
		}
		return writeDaiIgnoreWithPrompt(path, flagIgnoreForce, flagIgnoreYes)
	},
}

// promptCreateDaiIgnore is a helper you can call from `dai init` flow.
func promptCreateDaiIgnore(projectRoot string) error {
	path := filepath.Join(projectRoot, ".daiignore")

	// Ask if user wants to create it
	var create bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Create a default .daiignore file to customize triage ignores?",
		Default: true,
	}, &create); err != nil {
		return err
	}
	if !create {
		fmt.Println("Skipped creating .daiignore.")
		return nil
	}

	// If exists, ask to overwrite
	force := false
	if _, err := os.Stat(path); err == nil {
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf(".daiignore already exists at %s. Overwrite?", relPath(projectRoot, path)),
			Default: false,
		}, &force); err != nil {
			return err
		}
	}

	if err := writeDaiIgnore(path, ignore.DefaultDaiIgnore(), force); err != nil {
		return err
	}
	fmt.Printf("✓ Wrote %s\n", relPath(projectRoot, path))
	return nil
}

func writeDaiIgnoreWithPrompt(absPath string, force bool, yes bool) error {
	if !filepath.IsAbs(absPath) {
		return fmt.Errorf("path must be absolute: %s", absPath)
	}
	root := filepath.Dir(absPath) // best effort for pretty print

	if _, err := os.Stat(absPath); err == nil && !force {
		if yes {
			force = true
		} else {
			var ok bool
			if err := survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("%s already exists. Overwrite?", relPath(root, absPath)),
				Default: false,
			}, &ok); err != nil {
				return err
			}
			if !ok {
				fmt.Println("Aborted. No changes made.")
				return nil
			}
			force = true
		}
	}

	if err := writeDaiIgnore(absPath, ignore.DefaultDaiIgnore(), force); err != nil {
		return err
	}
	fmt.Printf("✓ Wrote %s\n", relPath(root, absPath))
	return nil
}

func writeDaiIgnore(absPath string, content string, force bool) error {
	if _, err := os.Stat(absPath); err == nil && !force {
		return fmt.Errorf("file exists: %s (use --force to overwrite)", absPath)
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(absPath, []byte(content), 0o644)
}

func relPath(base, full string) string {
	if r, err := filepath.Rel(base, full); err == nil && !strings.HasPrefix(r, "..") {
		return r
	}
	return full
}
