package recurrence_test

import (
	"testing"
	"time"

	"github.com/tanuvnair/unfold/internal/recurrence"
	"github.com/tanuvnair/unfold/internal/txn"
)

func TestNormalizePayee_StripsUPINoise(t *testing.T) {
	got := recurrence.NormalizePayee("UPI/Netflix/530812380036/Payment")
	if got != "NETFLIX" {
		t.Fatalf("got %q, want NETFLIX", got)
	}
}

func TestNormalizePayee_StripsDatesAndRefs(t *testing.T) {
	got := recurrence.NormalizePayee("NEFT/SPOTIFY INDIA/12-01-2026/REF99887766")
	if got != "SPOTIFY INDIA" {
		t.Fatalf("got %q, want SPOTIFY INDIA", got)
	}
}

func TestGroupByPayee(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI/Netflix/1111/Payment"},
		{Description: "UPI/Netflix/2222/Payment"},
		{Description: "UPI/Spotify/3333/Payment"},
	}
	groups := recurrence.GroupByPayee(txns)
	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(groups))
	}
	byToken := map[string]int{}
	for _, g := range groups {
		byToken[g.PayeeToken] = len(g.Transactions)
	}
	if byToken["NETFLIX"] != 2 || byToken["SPOTIFY"] != 1 {
		t.Fatalf("groups=%v", byToken)
	}
}

func monthly(payee string, amount float64, months int) []txn.Transaction {
	var out []txn.Transaction
	for i := 0; i < months; i++ {
		out = append(out, txn.Transaction{
			Description: "UPI/" + payee + "/8823" + string(rune('0'+i)) + "/Payment",
			Amount:      amount,
			Date:        time.Date(2025, time.Month(1+i), 5, 10, 0, 0, 0, time.UTC),
			Type:        "DR",
			Fields: txn.Fields{
				{Name: "Description", Value: "UPI/" + payee + "/Payment"},
				{Name: "Amount", Value: "199.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		})
	}
	return out
}

func TestDetectRecurring_HighConfidenceThreeHits(t *testing.T) {
	groups := recurrence.GroupByPayee(monthly("Netflix", 199, 3))
	got := recurrence.DetectRecurring(groups)
	if len(got) != 1 {
		t.Fatalf("got %d detections, want 1", len(got))
	}
	if got[0].Confidence != "high" {
		t.Fatalf("confidence=%q want high", got[0].Confidence)
	}
}

func TestDetectRecurring_MediumConfidenceTwoHits(t *testing.T) {
	groups := recurrence.GroupByPayee(monthly("Netflix", 199, 2))
	got := recurrence.DetectRecurring(groups)
	if len(got) != 1 {
		t.Fatalf("got %d detections, want 1", len(got))
	}
	if got[0].Confidence != "medium" {
		t.Fatalf("confidence=%q want medium", got[0].Confidence)
	}
}

func TestDetectRecurring_RejectsOneOff(t *testing.T) {
	txns := monthly("Friend", 500, 1)
	groups := recurrence.GroupByPayee(txns)
	got := recurrence.DetectRecurring(groups)
	if len(got) != 0 {
		t.Fatalf("got %d detections, want 0", len(got))
	}
}

func TestDetectRecurring_RejectsUnstableAmount(t *testing.T) {
	txns := monthly("Shop", 100, 3)
	txns[1].Amount = 500
	groups := recurrence.GroupByPayee(txns)
	got := recurrence.DetectRecurring(groups)
	if len(got) != 0 {
		t.Fatalf("got %d detections, want 0 for unstable amount", len(got))
	}
}
