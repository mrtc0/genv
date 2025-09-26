package secretmanager

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/secretutil"
	"google.golang.org/api/option"
)

var _ provider.SecretClient = &SecretManagerClient{}

type SecretManagerClientInterface interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
}

type SecretManagerClient struct {
	projectID string
	location  string
	Client    SecretManagerClientInterface
}

func NewSecretManager(ctx context.Context, projectID, location string) (*SecretManagerClient, error) {
	var opts []option.ClientOption
	if location != "" {
		endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", location)
		opts = append(opts, option.WithEndpoint(endpoint))
	}

	client, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &SecretManagerClient{projectID: projectID, location: location, Client: client}, nil
}

func (s *SecretManagerClient) GetSecret(ctx context.Context, ref provider.SecretRef) ([]byte, error) {
	result, err := s.Client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: s.buildResourceName(ref.Key),
	})
	if err != nil {
		return nil, err
	}

	if ref.Property == "" {
		return result.Payload.Data, nil
	}

	val, err := secretutil.GetValueFromJSON(result.Payload.Data, ref.Property)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s *SecretManagerClient) buildResourceName(secretName string) string {
	if s.location != "" {
		return fmt.Sprintf("projects/%s/locations/%s/secrets/%s/versions/latest", s.projectID, s.location, secretName)
	}

	return fmt.Sprintf("projects/%s/secrets/%s/versions/latest", s.projectID, secretName)
}
