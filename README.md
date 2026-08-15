# unfold
Find hidden subscriptions in seconds.

A local Go tool that parses bank statement CSVs, matches autopay/mandate
keywords, and reports the hits as JSON. Use the CLI, HTTP API, or web UI.

## Layout

| Path | Role |
|------|------|
| `cmd/cli` | CLI (`unfold` binary) |
| `cmd/api` | Local HTTP API |
| `internal/` | Shared parse → match → report pipeline |
| `configs/` | Bank profiles |
| `web/` | Vite + React + TanStack Router + TanStack Table + shadcn UI |

## CLI

```bash
make cli
./unfold configs/banks.json path/to/statement.csv
```

Useful flags: `--bank`, `-v` / `--verbose`, `--diff`, `--dry-run`.

Writes `autopay_report.json` next to the statement CSV (unless `--dry-run`).

## API + UI

**Production / single binary (preferred):** Vite build is copied into `internal/webui/dist` and embedded. One process serves `/` and `/api` on the same origin — leave `UNFOLD_WEB_API_BASE` empty.

```bash
make serve
# open http://127.0.0.1:8080
```

**Local development** (hot reload): API and Vite in two terminals.

```bash
make api
# or: UNFOLD_API_PORT=9090 make api
```

```bash
make web
# or: UNFOLD_WEB_PORT=3000 UNFOLD_API_PORT=9090 make web
```

Open http://localhost:5173 — Vite proxies `/api` to the Go server.

Copy [`.env.example`](.env.example) to `.env` at the repo root (Vite also reads it). Env vars:

| Variable | Default | Used by |
|----------|---------|---------|
| `UNFOLD_API_HOST` | `127.0.0.1` | API listen host; Vite proxy |
| `UNFOLD_API_PORT` | `8080` | API listen port; Vite proxy |
| `UNFOLD_API_ADDR` | *(host:port)* | Full API listen address (overrides host+port) |
| `UNFOLD_API_CONFIG` | `configs/banks.json` | API bank profiles |
| `UNFOLD_API_CORS_ORIGIN` | `http://localhost:5173` | API CORS allowlist (UI origin) |
| `UNFOLD_WEB_HOST` | `localhost` | Vite dev server host |
| `UNFOLD_WEB_PORT` / `PORT` | `5173` | Vite dev server port |
| `UNFOLD_WEB_API_BASE` | *(empty)* | Absolute API origin for split UI/API; empty = same-origin `/api` |
| `UNFOLD_API_PROXY_TARGET` | `http://$API_HOST:$API_PORT` | Vite `/api` proxy (dev, when API base empty) |

CLI flags (`-addr`, `-config`, `-cors-origin`) still override the env defaults when passed.

### API endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/health` | Liveness |
| `GET` | `/api/banks` | Configured bank profiles |
| `POST` | `/api/analyze` | Multipart form: `file` (CSV), optional `bank` |
| `GET` | `/api/reports/{id}/transactions` | Page/filter the in-memory report from analyze |

**Analyze** stores the full matched report in API process memory (not on disk)
and returns:

```json
{ "id": "<hex>", "bank_name": "Kotak Mahindra Bank", "transaction_count": 52 }
```

Restarting the API drops stored reports. The CLI is unchanged: it still writes
the full `autopay_report.json` next to the statement.

**Transactions** query params:

| Param | Default | Notes |
|-------|---------|-------|
| `q` | *(empty)* | Case-insensitive substring on Description |
| `type` | *(all)* | `DR` or `CR` |
| `page` | `0` | **0-based** page index |
| `page_size` | `10` | Clamped to 1–100 |

Response:

```json
{
  "rows": [
    {
      "id": "0",
      "transaction_date": "03-01-2026 08:16:49",
      "description": "UPI/Netflix/…/MandateExecute",
      "amount": "199.00",
      "type": "DR"
    }
  ],
  "row_count": 12,
  "page": 0,
  "page_size": 10
}
```

`row_count` is the filtered total (not just the current page). `id` on each row
is the original index in the matched report, stable across pages.

### Results UI

The web results table is a shadcn DataTable with **server-side** search, DR/CR
filter, and numbered pagination. Table state is in the URL (1-based `page`):

```
/results
/results?q=netflix
/results?type=DR&page=2&pageSize=25
```

The browser keeps `{ id, bank_name, transaction_count }` in session memory and
fetches each page from `/api/reports/{id}/transactions`. Refreshing the tab or
restarting the API requires analyzing again.

Home and results both align content to the **top** of the viewport (same
`Page` padding and stack).

## Tests

```bash
make test
```
