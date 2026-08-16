package parser

import (
	"fmt"
	"strings"
)

// ResolveColumnIndex finds columnName in header (case-insensitive, trimmed).
// Returns an error listing available headers when the column is missing.
func ResolveColumnIndex(header []string, columnName string) (int, error) {
	want := strings.ToUpper(strings.TrimSpace(columnName))
	if want == "" {
		return -1, fmt.Errorf("description column name is empty")
	}

	available := make([]string, 0, len(header))
	for i, col := range header {
		trimmed := strings.TrimSpace(col)
		available = append(available, trimmed)
		if strings.ToUpper(trimmed) == want {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column %q not found in header (available: %v)", columnName, available)
}
