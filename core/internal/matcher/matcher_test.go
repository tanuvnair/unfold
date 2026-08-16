package matcher_test

import (
	"testing"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/matcher"
	"github.com/tanuvnair/unfold/internal/txn"
)

func TestFilter_IncludeHit(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI payment to friend"},
		{Description: "NACH debit Mutual Fund"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, nil)
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1", len(got))
	}
	if got[0].Transaction.Description != "NACH debit Mutual Fund" {
		t.Fatalf("unexpected match: %q", got[0].Transaction.Description)
	}
	if got[0].Confidence != "high" || got[0].MatchedTerm != "NACH" {
		t.Fatalf("got confidence=%q term=%q", got[0].Confidence, got[0].MatchedTerm)
	}
}

func TestFilter_ExcludeSuppressesFalsePositive(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "NACH-REF-999 one-off transfer"},
		{Description: "NACH SIP mutual fund"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, []string{"ONE-OFF"})
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1", len(got))
	}
	if got[0].Transaction.Description != "NACH SIP mutual fund" {
		t.Fatalf("unexpected match: %q", got[0].Transaction.Description)
	}
}

func TestFilter_NoIncludeNoMatch(t *testing.T) {
	txns := []txn.Transaction{{Description: "grocery store"}}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, nil)
	if len(got) != 0 {
		t.Fatalf("got %d matches, want 0", len(got))
	}
}

func TestFilter_HighestTierWins(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "NACH RECURRING SIP debit"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{
		{Term: "RECURRING", Tier: "medium"},
		{Term: "NACH", Tier: "high"},
	}, nil)
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1", len(got))
	}
	if got[0].Confidence != "high" || got[0].MatchedTerm != "NACH" {
		t.Fatalf("got confidence=%q term=%q, want high/NACH", got[0].Confidence, got[0].MatchedTerm)
	}
}

func TestFilter_ExcludeWinsOverHighTierInclude(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "NACH one-off transfer"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, []string{"ONE-OFF"})
	if len(got) != 0 {
		t.Fatalf("got %d matches, want 0 (exclude should win)", len(got))
	}
}
