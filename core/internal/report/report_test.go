package report_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tanuvnair/unfold/internal/report"
	"github.com/tanuvnair/unfold/internal/txn"
)

func entry(desc, amount, confidence string) report.Entry {
	return report.Entry{
		Confidence: confidence,
		Fields: txn.Fields{
			{Name: "Description", Value: desc},
			{Name: "Amount", Value: amount},
		},
	}
}

func TestDiff_AddedRemovedUnchanged(t *testing.T) {
	prev := report.Report{Transactions: []report.Entry{
		entry("NACH A", "10", "high"),
		entry("NACH B", "20", "high"),
	}}
	next := report.Report{Transactions: []report.Entry{
		entry("NACH B", "20", "medium"),
		entry("NACH C", "30", "high"),
	}}

	d := report.Diff(prev, next)
	if len(d.Added) != 1 || report.EntryDescription(d.Added[0]) != "NACH C" {
		t.Fatalf("Added=%v", d.Added)
	}
	if len(d.Removed) != 1 || report.EntryDescription(d.Removed[0]) != "NACH A" {
		t.Fatalf("Removed=%v", d.Removed)
	}
	if d.Unchanged != 1 {
		t.Fatalf("Unchanged=%d want 1", d.Unchanged)
	}
}

func TestDiff_OrderIndependentFingerprint(t *testing.T) {
	prev := report.Report{Transactions: []report.Entry{{
		Fields: txn.Fields{
			{Name: "Amount", Value: "10"},
			{Name: "Description", Value: "NACH"},
		},
	}}}
	next := report.Report{Transactions: []report.Entry{{
		Fields: txn.Fields{
			{Name: "Description", Value: "NACH"},
			{Name: "Amount", Value: "10"},
		},
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
		Transactions: []report.Entry{{
			Confidence:  "high",
			MatchedTerm: "AUTOPAY",
			Fields: txn.Fields{
				{Name: "Sl. No.", Value: "1"},
				{Name: "Description", Value: "AUTOPAY"},
			},
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
	if got.Transactions[0].Confidence != "high" || got.Transactions[0].MatchedTerm != "AUTOPAY" {
		t.Fatalf("metadata lost: %+v", got.Transactions[0])
	}
	if got.Transactions[0].Fields[0].Name != "Sl. No." || got.Transactions[0].Fields[1].Name != "Description" {
		t.Fatalf("column order lost: %+v", got.Transactions[0].Fields)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !containsInOrder(string(raw), `"confidence"`, `"matched_term"`, `"Sl. No."`, `"Description"`) {
		t.Fatalf("raw JSON lost order: %s", raw)
	}
}

func TestBuild_FromMatches(t *testing.T) {
	got := report.Build("Kotak Mahindra Bank", []report.Entry{{
		Confidence:      "high",
		MatchedTerm:     "NACH",
		DetectionSource: report.SourceKeyword,
		Fields:          txn.Fields{{Name: "Description", Value: "NACH debit"}},
	}})
	if got.TransactionCount != 1 || got.Transactions[0].Confidence != "high" {
		t.Fatalf("got %+v", got)
	}
}

func containsInOrder(s string, parts ...string) bool {
	pos := 0
	for _, p := range parts {
		i := indexOf(s[pos:], p)
		if i < 0 {
			return false
		}
		pos += i + len(p)
	}
	return true
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
