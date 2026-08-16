package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
	from, to, err := parseDateRange(q.Get("from"), q.Get("to"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reportquery.Apply(rpt.Transactions, reportquery.Query{
		Search:     q.Get("q"),
		Type:       q.Get("type"),
		Confidence: q.Get("confidence"),
		Source:     q.Get("source"),
		Payee:      q.Get("payee"),
		DateFrom:   from,
		DateTo:     to,
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
	from, to, err := parseDateRange(q.Get("from"), q.Get("to"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reportquery.Summary(rpt.Transactions, reportquery.SummaryQuery{
		Search:     q.Get("q"),
		Type:       q.Get("type"),
		Confidence: q.Get("confidence"),
		Source:     q.Get("source"),
		DateFrom:   from,
		DateTo:     to,
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

func parseDateRange(fromRaw, toRaw string) (from, to time.Time, err error) {
	from, err = parseDateParam(fromRaw)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err = parseDateParam(toRaw)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return from, to, nil
}

func parseDateParam(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, nil
	}
	t, err := time.ParseInLocation("2006-01-02", raw, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
