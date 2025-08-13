package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/gorankrgovic/dai/internal/config"
)

var (
	flagAuthToken  string
	flagAuthShow   bool
	flagAuthDelete bool
)

// authCmd handles authentication-related actions (currently GitHub PAT).
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate DAI with external services (GitHub token for now)",
	Long: `Authenticate DAI with external services.

This stores your GitHub Personal Access Token (PAT) in ~/.dai/github_token with 0600 permissions.
It is separate from 'dai config' on purpose (future: cloud/local modes).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default behaviour when no subcommand is provided
		if flagAuthShow {
			return runAuthShow()
		}
		if flagAuthDelete {
			return runAuthDelete()
		}
		return runAuthInteractive()
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show whether a GitHub token is stored",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAuthShow()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVar(&flagAuthToken, "token", "", "GitHub Personal Access Token (non-interactive)")
	authCmd.Flags().BoolVar(&flagAuthShow, "show", false, "Show whether a token is stored (does not print the token)")
	authCmd.Flags().BoolVar(&flagAuthDelete, "delete", false, "Delete stored GitHub token")

	// Add subcommand: dai auth status
	authCmd.AddCommand(authStatusCmd)
}

func runAuthShow() error {
	exists, err := config.GitHubTokenExists()
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("✓ GitHub token is stored (in ~/.dai/github_token).")
	} else {
		fmt.Println("✗ No GitHub token stored.")
	}
	return nil
}

func runAuthDelete() error {
	if err := confirmDanger("This will delete the stored GitHub token. Continue?"); err != nil {
		fmt.Println("Aborted.")
		return nil
	}
	if err := config.DeleteGitHubToken(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No stored token to delete.")
			return nil
		}
		return err
	}
	fmt.Println("Deleted stored GitHub token.")
	return nil
}

func runAuthInteractive() error {
	token := strings.TrimSpace(flagAuthToken)
	if token == "" {
		var input string
		q := &survey.Password{
			Message: "Paste your GitHub Personal Access Token (will be hidden):",
		}
		if err := survey.AskOne(q, &input, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
		token = strings.TrimSpace(input)
	}

	if err := validateGitHubToken(token); err != nil {
		return err
	}

	if err := confirm("Save token to ~/.dai/github_token?"); err != nil {
		fmt.Println("Aborted.")
		return nil
	}

	if err := config.SaveGitHubToken(token); err != nil {
		return err
	}

	fmt.Println("✓ GitHub token saved to ~/.dai/github_token (permissions 0600).")
	return nil
}

func validateGitHubToken(tok string) error {
	if tok == "" {
		return errors.New("empty token")
	}
	if len(tok) < 20 {
		return errors.New("token looks too short")
	}
	if strings.ContainsAny(tok, " \t\r\n") {
		return errors.New("token must not contain whitespace")
	}
	var typical = regexp.MustCompile(`^(gh[pousr]_|github_pat_)`)
	if !typical.MatchString(tok) {
		fmt.Println("! Note: token doesn't match typical GitHub prefixes (ghp_/github_pat_/gho_/ghu_/ghs_/ghr_). Continuing anyway.")
	}
	return nil
}

func confirm(msg string) error {
	var ok bool
	p := &survey.Confirm{
		Message: msg,
		Default: true,
	}
	if err := survey.AskOne(p, &ok); err != nil {
		return err
	}
	if !ok {
		return errors.New("user declined")
	}
	return nil
}

func confirmDanger(msg string) error {
	var ok bool
	p := &survey.Confirm{
		Message: msg,
		Default: false,
	}
	if err := survey.AskOne(p, &ok); err != nil {
		return err
	}
	if !ok {
		return errors.New("user declined")
	}
	return nil
}
