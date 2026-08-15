package kotak

import (
	"strings"
	"testing"

	"github.com/tanuvnair/unfold/internal/config"
)

func TestParse_PreservesDuplicateDrCrColumns(t *testing.T) {
	const csv = `"Account Statement"

"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","03-01-2026 08:16:49","03-01-2026","UPI/Netflix/MandateExecute","UPI-1","199.00","DR","65,330.00","CR"
"2","05-01-2026 05:35:37","05-01-2026","IB:MONTHLY INVESTMENT","0001","5,000.00","DR","91,822.85","CR"
`

	cfg := config.Config{
		BankName:          "Kotak Mahindra Bank",
		DescriptionColumn: "Description",
		Keywords:          []string{"MANDATE"},
	}

	txns, err := (&Parser{}).Parse(strings.NewReader(csv), cfg)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(txns) != 2 {
		t.Fatalf("got %d transactions, want 2", len(txns))
	}

	first := txns[0]
	if got := first.Fields.Lookup("Dr / Cr"); got != "DR" {
		t.Errorf("Amount Dr / Cr = %q, want DR", got)
	}
	if got := first.Fields.Lookup("Balance Dr / Cr"); got != "CR" {
		t.Errorf("Balance Dr / Cr = %q, want CR", got)
	}
	if got := first.Fields.Lookup("Amount"); got != "199.00" {
		t.Errorf("Amount = %q, want 199.00", got)
	}
	if first.Description != "UPI/Netflix/MandateExecute" {
		t.Errorf("Description = %q", first.Description)
	}
}

func TestParse_SkipsFooterRows(t *testing.T) {
	const csv = `"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","01-01-2026","01-01-2026","UPI/Test/Mandate","UPI-1","10.00","DR","100.00","CR"
"Opening Balance","","","","","","","",""
"Notes: something","","","","","","","",""
`

	cfg := config.Config{
		BankName:          "Kotak Mahindra Bank",
		DescriptionColumn: "Description",
		Keywords:          []string{"MANDATE"},
	}

	txns, err := (&Parser{}).Parse(strings.NewReader(csv), cfg)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1", len(txns))
	}
}
