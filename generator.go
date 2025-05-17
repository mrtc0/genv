package genv

import (
	"context"
	"fmt"
)

type DotenvGeneratorConfig struct {
	OutputFilePath string
	Config         *Config
}

type DotenvGenerator struct {
	OutputFilePath        string
	Config                *Config
	SecretProviderService *SecretProviderService
}

func NewDotenvGenerator(ctx context.Context, config DotenvGeneratorConfig) (*DotenvGenerator, error) {
	svc, err := NewSecretProviderService(ctx, config.Config.SecretProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret provider, %w", err)
	}

	return &DotenvGenerator{
		OutputFilePath:        config.OutputFilePath,
		Config:                config.Config,
		SecretProviderService: svc,
	}, nil
}

func (d *DotenvGenerator) FetchSecrets(ctx context.Context) (map[string]string, error) {
	envMap := make(map[string]string)

	for key, envValue := range d.Config.Envs {
		if envValue.Value != "" {
			envMap[key] = envValue.Value
			continue
		}

		if envValue.SecretRef != nil {
			secret, err := d.SecretProviderService.GetSecret(ctx, envValue.SecretRef.Provider, GetSecretInput{
				Key:      envValue.SecretRef.Key,
				Property: envValue.SecretRef.Property,
			})

			if err != nil {
				return nil, fmt.Errorf("failed to get secret %s: %w", key, err)
			}

			envMap[key] = string(secret)
			continue
		}

		envMap[key] = ""
	}
	return envMap, nil
}
