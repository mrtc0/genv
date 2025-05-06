package genv_test

import (
	"context"
	"os"
	"testing"

	"github.com/mrtc0/genv"
	"github.com/mrtc0/genv/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDotenvGenerator_Generate(t *testing.T) {
	t.Parallel()

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

	file, err := os.CreateTemp("", ".env")
	require.NoError(t, err)

	defer os.Remove(file.Name())

	generator := &genv.DotenvGenerator{
		OutputFilePath: file.Name(),
		Config:         config,
	}

	generator.AddSecretProviderClient("example-account", &mockSecretClient{})
	err = generator.Generate(context.Background())
	assert.NoError(t, err)

	dotenv, err := os.ReadFile(file.Name())
	require.NoError(t, err)

	expect := `EXAMPLE_ENV=example-value
EXAMPLE_SECRET=secret-value
`

	assert.Equal(t, expect, string(dotenv))
}

type mockSecretClient struct{}

func (m *mockSecretClient) GetSecret(ctx context.Context, ref provider.SecretRef) ([]byte, error) {
	return []byte("secret-value"), nil
}
