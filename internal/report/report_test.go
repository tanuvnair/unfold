package report_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tanuvnair/unfold/internal/report"
	"github.com/tanuvnair/unfold/internal/txn"
)

func row(desc, amount string) txn.Fields {
	return txn.Fields{
		{Name: "Description", Value: desc},
		{Name: "Amount", Value: amount},
	}
}

func TestDiff_AddedRemovedUnchanged(t *testing.T) {
	prev := report.Report{Transactions: []txn.Fields{
		row("NACH A", "10"),
		row("NACH B", "20"),
	}}
	next := report.Report{Transactions: []txn.Fields{
		row("NACH B", "20"),
		row("NACH C", "30"),
	}}

	d := report.Diff(prev, next)
	if len(d.Added) != 1 || report.DescriptionOf(d.Added[0]) != "NACH C" {
		t.Fatalf("Added=%v", d.Added)
	}
	if len(d.Removed) != 1 || report.DescriptionOf(d.Removed[0]) != "NACH A" {
		t.Fatalf("Removed=%v", d.Removed)
	}
	if d.Unchanged != 1 {
		t.Fatalf("Unchanged=%d want 1", d.Unchanged)
	}
}

func TestDiff_OrderIndependentFingerprint(t *testing.T) {
	prev := report.Report{Transactions: []txn.Fields{{
		{Name: "Amount", Value: "10"},
		{Name: "Description", Value: "NACH"},
	}}}
	next := report.Report{Transactions: []txn.Fields{{
		{Name: "Description", Value: "NACH"},
		{Name: "Amount", Value: "10"},
	}}}
	d := report.Diff(prev, next)
	if len(d.Added) != 0 || len(d.Removed) != 0 || d.Unchanged != 1 {
		t.Fatalf("diff=%+v", d)
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "autopay_report.json")
	want := report.Report{
		BankName:         "Kotak Mahindra Bank",
		TransactionCount: 1,
		Transactions: []txn.Fields{{
			{Name: "Sl. No.", Value: "1"},
			{Name: "Description", Value: "AUTOPAY"},
		}},
	}
	if err := report.Write(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := report.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.BankName != want.BankName || got.TransactionCount != 1 {
		t.Fatalf("got %+v", got)
	}
	if got.Transactions[0][0].Name != "Sl. No." || got.Transactions[0][1].Name != "Description" {
		t.Fatalf("column order lost: %+v", got.Transactions[0])
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !containsInOrder(string(raw), `"Sl. No."`, `"Description"`) {
		t.Fatalf("raw JSON lost order: %s", raw)
	}
}

func containsInOrder(s, a, b string) bool {
	i := indexOf(s, a)
	j := indexOf(s, b)
	return i >= 0 && j >= 0 && i < j
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
