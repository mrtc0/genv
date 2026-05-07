package exec_test

import (
	"context"
	"testing"

	"github.com/mrtc0/genv/provider"
	execprovider "github.com/mrtc0/genv/provider/exec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetSecret(t *testing.T) {
	t.Parallel()

	type want struct {
		secret []byte
		errMsg string // empty means no error expected
	}

	testCases := map[string]struct {
		command []string
		ref     provider.SecretRef
		want    want
	}{
		"simple key lookup": {
			// command outputs a flat JSON object; ref.Key picks a top-level field
			command: []string{"sh", "-c", `echo '{"api_key":"my-secret-value"}'`},
			ref: provider.SecretRef{
				Key: "api_key",
			},
			want: want{
				secret: []byte("my-secret-value"),
			},
		},
		"nested key with gjson dot notation": {
			// gjson supports dot-notation for nested fields
			command: []string{"sh", "-c", `echo '{"data":{"db_url":"postgres://localhost/mydb"}}'`},
			ref: provider.SecretRef{
				Key: "data.db_url",
			},
			want: want{
				secret: []byte("postgres://localhost/mydb"),
			},
		},
		"with Property for double-encoded JSON": {
			// The secret value itself is a JSON string (e.g. AWS Secrets Manager style).
			// ref.Key fetches the raw string, ref.Property drills into the parsed result.
			command: []string{"sh", "-c", `echo '{"secret":"{\"token\":\"abc123\"}"}'`},
			ref: provider.SecretRef{
				Key:      "secret",
				Property: "token",
			},
			want: want{
				secret: []byte("abc123"),
			},
		},
		"key not found in output": {
			command: []string{"sh", "-c", `echo '{"api_key":"value"}'`},
			ref: provider.SecretRef{
				Key: "nonexistent",
			},
			want: want{
				errMsg: `exec provider: key "nonexistent" not found in output`,
			},
		},
		"property not found in value": {
			command: []string{"sh", "-c", `echo '{"secret":"{\"token\":\"abc123\"}"}'`},
			ref: provider.SecretRef{
				Key:      "secret",
				Property: "nonexistent",
			},
			want: want{
				errMsg: `exec provider: property "nonexistent" not found`,
			},
		},
		"command exits with non-zero status": {
			command: []string{"sh", "-c", "exit 1"},
			ref: provider.SecretRef{
				Key: "key",
			},
			want: want{
				errMsg: "exec provider: command failed",
			},
		},
		"command not found": {
			command: []string{"__nonexistent_command_genv__"},
			ref: provider.SecretRef{
				Key: "key",
			},
			want: want{
				errMsg: "exec provider: command failed",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg := &execprovider.ExecProviderConfig{
				ID:      "test",
				Command: tc.command,
			}
			p := execprovider.NewProvider(cfg)
			client, err := p.NewClient(context.Background())
			require.NoError(t, err)

			ctx := context.Background()
			got, err := client.GetSecret(ctx, tc.ref)

			if tc.want.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.want.errMsg)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want.secret, got)
		})
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	cfg := &execprovider.ExecProviderConfig{
		ID:      "test",
		Command: []string{},
	}
	p := execprovider.NewProvider(cfg)
	_, err := p.NewClient(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "command must not be empty")
}
