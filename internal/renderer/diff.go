package renderer

import (
	"fmt"
	"strings"

	"github.com/mrtc0/genv/diff"
)

// RenderDiff formats the diff output for the command line.
// It returns a string with the added, removed, and changed env variables.
//
// The output is formatted as follows:
// - Added variables are prefixed with "+ "
// - Removed variables are prefixed with "- "
// - Changed variables are prefixed with "~ "
func RenderDiff(diff *diff.Diff) string {
	// The output is aligned with the longest key length for better readability.
	// Each line contains the key, padding, and value.
	paddingBase := maxKeyLength(diff)

	result := ""
	for key, value := range diff.Added {
		padding := strings.Repeat(" ", paddingBase-len(key)+2)
		result += fmt.Sprintf("+ %s%s=  %q\n", key, padding, value)
	}

	for key, value := range diff.Removed {
		padding := strings.Repeat(" ", paddingBase-len(key)+2)
		result += fmt.Sprintf("- %s%s=  %q\n", key, padding, value)
	}

	for key, value := range diff.Changed {
		padding := strings.Repeat(" ", paddingBase-len(key)+2)
		result += fmt.Sprintf("~ %s%s=  %q => %q\n", key, padding, value.OldValue, value.NewValue)
	}

	return result
}

func maxKeyLength(diff *diff.Diff) int {
	max := 0
	for key := range diff.Added {
		if len(key) > max {
			max = len(key)
		}
	}
	for key := range diff.Removed {
		if len(key) > max {
			max = len(key)
		}
	}
	for key := range diff.Changed {
		if len(key) > max {
			max = len(key)
		}
	}
	return max
}
