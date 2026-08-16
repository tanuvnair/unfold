package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tanuvnair/unfold/internal/report"
	"github.com/tanuvnair/unfold/internal/reportquery"
	"github.com/tanuvnair/unfold/internal/txn"
)

func TestHandleReportTransactions(t *testing.T) {
	store := reportquery.NewStore()
	id, err := store.Put(report.Report{
		BankName:         "Kotak Mahindra Bank",
		TransactionCount: 2,
		Transactions: []txn.Fields{
			{
				{Name: "Transaction Date", Value: "03-01-2026 08:16:49"},
				{Name: "Description", Value: "UPI/Netflix/MandateExecute"},
				{Name: "Amount", Value: "199.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
			{
				{Name: "Transaction Date", Value: "05-01-2026 05:24:48"},
				{Name: "Description", Value: "IB:MONTHLY INVESTMENT"},
				{Name: "Amount", Value: "5,000.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	srv := &server{reports: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/reports/{id}/transactions", srv.handleReportTransactions)

	req := httptest.NewRequest(http.MethodGet, "/api/reports/"+id+"/transactions?q=netflix&page=0&page_size=10", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}

	var got reportquery.Result
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.RowCount != 1 || len(got.Rows) != 1 {
		t.Fatalf("got=%+v", got)
	}
	if got.Rows[0].Description != "UPI/Netflix/MandateExecute" {
		t.Fatalf("row=%+v", got.Rows[0])
	}
}

func TestHandleReportTransactions_NotFound(t *testing.T) {
	srv := &server{reports: reportquery.NewStore()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/reports/{id}/transactions", srv.handleReportTransactions)

	req := httptest.NewRequest(http.MethodGet, "/api/reports/missing/transactions", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rr.Code)
	}
}
