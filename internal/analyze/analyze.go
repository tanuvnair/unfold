// Package analyze runs the shared statement pipeline used by the CLI and API:
// parse → keyword match → report.
package analyze

import (
	"fmt"
	"io"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/matcher"
	"github.com/tanuvnair/unfold/internal/parser"
	"github.com/tanuvnair/unfold/internal/report"
)

// Run parses statement with the bank parser for cfg, filters autopay
// matches, and builds a report. It does not write anything to disk.
func Run(cfg config.Config, statement io.Reader) (report.Report, error) {
	bankParser, err := parser.Get(cfg.BankKey())
	if err != nil {
		return report.Report{}, err
	}

	transactions, err := bankParser.Parse(statement, cfg)
	if err != nil {
		return report.Report{}, fmt.Errorf("parse statement: %w", err)
	}

	matched := matcher.Filter(
		transactions,
		cfg.NormalizedKeywords(),
		cfg.NormalizedExcludeKeywords(),
	)
	return report.Build(cfg.BankName, matched), nil
}
