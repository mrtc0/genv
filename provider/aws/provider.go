package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/aws/secretsmanager"
)

type awsSecretService string

const (
	AWSSecretsManager awsSecretService = "SecretsManager"
)

var _ provider.Provider = &Provider{}

type AwsProviderConfig struct {
	ID      string           `yaml:"id"`
	Service awsSecretService `yaml:"service"`
	Region  string           `yaml:"region,omitempty"`
	Auth    AwsAuth          `yaml:"auth,omitempty"`
}

type AwsAuth struct {
	Profile string `yaml:"profile,omitempty"`
}

func (c *AwsProviderConfig) ProviderID() string {
	return c.ID
}

// Provider satisfies the provider interface
type Provider struct {
	Config *AwsProviderConfig
}

func NewProvider(cfg *AwsProviderConfig) provider.Provider {
	return &Provider{Config: cfg}
}

func (p *Provider) NewClient(ctx context.Context) (provider.SecretClient, error) {
	return newClient(ctx, p.Config)
}

func newClient(ctx context.Context, providerConfig *AwsProviderConfig) (provider.SecretClient, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedConfigProfile(providerConfig.Auth.Profile),
		config.WithRegion(providerConfig.Region),
	)
	if err != nil {
		return nil, err
	}

	switch providerConfig.Service {
	case AWSSecretsManager:
		return secretsmanager.NewSecretsManager(cfg), nil
	default:
		return nil, errors.New("unsupported AWS service")
	}
}
