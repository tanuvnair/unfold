package webui

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHasUI_MatchesEmbeddedIndex(t *testing.T) {
	_, err := dist.Open("dist/index.html")
	want := err == nil
	if got := HasUI(); got != want {
		t.Fatalf("HasUI() = %v, want %v", got, want)
	}
}

func TestHandler_UnknownPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/results", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)

	if HasUI() {
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 SPA fallback", rec.Code)
		}
		return
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 without embedded UI", rec.Code)
	}
}
