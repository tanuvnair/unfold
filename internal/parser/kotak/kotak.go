// Package kotak implements the parser.Parser interface for Kotak Mahindra
// Bank statement exports. Its quirks — where the header row lands, how
// footer rows are distinguished from transactions — are specific to Kotak's
// export format and don't belong anywhere more general.
package kotak

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

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
		transactions = append(transactions, toTransaction(header, record, descIdx))
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

func toTransaction(header, record []string, descIdx int) txn.Transaction {
	fields := make(txn.Fields, 0, len(header))
	indexByName := make(map[string]int, len(header))
	for i, col := range header {
		col = strings.TrimSpace(col)
		if col == "" {
			col = fmt.Sprintf("column_%d", i+1)
		}
		value := ""
		if i < len(record) {
			value = strings.TrimSpace(record[i])
		}
		// Duplicate headers: last wins (mirrors the old map behavior).
		if idx, ok := indexByName[col]; ok {
			fields[idx].Value = value
			continue
		}
		indexByName[col] = len(fields)
		fields = append(fields, txn.Field{Name: col, Value: value})
	}

	description := ""
	if descIdx >= 0 && descIdx < len(record) {
		description = strings.TrimSpace(record[descIdx])
	}

	return txn.Transaction{
		Description: description,
		Fields:      fields,
	}
}
