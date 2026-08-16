// Package report builds and writes the JSON output file.
package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tanuvnair/unfold/internal/txn"
)

// Detection source values for Entry.DetectionSource.
const (
	SourceKeyword    = "keyword"
	SourceRecurrence = "recurrence"
	SourceBoth       = "both"
)

// Reserved Entry JSON keys that are detection metadata, not bank columns.
const (
	jsonKeyConfidence          = "confidence"
	jsonKeyMatchedTerm         = "matched_term"
	jsonKeyDetectionSource     = "detection_source"
	jsonKeyPayeeToken          = "payee_token"
	jsonKeyAmountVarianceRatio = "amount_variance_ratio"
	jsonKeyAvgIntervalDays     = "avg_interval_days"
)

// Report is the top-level shape written to autopay_report.json.
type Report struct {
	BankName         string  `json:"bank_name"`
	TransactionCount int     `json:"transaction_count"`
	Transactions     []Entry `json:"transactions"`
}

// Entry is one matched transaction: detection metadata plus the full bank row.
// JSON flattens detection keys alongside bank columns.
// Date is set at analyze time for server-side filtering and is not written to
// autopay_report.json (same as HasRecurrenceMetrics).
type Entry struct {
	Confidence           string
	MatchedTerm          string
	DetectionSource      string
	PayeeToken           string
	AmountVarianceRatio  float64
	AvgIntervalDays      float64
	HasRecurrenceMetrics bool // controls whether variance/interval are written
	Date                 time.Time
	Fields               txn.Fields
}

// DiffResult summarizes what changed between two reports.
type DiffResult struct {
	Added     []Entry
	Removed   []Entry
	Unchanged int
}

// Build assembles a report from already-merged entries.
func Build(bankName string, entries []Entry) Report {
	return Report{
		BankName:         bankName,
		TransactionCount: len(entries),
		Transactions:     entries,
	}
}

// MarshalJSON encodes Entry as a flat object: detection keys then bank columns.
func (e Entry) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	writeKV := func(key, value string, omitEmpty bool) error {
		if omitEmpty && value == "" {
			return nil
		}
		if !first {
			buf.WriteByte(',')
		}
		first = false
		k, err := json.Marshal(key)
		if err != nil {
			return err
		}
		v, err := json.Marshal(value)
		if err != nil {
			return err
		}
		buf.Write(k)
		buf.WriteByte(':')
		buf.Write(v)
		return nil
	}
	writeFloat := func(key string, value float64) error {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		k, err := json.Marshal(key)
		if err != nil {
			return err
		}
		buf.Write(k)
		buf.WriteByte(':')
		buf.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
		return nil
	}
	if err := writeKV(jsonKeyConfidence, e.Confidence, true); err != nil {
		return nil, err
	}
	if err := writeKV(jsonKeyMatchedTerm, e.MatchedTerm, true); err != nil {
		return nil, err
	}
	if err := writeKV(jsonKeyDetectionSource, e.DetectionSource, true); err != nil {
		return nil, err
	}
	if err := writeKV(jsonKeyPayeeToken, e.PayeeToken, true); err != nil {
		return nil, err
	}
	if e.HasRecurrenceMetrics {
		if err := writeFloat(jsonKeyAmountVarianceRatio, e.AmountVarianceRatio); err != nil {
			return nil, err
		}
		if err := writeFloat(jsonKeyAvgIntervalDays, e.AvgIntervalDays); err != nil {
			return nil, err
		}
	}
	for _, field := range e.Fields {
		if isReservedEntryKey(field.Name) {
			continue
		}
		if err := writeKV(field.Name, field.Value, false); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON decodes a flat object into detection metadata + Fields.
func (e *Entry) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return fmt.Errorf("report.Entry: expected JSON object, got %v", tok)
	}

	var (
		entry  Entry
		fields txn.Fields
	)
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := keyTok.(string)
		if !ok {
			return fmt.Errorf("report.Entry: expected string key, got %T", keyTok)
		}
		switch key {
		case jsonKeyAmountVarianceRatio, jsonKeyAvgIntervalDays:
			var value float64
			if err := dec.Decode(&value); err != nil {
				return err
			}
			entry.HasRecurrenceMetrics = true
			if key == jsonKeyAmountVarianceRatio {
				entry.AmountVarianceRatio = value
			} else {
				entry.AvgIntervalDays = value
			}
		default:
			var value string
			if err := dec.Decode(&value); err != nil {
				return err
			}
			switch key {
			case jsonKeyConfidence:
				entry.Confidence = value
			case jsonKeyMatchedTerm:
				entry.MatchedTerm = value
			case jsonKeyDetectionSource:
				entry.DetectionSource = value
			case jsonKeyPayeeToken:
				entry.PayeeToken = value
			default:
				fields = append(fields, txn.Field{Name: key, Value: value})
			}
		}
	}
	tok, err = dec.Token()
	if err != nil {
		return err
	}
	delim, ok = tok.(json.Delim)
	if !ok || delim != '}' {
		return fmt.Errorf("report.Entry: expected end of object, got %v", tok)
	}
	entry.Fields = fields
	*e = entry
	return nil
}

func isReservedEntryKey(name string) bool {
	switch name {
	case jsonKeyConfidence, jsonKeyMatchedTerm, jsonKeyDetectionSource,
		jsonKeyPayeeToken, jsonKeyAmountVarianceRatio, jsonKeyAvgIntervalDays:
		return true
	default:
		return false
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
// fingerprints. Detection metadata is ignored for identity.
func Diff(prev, next Report) DiffResult {
	prevCounts := countFingerprints(prev.Transactions)

	var added []Entry
	unchanged := 0
	for _, row := range next.Transactions {
		fp := FieldsFingerprint(row.Fields)
		if prevCounts[fp] > 0 {
			prevCounts[fp]--
			unchanged++
			continue
		}
		added = append(added, row)
	}

	var removed []Entry
	for _, row := range prev.Transactions {
		fp := FieldsFingerprint(row.Fields)
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

func countFingerprints(rows []Entry) map[string]int {
	out := make(map[string]int, len(rows))
	for _, row := range rows {
		out[FieldsFingerprint(row.Fields)]++
	}
	return out
}

// FieldsFingerprint is a stable identity for a bank row (used by Diff and merge).
func FieldsFingerprint(fields txn.Fields) string {
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
	for _, f := range fields {
		if strings.EqualFold(f.Name, "Description") {
			return f.Value
		}
	}
	return ""
}

// EntryDescription returns the Description column for an Entry.
func EntryDescription(e Entry) string {
	return DescriptionOf(e.Fields)
}
