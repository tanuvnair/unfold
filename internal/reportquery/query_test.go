package reportquery

import (
	"testing"

	"github.com/tanuvnair/unfold/internal/txn"
)

func sampleRows() []txn.Fields {
	return []txn.Fields{
		{
			{Name: "Transaction Date", Value: "03-01-2026 08:16:49"},
			{Name: "Description", Value: "UPI/Netflix/530812380036/MandateExecute"},
			{Name: "Amount", Value: "199.00"},
			{Name: "Dr / Cr", Value: "DR"},
		},
		{
			{Name: "Transaction Date", Value: "05-01-2026 05:24:48"},
			{Name: "Description", Value: "IB:MONTHLY INVESTMENT"},
			{Name: "Amount", Value: "5,000.00"},
			{Name: "Dr / Cr", Value: "DR"},
		},
		{
			{Name: "Transaction Date", Value: "10-01-2026 12:00:00"},
			{Name: "Description", Value: "UPI/APPLE MEDIA SER/102431795941/UPI Mandate"},
			{Name: "Amount", Value: "99.00"},
			{Name: "Dr / Cr", Value: "CR"},
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
