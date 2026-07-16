package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: unfold <config.json> <statement.csv>")
	}

	config, err := loadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	file, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatalf("open statement: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Bank exports mix metadata, transactions, and footer notes.
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	if err := skipRows(reader, config.SkipRows); err != nil {
		log.Fatalf("skip header rows: %v", err)
	}

	header, err := findHeader(reader)
	if err != nil {
		log.Fatalf("find column header: %v", err)
	}

	outputFile, err := os.Create("autopay_report.csv")
	if err != nil {
		log.Fatalf("create output: %v", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		log.Fatalf("write header: %v", err)
	}

	keywords := normalizeKeywords(config.Keywords)
	fmt.Printf("Unfolding %s account...\n", config.BankName)

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("read statement row: %v", err)
		}
		if !isTransactionRow(record, len(header)) {
			continue
		}
		if !matchesAutopay(record, config.DescriptionColumnIndex, keywords) {
			continue
		}
		if err := writer.Write(record); err != nil {
			log.Fatalf("write match: %v", err)
		}
		count++
	}

	if err := writer.Error(); err != nil {
		log.Fatalf("flush output: %v", err)
	}

	fmt.Printf("Audit complete. Found %d autopay transactions.\n", count)
}

func loadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
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

func normalizeKeywords(keywords []string) []string {
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

// isTransactionRow drops footer/note lines that appear after the transaction block.
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

func matchesAutopay(record []string, descIdx int, keywords []string) bool {
	text := searchText(record, descIdx)
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func searchText(record []string, descIdx int) string {
	if descIdx >= 0 && descIdx < len(record) {
		return strings.ToUpper(record[descIdx])
	}
	return strings.ToUpper(strings.Join(record, " "))
}
