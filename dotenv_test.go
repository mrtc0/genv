package genv_test

import (
	"testing"

	"github.com/mrtc0/genv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		envMap   map[string]string
		expected string
	}{
		"multi values": {
			envMap: map[string]string{
				"EXAMPLE_ENV":    "example-value",
				"EXAMPLE_SECRET": "secret-value",
			},
			expected: `EXAMPLE_ENV="example-value"
EXAMPLE_SECRET="secret-value"`,
		},
		"empty map": {
			envMap:   map[string]string{},
			expected: ``,
		},
		"empty value": {
			envMap: map[string]string{
				"EXAMPLE_ENV": "",
			},
			expected: `EXAMPLE_ENV=""`,
		},
		"integer value": {
			envMap: map[string]string{
				"EXAMPLE_ENV": "12345",
			},
			expected: `EXAMPLE_ENV=12345`,
		},
		"value has LF": {
			envMap: map[string]string{
				"LF": "test1\ntest2\ntest3",
			},
			expected: "LF=\"test1\\ntest2\\ntest3\"",
		},
		"value has CR": {
			envMap: map[string]string{
				"CR": "test1\rtest2\rtest3",
			},
			expected: "CR=\"test1\\rtest2\\rtest3\"",
		},
		"values has CRLF": {
			envMap: map[string]string{
				"CRLF": "test1\r\ntest2\r\ntest3",
			},
			expected: "CRLF=\"test1\\r\\ntest2\\r\\ntest3\"",
		},
		"value has double quotes": {
			envMap: map[string]string{
				"EXAMPLE_ENV": `example"value`,
			},
			expected: `EXAMPLE_ENV="example\"value"`,
		},
		"value has backslash": {
			envMap: map[string]string{
				"EXAMPLE_ENV": `example\value`,
			},
			expected: `EXAMPLE_ENV="example\\value"`,
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := genv.Marshal(tt.envMap)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		content  []byte
		expected map[string]string
	}{
		"multi values": {
			content: []byte(`EXAMPLE_ENV="example-value"
EXAMPLE_SECRET="secret-value"`),
			expected: map[string]string{
				"EXAMPLE_ENV":    "example-value",
				"EXAMPLE_SECRET": "secret-value",
			},
		},
		"empty map": {
			content:  []byte(``),
			expected: map[string]string{},
		},
		"empty value": {
			content: []byte(`EXAMPLE_ENV=`),
			expected: map[string]string{
				"EXAMPLE_ENV": "",
			},
		},
		"empty value with double quote": {
			content: []byte(`EXAMPLE_ENV=""`),
			expected: map[string]string{
				"EXAMPLE_ENV": "",
			},
		},
		"empty value with single quote": {
			content: []byte(`EXAMPLE_ENV=''`),
			expected: map[string]string{
				"EXAMPLE_ENV": "",
			},
		},
		"integer value": {
			content: []byte(`EXAMPLE_ENV=12345`),
			expected: map[string]string{
				"EXAMPLE_ENV": "12345",
			},
		},
		"value has LF": {
			content: []byte(`LF="test1\ntest2\ntest3"`),
			expected: map[string]string{
				"LF": "test1\ntest2\ntest3",
			},
		},
		"value has CR": {
			content: []byte(`CR="test1\rtest2\rtest3"`),
			expected: map[string]string{
				"CR": "test1\rtest2\rtest3",
			},
		},
		"value has CRLF": {
			content: []byte(`CRLF="test1\r\ntest2\r\ntest3"`),
			expected: map[string]string{
				"CRLF": "test1\r\ntest2\r\ntest3",
			},
		},
		"value has double quotes": {
			content: []byte(`EXAMPLE_ENV="example\"value"`),
			expected: map[string]string{
				"EXAMPLE_ENV": `example"value`,
			},
		},
		"value has backslash": {
			content: []byte(`EXAMPLE_ENV="example\\value"`),
			expected: map[string]string{
				"EXAMPLE_ENV": `example\value`,
			},
		},
		"comment line": {
			content: []byte(`# comment
EXAMPLE_ENV="example-value"`),
			expected: map[string]string{
				"EXAMPLE_ENV": "example-value",
			},
		},
		"empty line": {
			content: []byte(`
EXAMPLE_ENV="example-value"`),
			expected: map[string]string{
				"EXAMPLE_ENV": "example-value",
			},
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := genv.Unmarshal(tt.content)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
