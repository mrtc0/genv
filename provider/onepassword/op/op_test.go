package op_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/onepassword/op"
	"github.com/stretchr/testify/assert"
)

func TestOPClient_GetSecret(t *testing.T) {
	t.Parallel()

	type want struct {
		secret []byte
		err    error
	}

	testCases := map[string]struct {
		account      string
		ref          provider.SecretRef
		mockExecutor op.OPCommandExecutor
		want         want
	}{
		"without account": {
			account: "",
			ref: provider.SecretRef{
				Key: "op://vault/item/field",
			},
			mockExecutor: &MockOPCommandExecutor{
				ExecFunc: func(ctx context.Context, args []string) ([]byte, error) {
					assert.Equal(t, []string{"read", "op://vault/item/field"}, args)
					return []byte("secret-value"), nil
				},
			},
			want: want{
				secret: []byte("secret-value"),
				err:    nil,
			},
		},
		"with account": {
			account: "my-account",
			ref: provider.SecretRef{
				Key: "op://vault/item/field",
			},
			mockExecutor: &MockOPCommandExecutor{
				ExecFunc: func(ctx context.Context, args []string) ([]byte, error) {
					assert.Equal(t, []string{"--account", "my-account", "read", "op://vault/item/field"}, args)
					return []byte("secret-value"), nil
				},
			},
			want: want{
				secret: []byte("secret-value"),
				err:    nil,
			},
		},
		"command error": {
			account: "my-account",
			ref: provider.SecretRef{
				Key: "op://vault/item/field",
			},
			mockExecutor: &MockOPCommandExecutor{
				ExecFunc: func(ctx context.Context, args []string) ([]byte, error) {
					return []byte("error output"), assert.AnError
				},
			},
			want: want{
				secret: nil,
				err:    errors.New("op command failed: " + assert.AnError.Error() + ": error output"),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := op.NewOPClient(tc.account)
			client.Executor = tc.mockExecutor

			got, err := client.GetSecret(context.Background(), tc.ref)

			if tc.want.err != nil {
				assert.ErrorContains(t, err, tc.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.secret, got)
		})
	}
}

type MockOPCommandExecutor struct {
	ExecFunc func(ctx context.Context, args []string) ([]byte, error)
	Called   bool
}

func (m *MockOPCommandExecutor) Exec(ctx context.Context, args []string) ([]byte, error) {
	m.Called = true
	return m.ExecFunc(ctx, args)
}
