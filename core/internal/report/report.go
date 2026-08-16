// Package report builds and writes the JSON output file.
package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tanuvnair/unfold/internal/txn"
)

// Report is the top-level shape written to autopay_report.json.
type Report struct {
	BankName         string       `json:"bank_name"`
	TransactionCount int          `json:"transaction_count"`
	Transactions     []txn.Fields `json:"transactions"`
}

// DiffResult summarizes what changed between two reports.
type DiffResult struct {
	Added     []txn.Fields
	Removed   []txn.Fields
	Unchanged int
}

// Build converts matched transactions into the report shape.
func Build(bankName string, transactions []txn.Transaction) Report {
	rows := make([]txn.Fields, 0, len(transactions))
	for _, t := range transactions {
		rows = append(rows, t.Fields)
	}
	return Report{
		BankName:         bankName,
		TransactionCount: len(rows),
		Transactions:     rows,
	}
}

// PathFor places autopay_report.json in the same directory as the input
// statement CSV, rather than the current working directory.
func PathFor(statementPath string) string {
	dir := filepath.Dir(statementPath)
	return filepath.Join(dir, "autopay_report.json")
}

// Read loads a previously written report from path.
func Read(path string) (Report, error) {
	f, err := os.Open(path)
	if err != nil {
		return Report{}, err
	}
	defer f.Close()

	var r Report
	if err := json.NewDecoder(f).Decode(&r); err != nil {
		return Report{}, fmt.Errorf("decode report: %w", err)
	}
	return r, nil
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

// Diff compares previous and next reports by order-independent field
// fingerprints so older map-shaped JSON and ordered Fields both work.
func Diff(prev, next Report) DiffResult {
	prevCounts := countFingerprints(prev.Transactions)

	var added []txn.Fields
	unchanged := 0
	for _, row := range next.Transactions {
		fp := fingerprint(row)
		if prevCounts[fp] > 0 {
			prevCounts[fp]--
			unchanged++
			continue
		}
		added = append(added, row)
	}

	var removed []txn.Fields
	for _, row := range prev.Transactions {
		fp := fingerprint(row)
		if prevCounts[fp] > 0 {
			removed = append(removed, row)
			prevCounts[fp]--
		}
	}

	return DiffResult{
		Added:     added,
		Removed:   removed,
		Unchanged: unchanged,
	}
}

func countFingerprints(rows []txn.Fields) map[string]int {
	out := make(map[string]int, len(rows))
	for _, row := range rows {
		out[fingerprint(row)]++
	}
	return out
}

func fingerprint(fields txn.Fields) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		parts = append(parts, f.Name+"\x00"+f.Value)
	}
	sort.Strings(parts)
	return strings.Join(parts, "\x01")
}

// DescriptionOf returns the Description column value for a row, if present.
func DescriptionOf(fields txn.Fields) string {
	if v := fields.Lookup("Description"); v != "" {
		return v
	}
	// Case-insensitive fallback for bank exports that vary capitalization.
	for _, f := range fields {
		if strings.EqualFold(f.Name, "Description") {
			return f.Value
		}
	}
	return ""
}
