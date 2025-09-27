package googlecloud

import (
	"context"
	"errors"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/googlecloud/secretmanager"
)

var _ provider.Provider = &Provider{}

type GoogleCloudProvider struct {
	ID        string
	Service   string
	ProjectID string
	Location  string
}

type Provider struct {
	Config *GoogleCloudProvider
}

func NewProvider(cfg *GoogleCloudProvider) provider.Provider {
	return &Provider{
		Config: cfg,
	}
}

func (p *Provider) NewClient(ctx context.Context) (provider.SecretClient, error) {
	switch p.Config.Service {
	case "SecretManager":
		return newClient(ctx, p.Config)
	default:
		return nil, errors.New("unsupported Google Cloud service: " + p.Config.Service)
	}
}

func newClient(ctx context.Context, cfg *GoogleCloudProvider) (provider.SecretClient, error) {
	return secretmanager.NewSecretManagerClient(
		ctx,
		cfg.ProjectID,
		cfg.Location,
	)
}
