package genv_test

import (
	"context"
	"testing"

	"github.com/mrtc0/genv"
	"github.com/mrtc0/genv/provider"
	"github.com/stretchr/testify/assert"
)

func TestSecretProviderService_GetSecret(t *testing.T) {
	t.Parallel()

	type arrange struct {
		providerID string
		client     provider.SecretClient
	}

	testCases := map[string]struct {
		arrange arrange
		args    genv.GetSecretInput
		want    []byte
		wantErr bool
	}{
		"when secret provider client is not found": {
			arrange: arrange{},
			args: genv.GetSecretInput{
				Key:      "example-key",
				Property: "example-property",
			},
			want:    nil,
			wantErr: true,
		},
		"when secret provider client is found": {
			arrange: arrange{
				providerID: "example-account",
				client: &mockSecretClient{
					returnSecretValue: []byte("example-secret"),
				},
			},
			args: genv.GetSecretInput{
				Key:      "example-key",
				Property: "example-property",
			},
			want:    []byte("example-secret"),
			wantErr: false,
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			s := &genv.SecretProviderService{}
			if tt.arrange.client != nil {
				s.AddSecretProviderClient(tt.arrange.providerID, tt.arrange.client)
			}

			got, err := s.GetSecret(ctx, tt.arrange.providerID, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
