package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tanuvnair/unfold/internal/analyze"
	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/parser"
	_ "github.com/tanuvnair/unfold/internal/parser/kotak"
	"github.com/tanuvnair/unfold/internal/reportquery"
	"github.com/tanuvnair/unfold/internal/webui"
)

const maxUploadBytes = 10 << 20 // 10 MiB

func main() {
	defaultAddr := listenAddrFromEnv()
	// Resolved relative to core/ when launched via Makefile (api/serve).
	defaultConfig := envOr("UNFOLD_API_CONFIG", "configs/banks.json")
	defaultCORS := envOr("UNFOLD_API_CORS_ORIGIN", "http://localhost:5173")

	addr := flag.String("addr", defaultAddr, "listen address (env: UNFOLD_API_ADDR or UNFOLD_API_HOST+UNFOLD_API_PORT)")
	configPath := flag.String("config", defaultConfig, "path to bank profiles config (env: UNFOLD_API_CONFIG)")
	corsOrigin := flag.String("cors-origin", defaultCORS, "allowed CORS origin (env: UNFOLD_API_CORS_ORIGIN)")
	flag.Parse()

	file, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	srv := &server{profiles: file, reports: reportquery.NewStore()}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", srv.handleHealth)
	mux.HandleFunc("GET /api/banks", srv.handleBanks)
	mux.HandleFunc("POST /api/analyze", srv.handleAnalyze)
	mux.HandleFunc("GET /api/reports/{id}/transactions", srv.handleReportTransactions)
	if webui.HasUI() {
		mux.Handle("/", webui.Handler())
	}

	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           corsMiddleware(*corsOrigin, mux),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
	}

	uiNote := "api-only"
	if webui.HasUI() {
		uiNote = "ui+api"
	}
	log.Printf("unfold API listening on http://%s (%s config=%s cors=%s)", *addr, uiNote, *configPath, *corsOrigin)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

// listenAddrFromEnv resolves UNFOLD_API_ADDR, or UNFOLD_API_HOST + UNFOLD_API_PORT.
func listenAddrFromEnv() string {
	if addr := strings.TrimSpace(os.Getenv("UNFOLD_API_ADDR")); addr != "" {
		return addr
	}
	host := envOr("UNFOLD_API_HOST", "127.0.0.1")
	port := envOr("UNFOLD_API_PORT", "8080")
	return net.JoinHostPort(host, port)
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

type server struct {
	profiles config.File
	reports  *reportquery.Store
}

type analyzeResponse struct {
	ID               string `json:"id"`
	BankName         string `json:"bank_name"`
	TransactionCount int    `json:"transaction_count"`
}

type bankInfo struct {
	Key       string `json:"key"`
	BankName  string `json:"bank_name"`
	HasParser bool   `json:"has_parser"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) handleBanks(w http.ResponseWriter, _ *http.Request) {
	registered := make(map[string]struct{}, len(parser.Keys()))
	for _, k := range parser.Keys() {
		registered[k] = struct{}{}
	}

	banks := make([]bankInfo, 0, len(s.profiles.Profiles))
	for _, p := range s.profiles.Profiles {
		key := p.BankKey()
		_, ok := registered[key]
		banks = append(banks, bankInfo{
			Key:       key,
			BankName:  p.BankName,
			HasParser: ok,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"banks": banks})
}

func (s *server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, "upload too large or invalid multipart form (max 10MB)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing form field \"file\"")
		return
	}
	defer file.Close()

	name := strings.ToLower(header.Filename)
	if !strings.HasSuffix(name, ".csv") {
		writeError(w, http.StatusBadRequest, "file must be a .csv statement export")
		return
	}
	ct := header.Header.Get("Content-Type")
	if ct != "" && !isAllowedCSVContentType(ct) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("unsupported content type %q", ct))
		return
	}

	bank := r.FormValue("bank")
	cfg, err := s.profiles.Select(bank)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	rpt, err := analyze.Run(cfg, file)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	id, err := s.reports.Put(rpt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not store report")
		return
	}
	writeJSON(w, http.StatusOK, analyzeResponse{
		ID:               id,
		BankName:         rpt.BankName,
		TransactionCount: rpt.TransactionCount,
	})
}

func isAllowedCSVContentType(ct string) bool {
	ct = strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	switch ct {
	case "text/csv", "application/csv", "application/vnd.ms-excel", "text/plain", "application/octet-stream":
		return true
	default:
		return false
	}
}

func corsMiddleware(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
