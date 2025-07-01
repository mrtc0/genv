package genv

import (
	"io"
	"os"
	"os/exec"
)

type CommandRunner interface {
	// Run executes the command with the provided environment variables.
	Run() error
}

type commandRunner struct {
	cmd *exec.Cmd
}

type CommandRunnerConfig struct {
	Name   string
	Args   []string
	Envs   map[string]string
	Stdout io.Writer
	Stderr io.Writer
}

func NewCommandRunner(cfg CommandRunnerConfig) (CommandRunner, error) {
	envs := make([]string, 0, len(cfg.Envs))
	for k, v := range cfg.Envs {
		envs = append(envs, k+"="+v)
	}

	command := exec.Command(cfg.Name, cfg.Args...)
	command.Env = append(os.Environ(), envs...)
	command.Stdout = cfg.Stdout
	command.Stderr = cfg.Stderr

	return &commandRunner{
		cmd: command,
	}, nil
}

func (c *commandRunner) Run() error {
	return c.cmd.Run()
}
