package reportquery

import (
	"testing"

	"github.com/tanuvnair/unfold/internal/report"
	"github.com/tanuvnair/unfold/internal/txn"
)

func TestSummary_AggregatesByPayee(t *testing.T) {
	rows := []report.Entry{
		{
			Confidence:           "high",
			DetectionSource:      report.SourceRecurrence,
			PayeeToken:           "NETFLIX",
			HasRecurrenceMetrics: true,
			AvgIntervalDays:      30,
			Fields: txn.Fields{
				{Name: "Transaction Date", Value: "05-01-2025 10:00:00"},
				{Name: "Description", Value: "UPI/Netflix/1/Payment"},
				{Name: "Amount", Value: "199.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
		{
			Confidence:           "high",
			DetectionSource:      report.SourceRecurrence,
			PayeeToken:           "NETFLIX",
			HasRecurrenceMetrics: true,
			AvgIntervalDays:      30,
			Fields: txn.Fields{
				{Name: "Transaction Date", Value: "05-02-2025 10:00:00"},
				{Name: "Description", Value: "UPI/Netflix/2/Payment"},
				{Name: "Amount", Value: "199.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
		{
			Confidence:      "medium",
			DetectionSource: report.SourceKeyword,
			PayeeToken:      "MUTUAL FUND",
			Fields: txn.Fields{
				{Name: "Transaction Date", Value: "01-01-2025 10:00:00"},
				{Name: "Description", Value: "NACH MUTUAL FUND"},
				{Name: "Amount", Value: "5,000.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
	}

	got := Summary(rows, SummaryQuery{})
	if got.GroupCount != 2 {
		t.Fatalf("GroupCount=%d want 2", got.GroupCount)
	}

	var netflix *SummaryGroup
	for i := range got.Groups {
		if got.Groups[i].Payee == "NETFLIX" {
			netflix = &got.Groups[i]
		}
	}
	if netflix == nil {
		t.Fatal("missing NETFLIX group")
	}
	if netflix.OccurrenceCount != 2 || netflix.TotalAmount != 398 {
		t.Fatalf("netflix=%+v", netflix)
	}
	if !netflix.IsMonthlyCadence || netflix.MonthlyEstimate != 199 {
		t.Fatalf("netflix monthly=%+v", netflix)
	}
	if got.EstimatedMonthlyTotal != 199 {
		t.Fatalf("EstimatedMonthlyTotal=%v want 199", got.EstimatedMonthlyTotal)
	}
}

func TestApply_PayeeFilter(t *testing.T) {
	rows := []report.Entry{
		{
			PayeeToken: "NETFLIX",
			Fields: txn.Fields{
				{Name: "Description", Value: "UPI/Netflix/1"},
				{Name: "Amount", Value: "199.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
		{
			PayeeToken: "SPOTIFY",
			Fields: txn.Fields{
				{Name: "Description", Value: "UPI/Spotify/1"},
				{Name: "Amount", Value: "99.00"},
				{Name: "Dr / Cr", Value: "DR"},
			},
		},
	}
	got := Apply(rows, Query{Payee: "NETFLIX", PageSize: 10})
	if got.RowCount != 1 || got.Rows[0].Description != "UPI/Netflix/1" {
		t.Fatalf("got=%+v", got)
	}
}
