package genv_test

import (
	"context"
	"testing"

	"github.com/mrtc0/genv"
	"github.com/mrtc0/genv/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDotenvGenerator_FetchSecrets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config := &genv.Config{
		SecretProvider: genv.SecretProvider{
			Aws: []genv.AwsProvider{
				{
					ID:      "example-account",
					Service: "secretsmanager",
					Region:  "us-east-1",
					Auth: genv.AwsAuth{
						Profile: "default",
					},
				},
			},
		},
		Envs: map[string]genv.EnvValue{
			"EXAMPLE_ENV": {
				Value: "example-value",
			},
			"EXAMPLE_SECRET": {
				SecretRef: &genv.SecretRef{
					Provider: "example-account",
					Key:      "example-key",
				},
			},
		},
	}

	svc, err := genv.NewSecretProviderService(ctx, config.SecretProvider)
	require.NoError(t, err)
	svc.AddSecretProviderClient("example-account", &mockSecretClient{
		returnSecretValue: []byte("secret-value"),
	})

	generator := &genv.DotenvGenerator{
		Config:                config,
		SecretProviderService: svc,
	}

	secrets, err := generator.FetchSecrets(ctx)
	assert.NoError(t, err)

	expected := map[string]string{
		"EXAMPLE_ENV":    "example-value",
		"EXAMPLE_SECRET": "secret-value",
	}

	assert.Equal(t, expected, secrets)
}

type mockSecretClient struct {
	returnSecretValue []byte
}

func (m *mockSecretClient) GetSecret(ctx context.Context, ref provider.SecretRef) ([]byte, error) {
	return m.returnSecretValue, nil
}
