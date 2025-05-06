package cmd

import (
	"fmt"

	"github.com/mrtc0/genv"
	"github.com/spf13/cobra"
)

const (
	defaultGenvFilePath   = ".genv.yaml"
	defaultOutputFilePath = ".env"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate .env file",
	Long:  `Generate .env file from secret providers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := genv.LoadConfig(".genv.yaml")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		generator := genv.NewDotenvGenerator(genv.DotenvGeneratorConfig{
			Config:         cfg,
			OutputFilePath: defaultOutputFilePath,
		})

		if err := generator.Generate(ctx); err != nil {
			return fmt.Errorf("failed to generate .env file: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
