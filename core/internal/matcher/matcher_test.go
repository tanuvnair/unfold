package matcher_test

import (
	"testing"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/matcher"
	"github.com/tanuvnair/unfold/internal/txn"
)

func TestFilter_IncludeHit(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI payment to friend", Type: "DR"},
		{Description: "NACH debit Mutual Fund", Type: "DR"},
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

func TestFilter_SkipsCredits(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "NACH MANDATE REFUND Apple", Type: "CR"},
		{Description: "NACH SIP mutual fund", Type: "DR"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, nil)
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1 (credits skipped)", len(got))
	}
	if got[0].Transaction.Description != "NACH SIP mutual fund" {
		t.Fatalf("unexpected match: %q", got[0].Transaction.Description)
	}
}

func TestFilter_ExcludeSuppressesFalsePositive(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "NACH-REF-999 one-off transfer", Type: "DR"},
		{Description: "NACH SIP mutual fund", Type: "DR"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, []string{"ONE-OFF"})
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1", len(got))
	}
	if got[0].Transaction.Description != "NACH SIP mutual fund" {
		t.Fatalf("unexpected match: %q", got[0].Transaction.Description)
	}
}

func TestFilter_ExcludeRefund(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI/APPLE MEDIA SER REFUND/1024/Mandate", Type: "DR"},
	}
	got := matcher.Filter(
		txns,
		[]config.KeywordRule{{Term: "MANDATE", Tier: "high"}},
		[]string{"REFUND"},
	)
	if len(got) != 0 {
		t.Fatalf("got %d matches, want 0 for refund", len(got))
	}
}

func TestFilter_NoIncludeNoMatch(t *testing.T) {
	txns := []txn.Transaction{{Description: "grocery store", Type: "DR"}}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, nil)
	if len(got) != 0 {
		t.Fatalf("got %d matches, want 0", len(got))
	}
}

func TestFilter_HighestTierWins(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "NACH RECURRING SIP debit", Type: "DR"},
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
		{Description: "NACH one-off transfer", Type: "DR"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{{Term: "NACH", Tier: "high"}}, []string{"ONE-OFF"})
	if len(got) != 0 {
		t.Fatalf("got %d matches, want 0 (exclude should win)", len(got))
	}
}

func TestFilter_SkipsIBSelfTransfer(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "IB:MONTHLY INVESTMENT", Type: "DR"},
		{Description: "NACH SIP mutual fund", Type: "DR"},
	}
	got := matcher.Filter(txns, []config.KeywordRule{
		{Term: "NACH", Tier: "high"},
		{Term: "MONTHLY INVESTMENT", Tier: "medium"},
	}, nil)
	if len(got) != 1 {
		t.Fatalf("got %d matches, want 1", len(got))
	}
	if got[0].Transaction.Description != "NACH SIP mutual fund" {
		t.Fatalf("unexpected match: %q", got[0].Transaction.Description)
	}
}

func TestFilter_ExcludeMonthlyInvestmentPhrase(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI/FOO/MONTHLY INVESTMENT TO ICICI ACCOUNT", Type: "DR"},
	}
	got := matcher.Filter(
		txns,
		[]config.KeywordRule{{Term: "MANDATE", Tier: "high"}},
		[]string{"MONTHLY INVESTMENT", "TO ICICI ACCOUNT"},
	)
	if len(got) != 0 {
		t.Fatalf("got %d matches, want 0", len(got))
	}
}

func TestLooksLikeSelfTransfer(t *testing.T) {
	if !matcher.LooksLikeSelfTransfer("IB:MONTHLY INVESTMENT") {
		t.Fatal("expected IB: prefix to be self-transfer")
	}
	if matcher.LooksLikeSelfTransfer("UPI/Netflix/111/Payment") {
		t.Fatal("Netflix UPI should not be self-transfer")
	}
}
