package reportquery

import (
	"strings"
	"testing"
	"time"

	"github.com/tanuvnair/unfold/internal/report"
	"github.com/tanuvnair/unfold/internal/txn"
)

func sampleRows() []report.Entry {
	return []report.Entry{
		{
			Confidence:      "high",
			DetectionSource: report.SourceKeyword,
			Date:            time.Date(2026, 1, 3, 8, 16, 49, 0, time.UTC),
			Fields: txn.Fields{
				{Name: "Transaction Date", Value: "03-01-2026 08:16:49"},
				{Name: "Description", Value: "UPI/Netflix/530812380036/MandateExecute"},
				{Name: "Amount", Value: "199.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
		{
			Confidence:      "medium",
			DetectionSource: report.SourceRecurrence,
			Date:            time.Date(2026, 1, 5, 5, 24, 48, 0, time.UTC),
			Fields: txn.Fields{
				{Name: "Transaction Date", Value: "05-01-2026 05:24:48"},
				{Name: "Description", Value: "IB:MONTHLY INVESTMENT"},
				{Name: "Amount", Value: "5,000.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
		{
			Confidence:      "high",
			DetectionSource: report.SourceBoth,
			Date:            time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC),
			Fields: txn.Fields{
				{Name: "Transaction Date", Value: "10-01-2026 12:00:00"},
				{Name: "Description", Value: "UPI/APPLE MEDIA SER/102431795941/UPI Mandate"},
				{Name: "Amount", Value: "99.00"},
				{Name: "Dr / Cr", Value: "CR"},
			},
		},
	}
}

func TestApply_SearchDescription(t *testing.T) {
	got := Apply(sampleRows(), Query{Search: "netflix", PageSize: 10})
	if got.RowCount != 1 {
		t.Fatalf("RowCount=%d want 1", got.RowCount)
	}
	if got.Rows[0].Description != "UPI/Netflix/530812380036/MandateExecute" {
		t.Fatalf("row=%+v", got.Rows[0])
	}
}

func TestApply_TypeFilter(t *testing.T) {
	got := Apply(sampleRows(), Query{Type: "CR", PageSize: 10})
	if got.RowCount != 1 || got.Rows[0].Type != "CR" {
		t.Fatalf("got=%+v", got)
	}
}

func TestApply_ConfidenceFilter(t *testing.T) {
	got := Apply(sampleRows(), Query{Confidence: "medium", PageSize: 10})
	if got.RowCount != 1 || got.Rows[0].Confidence != "medium" {
		t.Fatalf("got=%+v", got)
	}
}

func TestApply_Pagination(t *testing.T) {
	got := Apply(sampleRows(), Query{Page: 1, PageSize: 2})
	if got.RowCount != 3 {
		t.Fatalf("RowCount=%d want 3", got.RowCount)
	}
	if len(got.Rows) != 1 || got.Rows[0].ID != "2" {
		t.Fatalf("page=%+v", got.Rows)
	}
}

func TestApply_EmptyPagePastEnd(t *testing.T) {
	got := Apply(sampleRows(), Query{Page: 9, PageSize: 10})
	if got.RowCount != 3 || len(got.Rows) != 0 {
		t.Fatalf("got=%+v", got)
	}
}

func TestApply_AllTypeIsUnfiltered(t *testing.T) {
	got := Apply(sampleRows(), Query{Type: "all", PageSize: 10})
	if got.RowCount != 3 {
		t.Fatalf("RowCount=%d want 3", got.RowCount)
	}
}

func TestApply_StableIDsAreOriginalIndexes(t *testing.T) {
	got := Apply(sampleRows(), Query{Search: "apple", PageSize: 10})
	if len(got.Rows) != 1 || got.Rows[0].ID != "2" {
		t.Fatalf("got=%+v", got.Rows)
	}
}

func TestApply_SourceFilter(t *testing.T) {
	got := Apply(sampleRows(), Query{Source: "recurrence", PageSize: 10})
	if got.RowCount != 1 || got.Rows[0].Source != report.SourceRecurrence {
		t.Fatalf("got=%+v", got)
	}
}

func TestApply_DateFrom(t *testing.T) {
	from := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	got := Apply(sampleRows(), Query{DateFrom: from, PageSize: 10})
	if got.RowCount != 2 {
		t.Fatalf("RowCount=%d want 2", got.RowCount)
	}
}

func TestApply_DateTo(t *testing.T) {
	to := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	got := Apply(sampleRows(), Query{DateTo: to, PageSize: 10})
	if got.RowCount != 2 {
		t.Fatalf("RowCount=%d want 2", got.RowCount)
	}
}

func TestApply_DateRange(t *testing.T) {
	from := time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC)
	got := Apply(sampleRows(), Query{DateFrom: from, DateTo: to, PageSize: 10})
	if got.RowCount != 1 || !strings.Contains(got.Rows[0].Description, "MONTHLY") {
		t.Fatalf("got=%+v", got)
	}
}
