// Package config loads and validates the JSON config that tells unfold
// which bank profile to use and which keywords count as autopay.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// SupportedVersion is the only config_version this binary accepts.
const SupportedVersion = 1

// Valid keyword confidence tiers.
const (
	TierHigh   = "high"
	TierMedium = "medium"
	TierLow    = "low"
)

// File is the on-disk multi-profile config document.
type File struct {
	ConfigVersion int       `json:"config_version"`
	Profiles      []Profile `json:"profiles"`
}

// KeywordRule is one include keyword with a confidence tier.
type KeywordRule struct {
	Term string `json:"term"`
	Tier string `json:"tier"`
}

// Profile is one bank's parsing + matching settings inside a File.
type Profile struct {
	BankName          string        `json:"bank_name"`
	SkipRows          int           `json:"skip_rows"`
	DescriptionColumn string        `json:"description_column"`
	AmountColumn      string        `json:"amount_column"`
	DateColumn        string        `json:"date_column"`
	DateFormat        string        `json:"date_format"`
	TypeColumn        string        `json:"type_column"`
	Keywords          []KeywordRule `json:"keywords"`
	ExcludeKeywords   []string      `json:"exclude_keywords"`
}

// Config is a selected Profile — what parsers and the matcher consume.
// It is intentionally the same shape as Profile so call sites stay simple.
type Config = Profile

// Load reads and validates a config file from disk.
func Load(path string) (File, error) {
	f, err := os.Open(path)
	if err != nil {
		return File{}, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var file File
	if err := json.NewDecoder(f).Decode(&file); err != nil {
		return File{}, fmt.Errorf("parse config: %w", err)
	}

	if err := file.Validate(); err != nil {
		return File{}, fmt.Errorf("invalid config: %w", err)
	}

	return file, nil
}

// Validate catches config mistakes before they turn into a confusing empty
// report or a panic partway through parsing.
func (f File) Validate() error {
	if f.ConfigVersion != SupportedVersion {
		return fmt.Errorf("config_version must be %d, got %d", SupportedVersion, f.ConfigVersion)
	}
	if len(f.Profiles) == 0 {
		return fmt.Errorf("profiles must not be empty")
	}

	seen := make(map[string]struct{}, len(f.Profiles))
	for i, p := range f.Profiles {
		if err := p.Validate(); err != nil {
			return fmt.Errorf("profiles[%d]: %w", i, err)
		}
		key := p.BankKey()
		if _, ok := seen[key]; ok {
			return fmt.Errorf("duplicate bank profile %q", key)
		}
		seen[key] = struct{}{}
	}
	return nil
}

// Validate checks a single profile's fields.
func (c Config) Validate() error {
	if strings.TrimSpace(c.BankName) == "" {
		return fmt.Errorf("bank_name is required")
	}
	if c.SkipRows < 0 {
		return fmt.Errorf("skip_rows must be >= 0, got %d", c.SkipRows)
	}
	if strings.TrimSpace(c.DescriptionColumn) == "" {
		return fmt.Errorf("description_column is required")
	}
	if strings.TrimSpace(c.AmountColumn) == "" {
		return fmt.Errorf("amount_column is required")
	}
	if strings.TrimSpace(c.DateColumn) == "" {
		return fmt.Errorf("date_column is required")
	}
	if strings.TrimSpace(c.DateFormat) == "" {
		return fmt.Errorf("date_format is required")
	}
	if strings.TrimSpace(c.TypeColumn) == "" {
		return fmt.Errorf("type_column is required")
	}
	if len(c.Keywords) == 0 {
		return fmt.Errorf("keywords must not be empty")
	}
	for i, kw := range c.Keywords {
		if strings.TrimSpace(kw.Term) == "" {
			return fmt.Errorf("keywords[%d]: term must not be empty", i)
		}
		if err := validateTier(kw.Tier); err != nil {
			return fmt.Errorf("keywords[%d]: %w", i, err)
		}
	}
	return nil
}

func validateTier(tier string) error {
	switch strings.TrimSpace(tier) {
	case TierHigh, TierMedium, TierLow:
		return nil
	default:
		return fmt.Errorf("tier must be %q, %q, or %q, got %q", TierHigh, TierMedium, TierLow, tier)
	}
}

// Select returns the profile matching bank (BankKey or case-insensitive
// bank_name). If bank is empty and there is exactly one profile, that
// profile is returned; otherwise bank is required.
func (f File) Select(bank string) (Config, error) {
	bank = strings.TrimSpace(bank)
	if bank == "" {
		if len(f.Profiles) == 1 {
			return f.Profiles[0], nil
		}
		return Config{}, fmt.Errorf("--bank is required when config has %d profiles (known: %v)", len(f.Profiles), f.knownKeys())
	}

	wantKey := normalizeBankKey(bank)
	wantName := strings.ToLower(strings.TrimSpace(bank))
	for _, p := range f.Profiles {
		if p.BankKey() == wantKey || strings.ToLower(strings.TrimSpace(p.BankName)) == wantName {
			return p, nil
		}
	}
	return Config{}, fmt.Errorf("no profile for bank %q (known: %v)", bank, f.knownKeys())
}

func (f File) knownKeys() []string {
	keys := make([]string, 0, len(f.Profiles))
	for _, p := range f.Profiles {
		keys = append(keys, p.BankKey())
	}
	return keys
}

// NormalizedKeywords returns include keyword rules with Term trimmed and
// upper-cased and blanks removed — ready for substring matching.
func (c Config) NormalizedKeywords() []KeywordRule {
	out := make([]KeywordRule, 0, len(c.Keywords))
	for _, kw := range c.Keywords {
		term := strings.TrimSpace(kw.Term)
		if term == "" {
			continue
		}
		out = append(out, KeywordRule{
			Term: strings.ToUpper(term),
			Tier: strings.TrimSpace(kw.Tier),
		})
	}
	return out
}

// NormalizedExcludeKeywords returns exclude keywords ready for matching.
func (c Config) NormalizedExcludeKeywords() []string {
	return normalizeKeywordList(c.ExcludeKeywords)
}

func normalizeKeywordList(keywords []string) []string {
	out := make([]string, 0, len(keywords))
	for _, kw := range keywords {
		kw = strings.TrimSpace(kw)
		if kw == "" {
			continue
		}
		out = append(out, strings.ToUpper(kw))
	}
	return out
}

// BankKey normalizes bank_name into the registry key parsers register
// themselves under (lowercase, spaces collapsed to hyphens).
func (c Config) BankKey() string {
	return normalizeBankKey(c.BankName)
}

func normalizeBankKey(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(name), "-"))
}
