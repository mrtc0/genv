package op

import (
	"context"
	"errors"

	"github.com/mrtc0/genv/provider"
)

var _ provider.SecretClient = &OPClient{}

type OPClient struct {
	account  string
	Executor OPCommandExecutor
}

func NewOPClient(account string) *OPClient {
	return &OPClient{
		account:  account,
		Executor: &DefaultOPCommandExecutor{},
	}
}

func (c *OPClient) GetSecret(ctx context.Context, ref provider.SecretRef) ([]byte, error) {
	out, err := c.Executor.Exec(ctx, c.buildArgs(c.account, ref))
	if err != nil {
		return nil, errors.New("op command failed: " + err.Error() + ": " + string(out))
	}
	return out, nil
}

func (c *OPClient) buildArgs(account string, ref provider.SecretRef) []string {
	if account != "" {
		return []string{"--account", account, "read", ref.Key}
	}
	return []string{"read", ref.Key}
}
