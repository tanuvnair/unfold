# unfold

Find hidden subscriptions in seconds.

Local, private matching for Indian bank statement CSVs. unfold parses the
export, matches autopay / mandate / NACH language on the description, and
returns the hits — no bank login, no cloud sync, no accounts. Use the CLI,
local HTTP API, or web UI; all three share one parse → match → report pipeline.

**Today:** Kotak Mahindra Bank CSV profiles. Other banks are not promised.

## Quick start

```bash
cp .env.example .env   # optional; defaults are fine for local use

# Single binary (API + embedded UI)
make serve
# open http://127.0.0.1:8080

# Or hot-reload API + Vite together
make dev
# UI:  http://localhost:5173
# API: http://127.0.0.1:8080
```

CLI:

```bash
make cli-build
./dist/unfold-cli core/configs/banks.json path/to/statement.csv
```

Useful flags: `--bank`, `-v` / `--verbose`, `--diff`, `--dry-run`.

Writes `autopay_report.json` next to the statement CSV (unless `--dry-run`).

## Layout

| Path | Role |
|------|------|
| `core/` | Go module (CLI, API, shared pipeline, bank configs) |
| `core/cmd/cli` | CLI (`unfold-cli`) |
| `core/cmd/api` | Local HTTP API (`unfold-api`; embeds the web UI when built) |
| `core/internal/` | Shared parse → match → report pipeline |
| `core/configs/` | Bank profiles (`banks.json`) |
| `web/` | Vite + React + TanStack Router / Query / Table + shadcn/ui |

## Requirements

- Go (module in `core/`)
- Node.js + npm (web UI)
- Make

## API + UI

**Production / single binary (preferred):** Vite build is copied into
`core/internal/webui/dist` and embedded. One process serves `/` and `/api` on
the same origin — leave `UNFOLD_WEB_API_BASE` empty. `make api` / `make serve`
run with CWD=`core/` so `UNFOLD_API_CONFIG` defaults to `configs/banks.json`
there.

```bash
make serve
# open http://127.0.0.1:8080
```

**Local development** (hot reload): API and Vite together, or in two terminals.

```bash
make dev
# or separately:
make api          # UNFOLD_API_PORT=9090 make api
make web          # UNFOLD_WEB_PORT=3000 UNFOLD_API_PORT=9090 make web
```

Open http://localhost:5173 — Vite proxies `/api` to the Go server.

Copy [`.env.example`](.env.example) to `.env` at the repo root (Vite also
reads it). Env vars:

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

CLI flags (`-addr`, `-config`, `-cors-origin`) still override the env defaults
when passed. The default listen address is loopback — this is a local tool,
not a hosted multi-tenant service.

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

Uploads must be CSV; size is capped at 10 MiB.

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

Home and results both align content to the **top** of the viewport.

## Make targets

| Target | What it does |
|--------|----------------|
| `make serve` | Build and run API + embedded UI |
| `make dev` | API from source + Vite (Ctrl+C stops both) |
| `make cli-build` | Build `./dist/unfold-cli` |
| `make api-build` | Build `./dist/unfold-api` (embeds UI) |
| `make test` | Go tests |
| `make check` | fmt, vet, test |
| `make help` | List targets |

## Privacy

Statements stay on this machine. The web UI does not write statements to disk;
reports live in one Go process’s RAM until restart. There is no bank connect,
no cloud sync, and no stored statement history.

## License

[MIT](LICENSE)
