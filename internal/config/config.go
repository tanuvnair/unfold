// Package config loads and validates the JSON config that tells unfold
// which bank profile to use and which keywords count as autopay.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config describes how to parse a bank statement CSV for autopay rows.
type Config struct {
	BankName               string   `json:"bank_name"`
	SkipRows               int      `json:"skip_rows"`
	DescriptionColumnIndex int      `json:"description_column_index"`
	Keywords               []string `json:"keywords"`
}

// Load reads and validates a config file from disk.
func Load(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate catches config mistakes before they turn into a confusing empty
// report or a panic partway through parsing.
func (c Config) Validate() error {
	if strings.TrimSpace(c.BankName) == "" {
		return fmt.Errorf("bank_name is required")
	}
	if c.SkipRows < 0 {
		return fmt.Errorf("skip_rows must be >= 0, got %d", c.SkipRows)
	}
	if c.DescriptionColumnIndex < 0 {
		return fmt.Errorf("description_column_index must be >= 0, got %d", c.DescriptionColumnIndex)
	}
	if len(c.Keywords) == 0 {
		return fmt.Errorf("keywords must not be empty")
	}
	return nil
}

// NormalizedKeywords returns keywords trimmed, upper-cased, and with blanks
// removed — ready for substring matching.
func (c Config) NormalizedKeywords() []string {
	out := make([]string, 0, len(c.Keywords))
	for _, kw := range c.Keywords {
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
	return strings.ToLower(strings.Join(strings.Fields(c.BankName), "-"))
}
