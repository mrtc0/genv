package genv_test

import (
	"context"
	"testing"

	"github.com/mrtc0/genv"
	"github.com/mrtc0/genv/diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffEnvName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		cfg      *genv.Config
		dotenv   map[string]string
		expected diff.Diff
	}{
		"only the names of environment variables can be differenced": {
			cfg: &genv.Config{
				Envs: map[string]genv.EnvValue{
					"NOT_CHANGED_ENV":   {Value: "not-changed-value"},
					"NEW_DEFINED_ENV":   {Value: "new-defined-value"},
					"VALUE_CHANGED_ENV": {Value: "new-value"},
				},
			},
			dotenv: map[string]string{
				"NOT_CHANGED_ENV":   "not-changed-value",
				"REMOVED_ENV":       "removed-value",
				"VALUE_CHANGED_ENV": "old-value",
			},
			expected: diff.Diff{
				Added:   map[string]string{"NEW_DEFINED_ENV": "(value not fetched)"},
				Removed: map[string]string{"REMOVED_ENV": "(value not fetched)"},
				// The value has changed, but since only the name is used to take the difference,
				// the value is not included in the difference
				Changed: map[string]diff.ChangeValue{},
			},
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := genv.DiffEnvName(context.Background(), tt.cfg, tt.dotenv)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, *actual)
		})
	}
}
