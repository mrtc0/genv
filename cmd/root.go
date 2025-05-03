package cmd

import (
	"os"

	"github.com/mrtc0/genv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "genv",
	Short:        "genv is a dotenv generator",
	Long:         `genv is a dotenv generator that generates .env files from various secret providers.`,
	Version:      genv.Version,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
