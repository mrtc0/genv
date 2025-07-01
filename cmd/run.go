package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mrtc0/genv/dotenv"
	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "run [options] [COMMAND [ARG...]]",
	Short: "Run a command with environment variables from .env file",
	Long:  `Run a command with environment variables loaded from a .env file.`,
	Example: `genv run some-command
genv run --envfile /path/to/.env some-command`,
	RunE: run,
}

func init() {
	runCommand.Flags().StringP("envfile", "e", ".env", "Path to the .env file")

	rootCmd.AddCommand(runCommand)
}

func run(cmd *cobra.Command, args []string) error {
	envFile, err := cmd.Flags().GetString("envfile")
	if err != nil {
		return err
	}

	envMap, err := dotenv.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	for k, v := range envMap {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", k, err)
		}
	}

	command := exec.Command(args[0], args[1:]...)
	command.Env = os.Environ()
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}
	return nil
}
