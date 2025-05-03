package secretsmanager

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/mrtc0/genv/provider"
	"github.com/tidwall/gjson"
)

var _ provider.SecretClient = &SecretsManager{}

// SecretsManager is a client for AWS Secrets Manager.
type SecretsManager struct {
	client *awssm.Client
}

// NewSecretsManager creates a new SecretsManager client.
func NewSecretsManager(cfg aws.Config) *SecretsManager {
	client := awssm.NewFromConfig(cfg)
	return &SecretsManager{
		client: client,
	}
}

// GetSecret retrieves a secret from AWS Secrets Manager.
func (s *SecretsManager) GetSecret(ctx context.Context, ref provider.SecretRef) ([]byte, error) {
	secret, err := s.fetch(ctx, ref.Key)
	if err != nil {
		return nil, err
	}

	if ref.Property == "" {
		return secret, nil
	}

	val, err := s.getValueFromJSON(secret, ref.Property)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s *SecretsManager) fetch(ctx context.Context, key string) ([]byte, error) {
	version := "AWSCURRENT"

	input := &awssm.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String(version),
	}

	result, err := s.client.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.SecretString != nil {
		return []byte(*result.SecretString), nil
	}

	return result.SecretBinary, nil
}

func (s *SecretsManager) getValueFromJSON(secret []byte, property string) ([]byte, error) {
	result := gjson.Get(string(secret), property)
	if !result.Exists() {
		return nil, errors.New("property not found in secret")
	}

	return []byte(result.String()), nil
}
