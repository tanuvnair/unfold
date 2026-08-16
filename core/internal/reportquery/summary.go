package reportquery

import (
	"sort"
	"strconv"
	"strings"

	"github.com/tanuvnair/unfold/internal/report"
)

const (
	monthlyMinDays = 25.0
	monthlyMaxDays = 35.0
)

// SummaryQuery filters entries before grouping (same filters as Query, no page).
type SummaryQuery struct {
	Search     string
	Type       string
	Confidence string
	Source     string
}

// SummaryGroup is one payee aggregate for the grouped results view.
type SummaryGroup struct {
	Payee            string  `json:"payee"`
	OccurrenceCount  int     `json:"occurrence_count"`
	TotalAmount      float64 `json:"total_amount"`
	AvgAmount        float64 `json:"avg_amount"`
	LatestAmount     float64 `json:"latest_amount"`
	FirstSeen        string  `json:"first_seen"`
	LastSeen         string  `json:"last_seen"`
	Confidence       string  `json:"confidence"`
	Source           string  `json:"source"`
	MonthlyEstimate  float64 `json:"monthly_estimate"`
	IsMonthlyCadence bool    `json:"is_monthly_cadence"`
}

// SummaryResult is the grouped payee view plus a monthly spend estimate.
type SummaryResult struct {
	Groups                []SummaryGroup `json:"groups"`
	EstimatedMonthlyTotal float64        `json:"estimated_monthly_total"`
	GroupCount            int            `json:"group_count"`
}

// Summary aggregates filtered entries by PayeeToken.
func Summary(rows []report.Entry, q SummaryQuery) SummaryResult {
	nq := Query{
		Search:     q.Search,
		Type:       q.Type,
		Confidence: q.Confidence,
		Source:     q.Source,
	}.normalized()

	filtered := filterEntries(rows, nq)
	byPayee := make(map[string][]report.Entry)
	order := make([]string, 0)
	for _, e := range filtered {
		token := strings.TrimSpace(e.PayeeToken)
		if token == "" {
			token = report.DescriptionOf(e.Fields)
		}
		if token == "" {
			token = "(unknown)"
		}
		if _, ok := byPayee[token]; !ok {
			order = append(order, token)
		}
		byPayee[token] = append(byPayee[token], e)
	}

	groups := make([]SummaryGroup, 0, len(order))
	var monthlyTotal float64
	for _, token := range order {
		g := buildSummaryGroup(token, byPayee[token])
		groups = append(groups, g)
		if g.IsMonthlyCadence {
			monthlyTotal += g.MonthlyEstimate
		}
	}

	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].OccurrenceCount != groups[j].OccurrenceCount {
			return groups[i].OccurrenceCount > groups[j].OccurrenceCount
		}
		return groups[i].Payee < groups[j].Payee
	})

	return SummaryResult{
		Groups:                groups,
		EstimatedMonthlyTotal: monthlyTotal,
		GroupCount:            len(groups),
	}
}

func filterEntries(rows []report.Entry, q Query) []report.Entry {
	search := strings.ToLower(strings.TrimSpace(q.Search))
	typeFilter := strings.ToUpper(strings.TrimSpace(q.Type))
	confidenceFilter := strings.ToLower(strings.TrimSpace(q.Confidence))
	sourceFilter := strings.ToLower(strings.TrimSpace(q.Source))

	out := make([]report.Entry, 0, len(rows))
	for _, row := range rows {
		if search != "" && !strings.Contains(strings.ToLower(report.DescriptionOf(row.Fields)), search) {
			continue
		}
		if typeFilter != "" {
			if strings.ToUpper(strings.TrimSpace(lookup(row.Fields, "Dr / Cr"))) != typeFilter {
				continue
			}
		}
		if confidenceFilter != "" {
			if strings.ToLower(strings.TrimSpace(row.Confidence)) != confidenceFilter {
				continue
			}
		}
		if sourceFilter != "" {
			if strings.ToLower(strings.TrimSpace(row.DetectionSource)) != sourceFilter {
				continue
			}
		}
		out = append(out, row)
	}
	return out
}

func buildSummaryGroup(payee string, entries []report.Entry) SummaryGroup {
	sort.SliceStable(entries, func(i, j int) bool {
		return lookup(entries[i].Fields, "Transaction Date") < lookup(entries[j].Fields, "Transaction Date")
	})

	var total float64
	var amounts []float64
	confidence := ""
	sources := map[string]struct{}{}
	var avgInterval float64
	var hasInterval bool

	for _, e := range entries {
		amt := parseAmountString(lookup(e.Fields, "Amount"))
		total += amt
		amounts = append(amounts, amt)
		confidence = maxConfidence(confidence, e.Confidence)
		if e.DetectionSource != "" {
			sources[e.DetectionSource] = struct{}{}
		}
		if e.HasRecurrenceMetrics && e.AvgIntervalDays > 0 {
			avgInterval = e.AvgIntervalDays
			hasInterval = true
		}
	}

	latest := 0.0
	if len(amounts) > 0 {
		latest = amounts[len(amounts)-1]
	}
	avg := 0.0
	if len(amounts) > 0 {
		avg = total / float64(len(amounts))
	}

	monthly := false
	if hasInterval && avgInterval >= monthlyMinDays && avgInterval <= monthlyMaxDays {
		monthly = true
	} else if len(entries) >= 2 && !hasInterval {
		// Keyword-only multi-hit groups: treat as monthly if ≥2 occurrences.
		monthly = true
	}

	estimate := 0.0
	if monthly {
		estimate = latest
	}

	return SummaryGroup{
		Payee:            payee,
		OccurrenceCount:  len(entries),
		TotalAmount:      total,
		AvgAmount:        avg,
		LatestAmount:     latest,
		FirstSeen:        lookup(entries[0].Fields, "Transaction Date"),
		LastSeen:         lookup(entries[len(entries)-1].Fields, "Transaction Date"),
		Confidence:       confidence,
		Source:           mergeSources(sources),
		MonthlyEstimate:  estimate,
		IsMonthlyCadence: monthly,
	}
}

func parseAmountString(raw string) float64 {
	cleaned := strings.ReplaceAll(strings.TrimSpace(raw), ",", "")
	cleaned = strings.ReplaceAll(cleaned, "₹", "")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return 0
	}
	n, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}
	return n
}

func mergeSources(sources map[string]struct{}) string {
	_, kw := sources[report.SourceKeyword]
	_, rec := sources[report.SourceRecurrence]
	_, both := sources[report.SourceBoth]
	if both || (kw && rec) {
		return report.SourceBoth
	}
	if kw {
		return report.SourceKeyword
	}
	if rec {
		return report.SourceRecurrence
	}
	return ""
}

func maxConfidence(a, b string) string {
	rank := func(t string) int {
		switch strings.ToLower(t) {
		case "high":
			return 3
		case "medium":
			return 2
		case "low":
			return 1
		default:
			return 0
		}
	}
	if rank(b) > rank(a) {
		return b
	}
	return a
}
