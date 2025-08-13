package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/gorankrgovic/dai/internal/gitutil"
	"github.com/gorankrgovic/dai/internal/project"
)

var (
	initFlagForce   bool
	initFlagVerbose bool
	initFlagPath    string
)

func init() {
	initCmd.Flags().BoolVarP(&initFlagForce, "force", "f", false, "overwrite existing .dai/project.yaml without prompt")
	initCmd.Flags().BoolVarP(&initFlagVerbose, "verbose", "v", false, "print detected git info")
	initCmd.Flags().StringVarP(&initFlagPath, "path", "p", ".", "project path (defaults to current directory)")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [optional_project_path]",
	Short: "Initialize DAI for this project (.dai/project.yaml)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := initFlagPath
		if len(args) == 1 {
			target = args[0]
		}
		absTarget, err := filepath.Abs(target)
		if err != nil {
			return err
		}

		// 1) git repo?
		if !gitutil.IsRepo(absTarget) {
			return fmt.Errorf("git is not initialized in this directory.\nPath: %s\nRun: git init && git remote add origin <url>", absTarget)
		}

		// 2) remote origin?
		remote, ok := gitutil.RemoteOriginURL(absTarget)
		if !ok {
			return fmt.Errorf("git remote 'origin' not found at %s.\nRun: git remote add origin <url>", absTarget)
		}

		// uvek radi u git root-u (podržava i podfoldere u repo-u)
		repoRoot, _ := gitutil.RepoRoot(absTarget)
		baseDir := repoRoot
		if baseDir == "" {
			baseDir = absTarget
		}

		if initFlagVerbose {
			fmt.Println("Repo root:", baseDir)
			fmt.Println("Detected origin:", remote)
		}

		owner, repo := gitutil.ParseOwnerRepo(remote)
		if owner == "" || repo == "" {
			return fmt.Errorf("could not parse owner/repo from origin URL: %s\nTip: use a GitHub URL (ssh or https), e.g. git@github.com:ORG/REPO.git", remote)
		}

		projectPath := filepath.Join(baseDir, ".dai", "project.yaml")
		if _, err := os.Stat(projectPath); err == nil && !initFlagForce {
			confirm := false
			if err := survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("%s already exists. Overwrite?", projectPath),
				Default: false,
			}, &confirm); err != nil {
				return err
			}
			if !confirm {
				fmt.Println("Aborted.")
				return nil
			}
		}

		// 3) sačuvaj project config u ROOT-u
		p := &project.Project{Provider: "github", Owner: owner, Repo: repo}
		if err := project.Save(baseDir, p); err != nil {
			return err
		}

		// 4) dodaj .dai u .gitignore u ROOT-u (ako nije tamo)
		gitignorePath := filepath.Join(baseDir, ".gitignore")
		added, err := ensureGitignoreHasEntry(gitignorePath, ".dai")
		if err != nil {
			return err
		}
		if initFlagVerbose {
			fmt.Println("Gitignore path:", gitignorePath)
			if added {
				fmt.Println("Added `.dai` to .gitignore")
			} else {
				fmt.Println("`.dai` already present in .gitignore (or equivalent)")
			}
		}

		// ▼ DODAJ OVO: ponudi kreiranje .daiignore u project root-u
		if err := promptCreateDaiIgnore(baseDir); err != nil {
			return err
		}

		fmt.Printf("DAI project initialized at %s\n", projectPath)
		fmt.Printf("Detected repo: %s/%s\n", owner, repo)
		fmt.Println("If you change origin, run 'dai init' again to update project config.")
		return nil
	},
}

// pametniji check za prisustvo .dai
func ensureGitignoreHasEntry(path string, entry string) (bool, error) {
	// ako .gitignore ne postoji → kreiraj ga sa entry
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(entry+"\n"), 0o644); err != nil {
			return false, err
		}
		return true, nil
	}

	// proveri par tipičnih varijanti
	equivalents := map[string]struct{}{
		entry:         {},
		entry + "/":   {},
		"/" + entry:   {},
		"**/" + entry: {},
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if _, ok := equivalents[line]; ok {
			return false, nil // već postoji
		}
	}
	if err := s.Err(); err != nil {
		return false, err
	}

	// dodaj na kraj
	af, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return false, err
	}
	defer af.Close()
	if _, err := af.WriteString(entry + "\n"); err != nil {
		return false, err
	}
	return true, nil
}
