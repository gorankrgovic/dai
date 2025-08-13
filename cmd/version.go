package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gorankrgovic/dai/internal/buildinfo"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print DAI version info",
	Run: func(cmd *cobra.Command, args []string) {
		v := buildinfo.Version
		if v == "" {
			v = "dev"
		}
		fmt.Printf("dai %s", v)
		if buildinfo.Commit != "" {
			fmt.Printf(" (%s)", buildinfo.Commit[:8])
		}
		if buildinfo.Date != "" {
			fmt.Printf(" %s", buildinfo.Date)
		}
		fmt.Println()
	},
}
