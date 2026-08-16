# AGENTS.md

Guidance for coding agents working on unfold. Product truth lives in
[`PRODUCT.md`](PRODUCT.md). Visual system lives in [`DESIGN.md`](DESIGN.md).
How to run the repo lives in [`README.md`](README.md).

## Product constraints (do not drift)

- Local, private, India-specific bank-CSV matching. No bank login, no cloud
  sync, no accounts, no stored statement history.
- One parse → keyword-match → report pipeline shared by CLI, HTTP API, and UI.
- CLI still writes `autopay_report.json` next to the statement. The web UI does
  not persist statements to disk.
- Default API listen address is loopback. Do not imply hosted/multi-tenant use.
- Wordmark is lowercase **unfold**. System Blue is the only chromatic accent
  (commit actions). Do not invent testimonials or bank partnerships.

## Stack

- Go API/CLI (`core/cmd/api`, `core/cmd/cli`, `core/internal/`).
- Web: Vite + React + TanStack Router + TanStack Query + shadcn/ui
  (style `base-rhea`, Base UI primitives, Hugeicons).
- Table: TanStack Table v9 (`useTable`, `tableFeatures()`). Do not use the v8
  `useReactTable` API.

## Decisions on this branch

### Keep server-side table operations

Search, DR/CR filter, and pagination run on the **API**, not in the browser.

Typical reports are small enough for client-side filtering. That was considered
and rejected: keep the server-side table architecture. Do not silently move
filtering back to the client.

What this is **not**: a database, a session store on disk, or a multi-user
backend. Reports live in **one Go process’s RAM** (`core/internal/reportquery.Store`)
keyed by a random id. Restarting the API drops them. Refreshing the tab drops
the client-side id. The user must analyze again.

Do not add SQL, files, or cookies to “fix” that unless the product explicitly
changes.

### Analyze returns metadata; rows are queried

- `POST /api/analyze` stores the full report in RAM and returns
  `{ id, bank_name, transaction_count }` — not the transaction list.
- `GET /api/reports/{id}/transactions` applies `q`, `type`, `page`, `page_size`
  and returns one page plus `row_count`.
- The web session store (`web/src/lib/report-store.ts`) holds only that
  metadata. The DataTable fetches pages with TanStack Query.
- CLI output shape is unchanged. Do not put `id` into `autopay_report.json`.

### URL owns table UI state

`/results` search params (TanStack Router `validateSearch` +
`stripSearchParams`):

| Param | Default | Notes |
| --- | --- | --- |
| `q` | `""` | Description substring; input is debounced 300ms before the URL/query updates |
| `type` | `all` | `all` \| `DR` \| `CR` |
| `page` | `1` | **1-based in the URL** |
| `pageSize` | `10` | `10` \| `20` \| `25` \| `50` |

Defaults are omitted from the URL. Changing `q` or `type` resets `page` to 1.

The HTTP API `page` query is **0-based**. Convert in
`fetchReportTransactions` (`page: search.page - 1`). Do not mix the two.

Parser/defaults: `web/src/lib/results-search.ts`.

### DataTable UI

- Reusable table: `web/src/components/ui/data-table.tsx` with
  `manualPagination: true` (no client row models for filter/page).
- Features object: `web/src/components/ui/data-table-features.ts`. Register only
  `rowPaginationFeature`. `getVisibleCells` requires
  `columnVisibilityFeature`; this table uses `row.getAllCells()`.
- Columns: `web/src/components/matched-transactions/columns.tsx`
  (`transactionDate`, `description`, `amount`, `type`).
- Pagination control is the numbered shadcn Pagination (Previous, page links,
  ellipsis, Next) plus a rows-per-page Select — not Previous/Next only.
- DR/CR is a `ToggleGroup` (All / DR / CR). Search uses `InputGroup`.
- Empty filtered results use the `Empty` component; load errors use `Alert`.

Add shadcn pieces with `npx shadcn@latest` from `web/`. Do not `--overwrite`
existing UI files unless the user asks. Icon library is Hugeicons, not Lucide.

### Page layout

Home (`/`) and results (`/results`) both start at the **top** of the viewport.
Use `<Page className="flex flex-col gap-8">`. Do not vertically center the
home form with `min-h-svh justify-center`. Home is `size="md"` (`max-w-3xl`);
results is `size="lg"` (`max-w-5xl`).

## Conventions

- TypeScript/JavaScript: Google style. Go: Google style.
- Commit messages: Conventional Commits.
- No emojis in generated UI copy or docs.
- shadcn: `className` for layout not color; `gap-*` not `space-y-*`; `Field` +
  `FieldGroup` for forms; `SelectItem` inside `SelectGroup`; Base UI `render`
  not Radix `asChild`.
- Do not invent bank parsers, PDF input, or hosted-cloud copy.
