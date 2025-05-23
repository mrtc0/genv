package genv

import (
	"context"
	"fmt"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/aws"
)

type SecretProviderService struct {
	clients map[string]provider.SecretClient
}

func NewSecretProviderService(ctx context.Context, sp SecretProvider) (*SecretProviderService, error) {
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

	return &SecretProviderService{
		clients: secretProviderClients,
	}, nil
}

func (s *SecretProviderService) AddSecretProviderClient(providerID string, client provider.SecretClient) {
	if s.clients == nil {
		s.clients = make(map[string]provider.SecretClient)
	}

	s.clients[providerID] = client
}

type GetSecretInput struct {
	Key      string
	Property string
}

func (s *SecretProviderService) GetSecret(ctx context.Context, providerID string, input GetSecretInput) ([]byte, error) {
	ref := provider.SecretRef{
		Key:      input.Key,
		Property: input.Property,
	}

	client, ok := s.clients[providerID]
	if !ok {
		return nil, fmt.Errorf("secret provider client not found for ID: %s", providerID)
	}

	return client.GetSecret(ctx, ref)
}
