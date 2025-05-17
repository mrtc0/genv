package genv

import (
	"context"
	"fmt"
	"os"
)

var (
	dotenv map[string]string = make(map[string]string)
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

func (d DotenvGenerator) Generate(ctx context.Context) error {
	for key, envValue := range d.Config.Envs {
		if envValue.Value != "" {
			dotenv[key] = envValue.Value
			continue
		}

		if envValue.SecretRef != nil {
			secret, err := d.SecretProviderService.GetSecret(ctx, envValue.SecretRef.Provider, GetSecretInput{
				Key:      envValue.SecretRef.Key,
				Property: envValue.SecretRef.Property,
			})

			if err != nil {
				return fmt.Errorf("failed to get secret %s: %w", key, err)
			}

			dotenv[key] = string(secret)
			continue
		}

		dotenv[key] = ""
	}

	if err := writeDotenvFile(d.OutputFilePath, dotenv); err != nil {
		return err
	}

	return nil
}

func writeDotenvFile(filePath string, dotenv map[string]string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", filePath, err)
	}
	defer f.Close()

	for key, value := range dotenv {
		if _, err := f.WriteString(key + "=" + value + "\n"); err != nil {
			return err
		}
	}

	return nil
}
