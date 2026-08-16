package reportquery

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/tanuvnair/unfold/internal/report"
)

// Store keeps analyzed reports in process memory so the UI can page and
// filter them without writing statements to disk.
type Store struct {
	mu      sync.RWMutex
	reports map[string]report.Report
}

// NewStore returns an empty in-memory report store.
func NewStore() *Store {
	return &Store{reports: make(map[string]report.Report)}
}

// Put stores r and returns a random id for later queries.
func (s *Store) Put(r report.Report) (string, error) {
	id, err := newID()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.reports[id] = r
	s.mu.Unlock()
	return id, nil
}

// Get returns the report for id, if present.
func (s *Store) Get(id string) (report.Report, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.reports[id]
	return r, ok
}

func newID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
