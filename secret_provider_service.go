package genv

import (
	"context"
	"fmt"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/aws"
	"github.com/mrtc0/genv/provider/googlecloud"
	"github.com/mrtc0/genv/provider/onepassword"
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

	for _, p := range sp.GoogleCloud {
		providerConfig := googlecloud.GoogleCloudProvider{
			ID:        p.ID,
			Service:   p.Service,
			ProjectID: p.ProjectID,
			Location:  p.Location,
		}

		gc := googlecloud.NewProvider(&providerConfig)
		client, err := gc.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create Google Cloud secret client: %w", err)
		}

		secretProviderClients[p.ID] = client
	}

	for _, p := range sp.OnePassword {
		opProvider := onepassword.NewProvider(
			onepassword.WithAccount(p.Auth.Account),
			onepassword.WithAuthMethod(p.Auth.Method),
		)
		client, err := opProvider.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create 1Password secret client: %w", err)
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
