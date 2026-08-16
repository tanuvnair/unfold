// Package kotak implements the parser.Parser interface for Kotak Mahindra
// Bank statement exports. Its quirks — where the header row lands, how
// footer rows are distinguished from transactions — are specific to Kotak's
// export format and don't belong anywhere more general.
package kotak

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/parser"
	"github.com/tanuvnair/unfold/internal/txn"
)

func init() {
	parser.Register("kotak-mahindra-bank", func() parser.Parser {
		return &Parser{}
	})
}

// Parser implements parser.Parser for Kotak exports.
type Parser struct{}

func (p *Parser) Parse(r io.Reader, cfg config.Config) ([]txn.Transaction, error) {
	reader := csv.NewReader(r)
	// Kotak exports mix metadata, transactions, and footer notes.
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	if err := skipRows(reader, cfg.SkipRows); err != nil {
		return nil, fmt.Errorf("skip header rows: %w", err)
	}

	header, err := findHeader(reader)
	if err != nil {
		return nil, fmt.Errorf("find column header: %w", err)
	}

	descIdx, err := parser.ResolveColumnIndex(header, cfg.DescriptionColumn)
	if err != nil {
		return nil, fmt.Errorf("resolve description column: %w", err)
	}
	amountIdx, err := parser.ResolveColumnIndex(header, cfg.AmountColumn)
	if err != nil {
		return nil, fmt.Errorf("resolve amount column: %w", err)
	}
	dateIdx, err := parser.ResolveColumnIndex(header, cfg.DateColumn)
	if err != nil {
		return nil, fmt.Errorf("resolve date column: %w", err)
	}
	typeIdx, err := parser.ResolveColumnIndex(header, cfg.TypeColumn)
	if err != nil {
		return nil, fmt.Errorf("resolve type column: %w", err)
	}

	var transactions []txn.Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read statement row: %w", err)
		}
		if !isTransactionRow(record, len(header)) {
			continue
		}
		t, err := toTransaction(header, record, descIdx, amountIdx, dateIdx, typeIdx, cfg.DateFormat)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func skipRows(reader *csv.Reader, n int) error {
	for i := 0; i < n; i++ {
		if _, err := reader.Read(); err != nil {
			return err
		}
	}
	return nil
}

// findHeader reads rows until it finds the transaction column header.
// encoding/csv silently skips blank lines, so a fixed skip_rows count alone
// is brittle across bank exports.
func findHeader(reader *csv.Reader) ([]string, error) {
	for {
		record, err := reader.Read()
		if err != nil {
			return nil, err
		}
		if looksLikeHeader(record) {
			return record, nil
		}
	}
}

func looksLikeHeader(record []string) bool {
	joined := strings.ToUpper(strings.Join(record, " "))
	return strings.Contains(joined, "DESCRIPTION") &&
		(strings.Contains(joined, "TRANSACTION DATE") ||
			strings.Contains(joined, "VALUE DATE") ||
			strings.Contains(joined, "SL. NO"))
}

// isTransactionRow drops footer/note lines that appear after the
// transaction block.
func isTransactionRow(record []string, expectedCols int) bool {
	if expectedCols > 0 && len(record) < expectedCols {
		return false
	}
	if len(record) == 0 || strings.TrimSpace(record[0]) == "" {
		return false
	}
	for _, r := range record[0] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func toTransaction(header, record []string, descIdx, amountIdx, dateIdx, typeIdx int, dateFormat string) (txn.Transaction, error) {
	fields := make(txn.Fields, 0, len(header))
	seen := make(map[string]struct{}, len(header))
	for i := range header {
		name := uniqueColumnName(header, i, seen)
		seen[name] = struct{}{}
		value := ""
		if i < len(record) {
			value = strings.TrimSpace(record[i])
		}
		fields = append(fields, txn.Field{Name: name, Value: value})
	}

	description := cell(record, descIdx)
	amountRaw := cell(record, amountIdx)
	amount, err := parseAmount(amountRaw)
	if err != nil {
		return txn.Transaction{}, fmt.Errorf("parse amount %q: %w", amountRaw, err)
	}
	dateRaw := cell(record, dateIdx)
	date, err := time.Parse(dateFormat, dateRaw)
	if err != nil {
		return txn.Transaction{}, fmt.Errorf("parse date %q with layout %q: %w", dateRaw, dateFormat, err)
	}
	txnType := strings.ToUpper(strings.TrimSpace(cell(record, typeIdx)))

	return txn.Transaction{
		Description: description,
		Amount:      amount,
		Date:        date,
		Type:        txnType,
		Fields:      fields,
	}, nil
}

func cell(record []string, idx int) string {
	if idx < 0 || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}

// parseAmount converts Indian bank amount strings (commas, optional currency
// symbols) into float64. Empty or non-numeric values are errors.
func parseAmount(raw string) (float64, error) {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "₹", "")
	cleaned = strings.TrimSpace(cleaned)
	if strings.HasPrefix(strings.ToUpper(cleaned), "RS") {
		cleaned = strings.TrimSpace(cleaned[2:])
		cleaned = strings.TrimLeft(cleaned, ". ")
	}
	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, cleaned)
	if cleaned == "" {
		return 0, fmt.Errorf("empty amount")
	}
	return strconv.ParseFloat(cleaned, 64)
}

// uniqueColumnName keeps the first occurrence of a header as-is. Duplicates
// are renamed using the previous column (e.g. second "Dr / Cr" after
// "Balance" becomes "Balance Dr / Cr") so amount and balance sides are both
// preserved.
func uniqueColumnName(header []string, i int, seen map[string]struct{}) string {
	col := strings.TrimSpace(header[i])
	if col == "" {
		return fmt.Sprintf("column_%d", i+1)
	}
	if _, exists := seen[col]; !exists {
		return col
	}
	if i > 0 {
		prev := strings.TrimSpace(header[i-1])
		if prev != "" {
			candidate := prev + " " + col
			if _, exists := seen[candidate]; !exists {
				return candidate
			}
		}
	}
	for n := 2; ; n++ {
		candidate := fmt.Sprintf("%s (%d)", col, n)
		if _, exists := seen[candidate]; !exists {
			return candidate
		}
	}
}
