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
		Search:   q.Get("q"),
		Type:     q.Get("type"),
		Page:     parseIntParam(q.Get("page"), 0),
		PageSize: parseIntParam(q.Get("page_size"), 10),
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
