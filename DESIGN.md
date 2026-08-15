---
name: unfold
description: Find hidden subscriptions in seconds.
colors:
  pine: "oklch(0.508 0.118 165.612)"
  pine-foreground: "oklch(0.979 0.021 166.113)"
  pine-deep: "oklch(0.432 0.095 166.913)"
  paper: "oklch(1 0 0)"
  ink: "oklch(0.145 0 0)"
  muted: "oklch(0.97 0 0)"
  muted-ink: "oklch(0.556 0 0)"
  cool-wash: "oklch(0.967 0.001 286.375)"
  cool-wash-ink: "oklch(0.21 0.006 285.885)"
  rule: "oklch(0.922 0 0)"
  field: "oklch(0.922 0 0)"
  ring: "oklch(0.708 0 0)"
  destructive: "oklch(0.577 0.245 27.325)"
typography:
  headline:
    fontFamily: "Inter Variable, sans-serif"
    fontSize: "1.875rem"
    fontWeight: 600
    lineHeight: 1.25
    letterSpacing: "-0.025em"
  title:
    fontFamily: "Inter Variable, sans-serif"
    fontSize: "1.125rem"
    fontWeight: 500
    lineHeight: 1.4
    letterSpacing: "-0.025em"
  body:
    fontFamily: "Inter Variable, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 400
    lineHeight: 1.5
    letterSpacing: "normal"
  label:
    fontFamily: "Inter Variable, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 500
    lineHeight: 1.4
    letterSpacing: "normal"
rounded:
  sm: "0.375rem"
  md: "0.5rem"
  lg: "0.625rem"
  xl: "0.875rem"
  2xl: "1.125rem"
  3xl: "1.375rem"
  4xl: "1.625rem"
  card: "24px"
spacing:
  sm: "0.5rem"
  md: "1.25rem"
  lg: "1.5rem"
  xl: "2rem"
  2xl: "2.5rem"
components:
  button-primary:
    backgroundColor: "{colors.pine}"
    textColor: "{colors.pine-foreground}"
    typography: "{typography.label}"
    rounded: "{rounded.2xl}"
    padding: "0 1rem"
    height: "2.25rem"
  button-primary-hover:
    backgroundColor: "color-mix(in oklch, oklch(0.508 0.118 165.612) 80%, transparent)"
    textColor: "{colors.pine-foreground}"
    rounded: "{rounded.2xl}"
    padding: "0 1rem"
    height: "2.25rem"
  button-outline:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    typography: "{typography.label}"
    rounded: "{rounded.2xl}"
    padding: "0 0.75rem"
    height: "2rem"
  button-outline-hover:
    backgroundColor: "{colors.muted}"
    textColor: "{colors.ink}"
    rounded: "{rounded.2xl}"
    padding: "0 0.75rem"
    height: "2rem"
  badge-secondary:
    backgroundColor: "{colors.cool-wash}"
    textColor: "{colors.cool-wash-ink}"
    typography: "{typography.label}"
    rounded: "{rounded.2xl}"
    padding: "0.125rem 0.5rem"
    height: "1.25rem"
  card:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.card}"
    padding: "1.25rem 0"
  input:
    backgroundColor: "color-mix(in oklch, oklch(0.922 0 0) 50%, transparent)"
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.2xl}"
    padding: "0.25rem 0.625rem"
    height: "2rem"
  select-trigger:
    backgroundColor: "color-mix(in oklch, oklch(0.922 0 0) 50%, transparent)"
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.2xl}"
    padding: "0.5rem 0.75rem"
    height: "2rem"
  nav-wordmark:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    typography: "{typography.title}"
    height: "3.5rem"
---

# Design System: unfold

## Overview

**Creative North Star: "The Sealed Ledger"**

unfold's interface treats privacy as a material, not a slogan. The screen is a sealed paper field: white, hairline-edged, and chromatically silent until the user must act. Then Pine fires once, on the control that commits the statement. That tap is tactile and confident; everything around it stays quiet.

The system is Inter Variable at every rank, shadcn base-rhea geometry with pill-rounded controls and larger sealed cards, and depth that is tonal rather than theatrical. Cards sit on the page through a 1px ink-whisper ring and a faint shadow. Nothing floats. Nothing ornaments. The job is a short local audit, so the layout stays a single centered column, not a dashboard.

**Key Characteristics:**

- Inter Variable for heading and body; no second family
- Pine as the sole chromatic accent, used on committing actions
- Soft-seal radii: 18px controls, 24px cards
- Hairline ring plus faint shadow; no floating chrome
- Light paper default; dark tokens exist as a theme, not the primary chrome
- Lowercase wordmark "unfold"

## Colors

A near-achromatic paper ledger with one mid teal-green voice. Accent and primary are the same token.

### Primary

- **Pine** (`colors.pine`): the only chromatic fill. Primary buttons, and the same value as `accent` for selected menu rows. Its job is to mark the commit, not to paint the chrome.
- **Pine Ink** (`colors.pine-foreground`): text and icons on Pine fills. Near-white with a faint teal cast.
- **Pine Deep** (`colors.pine-deep`): dark-theme primary. Do not use on the light paper default.

### Neutral

- **Paper** (`colors.paper`): page, card, and popover ground in light mode.
- **Ink** (`colors.ink`): primary text and the wordmark.
- **Muted** (`colors.muted`): hover washes, table row hover, empty-state icon wells.
- **Muted Ink** (`colors.muted-ink`): supporting copy, field descriptions, placeholders, header subtitles.
- **Cool Wash** (`colors.cool-wash`): secondary badge fill; a near-gray with a whisper of violet. Count, not alarm.
- **Cool Wash Ink** (`colors.cool-wash-ink`): text on Cool Wash.
- **Rule** (`colors.rule`): default borders, header rule, table row rules, input token before the 50% wash.
- **Field** (`colors.field`): same value as Rule; fields paint it at 50% over Paper so the control reads as a wash, not a box.
- **Ring** (`colors.ring`): focus border and the 3px focus halo at 30% opacity.

### Destructive

- **Signal Red** (`colors.destructive`): error alerts, invalid fields, destructive button text. Alerts keep a Paper ground and color the type; they do not flood the card.

**The Pine Tap Rule.** Pine appears on the action that commits, and almost nowhere else. Filling a header, a page wash, or decorative shapes with Pine breaks the seal.

**The One Accent Rule.** Primary and accent are the same Pine token. Do not introduce a second brand hue.

## Typography

**Display Font:** Inter Variable (with sans-serif)
**Body Font:** Inter Variable (with sans-serif)

**Character:** One family, two weights. Semibold tracking-tight for page titles; medium for card and empty titles; regular for reading; medium for labels and buttons. No serif, no mono, no display face.

### Hierarchy

- **Headline** (semibold / 1.875rem / tight tracking): page titles ("Find hidden subscriptions", "Autopay matches"). The largest type the product uses; there is no poster display size.
- **Title** (medium / 1.125rem / tight tracking): empty-state titles. Card titles step down to 1rem medium on the same family.
- **Body** (regular / 0.875rem / 1.5): descriptions, table cells, supporting copy. Inputs use 1rem on small viewports, then 0.875rem from the `md` breakpoint so iOS does not zoom.
- **Label** (medium / 0.875rem): field labels, buttons, badges (badges also use 0.75rem). Sentence case. No uppercase tracking.

**The One Voice Rule.** Heading and body share Inter Variable. `--font-heading` is an alias of `--font-sans`, not a second face.

## Layout

A single centered column. The form page is `max-w-3xl`; the results page is `max-w-5xl`. Horizontal inset is 1.5rem; vertical page padding is 2.5rem; stacks between title and card use 2rem. The header is 3.5rem tall, the same 5xl cap, with a bottom Rule.

Cards use an internal 1.25rem spacing token (`--card-spacing`). Field groups stack at 1.5rem. This is a short audit, not a 12-column dashboard: no sidebar, no app rail, no widget grid.

**The Short Audit Rule.** If a screen needs a second column to feel complete, the layout has drifted.

## Elevation & Depth

Surfaces are paper, seated with a hairline, not lifted into floating chrome. Cards use a faint `shadow-sm` plus `ring-1 ring-foreground/5` (10% in dark). Select popovers add `shadow-lg` because they leave the page plane; that is the only structural lift. Inputs have no shadow and no visible border at rest: they are a 50% Field wash.

**The Hairline Seat Rule.** Resting surfaces sit on Paper via a 1px ink-whisper ring and a faint shadow. Do not add large drop shadows, glass, or gradient scrims.

## Shapes

Controls are soft seals, not sharp tools and not stadium pills. Buttons, inputs, badges, alerts, and select triggers use 18px (`rounded.2xl`). Select items ease to 14px. Empty states use 22px. Cards use a 24px cap (`min(radius-4xl, 24px)`), the largest corner in the system.

Borders on controls are transparent at rest. Focus draws a Ring border plus a 3px halo at 30% opacity. Active buttons press `1px` down. Empty states, when they need an outline, use a dashed Rule, not a heavy stroke.

**The Soft Seal Rule.** Interactive chrome is 18px. Cards are 24px. Do not mix sharp 4px rectangles into this vocabulary, and do not fully pill large surfaces.

## Components

### Buttons

Tactile and confident: the primary fill is the one firm tap on an otherwise paper screen.

- **Shape:** soft seal (18px). Medium Inter, no uppercase.
- **Primary:** Pine fill, Pine Ink type. Default height 2rem / 12px horizontal padding; the analyze commit uses large (2.25rem / 16px). Hover: Pine at 80%. Focus: Ring border and 3px halo. Active: 1px press.
- **Outline:** Paper fill, Rule-transparent border that reads as `border-border`, Ink type. Hover: Muted wash. Used for "Analyze another".
- **Ghost / Secondary / Destructive / Link:** available in the kit; the shipped screens use default and outline only. Destructive is a Signal Red tint, not a solid red brick.
- **Disabled:** 50% opacity, no pointer.

### Chips

- **Style:** Badge, 20px tall, 18px corners, 12px medium type. Secondary (Cool Wash) is the count chip on results. Default (Pine) exists but is not the results treatment.
- **State:** static count, not a filter chip.

### Cards / Containers

The sealed unit of work: one card per job (analyze form, matched table).

- **Corner Style:** 24px seal
- **Background:** Paper
- **Shadow Strategy:** Hairline Seat (faint shadow + 5% ink ring)
- **Border:** none; the ring is the edge
- **Internal Padding:** 1.25rem, including a bordered footer with the primary action right-aligned

### Inputs / Fields

- **Style:** 18px seal, 2rem tall, transparent border, Field at 50% over Paper. File upload sits in an Input Group with a muted icon addon; the inner control has no own border.
- **Focus:** Ring border and 3px halo on the group, 200ms color/shadow.
- **Error:** Signal Red border and halo. Field error copy is Signal Red at body size.
- **Select:** same wash and height as input; 18px trigger; 18px popover with `shadow-lg` and the same hairline ring. Items are 14px corners; selected row uses Pine as accent fill.

### Navigation

- **Style:** 3.5rem Paper bar, bottom Rule, no actions except the wordmark.
- **Wordmark:** lowercase "unfold", title-sized Inter semibold, tight tracking, Ink. It is a home link, not a logo lockup. No mark, no tagline in the bar.

### Empty state

Dashed 22px seal, 3rem padding, centered. Icon well is 2.5rem, 12px corners, Muted fill. Title uses the Title role; description is Muted Ink.

### Table

Body-sized, 2.5rem header cells, 0.5rem cell padding, Rule row dividers. Row hover is Muted at 50%. Description cells may wrap; other columns stay nowrap.

### Alerts

18px seal, 16px/12px padding, Paper ground, 1px Rule. Destructive colors the title and description, not the background.

## Do's and Don'ts

### Do:

- **Do** keep the wordmark lowercase "unfold" in Inter semibold.
- **Do** spend Pine on the committing control; let Paper and Ink do the rest.
- **Do** seat cards with a hairline ring and a faint shadow.
- **Do** use 18px seals on controls and 24px on cards.
- **Do** keep new screens in the centered 3xl/5xl column.

### Don't:

- **Don't** introduce a second typeface or a display/serif pairing.
- **Don't** paint headers, page backgrounds, or illustration with Pine.
- **Don't** add accounts, avatars, cloud glyphs, or social proof; those break the sealed, non-invasive product.
- **Don't** use heavy drop shadows, glass, or gradient scrims.
- **Don't** treat the unused isometric `hero.png` as brand art; it is not in this system.
