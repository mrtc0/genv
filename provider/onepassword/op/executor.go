package op

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
)

type OPCommandExecutor interface {
	Exec(ctx context.Context, args []string) ([]byte, error)
}

type DefaultOPCommandExecutor struct{}

func (e *DefaultOPCommandExecutor) Exec(ctx context.Context, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "op", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New("op command failed: " + err.Error() + ": " + string(out))
	}

	return bytes.TrimSuffix(out, []byte{'\n'}), nil
}
