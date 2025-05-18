package diff_test

import (
	"testing"

	"github.com/mrtc0/genv/diff"
	"github.com/stretchr/testify/assert"
)

func TestDiffEnvMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		old, new map[string]string
		expected diff.Diff
	}{
		"empty": {
			old: map[string]string{},
			new: map[string]string{},
			expected: diff.Diff{
				Added:   map[string]string{},
				Removed: map[string]string{},
				Changed: map[string]diff.ChangeValue{},
			},
		},
		"added": {
			old: map[string]string{"A": "1"},
			new: map[string]string{"A": "1", "B": "2"},
			expected: diff.Diff{
				Added:   map[string]string{"B": "2"},
				Removed: map[string]string{},
				Changed: map[string]diff.ChangeValue{},
			},
		},
		"removed": {
			old: map[string]string{"A": "1", "B": "2"},
			new: map[string]string{"A": "1"},
			expected: diff.Diff{
				Added:   map[string]string{},
				Removed: map[string]string{"B": "2"},
				Changed: map[string]diff.ChangeValue{},
			},
		},
		"changed": {
			old: map[string]string{"A": "1", "B": "2"},
			new: map[string]string{"A": "1", "B": "3"},
			expected: diff.Diff{
				Added:   map[string]string{},
				Removed: map[string]string{},
				Changed: map[string]diff.ChangeValue{
					"B": {NewValue: "3", OldValue: "2"},
				},
			},
		},
		"added, removed, changed": {
			old: map[string]string{"A": "1", "B": "2"},
			new: map[string]string{"A": "2", "C": "3"},
			expected: diff.Diff{
				Added:   map[string]string{"C": "3"},
				Removed: map[string]string{"B": "2"},
				Changed: map[string]diff.ChangeValue{
					"A": {NewValue: "2", OldValue: "1"},
				},
			},
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := diff.DiffEnvMap(tt.old, tt.new)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestDiff_IsChanged(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		diff     diff.Diff
		expected bool
	}{
		"empty": {
			diff: diff.Diff{
				Added:   make(map[string]string),
				Removed: make(map[string]string),
				Changed: make(map[string]diff.ChangeValue),
			},
			expected: false,
		},
		"added": {
			diff: diff.Diff{
				Added: map[string]string{"A": "1"},
			},
			expected: true,
		},
		"removed": {
			diff: diff.Diff{
				Removed: map[string]string{"A": "1"},
			},
			expected: true,
		},
		"changed": {
			diff: diff.Diff{
				Changed: map[string]diff.ChangeValue{"A": {NewValue: "2", OldValue: "1"}},
			},
			expected: true,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.diff.IsChanged()
			assert.Equal(t, tt.expected, got)
		})
	}
}
