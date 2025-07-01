package cmd

import (
	"fmt"
	"os"

	"github.com/mrtc0/genv"
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

	runner, err := genv.NewCommandRunner(genv.CommandRunnerConfig{
		Name:   args[0],
		Args:   args[1:],
		Envs:   envMap,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("failed to create command runner: %w", err)
	}

	if err := runner.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}
	return nil
}
