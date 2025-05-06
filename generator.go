package genv

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/aws"
)

var (
	dotenv map[string]string = make(map[string]string)
)

const (
	defaultGenvFilePath   = ".genv.yaml"
	defaultOutputFilePath = ".env"
)

type DotenvGenerator interface {
	Generate(ctx context.Context) error
}

type DotenvGeneratorConfig struct {
	OutputFilePath string
	Config         *Config
}

type dotnevGenerator struct {
	OutputFilePath string
	Config         *Config
}

func NewDotenvGenerator(config DotenvGeneratorConfig) DotenvGenerator {
	return &dotnevGenerator{
		OutputFilePath: config.OutputFilePath,
		Config:         config.Config,
	}
}

func (d *dotnevGenerator) Generate(ctx context.Context) error {
	for _, p := range d.Config.SecretProvider.Aws {
		providerConfig := &aws.AwsProviderConfig{
			ID:      p.ID,
			Service: aws.AWSSecretsManager,
			Region:  p.Region,
			Auth: aws.AwsAuth{
				Profile: p.Auth.Profile,
			},
		}

		awsProvider := aws.NewProvider(providerConfig)
		provider.SecretProviders[p.ID] = awsProvider
	}

	for key, envValue := range d.Config.Envs {
		if envValue.Value != "" {
			dotenv[key] = envValue.Value
			continue
		}

		if envValue.SecretRef != nil {
			providerID := envValue.SecretRef.Provider

			p, ok := provider.SecretProviders[providerID]
			if !ok {
				return errors.New("provider not found")
			}

			client, err := p.NewClient(ctx)
			if err != nil {
				return err
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
