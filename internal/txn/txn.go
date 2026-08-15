// Package txn defines the canonical transaction shape that every bank
// parser produces. Downstream code (matcher, report) only ever depends on
// this type, never on a specific bank's column layout.
package txn

// Transaction is one parsed row from a bank statement.
//
// Description is pulled out explicitly because it's the one field every
// matcher needs to inspect, regardless of which bank produced it. Fields
// holds every column (including Description) keyed by its header name, so
// the report can still show the full row.
type Transaction struct {
	Description string
	Fields      map[string]string
}
