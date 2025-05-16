package genv

import (
	"context"
	"fmt"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/aws"
)

func NewSecretProviderClient(ctx context.Context, sp SecretProvider) (map[string]provider.SecretClient, error) {
	secretProviderClients := make(map[string]provider.SecretClient)

	for _, p := range sp.Aws {
		providerConfig := &aws.AwsProviderConfig{
			ID:      p.ID,
			Service: aws.AWSSecretsManager,
			Region:  p.Region,
			Auth: aws.AwsAuth{
				Profile: p.Auth.Profile,
			},
		}

		awsProvider := aws.NewProvider(providerConfig)
		client, err := awsProvider.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create secret client: %w", err)
		}

		secretProviderClients[p.ID] = client
	}

	return secretProviderClients, nil
}
