// Package analyze runs the shared statement pipeline used by the CLI and API:
// parse → keyword match + recurrence → merged report.
package analyze

import (
	"fmt"
	"io"
	"strings"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/matcher"
	"github.com/tanuvnair/unfold/internal/parser"
	"github.com/tanuvnair/unfold/internal/recurrence"
	"github.com/tanuvnair/unfold/internal/report"
	"github.com/tanuvnair/unfold/internal/txn"
)

// Run parses statement with the bank parser for cfg, runs keyword matching
// and recurrence detection independently, merges hits, and builds a report.
// It does not write anything to disk.
func Run(cfg config.Config, statement io.Reader) (report.Report, error) {
	bankParser, err := parser.Get(cfg.BankKey())
	if err != nil {
		return report.Report{}, err
	}

	transactions, err := bankParser.Parse(statement, cfg)
	if err != nil {
		return report.Report{}, fmt.Errorf("parse statement: %w", err)
	}

	keywordHits := matcher.Filter(
		transactions,
		cfg.NormalizedKeywords(),
		cfg.NormalizedExcludeKeywords(),
	)
	recHits := recurrence.DetectRecurring(
		recurrence.GroupByPayee(
			excludeDescriptions(transactions, cfg.NormalizedExcludeKeywords()),
		),
	)

	entries := merge(keywordHits, recHits)
	return report.Build(cfg.BankName, entries), nil
}

func excludeDescriptions(transactions []txn.Transaction, exclude []string) []txn.Transaction {
	if len(exclude) == 0 {
		return transactions
	}
	out := make([]txn.Transaction, 0, len(transactions))
	for _, t := range transactions {
		desc := strings.ToUpper(t.Description)
		skip := false
		for _, kw := range exclude {
			if strings.Contains(desc, kw) {
				skip = true
				break
			}
		}
		if !skip {
			out = append(out, t)
		}
	}
	return out
}

type pending struct {
	txn                  txn.Transaction
	keywordConfidence    string
	matchedTerm          string
	hasKeyword           bool
	recurrenceConfidence string
	amountVariance       float64
	avgIntervalDays      float64
	hasRecurrence        bool
}

func merge(keywordHits []matcher.Match, recHits []recurrence.Detection) []report.Entry {
	byFP := make(map[string]*pending)
	order := make([]string, 0)

	get := func(t txn.Transaction) *pending {
		fp := report.FieldsFingerprint(t.Fields)
		if p, ok := byFP[fp]; ok {
			return p
		}
		p := &pending{txn: t}
		byFP[fp] = p
		order = append(order, fp)
		return p
	}

	for _, m := range keywordHits {
		p := get(m.Transaction)
		p.hasKeyword = true
		p.keywordConfidence = m.Confidence
		p.matchedTerm = m.MatchedTerm
	}

	for _, d := range recHits {
		for _, t := range d.Group.Transactions {
			p := get(t)
			p.hasRecurrence = true
			p.recurrenceConfidence = d.Confidence
			p.amountVariance = d.AmountVariance
			p.avgIntervalDays = d.AvgIntervalDays
		}
	}

	out := make([]report.Entry, 0, len(order))
	for _, fp := range order {
		p := byFP[fp]
		out = append(out, toEntry(p))
	}
	return out
}

func toEntry(p *pending) report.Entry {
	e := report.Entry{
		Fields:     p.txn.Fields,
		PayeeToken: recurrence.NormalizePayee(p.txn.Description),
	}
	switch {
	case p.hasKeyword && p.hasRecurrence:
		e.DetectionSource = report.SourceBoth
		e.Confidence = maxConfidence(p.keywordConfidence, p.recurrenceConfidence)
		e.MatchedTerm = p.matchedTerm
		e.AmountVarianceRatio = p.amountVariance
		e.AvgIntervalDays = p.avgIntervalDays
		e.HasRecurrenceMetrics = true
	case p.hasKeyword:
		e.DetectionSource = report.SourceKeyword
		e.Confidence = p.keywordConfidence
		e.MatchedTerm = p.matchedTerm
	case p.hasRecurrence:
		e.DetectionSource = report.SourceRecurrence
		e.Confidence = p.recurrenceConfidence
		e.AmountVarianceRatio = p.amountVariance
		e.AvgIntervalDays = p.avgIntervalDays
		e.HasRecurrenceMetrics = true
	}
	return e
}

func maxConfidence(a, b string) string {
	if tierRank(a) >= tierRank(b) {
		return a
	}
	return b
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
