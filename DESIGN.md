---
name: unfold
description: Find hidden subscriptions in seconds.
colors:
  system-blue: "oklch(0.603 0.218 257.424)"
  system-blue-foreground: "oklch(1 0 0)"
  system-blue-dark: "oklch(0.624 0.206 255.486)"
  grouped: "oklch(0.963 0.007 286.274)"
  paper: "oklch(1 0 0)"
  ink: "oklch(0.232 0.004 286.099)"
  muted: "oklch(0.923 0.007 286.267)"
  muted-ink: "oklch(0.648 0.007 286.193)"
  cool-wash: "oklch(0.923 0.007 286.267)"
  cool-wash-ink: "oklch(0.232 0.004 286.099)"
  rule: "oklch(0.827 0.003 286.339)"
  field: "oklch(0.923 0.007 286.267)"
  ring: "oklch(0.603 0.218 257.424)"
  destructive: "oklch(0.654 0.232 28.659)"
  success: "oklch(0.73 0.194 147.444)"
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
    backgroundColor: "{colors.system-blue}"
    textColor: "{colors.system-blue-foreground}"
    typography: "{typography.label}"
    rounded: "{rounded.2xl}"
    padding: "0 1rem"
    height: "2.25rem"
  button-primary-hover:
    backgroundColor: "color-mix(in oklch, oklch(0.603 0.218 257.424) 80%, transparent)"
    textColor: "{colors.system-blue-foreground}"
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
    backgroundColor: "color-mix(in oklch, oklch(0.923 0.007 286.267) 50%, transparent)"
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.2xl}"
    padding: "0.25rem 0.625rem"
    height: "2rem"
  select-trigger:
    backgroundColor: "color-mix(in oklch, oklch(0.923 0.007 286.267) 50%, transparent)"
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.2xl}"
    padding: "0.5rem 0.75rem"
    height: "2rem"
  nav-wordmark:
    backgroundColor: "{colors.grouped}"
    textColor: "{colors.ink}"
    typography: "{typography.title}"
    height: "3.5rem"
---

# Design System: unfold

## Overview

**Creative North Star: "The Sealed Ledger"**

unfold's interface treats privacy as a material, not a slogan. The screen is an iOS grouped field: system-gray canvas, white cards, hairline-edged, and chromatically silent until the user must act. Then System Blue fires once, on the control that commits the statement. That tap is tactile and confident; everything around it stays quiet.

The system is Inter Variable at every rank, shadcn base-rhea geometry with pill-rounded controls and larger sealed cards, and depth that is tonal rather than theatrical. Cards sit on the grouped canvas through a 1px ink-whisper ring and a faint shadow. Nothing floats. Nothing ornaments. The job is a short local audit, so the layout stays a single centered column, not a dashboard.

**Key Characteristics:**

- Inter Variable for heading and body; no second family
- System Blue as the sole chromatic accent, used on committing actions
- Soft-seal radii: 18px controls, 24px cards
- Hairline ring plus faint shadow; no floating chrome
- Light grouped-gray default; dark tokens exist as a theme, not the primary chrome
- Lowercase wordmark "unfold"

## Colors

An Apple HIG grouped canvas with one System Blue voice. Accent and primary are the same token.

### Primary

- **System Blue** (`colors.system-blue`): the only chromatic fill (`#007AFF`). Primary buttons, the mark tile, and the focus ring. Its job is to mark the commit, not to paint the chrome.
- **System Blue Ink** (`colors.system-blue-foreground`): text and icons on System Blue fills. White.
- **System Blue Dark** (`colors.system-blue-dark`): dark-theme primary (`#0A84FF`). Do not use on the light grouped default.

### Neutral

- **Grouped** (`colors.grouped`): page ground in light mode (`#F2F2F7`).
- **Paper** (`colors.paper`): card and popover ground. White.
- **Ink** (`colors.ink`): primary text and the wordmark (`#1d1d1f`).
- **Muted** (`colors.muted`): hover washes, table row hover, empty-state icon wells (`#E5E5EA`).
- **Muted Ink** (`colors.muted-ink`): supporting copy, field descriptions, placeholders, header subtitles (`#8E8E93`).
- **Cool Wash** (`colors.cool-wash`): secondary badge fill; system gray 5. Count, not alarm.
- **Cool Wash Ink** (`colors.cool-wash-ink`): text on Cool Wash.
- **Rule** (`colors.rule`): default borders, header rule, table row rules, opaque separator (`#C6C6C8`).
- **Field** (`colors.field`): same value as Muted; fields paint it at 50% over Paper so the control reads as a wash, not a box.
- **Ring** (`colors.ring`): focus border and the 3px focus halo at 30% opacity. Same value as System Blue.

### Destructive

- **Signal Red** (`colors.destructive`): error alerts, invalid fields, destructive button text (`#FF3B30`). Alerts keep a Paper ground and color the type; they do not flood the card.

**The Blue Tap Rule.** System Blue appears on the action that commits, the mark, and the focus ring, and almost nowhere else. Filling a header, a page wash, or decorative shapes with System Blue breaks the seal.

**The One Accent Rule.** Primary, ring, and the mark tile are the same System Blue token. Do not introduce a second brand hue on chrome.

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
- **Primary:** System Blue fill, white type. Default height 2rem / 12px horizontal padding; the analyze commit uses large (2.25rem / 16px). Hover: System Blue at 80%. Focus: Ring border and 3px halo. Active: 1px press.
- **Outline:** Paper fill, Rule-transparent border that reads as `border-border`, Ink type. Hover: Muted wash. Used for "Analyze another".
- **Ghost / Secondary / Destructive / Link:** available in the kit; the shipped screens use default and outline only. Destructive is a Signal Red tint, not a solid red brick.
- **Disabled:** 50% opacity, no pointer.

### Chips

- **Style:** Badge, 20px tall, 18px corners, 12px medium type. Secondary (Cool Wash) is the count chip on results. Default (System Blue) exists but is not the results treatment.
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
- **Select:** same wash and height as input; 18px trigger; 18px popover with `shadow-lg` and the same hairline ring. Items are 14px corners; selected row uses a muted wash, not a System Blue fill.

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
- **Do** spend System Blue on the committing control; let Grouped, Paper, and Ink do the rest.
- **Do** seat cards with a hairline ring and a faint shadow.
- **Do** use 18px seals on controls and 24px on cards.
- **Do** keep new screens in the centered 3xl/5xl column.

### Don't:

- **Don't** introduce a second typeface or a display/serif pairing.
- **Don't** paint headers, page backgrounds, or illustration with System Blue.
- **Don't** add accounts, avatars, cloud glyphs, or social proof; those break the sealed, non-invasive product.
- **Don't** use heavy drop shadows, glass, or gradient scrims.
- **Don't** treat the unused isometric `hero.png` as brand art; it is not in this system.
