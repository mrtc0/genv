package genv_test

import (
	"testing"

	"github.com/mrtc0/genv"
	"github.com/stretchr/testify/assert"
)

func TestDiffEnvMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		old, new map[string]string
		expected genv.Diff
	}{
		"empty": {
			old: map[string]string{},
			new: map[string]string{},
			expected: genv.Diff{
				Added:   map[string]string{},
				Removed: map[string]string{},
				Changed: map[string]genv.ChangeValue{},
			},
		},
		"added": {
			old: map[string]string{"A": "1"},
			new: map[string]string{"A": "1", "B": "2"},
			expected: genv.Diff{
				Added:   map[string]string{"B": "2"},
				Removed: map[string]string{},
				Changed: map[string]genv.ChangeValue{},
			},
		},
		"removed": {
			old: map[string]string{"A": "1", "B": "2"},
			new: map[string]string{"A": "1"},
			expected: genv.Diff{
				Added:   map[string]string{},
				Removed: map[string]string{"B": "2"},
				Changed: map[string]genv.ChangeValue{},
			},
		},
		"changed": {
			old: map[string]string{"A": "1", "B": "2"},
			new: map[string]string{"A": "1", "B": "3"},
			expected: genv.Diff{
				Added:   map[string]string{},
				Removed: map[string]string{},
				Changed: map[string]genv.ChangeValue{
					"B": {NewValue: "3", OldValue: "2"},
				},
			},
		},
		"added, removed, changed": {
			old: map[string]string{"A": "1", "B": "2"},
			new: map[string]string{"A": "2", "C": "3"},
			expected: genv.Diff{
				Added:   map[string]string{"C": "3"},
				Removed: map[string]string{"B": "2"},
				Changed: map[string]genv.ChangeValue{
					"A": {NewValue: "2", OldValue: "1"},
				},
			},
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := genv.DiffEnvMap(tt.old, tt.new)
			assert.Equal(t, tt.expected, got)
		})
	}
}
