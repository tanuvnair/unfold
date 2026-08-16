import {type VariantProps} from 'class-variance-authority';

import {type badgeVariants} from '@/components/ui/badge';

type BadgeVariant = NonNullable<VariantProps<typeof badgeVariants>['variant']>;

/**
 * Status chips on muted/hover surfaces. Avoid `secondary` — it shares the same
 * token as `muted`, so the fill disappears on ghost row hover. Prefer
 * `success` (tint) or `outline` (paper + border, matches outline Button).
 * Never `default` (System Blue is commit-only per DESIGN.md).
 */
export function confidenceBadgeVariant(value: string): BadgeVariant {
  switch (value) {
    case 'high':
      return 'success';
    case 'medium':
      return 'outline';
    case 'low':
      return 'outline';
    default:
      return 'outline';
  }
}

export function sourceBadgeVariant(value: string): BadgeVariant {
  switch (value) {
    case 'both':
      return 'outline';
    case 'recurrence':
      return 'outline';
    case 'keyword':
      return 'outline';
    default:
      return 'outline';
  }
}

export function confidenceLabel(value: string): string {
  switch (value) {
    case 'high':
      return 'High';
    case 'medium':
      return 'Medium';
    case 'low':
      return 'Low';
    default:
      return value || '—';
  }
}

export function sourceLabel(value: string): string {
  switch (value) {
    case 'keyword':
      return 'Keyword';
    case 'recurrence':
      return 'Pattern';
    case 'both':
      return 'Both';
    default:
      return value || '—';
  }
}

export function sourceTitle(value: string): string | undefined {
  if (value === 'recurrence') {
    return 'Flagged by repeating amount and monthly cadence, not keyword text.';
  }
  if (value === 'both') {
    return 'Matched by keyword language and by repeating amount/cadence.';
  }
  if (value === 'keyword') {
    return 'Matched by autopay/mandate keyword in the description.';
  }
  return undefined;
}
