package sdk

import (
	"context"
	"os"

	"github.com/1password/onepassword-sdk-go"
	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/version"
)

const (
	OnePasswordServiceAccountTokenEnv = "OP_SERVICE_ACCOUNT_TOKEN"
)

var _ provider.SecretClient = &OnePasswordClient{}

type OnePasswordClient struct {
	client *onepassword.Client
}

func NewOnePasswordClient() (*OnePasswordClient, error) {
	token := os.Getenv(OnePasswordServiceAccountTokenEnv)

	client, err := onepassword.NewClient(
		context.Background(),
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("genv", version.Version),
	)
	if err != nil {
		return nil, err
	}

	return &OnePasswordClient{client: client}, nil
}

func (c *OnePasswordClient) GetSecret(ctx context.Context, ref provider.SecretRef) ([]byte, error) {
	secret, err := c.client.Secrets().Resolve(ctx, ref.Key)
	if err != nil {
		return nil, err
	}

	return []byte(secret), nil
}
