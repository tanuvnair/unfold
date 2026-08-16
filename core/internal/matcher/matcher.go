// Package matcher filters canonical transactions by keyword. It never
// touches raw CSV columns or bank-specific layout — that separation is what
// lets the same matcher work across every bank parser.
package matcher

import (
	"strings"

	"github.com/tanuvnair/unfold/internal/txn"
)

// Filter returns transactions whose Description contains at least one
// include keyword and none of the exclude keywords. Both slices must already
// be normalized (upper-case, blanks removed).
func Filter(transactions []txn.Transaction, include, exclude []string) []txn.Transaction {
	var matched []txn.Transaction
	for _, t := range transactions {
		if matches(t, include, exclude) {
			matched = append(matched, t)
		}
	}
	return matched
}

func matches(t txn.Transaction, include, exclude []string) bool {
	description := strings.ToUpper(t.Description)
	if !containsAny(description, include) {
		return false
	}
	if containsAny(description, exclude) {
		return false
	}
	return true
}

func containsAny(description string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(description, kw) {
			return true
		}
	}
	return false
}
