// Package matcher filters canonical transactions by keyword. It never
// touches raw CSV columns or bank-specific layout — that separation is what
// lets the same matcher work across every bank parser.
package matcher

import (
	"strings"

	"github.com/tanuvnair/unfold/internal/txn"
)

// Filter returns only the transactions whose Description contains at least
// one of the given (already normalized, upper-case) keywords.
func Filter(transactions []txn.Transaction, normalizedKeywords []string) []txn.Transaction {
	var matched []txn.Transaction
	for _, t := range transactions {
		if matches(t, normalizedKeywords) {
			matched = append(matched, t)
		}
	}
	return matched
}

func matches(t txn.Transaction, keywords []string) bool {
	description := strings.ToUpper(t.Description)
	for _, kw := range keywords {
		if strings.Contains(description, kw) {
			return true
		}
	}
	return false
}
