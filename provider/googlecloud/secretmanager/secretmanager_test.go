package secretmanager_test

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/googlecloud/secretmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dummyProjectID = "dummy-project"
)

func TestGetSecret(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	type arrange struct {
		location   string
		mockClient *mockSecretManagerClient
	}

	type want struct {
		secret []byte
		err    bool
	}

	testCases := map[string]struct {
		arrange arrange
		ref     provider.SecretRef
		want    want
	}{
		"get secret without property": {
			arrange: arrange{
				mockClient: &mockSecretManagerClient{
					AccessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
						assert.Equal(t, req.Name, fmt.Sprintf("projects/%s/secrets/my-secret/versions/latest", dummyProjectID))
						return &secretmanagerpb.AccessSecretVersionResponse{
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte("my-secret-value"),
							},
						}, nil
					},
				},
			},
			ref: provider.SecretRef{
				Key: "my-secret",
			},
			want: want{
				secret: []byte("my-secret-value"),
			},
		},
		"get secret with property": {
			arrange: arrange{
				mockClient: &mockSecretManagerClient{
					AccessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
						assert.Equal(t, req.Name, fmt.Sprintf("projects/%s/secrets/my-json-secret/versions/latest", dummyProjectID))
						return &secretmanagerpb.AccessSecretVersionResponse{
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte(`{"username": "my-username", "password": "my-secret-password"}`),
							},
						}, nil
					},
				},
			},
			ref: provider.SecretRef{
				Key:      "my-json-secret",
				Property: "password",
			},
			want: want{
				secret: []byte("my-secret-password"),
			},
		},
		"when regional secret": {
			arrange: arrange{
				mockClient: &mockSecretManagerClient{
					AccessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
						assert.Equal(t, req.Name, fmt.Sprintf("projects/%s/locations/us-central1/secrets/my-regional-secret/versions/latest", dummyProjectID))
						return &secretmanagerpb.AccessSecretVersionResponse{
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte("my-regional-secret-value"),
							},
						}, nil
					},
				},
				location: "us-central1",
			},
			ref: provider.SecretRef{
				Key: "my-regional-secret",
			},
			want: want{
				secret: []byte("my-regional-secret-value"),
			},
		},
		"get non-existing secret": {
			arrange: arrange{
				mockClient: &mockSecretManagerClient{
					AccessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
						return nil, fmt.Errorf("not found")
					},
				},
			},
			ref: provider.SecretRef{
				Key: "non-existing-secret",
			},
			want: want{
				secret: nil,
				err:    true,
			},
		},
		"get secret with non-existing property": {
			arrange: arrange{
				mockClient: &mockSecretManagerClient{
					AccessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
						assert.Equal(t, req.Name, fmt.Sprintf("projects/%s/secrets/existing-secret/versions/latest", dummyProjectID))
						return &secretmanagerpb.AccessSecretVersionResponse{
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte(`{"username": "my-username", "password": "my-secret-password"}`),
							},
						}, nil
					},
				},
			},
			ref: provider.SecretRef{
				Key:      "existing-secret",
				Property: "non-existing-property",
			},
			want: want{
				secret: nil,
				err:    true,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &secretmanager.SecretManagerClient{
				ProjectID: dummyProjectID,
				Location:  tc.arrange.location,
				Client:    tc.arrange.mockClient,
			}

			secret, err := client.GetSecret(ctx, tc.ref)
			if tc.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want.secret, secret)
			assert.True(t, tc.arrange.mockClient.Called)
		})
	}
}

type mockSecretManagerClient struct {
	AccessSecretVersionFunc func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)

	Called bool
}

func (m *mockSecretManagerClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	if m.AccessSecretVersionFunc != nil {
		m.Called = true
		return m.AccessSecretVersionFunc(ctx, req, opts...)
	}

	panic("AccessSecretVersionFunc not implemented")
}
