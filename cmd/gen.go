package cmd

import (
	"github.com/mrtc0/genv"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate .env file",
	Long:  `Generate .env file from secret providers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := genv.LoadConfig(".genv.yaml")
		if err != nil {
			return err
		}

		generator := genv.NewDotenvGenerator(genv.DotenvGeneratorConfig{
			Config: cfg,
		})

		if err := generator.Generate(ctx); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
