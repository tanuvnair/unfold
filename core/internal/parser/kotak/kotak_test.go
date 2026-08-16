package kotak

import (
	"strings"
	"testing"
	"time"

	"github.com/tanuvnair/unfold/internal/config"
)

func testConfig() config.Config {
	return config.Config{
		BankName:          "Kotak Mahindra Bank",
		DescriptionColumn: "Description",
		AmountColumn:      "Amount",
		DateColumn:        "Transaction Date",
		DateFormat:        "02-01-2006 15:04:05",
		TypeColumn:        "Dr / Cr",
		Keywords:          []config.KeywordRule{{Term: "MANDATE", Tier: "high"}},
	}
}

func TestParse_PreservesDuplicateDrCrColumns(t *testing.T) {
	const csv = `"Account Statement"

"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","03-01-2026 08:16:49","03-01-2026","UPI/Netflix/MandateExecute","UPI-1","199.00","DR","65,330.00","CR"
"2","05-01-2026 05:35:37","05-01-2026 00:00:00","IB:MONTHLY INVESTMENT","0001","5,000.00","DR","91,822.85","CR"
`

	txns, err := (&Parser{}).Parse(strings.NewReader(csv), testConfig())
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
	if first.Amount != 199.00 {
		t.Errorf("Amount = %v, want 199", first.Amount)
	}
	if first.Type != "DR" {
		t.Errorf("Type = %q, want DR", first.Type)
	}
	wantDate := time.Date(2026, 1, 3, 8, 16, 49, 0, time.UTC)
	if !first.Date.Equal(wantDate) {
		t.Errorf("Date = %v, want %v", first.Date, wantDate)
	}
	if txns[1].Amount != 5000.00 {
		t.Errorf("second Amount = %v, want 5000", txns[1].Amount)
	}
}

func TestParse_SkipsFooterRows(t *testing.T) {
	const csv = `"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","01-01-2026 00:00:00","01-01-2026","UPI/Test/Mandate","UPI-1","10.00","DR","100.00","CR"
"Opening Balance","","","","","","","",""
"Notes: something","","","","","","","",""
`

	txns, err := (&Parser{}).Parse(strings.NewReader(csv), testConfig())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1", len(txns))
	}
}

func TestParse_CommaAndCurrencyAmount(t *testing.T) {
	const csv = `"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","01-01-2026 12:00:00","01-01-2026","UPI/Test/Mandate","UPI-1","₹1,234.50","DR","100.00","CR"
`

	txns, err := (&Parser{}).Parse(strings.NewReader(csv), testConfig())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if txns[0].Amount != 1234.50 {
		t.Fatalf("Amount = %v, want 1234.50", txns[0].Amount)
	}
}

func TestParse_MalformedDate(t *testing.T) {
	const csv = `"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","not-a-date","01-01-2026","UPI/Test/Mandate","UPI-1","10.00","DR","100.00","CR"
`

	_, err := (&Parser{}).Parse(strings.NewReader(csv), testConfig())
	if err == nil {
		t.Fatal("expected parse error for malformed date")
	}
	if !strings.Contains(err.Error(), "parse date") {
		t.Fatalf("error should mention parse date, got: %v", err)
	}
}

func TestParse_MalformedAmount(t *testing.T) {
	const csv = `"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Amount","Dr / Cr","Balance","Dr / Cr"
"1","01-01-2026 12:00:00","01-01-2026","UPI/Test/Mandate","UPI-1","not-money","DR","100.00","CR"
`

	_, err := (&Parser{}).Parse(strings.NewReader(csv), testConfig())
	if err == nil {
		t.Fatal("expected parse error for malformed amount")
	}
	if !strings.Contains(err.Error(), "parse amount") {
		t.Fatalf("error should mention parse amount, got: %v", err)
	}
}

func TestParse_MissingAmountColumn(t *testing.T) {
	const csv = `"Sl. No.","Transaction Date","Value Date","Description","Chq /Ref No.","Dr / Cr","Balance","Dr / Cr"
"1","01-01-2026 12:00:00","01-01-2026","UPI/Test/Mandate","UPI-1","DR","100.00","CR"
`

	_, err := (&Parser{}).Parse(strings.NewReader(csv), testConfig())
	if err == nil {
		t.Fatal("expected error for missing amount column")
	}
	if !strings.Contains(err.Error(), "amount column") {
		t.Fatalf("error should mention amount column, got: %v", err)
	}
}

func TestParseAmount(t *testing.T) {
	cases := []struct {
		raw  string
		want float64
	}{
		{"199.00", 199},
		{"5,000.00", 5000},
		{"₹99.00", 99},
		{"Rs. 1,250.25", 1250.25},
	}
	for _, tc := range cases {
		got, err := parseAmount(tc.raw)
		if err != nil {
			t.Fatalf("parseAmount(%q): %v", tc.raw, err)
		}
		if got != tc.want {
			t.Fatalf("parseAmount(%q)=%v want %v", tc.raw, got, tc.want)
		}
	}
}
