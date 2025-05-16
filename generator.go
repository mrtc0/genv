package genv

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mrtc0/genv/provider"
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
	secretProviderClients map[string]provider.SecretClient
}

func NewDotenvGenerator(ctx context.Context, config DotenvGeneratorConfig) (*DotenvGenerator, error) {
	secretProviderClients, err := NewSecretProviderClient(ctx, config.Config.SecretProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret provider client: %w", err)
	}

	return &DotenvGenerator{
		OutputFilePath:        config.OutputFilePath,
		Config:                config.Config,
		secretProviderClients: secretProviderClients,
	}, nil
}

func (d DotenvGenerator) Generate(ctx context.Context) error {
	for key, envValue := range d.Config.Envs {
		if envValue.Value != "" {
			dotenv[key] = envValue.Value
			continue
		}

		if envValue.SecretRef != nil {
			client, ok := d.secretProviderClients[envValue.SecretRef.Provider]
			if !ok {
				return errors.New("secret provider client not found")
			}

			secret, err := client.GetSecret(ctx, provider.SecretRef{
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

func (d *DotenvGenerator) AddSecretProviderClient(providerID string, client provider.SecretClient) {
	if d.secretProviderClients == nil {
		d.secretProviderClients = make(map[string]provider.SecretClient)
	}

	d.secretProviderClients[providerID] = client
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
