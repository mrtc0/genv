package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	ID       string           `yaml:"id"`
	Service  awsSecretService `yaml:"service"`
	Region   string           `yaml:"region,omitempty"`
	Endpoint string           `yaml:"endpoint,omitempty"`
	Auth     AwsAuth          `yaml:"auth,omitempty"`
}

type AwsAuth struct {
	Profile                string
	SharedCredentialsFiles []string
	SharedConfigFiles      []string
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
	cfg, err := GetAWSConfig(ctx, providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	switch providerConfig.Service {
	case AWSSecretsManager:
		return secretsmanager.NewSecretsManager(cfg), nil
	default:
		return nil, errors.New("unsupported AWS service")
	}
}

func GetAWSConfig(ctx context.Context, providerConfig *AwsProviderConfig) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, loadOptions(providerConfig)...)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

func loadOptions(c *AwsProviderConfig) []func(*config.LoadOptions) error {
	loadOptions := []func(*config.LoadOptions) error{}

	if c.Region != "" {
		loadOptions = append(loadOptions, config.WithRegion(c.Region))
	}

	if len(c.Auth.SharedConfigFiles) > 0 {
		loadOptions = append(loadOptions, config.WithSharedConfigFiles(c.Auth.SharedConfigFiles))
	}

	if len(c.Auth.SharedCredentialsFiles) > 0 {
		loadOptions = append(loadOptions, config.WithSharedCredentialsFiles(c.Auth.SharedCredentialsFiles))
	}

	if c.Auth.Profile != "" {
		loadOptions = append(loadOptions, config.WithSharedConfigProfile(c.Auth.Profile))
	}

	if c.Endpoint != "" {
		loadOptions = append(loadOptions, config.WithBaseEndpoint(c.Endpoint))
	}

	return loadOptions
}
