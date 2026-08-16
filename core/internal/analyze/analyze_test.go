package analyze_test

import (
	"strings"
	"testing"

	"github.com/tanuvnair/unfold/internal/analyze"
	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/report"

	_ "github.com/tanuvnair/unfold/internal/parser/kotak"
)

func TestRun_ClassifiesKeywordRecurrenceAndOneOff(t *testing.T) {
	// Synthetic statement:
	// (a) NACH keyword hit (single debit)
	// (b) UPI Netflix recurring 3x, no keyword language
	// (c) one-off UPI to a friend
	// (d) mandate refund credit — must not match
	const csv = `"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","05-01-2025 10:00:00","05-01-2025","NACH MUTUAL FUND SIP","REF1","2000.00","DR","100.00","CR"
"2","05-01-2025 11:00:00","05-01-2025","UPI/Netflix/1111111111/Payment","UPI-1","199.00","DR","100.00","CR"
"3","05-02-2025 11:00:00","05-02-2025","UPI/Netflix/2222222222/Payment","UPI-2","199.00","DR","100.00","CR"
"4","05-03-2025 11:00:00","05-03-2025","UPI/Netflix/3333333333/Payment","UPI-3","199.00","DR","100.00","CR"
"5","10-01-2025 12:00:00","10-01-2025","UPI/Friend/4444444444/Payment","UPI-4","50.00","DR","100.00","CR"
"6","12-01-2025 12:00:00","12-01-2025","UPI/APPLE MEDIA SER REFUND/5555/Mandate","UPI-5","5.00","CR","100.00","CR"
`

	cfg := config.Config{
		BankName:          "Kotak Mahindra Bank",
		DescriptionColumn: "Description",
		AmountColumn:      "Amount",
		DateColumn:        "Transaction Date",
		DateFormat:        "02-01-2006 15:04:05",
		TypeColumn:        "Dr / Cr",
		Keywords: []config.KeywordRule{
			{Term: "NACH", Tier: "high"},
			{Term: "MANDATE", Tier: "high"},
		},
		ExcludeKeywords: []string{"REFUND"},
	}

	rpt, err := analyze.Run(cfg, strings.NewReader(csv))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	byDesc := map[string]report.Entry{}
	for _, e := range rpt.Transactions {
		byDesc[report.EntryDescription(e)] = e
	}

	nach, ok := byDesc["NACH MUTUAL FUND SIP"]
	if !ok {
		t.Fatalf("missing NACH keyword hit; got %d rows", rpt.TransactionCount)
	}
	if nach.DetectionSource != report.SourceKeyword || nach.Confidence != "high" {
		t.Fatalf("NACH entry=%+v", nach)
	}

	netflixCount := 0
	for _, e := range rpt.Transactions {
		if strings.Contains(report.EntryDescription(e), "Netflix") {
			netflixCount++
			if e.DetectionSource != report.SourceRecurrence {
				t.Fatalf("Netflix source=%q want recurrence", e.DetectionSource)
			}
		}
	}
	if netflixCount != 3 {
		t.Fatalf("Netflix recurrence hits=%d want 3", netflixCount)
	}

	for _, e := range rpt.Transactions {
		desc := report.EntryDescription(e)
		if strings.Contains(desc, "Friend") {
			t.Fatalf("one-off Friend payment should not be flagged: %+v", e)
		}
		if strings.Contains(desc, "REFUND") {
			t.Fatalf("refund should not be flagged: %+v", e)
		}
	}
}
