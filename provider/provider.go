package provider

import "context"

const (
	// AWSProviderName is the name of the AWS provider.
	AWSProviderName ProviderName = "aws"
	// OnePasswordProviderName is the name of the 1Password provider.
	OnePasswordProviderName ProviderName = "1password"
)

type ProviderName string

// Provider is an interface for interacting with secret providers.
type Provider interface {
	NewClient(ctx context.Context) (SecretClient, error)
}

// SecretClient provides access to secrets.
type SecretClient interface {
	GetSecret(ctx context.Context, ref SecretRef) ([]byte, error)
}

// SecretRef represents a location of a secret value.
type SecretRef struct {
	Key      string
	Property string
}
