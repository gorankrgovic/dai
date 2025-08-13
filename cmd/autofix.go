package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(autofixCmd)
}

var autofixCmd = &cobra.Command{
	Use:   "autofix [issue_number]",
	Short: "Auto-fix an Issue by opening a PR (in development)",
	Long: `🚧 This command is currently in development (TODO).

Planned flow:
- Read the GitHub Issue
- Analyze diff/context with LLM
- Generate a patch
- Create a branch and open a PR`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚧 'dai autofix' is in development. No changes were made.")
		return nil
	},
}
