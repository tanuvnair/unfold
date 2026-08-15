package txn_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/tanuvnair/unfold/internal/txn"
)

func TestFieldsMarshalPreservesKeyOrder(t *testing.T) {
	fields := txn.Fields{
		{Name: "Sl. No.", Value: "1"},
		{Name: "Description", Value: "NACH SIP"},
		{Name: "Amount", Value: "100.00"},
	}
	raw, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	sl := strings.Index(got, `"Sl. No."`)
	desc := strings.Index(got, `"Description"`)
	amt := strings.Index(got, `"Amount"`)
	if sl < 0 || desc < 0 || amt < 0 {
		t.Fatalf("missing keys in %s", got)
	}
	if !(sl < desc && desc < amt) {
		t.Fatalf("key order not preserved in %s", got)
	}
}

func TestFieldsRoundTrip(t *testing.T) {
	want := txn.Fields{
		{Name: "Transaction Date", Value: "01-01-2026"},
		{Name: "Description", Value: "AUTOPAY rent"},
	}
	raw, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got txn.Fields
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("len=%d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("field[%d]=%v want %v", i, got[i], want[i])
		}
	}
}

func TestReportTransactionsKeepColumnOrder(t *testing.T) {
	type wrap struct {
		Transactions []txn.Fields `json:"transactions"`
	}
	w := wrap{Transactions: []txn.Fields{{
		{Name: "Zebra", Value: "z"},
		{Name: "Alpha", Value: "a"},
	}}}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(w); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Index(out, `"Zebra"`) > strings.Index(out, `"Alpha"`) {
		t.Fatalf("expected Zebra before Alpha in %s", out)
	}
}
