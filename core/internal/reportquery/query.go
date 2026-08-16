package reportquery

import (
	"strconv"
	"strings"

	"github.com/tanuvnair/unfold/internal/report"
	"github.com/tanuvnair/unfold/internal/txn"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// Query is the server-side search, type filter, and page window.
type Query struct {
	Search     string
	Type       string
	Confidence string
	Source     string
	Payee      string // exact PayeeToken match when set
	Page       int
	PageSize   int
}

// Row is one matched transaction in the shape the results table expects.
type Row struct {
	ID              string `json:"id"`
	TransactionDate string `json:"transaction_date"`
	Description     string `json:"description"`
	Amount          string `json:"amount"`
	Type            string `json:"type"`
	Confidence      string `json:"confidence"`
	Source          string `json:"source"`
}

// Result is one page of filtered transactions plus the filtered total.
type Result struct {
	Rows     []Row `json:"rows"`
	RowCount int   `json:"row_count"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// Apply filters and paginates transactions for the results table.
func Apply(rows []report.Entry, q Query) Result {
	q = q.normalized()

	filtered := make([]indexedRow, 0, len(rows))
	search := strings.ToLower(strings.TrimSpace(q.Search))
	typeFilter := strings.ToUpper(strings.TrimSpace(q.Type))
	confidenceFilter := strings.ToLower(strings.TrimSpace(q.Confidence))
	sourceFilter := strings.ToLower(strings.TrimSpace(q.Source))
	payeeFilter := strings.TrimSpace(q.Payee)

	for i, row := range rows {
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
		if payeeFilter != "" {
			token := strings.TrimSpace(row.PayeeToken)
			if !strings.EqualFold(token, payeeFilter) {
				continue
			}
		}
		filtered = append(filtered, indexedRow{index: i, entry: row})
	}

	rowCount := len(filtered)
	start := q.Page * q.PageSize
	if start >= rowCount {
		return Result{
			Rows:     []Row{},
			RowCount: rowCount,
			Page:     q.Page,
			PageSize: q.PageSize,
		}
	}
	end := min(start+q.PageSize, rowCount)

	out := make([]Row, 0, end-start)
	for _, ir := range filtered[start:end] {
		out = append(out, mapRow(ir.index, ir.entry))
	}
	return Result{
		Rows:     out,
		RowCount: rowCount,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
}

type indexedRow struct {
	index int
	entry report.Entry
}

func (q Query) normalized() Query {
	if q.Page < 0 {
		q.Page = 0
	}
	if q.PageSize <= 0 {
		q.PageSize = defaultPageSize
	}
	if q.PageSize > maxPageSize {
		q.PageSize = maxPageSize
	}
	q.Type = strings.ToUpper(strings.TrimSpace(q.Type))
	if q.Type == "ALL" {
		q.Type = ""
	}
	q.Confidence = strings.ToLower(strings.TrimSpace(q.Confidence))
	if q.Confidence == "all" {
		q.Confidence = ""
	}
	q.Source = strings.ToLower(strings.TrimSpace(q.Source))
	if q.Source == "all" {
		q.Source = ""
	}
	return q
}

func mapRow(index int, entry report.Entry) Row {
	return Row{
		ID:              strconv.Itoa(index),
		TransactionDate: lookup(entry.Fields, "Transaction Date"),
		Description:     report.DescriptionOf(entry.Fields),
		Amount:          lookup(entry.Fields, "Amount"),
		Type:            lookup(entry.Fields, "Dr / Cr"),
		Confidence:      entry.Confidence,
		Source:          entry.DetectionSource,
	}
}

func lookup(fields txn.Fields, name string) string {
	if v := fields.Lookup(name); v != "" {
		return v
	}
	for _, f := range fields {
		if strings.EqualFold(f.Name, name) {
			return f.Value
		}
	}
	return ""
}
