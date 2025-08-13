package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	survey "github.com/AlecAivazis/survey/v2"

	"github.com/gorankrgovic/dai/internal/config"
)

var modelChoices = []string{
	"gpt-4o",
	"gpt-4o-mini",
	"gpt-4.1",
	"o4-mini",
	"Custom…",
}

var configWizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Interactive setup (OpenAI key + model)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Key (sa maskiranjem)
		fmt.Print("Enter your OpenAI API key (input hidden): ")
		secret, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}
		if len(secret) == 0 {
			return fmt.Errorf("empty key")
		}

		// 2) Model (interaktivni select)
		var sel string
		err = survey.AskOne(&survey.Select{
			Message:  "Choose default model:",
			Options:  modelChoices,
			Default:  "gpt-4o-mini",
			PageSize: 7,
		}, &sel)
		if err != nil {
			return err
		}

		if sel == "Custom…" {
			var custom string
			if err := survey.AskOne(&survey.Input{
				Message: "Enter custom model name:",
				Help:    "Type the exact model id as in your provider (e.g. gpt-4o-mini-2025-05-xx)",
			}, &custom, survey.WithValidator(survey.Required)); err != nil {
				return err
			}
			sel = custom
		}

		cfg := &config.Config{
			OpenAIKey: string(secret),
			Model:     sel,
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("Config saved.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configWizardCmd)
	configCmd.AddCommand(configSetKeyCmd)
	configCmd.AddCommand(configSetModelCmd)
	configCmd.AddCommand(configShowCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure global DAI settings (~/.dai/config.yaml)",
}

var configSetKeyCmd = &cobra.Command{
	Use:   "set-key",
	Short: "Set OpenAI API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Enter your OpenAI API key (input hidden): ")
		secret, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}
		if len(secret) == 0 {
			return fmt.Errorf("empty key")
		}
		cfg, _ := config.Load()
		if cfg == nil {
			cfg = &config.Config{}
		}
		cfg.OpenAIKey = string(secret)
		if cfg.Model == "" {
			cfg.Model = "gpt-4o-mini"
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("OpenAI key saved.")
		return nil
	},
}

var configSetModelCmd = &cobra.Command{
	Use:   "set-model",
	Short: "Interactively choose default model",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("run 'dai config set-key' or 'dai config wizard' first: %w", err)
		}

		def := cfg.Model
		if def == "" {
			def = "gpt-4o-mini"
		}

		var sel string
		err = survey.AskOne(&survey.Select{
			Message:  "Choose default model:",
			Options:  modelChoices,
			Default:  def,
			PageSize: 7,
		}, &sel)
		if err != nil {
			return err
		}

		if sel == "Custom…" {
			var custom string
			if err := survey.AskOne(&survey.Input{
				Message: "Enter custom model name:",
				Help:    "Type the exact model id as in your provider",
			}, &custom, survey.WithValidator(survey.Required)); err != nil {
				return err
			}
			sel = custom
		}

		cfg.Model = sel
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("Model saved:", cfg.Model)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current global config (key masked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("no config — run 'dai config set-key'")
		}
		masked := "not set"
		if cfg.OpenAIKey != "" {
			if len(cfg.OpenAIKey) > 8 {
				masked = cfg.OpenAIKey[:4] + "..." + cfg.OpenAIKey[len(cfg.OpenAIKey)-4:]
			} else {
				masked = "****"
			}
		}
		fmt.Printf("Model: %s\nOpenAI Key: %s\n", cfg.Model, masked)
		return nil
	},
}
