package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tanuvnair/unfold/internal/config"
)

func validFile(profiles ...config.Profile) config.File {
	return config.File{
		ConfigVersion: config.SupportedVersion,
		Profiles:      profiles,
	}
}

func kotakProfile() config.Profile {
	return config.Profile{
		BankName:          "Kotak Mahindra Bank",
		SkipRows:          0,
		DescriptionColumn: "Description",
		AmountColumn:      "Amount",
		DateColumn:        "Transaction Date",
		DateFormat:        "02-01-2006 15:04:05",
		TypeColumn:        "Dr / Cr",
		Keywords: []config.KeywordRule{
			{Term: "NACH", Tier: "high"},
			{Term: "AUTOPAY", Tier: "high"},
		},
		ExcludeKeywords: []string{"ONE-OFF"},
	}
}

func TestFileValidate_RejectsWrongVersion(t *testing.T) {
	f := validFile(kotakProfile())
	f.ConfigVersion = 99
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for unsupported config_version")
	}
}

func TestFileValidate_RejectsEmptyProfiles(t *testing.T) {
	f := validFile()
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty profiles")
	}
}

func TestFileValidate_RejectsDuplicateBanks(t *testing.T) {
	f := validFile(kotakProfile(), kotakProfile())
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for duplicate bank profiles")
	}
}

func TestFileValidate_RejectsMissingDescriptionColumn(t *testing.T) {
	p := kotakProfile()
	p.DescriptionColumn = ""
	f := validFile(p)
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty description_column")
	}
}

func TestFileValidate_RejectsUnknownTier(t *testing.T) {
	p := kotakProfile()
	p.Keywords = []config.KeywordRule{{Term: "NACH", Tier: "urgent"}}
	f := validFile(p)
	err := f.Validate()
	if err == nil {
		t.Fatal("expected error for unknown tier")
	}
	if !strings.Contains(err.Error(), "tier") {
		t.Fatalf("error should mention tier, got: %v", err)
	}
}

func TestFileValidate_RejectsEmptyTerm(t *testing.T) {
	p := kotakProfile()
	p.Keywords = []config.KeywordRule{{Term: "  ", Tier: "high"}}
	f := validFile(p)
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty term")
	}
}

func TestSelect_SingleProfileWithoutBank(t *testing.T) {
	f := validFile(kotakProfile())
	cfg, err := f.Select("")
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if cfg.BankKey() != "kotak-mahindra-bank" {
		t.Fatalf("BankKey = %q", cfg.BankKey())
	}
}

func TestSelect_RequiresBankWhenMultipleProfiles(t *testing.T) {
	second := kotakProfile()
	second.BankName = "HDFC Bank"
	f := validFile(kotakProfile(), second)
	if _, err := f.Select(""); err == nil {
		t.Fatal("expected error when --bank omitted with multiple profiles")
	}
}

func TestSelect_ByKeyAndName(t *testing.T) {
	f := validFile(kotakProfile())
	for _, bank := range []string{"kotak-mahindra-bank", "Kotak Mahindra Bank", "KOTAK MAHINDRA BANK"} {
		cfg, err := f.Select(bank)
		if err != nil {
			t.Fatalf("Select(%q): %v", bank, err)
		}
		if cfg.BankName != "Kotak Mahindra Bank" {
			t.Fatalf("Select(%q) bank_name = %q", bank, cfg.BankName)
		}
	}
}

func TestSelect_UnknownBank(t *testing.T) {
	f := validFile(kotakProfile())
	if _, err := f.Select("sbi"); err == nil {
		t.Fatal("expected error for unknown bank")
	}
}

func TestNormalizedKeywords(t *testing.T) {
	p := kotakProfile()
	p.Keywords = []config.KeywordRule{
		{Term: "  nach ", Tier: "high"},
		{Term: "", Tier: "high"},
		{Term: "AutoPay", Tier: "medium"},
	}
	got := p.NormalizedKeywords()
	want := []config.KeywordRule{
		{Term: "NACH", Tier: "high"},
		{Term: "AUTOPAY", Tier: "medium"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "config_version": 1,
  "profiles": [
    {
      "bank_name": "Kotak Mahindra Bank",
      "skip_rows": 0,
      "description_column": "Description",
      "amount_column": "Amount",
      "date_column": "Transaction Date",
      "date_format": "02-01-2006 15:04:05",
      "type_column": "Dr / Cr",
      "keywords": [{"term": "NACH", "tier": "high"}],
      "exclude_keywords": ["FALSE"]
    }
  ]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	file, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	cfg, err := file.Select("")
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if cfg.DescriptionColumn != "Description" {
		t.Fatalf("DescriptionColumn = %q", cfg.DescriptionColumn)
	}
	if got := cfg.NormalizedExcludeKeywords(); len(got) != 1 || got[0] != "FALSE" {
		t.Fatalf("NormalizedExcludeKeywords = %v", got)
	}
	rules := cfg.NormalizedKeywords()
	if len(rules) != 1 || rules[0].Term != "NACH" || rules[0].Tier != "high" {
		t.Fatalf("NormalizedKeywords = %v", rules)
	}
}

func TestLoad_RejectsLegacyShape(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "legacy.json")
	// Old flat schema has no config_version / profiles.
	content := `{"bank_name":"Kotak","description_column_index":3,"keywords":["NACH"]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for legacy config")
	}
	if !strings.Contains(err.Error(), "config_version") {
		t.Fatalf("error should mention config_version, got: %v", err)
	}
}
