package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/secretutil"
)

var _ provider.Provider = &Provider{}
var _ provider.SecretClient = &Client{}

// ExecProviderConfig holds configuration for the exec provider.
type ExecProviderConfig struct {
	ID      string
	Command []string
}

// Provider satisfies the provider.Provider interface.
type Provider struct {
	Config *ExecProviderConfig
}

func NewProvider(cfg *ExecProviderConfig) provider.Provider {
	return &Provider{Config: cfg}
}

func (p *Provider) NewClient(ctx context.Context) (provider.SecretClient, error) {
	if len(p.Config.Command) == 0 {
		return nil, fmt.Errorf("exec provider %q: command must not be empty", p.Config.ID)
	}
	return &Client{command: p.Config.Command}, nil
}

// Client executes the configured command once and caches the JSON output.
type Client struct {
	command []string
	once    sync.Once
	output  []byte
	execErr error
}

// GetSecret runs the command (at most once), then extracts the value at
// ref.Key using gjson. If ref.Property is also set it is applied as a
// second gjson path on the result of ref.Key.
func (c *Client) GetSecret(ctx context.Context, ref provider.SecretRef) ([]byte, error) {
	c.once.Do(func() {
		c.output, c.execErr = c.run(ctx)
	})
	if c.execErr != nil {
		return nil, c.execErr
	}

	val, err := secretutil.GetValueFromJSON(c.output, ref.Key)
	if err != nil {
		return nil, fmt.Errorf("exec provider: key %q not found in output: %w", ref.Key, err)
	}

	if ref.Property == "" {
		return val, nil
	}

	val, err = secretutil.GetValueFromJSON(val, ref.Property)
	if err != nil {
		return nil, fmt.Errorf("exec provider: property %q not found: %w", ref.Property, err)
	}

	return val, nil
}

func (c *Client) run(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, c.command[0], c.command[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("exec provider: command failed: %w\nstderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
