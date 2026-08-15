package matcher_test

import (
	"testing"

	"github.com/tanuvnair/unfold/internal/matcher"
	"github.com/tanuvnair/unfold/internal/txn"
)

func TestFilter_IncludeHit(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI payment to friend"},
		{Description: "NACH debit Mutual Fund"},
	}
	got := matcher.Filter(txns, []string{"NACH"}, nil)
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1", len(got))
	}
	if got[0].Description != "NACH debit Mutual Fund" {
		t.Fatalf("unexpected match: %q", got[0].Description)
	}
}

func TestFilter_ExcludeSuppressesFalsePositive(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "NACH-REF-999 one-off transfer"},
		{Description: "NACH SIP mutual fund"},
	}
	got := matcher.Filter(txns, []string{"NACH"}, []string{"ONE-OFF"})
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1", len(got))
	}
	if got[0].Description != "NACH SIP mutual fund" {
		t.Fatalf("unexpected match: %q", got[0].Description)
	}
}

func TestFilter_NoIncludeNoMatch(t *testing.T) {
	txns := []txn.Transaction{{Description: "grocery store"}}
	got := matcher.Filter(txns, []string{"NACH"}, nil)
	if len(got) != 0 {
		t.Fatalf("got %d matches, want 0", len(got))
	}
}
