# unfold web UI

Vite + React app for the local unfold audit: pick a bank profile, upload a
statement CSV, read the autopay/mandate hits.

Run it from the **repo root** (`make web`), not from this directory, so env
and the API proxy match [`.env.example`](../.env.example). See the root
[`README.md`](../README.md).

## Stack

- Vite, React 19, TypeScript
- TanStack Router, TanStack Query, TanStack Table v9
- shadcn/ui (`base-rhea`, Base UI, Hugeicons)

## Screens

| Route | Role |
|-------|------|
| `/` | Bank select + CSV dropzone. Analyze posts to `/api/analyze`. |
| `/results` | Matched transactions DataTable. Requires a report id in session memory. |

Both screens use the shared `Page` layout and align to the **top** of the
viewport (`flex flex-col gap-8`). Home is `max-w-3xl`; results is `max-w-5xl`.

## Results table

Search, DR/CR, and pagination are **server-side**. The UI does not keep the
full row list. After analyze, it stores `{ id, bank_name, transaction_count }`
and loads pages from `GET /api/reports/{id}/transactions`.

Table state is in the URL (TanStack Router search params). Defaults are
omitted from the address bar.

| Param | Default | Notes |
|-------|---------|-------|
| `q` | *(empty)* | Description search (debounced 300ms) |
| `type` | `all` | `all`, `DR`, or `CR` |
| `page` | `1` | **1-based** (the API `page` param is 0-based) |
| `pageSize` | `10` | `10`, `20`, `25`, or `50` |

Pagination uses the numbered shadcn Pagination control (page links +
ellipsis), not Previous/Next only.

Key files:

- `src/components/ui/data-table.tsx` — generic DataTable (`manualPagination`)
- `src/components/matched-transactions/` — columns + results table wiring
- `src/lib/results-search.ts` — URL search parse/defaults
- `src/lib/report-store.ts` — session metadata only
- `src/lib/api.ts` — `/api/banks`, `/api/analyze`, `/api/reports/…`

## Scripts

From this directory, or via `make web` / `make web-build` at the repo root:

```bash
npm run dev      # Vite; proxies /api when UNFOLD_WEB_API_BASE is empty
npm run build    # tsc -b && vite build
npm run lint     # oxlint
```

Add shadcn components from **this** directory:

```bash
npx shadcn@latest add <component>
```

Do not overwrite existing `src/components/ui` files unless asked. Icons are
Hugeicons (`@hugeicons/core-free-icons`), not Lucide.
