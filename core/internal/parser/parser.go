// Package parser defines the seam between "raw bank export" and "canonical
// transactions". Each bank gets its own implementation of Parser, registered
// under a key derived from config.BankKey(). main.go and everything
// downstream of parsing never needs a bank-specific switch statement.
package parser

import (
	"fmt"
	"io"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/txn"
)

// Parser turns a raw statement export into canonical transactions.
// Implementations own all bank-specific quirks: header detection, footer
// rows, date formats, which column is the description, etc.
type Parser interface {
	// Parse reads a full statement export and returns every transaction row
	// it can find. It does NOT filter by keyword — that's the matcher's job.
	Parse(r io.Reader, cfg config.Config) ([]txn.Transaction, error)
}

// Factory constructs a fresh Parser instance. Parsers are constructed fresh
// per run rather than shared, so they're free to hold per-run state without
// worrying about concurrent reuse.
type Factory func() Parser

var registry = map[string]Factory{}

// Register makes a parser available under bankKey. Called from each bank
// package's init(), so importing the package for its side effect is what
// registers it — see cmd/unfold/main.go's blank imports.
func Register(bankKey string, factory Factory) {
	if _, exists := registry[bankKey]; exists {
		panic(fmt.Sprintf("parser: duplicate registration for %q", bankKey))
	}
	registry[bankKey] = factory
}

// Get looks up the parser registered for bankKey.
func Get(bankKey string) (Parser, error) {
	factory, ok := registry[bankKey]
	if !ok {
		return nil, fmt.Errorf("no parser registered for bank %q (known: %v)", bankKey, Keys())
	}
	return factory(), nil
}

// Keys returns the registered bank parser keys.
func Keys() []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}
