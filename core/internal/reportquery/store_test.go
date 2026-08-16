package reportquery

import (
	"testing"

	"github.com/tanuvnair/unfold/internal/report"
)

func TestStorePutGet(t *testing.T) {
	s := NewStore()
	id, err := s.Put(report.Report{BankName: "Kotak Mahindra Bank", TransactionCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("empty id")
	}
	got, ok := s.Get(id)
	if !ok {
		t.Fatal("missing report")
	}
	if got.BankName != "Kotak Mahindra Bank" {
		t.Fatalf("got %+v", got)
	}
	if _, ok := s.Get("missing"); ok {
		t.Fatal("expected miss")
	}
}
