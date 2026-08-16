// Package matcher filters canonical transactions by keyword. It never
// touches raw CSV columns or bank-specific layout — that separation is what
// lets the same matcher work across every bank parser.
package matcher

import (
	"strings"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/txn"
)

// Match is one transaction that hit an include keyword, with the highest
// confidence tier among all include hits and the term that produced it.
type Match struct {
	Transaction txn.Transaction
	Confidence  string
	MatchedTerm string
}

// Filter returns matches whose Description contains at least one include
// keyword and none of the exclude keywords. include must already be
// normalized (upper-case Term, blanks removed); exclude likewise.
func Filter(transactions []txn.Transaction, include []config.KeywordRule, exclude []string) []Match {
	var matched []Match
	for _, t := range transactions {
		if m, ok := match(t, include, exclude); ok {
			matched = append(matched, m)
		}
	}
	return matched
}

func match(t txn.Transaction, include []config.KeywordRule, exclude []string) (Match, bool) {
	// Autopay / mandate hits are outgoing. Credits (refunds) are noise.
	if strings.EqualFold(strings.TrimSpace(t.Type), "CR") {
		return Match{}, false
	}
	description := strings.ToUpper(t.Description)
	confidence, term, ok := bestIncludeHit(description, include)
	if !ok {
		return Match{}, false
	}
	if containsAny(description, exclude) {
		return Match{}, false
	}
	return Match{
		Transaction: t,
		Confidence:  confidence,
		MatchedTerm: term,
	}, true
}

func bestIncludeHit(description string, include []config.KeywordRule) (confidence, term string, ok bool) {
	bestRank := 0
	for _, kw := range include {
		if !strings.Contains(description, kw.Term) {
			continue
		}
		rank := tierRank(kw.Tier)
		if !ok || rank > bestRank {
			ok = true
			bestRank = rank
			confidence = kw.Tier
			term = kw.Term
		}
	}
	return confidence, term, ok
}

func tierRank(tier string) int {
	switch tier {
	case config.TierHigh:
		return 3
	case config.TierMedium:
		return 2
	case config.TierLow:
		return 1
	default:
		return 0
	}
}

func containsAny(description string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(description, kw) {
			return true
		}
	}
	return false
}
