package main

import (
	"net/http"
	"strconv"

	"github.com/tanuvnair/unfold/internal/reportquery"
)

func (s *server) handleReportTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rpt, ok := s.reports.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "report not found")
		return
	}

	q := r.URL.Query()
	writeJSON(w, http.StatusOK, reportquery.Apply(rpt.Transactions, reportquery.Query{
		Search:     q.Get("q"),
		Type:       q.Get("type"),
		Confidence: q.Get("confidence"),
		Source:     q.Get("source"),
		Payee:      q.Get("payee"),
		Page:       parseIntParam(q.Get("page"), 0),
		PageSize:   parseIntParam(q.Get("page_size"), 10),
	}))
}

func (s *server) handleReportSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rpt, ok := s.reports.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "report not found")
		return
	}

	q := r.URL.Query()
	writeJSON(w, http.StatusOK, reportquery.Summary(rpt.Transactions, reportquery.SummaryQuery{
		Search:     q.Get("q"),
		Type:       q.Get("type"),
		Confidence: q.Get("confidence"),
		Source:     q.Get("source"),
	}))
}

func parseIntParam(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}
