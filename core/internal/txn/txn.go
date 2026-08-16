// Package txn defines the canonical transaction shape that every bank
// parser produces. Downstream code (matcher, report) only ever depends on
// this type, never on a specific bank's column layout.
package txn

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Field is one statement column name/value pair.
type Field struct {
	Name  string
	Value string
}

// Fields is an ordered list of columns as they appeared in the statement
// header. It marshals to a JSON object so reports stay human-readable while
// preserving column order (unlike map[string]string).
type Fields []Field

// Transaction is one parsed row from a bank statement.
//
// Description is pulled out explicitly because it's the one field every
// matcher needs to inspect, regardless of which bank produced it. Fields
// holds every column (including Description) in header order so the report
// can show the full row.
type Transaction struct {
	Description string
	Fields      Fields
}

// MarshalJSON encodes Fields as a JSON object with keys in slice order.
func (f Fields) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, field := range f {
		if i > 0 {
			buf.WriteByte(',')
		}
		name, err := json.Marshal(field.Name)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(field.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(name)
		buf.WriteByte(':')
		buf.Write(value)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON decodes a JSON object into Fields, preserving key order
// from the input bytes via a Decoder token stream.
func (f *Fields) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return fmt.Errorf("txn.Fields: expected JSON object, got %v", tok)
	}

	var out Fields
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := keyTok.(string)
		if !ok {
			return fmt.Errorf("txn.Fields: expected string key, got %T", keyTok)
		}
		var value string
		if err := dec.Decode(&value); err != nil {
			return err
		}
		out = append(out, Field{Name: key, Value: value})
	}
	tok, err = dec.Token()
	if err != nil {
		return err
	}
	delim, ok = tok.(json.Delim)
	if !ok || delim != '}' {
		return fmt.Errorf("txn.Fields: expected end of object, got %v", tok)
	}
	*f = out
	return nil
}

// Lookup returns the value for name (exact match), or "" if missing.
func (f Fields) Lookup(name string) string {
	for _, field := range f {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}
