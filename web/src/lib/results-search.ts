import type {ConfidenceFilter, SourceFilter} from '@/lib/api';

export const RESULTS_PAGE_SIZES = [10, 20, 25, 50] as const;

export type ResultsView = 'grouped' | 'transactions';

export type ResultsSearch = {
  q: string;
  confidence: ConfidenceFilter;
  source: SourceFilter;
  from: string;
  to: string;
  view: ResultsView;
  page: number;
  pageSize: number;
};

export const resultsSearchDefaults: ResultsSearch = {
  q: '',
  confidence: 'all',
  source: 'all',
  from: '',
  to: '',
  view: 'grouped',
  page: 1,
  pageSize: 10,
};

function parsePositiveInt(value: unknown, fallback: number): number {
  const n = typeof value === 'number' ? value : Number(value);
  if (!Number.isFinite(n) || n < 1) {
    return fallback;
  }
  return Math.floor(n);
}

function parsePageSize(value: unknown): number {
  const n = parsePositiveInt(value, resultsSearchDefaults.pageSize);
  return (RESULTS_PAGE_SIZES as readonly number[]).includes(n)
    ? n
    : resultsSearchDefaults.pageSize;
}

function parseConfidence(value: unknown): ConfidenceFilter {
  if (
    value === 'high' ||
    value === 'medium' ||
    value === 'low' ||
    value === 'all'
  ) {
    return value;
  }
  return resultsSearchDefaults.confidence;
}

function parseSource(value: unknown): SourceFilter {
  if (
    value === 'keyword' ||
    value === 'recurrence' ||
    value === 'both' ||
    value === 'all'
  ) {
    return value;
  }
  return resultsSearchDefaults.source;
}

function parseView(value: unknown): ResultsView {
  if (value === 'grouped' || value === 'transactions') {
    return value;
  }
  return resultsSearchDefaults.view;
}

/** YYYY-MM-DD for date-range filter / API query params. */
function parseISODate(value: unknown): string {
  if (typeof value !== 'string') {
    return '';
  }
  const trimmed = value.trim();
  if (!/^\d{4}-\d{2}-\d{2}$/.test(trimmed)) {
    return '';
  }
  const t = Date.parse(`${trimmed}T00:00:00Z`);
  if (Number.isNaN(t)) {
    return '';
  }
  return trimmed;
}

export function parseResultsSearch(
  search: Record<string, unknown>,
): ResultsSearch {
  return {
    q: typeof search.q === 'string' ? search.q : resultsSearchDefaults.q,
    confidence: parseConfidence(search.confidence),
    source: parseSource(search.source),
    from: parseISODate(search.from),
    to: parseISODate(search.to),
    view: parseView(search.view),
    page: parsePositiveInt(search.page, resultsSearchDefaults.page),
    pageSize: parsePageSize(search.pageSize),
  };
}
