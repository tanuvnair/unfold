// Package recurrence detects recurring charges by grouping transactions on a
// normalized payee token and checking for stable amount + regular cadence,
// independent of keyword matching. This catches autopay charges whose
// description never contains NACH/MANDATE/AUTOPAY language.
package recurrence

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/txn"
)

const (
	minOccurrences     = 2
	highMinOccurrences = 3
	maxAmountVariance  = 0.05 // 5%
	minIntervalDays    = 25.0
	maxIntervalDays    = 35.0
	looseIntervalDays  = 10.0 // half-width around ~30 for medium confidence
)

var (
	digitRunRe   = regexp.MustCompile(`\d{4,}`)
	dateLikeRe   = regexp.MustCompile(`\b\d{1,2}[-/]\d{1,2}([-/]\d{2,4})?\b`)
	noiseTokens  = regexp.MustCompile(`\b(UPI|NEFT|IMPS|RTGS|ACH|NACH|MANDATE|AUTOPAY|PAYMENT|PAY|REF|TXN|TRANSFER)\b`)
	slashNoiseRe = regexp.MustCompile(`/+`)
)

// Group is transactions that share a normalized payee token.
type Group struct {
	PayeeToken   string
	Transactions []txn.Transaction
}

// Detection is a payee group flagged as recurring.
type Detection struct {
	Group           Group
	Confidence      string
	AmountVariance  float64
	AvgIntervalDays float64
}

// NormalizePayee strips UPI reference numbers, dates, and known noise tokens
// from a description. Fully offline — no merchant lookup.
func NormalizePayee(description string) string {
	s := strings.ToUpper(strings.TrimSpace(description))
	s = digitRunRe.ReplaceAllString(s, " ")
	s = dateLikeRe.ReplaceAllString(s, " ")
	s = noiseTokens.ReplaceAllString(s, " ")
	s = strings.Map(func(r rune) rune {
		switch r {
		case '/', '-', '_', '.', ',', '|', '\\', ':', ';', '(', ')', '[', ']', '{', '}':
			return ' '
		default:
			if unicode.IsSpace(r) {
				return ' '
			}
			return r
		}
	}, s)
	s = slashNoiseRe.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(s), " ")
}

// GroupByPayee clusters transactions by NormalizePayee(description).
// Empty tokens are skipped.
func GroupByPayee(transactions []txn.Transaction) []Group {
	buckets := make(map[string][]txn.Transaction)
	order := make([]string, 0)
	for _, t := range transactions {
		token := NormalizePayee(t.Description)
		if token == "" {
			continue
		}
		if _, ok := buckets[token]; !ok {
			order = append(order, token)
		}
		buckets[token] = append(buckets[token], t)
	}
	out := make([]Group, 0, len(order))
	for _, token := range order {
		out = append(out, Group{PayeeToken: token, Transactions: buckets[token]})
	}
	return out
}

// DetectRecurring flags groups with stable amounts and monthly-ish cadence.
func DetectRecurring(groups []Group) []Detection {
	var out []Detection
	for _, g := range groups {
		if d, ok := detectGroup(g); ok {
			out = append(out, d)
		}
	}
	return out
}

func detectGroup(g Group) (Detection, bool) {
	if len(g.Transactions) < minOccurrences {
		return Detection{}, false
	}

	txns := append([]txn.Transaction(nil), g.Transactions...)
	sort.Slice(txns, func(i, j int) bool {
		return txns[i].Date.Before(txns[j].Date)
	})

	variance := amountVarianceRatio(txns)
	if variance > maxAmountVariance {
		return Detection{}, false
	}

	avgInterval, tight := intervalStats(txns)
	if avgInterval <= 0 {
		return Detection{}, false
	}

	confidence := config.TierMedium
	if len(txns) >= highMinOccurrences && tight {
		confidence = config.TierHigh
	} else if !tight && !looseIntervalFit(avgInterval) {
		return Detection{}, false
	}

	return Detection{
		Group: Group{
			PayeeToken:   g.PayeeToken,
			Transactions: txns,
		},
		Confidence:      confidence,
		AmountVariance:  variance,
		AvgIntervalDays: avgInterval,
	}, true
}

func amountVarianceRatio(txns []txn.Transaction) float64 {
	if len(txns) == 0 {
		return 0
	}
	var sum float64
	for _, t := range txns {
		sum += math.Abs(t.Amount)
	}
	mean := sum / float64(len(txns))
	if mean == 0 {
		return 0
	}
	var maxDev float64
	for _, t := range txns {
		dev := math.Abs(math.Abs(t.Amount)-mean) / mean
		if dev > maxDev {
			maxDev = dev
		}
	}
	return maxDev
}

func intervalStats(txns []txn.Transaction) (avg float64, tight bool) {
	if len(txns) < 2 {
		return 0, false
	}
	var sum float64
	var intervals []float64
	for i := 1; i < len(txns); i++ {
		days := txns[i].Date.Sub(txns[i-1].Date).Hours() / 24
		if days <= 0 {
			continue
		}
		intervals = append(intervals, days)
		sum += days
	}
	if len(intervals) == 0 {
		return 0, false
	}
	avg = sum / float64(len(intervals))
	tight = true
	for _, d := range intervals {
		if d < minIntervalDays || d > maxIntervalDays {
			tight = false
			break
		}
	}
	return avg, tight
}

func looseIntervalFit(avg float64) bool {
	center := 30.0
	return avg >= center-looseIntervalDays && avg <= center+looseIntervalDays
}
