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
	// At least three spaced debits — two hits (or same-day retries) are too
	// easy to confuse with one-off UPI purchases.
	minOccurrences = 3
	maxAmountVariance = 0.05 // 5%
	minIntervalDays   = 25.0
	maxIntervalDays   = 35.0
	looseIntervalDays = 10.0 // half-width around ~30 for medium confidence
)

var (
	digitRunRe   = regexp.MustCompile(`\d{4,}`)
	dateLikeRe   = regexp.MustCompile(`\b\d{1,2}[-/]\d{1,2}([-/]\d{2,4})?\b`)
	noiseTokens  = regexp.MustCompile(`\b(UPI|NEFT|IMPS|RTGS|ACH|NACH|MANDATE|AUTOPAY|PAYMENT|PAY|REF|TXN|TRANSFER|REFUND|REVERSAL)\b`)
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

// GroupByPayee clusters debit transactions by NormalizePayee(description).
// Credits and empty tokens are skipped — autopay charges are outgoing.
func GroupByPayee(transactions []txn.Transaction) []Group {
	buckets := make(map[string][]txn.Transaction)
	order := make([]string, 0)
	for _, t := range transactions {
		if isCredit(t) {
			continue
		}
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
	txns := append([]txn.Transaction(nil), g.Transactions...)
	sort.Slice(txns, func(i, j int) bool {
		return txns[i].Date.Before(txns[j].Date)
	})
	txns = collapseSameCalendarDay(txns)
	if len(txns) < minOccurrences {
		return Detection{}, false
	}

	variance := amountVarianceRatio(txns)
	if variance > maxAmountVariance {
		return Detection{}, false
	}

	avgInterval, tight, loose := intervalStats(txns)
	if avgInterval <= 0 || !loose {
		return Detection{}, false
	}

	confidence := config.TierMedium
	if tight {
		confidence = config.TierHigh
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

func isCredit(t txn.Transaction) bool {
	return strings.EqualFold(strings.TrimSpace(t.Type), "CR")
}

// collapseSameCalendarDay keeps the first debit per UTC calendar day so a
// failed UPI retry a minute later is not treated as a second monthly hit.
func collapseSameCalendarDay(txns []txn.Transaction) []txn.Transaction {
	if len(txns) == 0 {
		return txns
	}
	out := make([]txn.Transaction, 0, len(txns))
	out = append(out, txns[0])
	for i := 1; i < len(txns); i++ {
		prev := out[len(out)-1].Date
		cur := txns[i].Date
		if prev.Year() == cur.Year() && prev.YearDay() == cur.YearDay() {
			continue
		}
		out = append(out, txns[i])
	}
	return out
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

func intervalStats(txns []txn.Transaction) (avg float64, tight, loose bool) {
	if len(txns) < 2 {
		return 0, false, false
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
		return 0, false, false
	}
	avg = sum / float64(len(intervals))
	tight = true
	loose = true
	center := 30.0
	for _, d := range intervals {
		if d < minIntervalDays || d > maxIntervalDays {
			tight = false
		}
		if d < center-looseIntervalDays || d > center+looseIntervalDays {
			loose = false
		}
	}
	return avg, tight, loose
}
