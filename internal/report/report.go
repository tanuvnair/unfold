// Package report builds and writes the JSON output file.
package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tanuvnair/unfold/internal/txn"
)

// Report is the top-level shape written to autopay_report.json.
type Report struct {
	BankName         string              `json:"bank_name"`
	TransactionCount int                 `json:"transaction_count"`
	Transactions     []map[string]string `json:"transactions"`
}

// Build converts matched transactions into the report shape.
func Build(bankName string, transactions []txn.Transaction) Report {
	fields := make([]map[string]string, 0, len(transactions))
	for _, t := range transactions {
		fields = append(fields, t.Fields)
	}
	return Report{
		BankName:         bankName,
		TransactionCount: len(fields),
		Transactions:     fields,
	}
}

// PathFor places autopay_report.json in the same directory as the input
// statement CSV, rather than the current working directory.
func PathFor(statementPath string) string {
	dir := filepath.Dir(statementPath)
	return filepath.Join(dir, "autopay_report.json")
}

// Write pretty-prints the report to path, human-readable and with HTML
// escaping disabled so merchant names with "&" don't turn into \u0026.
func Write(path string, r Report) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("encode report: %w", err)
	}
	return nil
}
