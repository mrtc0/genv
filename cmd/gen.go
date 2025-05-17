package cmd

import (
	"fmt"

	"github.com/mrtc0/genv"
	"github.com/spf13/cobra"
)

var (
	genvFilePath   string
	outputFilePath string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate .env file",
	Long:  `Generate .env file from secret providers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := genv.LoadConfig(genvFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		generator, err := genv.NewDotenvGenerator(ctx, genv.DotenvGeneratorConfig{
			Config:         cfg,
			OutputFilePath: outputFilePath,
		})
		if err != nil {
			return fmt.Errorf("failed to create dotenv generator: %w", err)
		}

		secrets, err := generator.FetchSecrets(ctx)
		if err != nil {
			return fmt.Errorf("failed to generate .env file: %w", err)
		}

		if err := genv.WriteDotenvFile(outputFilePath, secrets); err != nil {
			return fmt.Errorf("failed to write .env file: %w", err)
		}

		return nil
	},
}

func init() {
	genCmd.Flags().StringVar(&genvFilePath, "config", ".genv.yaml", "Path to the genv config file")
	genCmd.Flags().StringVar(&outputFilePath, "output", ".env", "Path to the output dotenv file")
	rootCmd.AddCommand(genCmd)
}
