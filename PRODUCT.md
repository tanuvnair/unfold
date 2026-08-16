# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

People in India reviewing their own bank statement to find hidden subscriptions, SIPs, and standing instructions. They already have (or can export) a CSV; they do not want to connect their bank to a third-party app.

## Product Purpose

unfold finds hidden recurring charges in a bank statement CSV in seconds. The user picks a bank profile, supplies the export, and gets the autopay/mandate hits that were buried in the file. Success is a trustworthy list they can act on (cancel, budget, ignore) without giving up bank login or sending the statement away.

## Positioning

Local and private, and India-specific bank-CSV matching: the statement stays on this machine, there is no bank login and no merchant/cloud database, and detection uses the language that actually appears on Indian exports (NACH, mandate, autopay, recurring, standing, monthly investment). Neighboring products (subscription managers that require bank connect, generic spreadsheet filters, bank apps) cannot honestly claim both.

## Operating Context

The job is a short local audit, not an ongoing finance dashboard. Typical flow: export CSV from the bank, choose the matching profile, run unfold (CLI, local API, or web UI), read the matched rows. The CLI writes `autopay_report.json` next to the statement (unless `--dry-run`); `--diff` compares against a previous report. The web UI keeps the latest report in session memory and does not persist statements. CLI, HTTP API, and UI share one parse → keyword-match → report pipeline.

## Capabilities and Constraints

Confirmed:

- Local Go tool with three surfaces: CLI (`unfold-cli`), HTTP API, web UI (Vite + React, embeddable in a single `unfold-api` binary).
- Input is a bank statement CSV. API rejects non-CSV uploads and caps size at 10 MiB.
- Matching is include/exclude keywords on the transaction description (config-driven per bank profile).
- One implemented parser today: Kotak Mahindra Bank (`kotak-mahindra-bank`). Other banks are not promised.
- Default listen address is loopback (`127.0.0.1`).
- Must be secure and not invasive: do not imply accounts, tracking, cloud sync, or sending statements off-device.

Undecided (do not invent in UI copy or flows):

- Additional bank parsers, PDF/screenshot input, live bank connect, user accounts, stored statement history, hosted/multi-tenant deployment.

Terminology: autopay, mandate, NACH, SIP / monthly investment, standing instruction, statement CSV, bank profile, report / matched transactions.

## Brand Commitments

- Product name: unfold (lowercase in the product UI).
- Existing promise: "Find hidden subscriptions in seconds."
- Personality constraint: secure and not invasive. Do not add growth-hacking, social proof fabrication, or copy that treats the statement as something to mine.

## Evidence on Hand

- No testimonials, case studies, press, pricing, or bank-partnership claims. Future work must not fabricate them.
- Real demonstration is the local pipeline itself (CSV in, matched rows out).
- Unused raster: `web/src/assets/hero.png` (isometric layered tiles). Not currently referenced by the UI.

## Product Principles

1. Privacy is the product: statements stay local; never imply login, sync, or off-device processing.
2. Ask only for the CSV, show only the matches: secure and not invasive.
3. Match Indian bank-export language as it appears, not a generic merchant catalog.
4. Seconds to a list, then get out of the way: unfold finds charges; canceling and budgeting happen elsewhere.
5. One truth: CLI, API, and UI report the same matches from the same pipeline.
