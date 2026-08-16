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

func TestNormalizePayee_StripsRefundNoise(t *testing.T) {
	got := recurrence.NormalizePayee("UPI/APPLE MEDIA SER REFUND/1024/Mandate")
	if got != "APPLE MEDIA SER" {
		t.Fatalf("got %q, want APPLE MEDIA SER", got)
	}
}

func TestGroupByPayee(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI/Netflix/1111/Payment", Type: "DR"},
		{Description: "UPI/Netflix/2222/Payment", Type: "DR"},
		{Description: "UPI/Spotify/3333/Payment", Type: "DR"},
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

func TestGroupByPayee_SkipsCredits(t *testing.T) {
	txns := []txn.Transaction{
		{Description: "UPI/Netflix/1111/Payment", Type: "DR"},
		{Description: "UPI/Netflix/2222/Refund", Type: "CR"},
		{Description: "UPI/Netflix/3333/Payment", Type: "DR"},
	}
	groups := recurrence.GroupByPayee(txns)
	if len(groups) != 1 || len(groups[0].Transactions) != 2 {
		t.Fatalf("got %+v, want 1 group with 2 DR rows", groups)
	}
}

func monthly(payee string, amount float64, months int) []txn.Transaction {
	var out []txn.Transaction
	for i := 0; i < months; i++ {
		out = append(out, txn.Transaction{
			Description: "UPI/" + payee + "/88230" + string(rune('0'+i)) + "/Payment",
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

func TestDetectRecurring_MediumConfidenceLooseIntervals(t *testing.T) {
	txns := []txn.Transaction{
		{
			Description: "UPI/Spotify/1111/Payment",
			Amount:      119,
			Date:        time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			Type:        "DR",
		},
		{
			Description: "UPI/Spotify/2222/Payment",
			Amount:      119,
			Date:        time.Date(2025, 1, 22, 10, 0, 0, 0, time.UTC),
			Type:        "DR",
		},
		{
			Description: "UPI/Spotify/3333/Payment",
			Amount:      119,
			Date:        time.Date(2025, 2, 28, 10, 0, 0, 0, time.UTC),
			Type:        "DR",
		},
	}
	groups := recurrence.GroupByPayee(txns)
	got := recurrence.DetectRecurring(groups)
	if len(got) != 1 {
		t.Fatalf("got %d detections, want 1", len(got))
	}
	if got[0].Confidence != "medium" {
		t.Fatalf("confidence=%q want medium", got[0].Confidence)
	}
}

func TestDetectRecurring_RejectsTwoHits(t *testing.T) {
	groups := recurrence.GroupByPayee(monthly("Netflix", 199, 2))
	got := recurrence.DetectRecurring(groups)
	if len(got) != 0 {
		t.Fatalf("got %d detections, want 0 for only two spaced hits", len(got))
	}
}

func TestDetectRecurring_RejectsSameDayRetries(t *testing.T) {
	// One purchase in April, then a same-day double tap in June — not monthly.
	txns := []txn.Transaction{
		{
			Description: "UPI/HEADPHONE ZONE/1111/PaymenttoHEADPH",
			Amount:      2199,
			Date:        time.Date(2026, 4, 5, 22, 23, 51, 0, time.UTC),
			Type:        "DR",
		},
		{
			Description: "UPI/HEADPHONE ZONE/2222/PaymenttoHEADPH",
			Amount:      2189,
			Date:        time.Date(2026, 6, 10, 23, 4, 11, 0, time.UTC),
			Type:        "DR",
		},
		{
			Description: "UPI/HEADPHONE ZONE/3333/PaymentToHEADPH",
			Amount:      2189,
			Date:        time.Date(2026, 6, 10, 23, 5, 9, 0, time.UTC),
			Type:        "DR",
		},
	}
	groups := recurrence.GroupByPayee(txns)
	got := recurrence.DetectRecurring(groups)
	if len(got) != 0 {
		t.Fatalf("got %d detections, want 0 for same-day UPI retries", len(got))
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

func TestDetectRecurring_RejectsWideGaps(t *testing.T) {
	txns := []txn.Transaction{
		{
			Description: "UPI/Shop/1111/Payment",
			Amount:      200,
			Date:        time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			Type:        "DR",
		},
		{
			Description: "UPI/Shop/2222/Payment",
			Amount:      200,
			Date:        time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
			Type:        "DR",
		},
		{
			Description: "UPI/Shop/3333/Payment",
			Amount:      200,
			Date:        time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC),
			Type:        "DR",
		},
	}
	groups := recurrence.GroupByPayee(txns)
	got := recurrence.DetectRecurring(groups)
	if len(got) != 0 {
		t.Fatalf("got %d detections, want 0 for non-monthly gaps", len(got))
	}
}

func TestGroupByPayee_StillGroupsIBWhenPassed(t *testing.T) {
	// analyze strips IB: before GroupByPayee; this asserts GroupByPayee itself
	// still accepts DR rows (exclusion is analyze/matcher responsibility).
	txns := []txn.Transaction{
		{Description: "IB:MONTHLY INVESTMENT", Type: "DR", Amount: 5000},
		{Description: "IB:MONTHLY INVESTMENT", Type: "DR", Amount: 5000},
	}
	groups := recurrence.GroupByPayee(txns)
	if len(groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(groups))
	}
}
