package genv_test

import (
	"bytes"
	"testing"

	"github.com/mrtc0/genv"
	"github.com/stretchr/testify/assert"
)

func TestCommandRunner_Run(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	runner, err := genv.NewCommandRunner(genv.CommandRunnerConfig{
		Name: "sh",
		Args: []string{"-c", "echo Hello \"$ENV_VAR_1\" \"$ENV_VAR_2\""},
		Envs: map[string]string{
			"ENV_VAR_1": "Value1",
			"ENV_VAR_2": "Value2",
		},
		Stdout: stdout,
		Stderr: stderr,
	})
	assert.NoError(t, err)

	err = runner.Run()
	assert.NoError(t, err)

	expectedOutput := "Hello Value1 Value2\n"
	assert.Equal(t, expectedOutput, stdout.String())

	assert.Empty(t, stderr.String())
}
