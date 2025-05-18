package cmd

import (
	"errors"
	"fmt"

	"github.com/mrtc0/genv"
	"github.com/mrtc0/genv/dotenv"
	"github.com/mrtc0/genv/internal/renderer"
	"github.com/spf13/cobra"
)

var (
	dotenvFilePath string
	nameOnly       bool
)

var outdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Show outdated envs in the dotenv file.",
	Long:  `Show the difference between the current genv environment variable definitions and the dotenv file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := genv.LoadConfig(genvFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		dotenvMap, err := dotenv.ReadFile(dotenvFilePath)
		if err != nil {
			return fmt.Errorf("failed to read dotenv file: %w", err)
		}

		diff, err := genv.Diff(ctx, cfg, dotenvMap, nameOnly)
		if err != nil {
			return fmt.Errorf("failed to diff envs: %w", err)
		}

		if !diff.IsChanged() {
			return nil
		}

		fmt.Printf("%s\n", renderer.RenderDiff(diff))

		return errors.New("outdated envs found")
	},
}

func init() {
	outdatedCmd.Flags().StringVar(&genvFilePath, "config", ".genv.yaml", "Path to the genv config file.")
	outdatedCmd.Flags().StringVar(&dotenvFilePath, "envfile", ".env", "Path to the dotenv file.")
	outdatedCmd.Flags().BoolVar(&nameOnly, "name-only", false, "Only the differences in the variable names of the environment variables are checked. No values are retrieved from remote credential providers.")
	rootCmd.AddCommand(outdatedCmd)
}
