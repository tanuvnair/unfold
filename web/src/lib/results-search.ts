import type {ConfidenceFilter, SourceFilter} from '@/lib/api';

export const RESULTS_PAGE_SIZES = [10, 20, 25, 50] as const;

export type ResultsView = 'grouped' | 'transactions';

export type ResultsSearch = {
  q: string;
  confidence: ConfidenceFilter;
  source: SourceFilter;
  view: ResultsView;
  page: number;
  pageSize: number;
};

export const resultsSearchDefaults: ResultsSearch = {
  q: '',
  confidence: 'all',
  source: 'all',
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

export function parseResultsSearch(
  search: Record<string, unknown>,
): ResultsSearch {
  return {
    q: typeof search.q === 'string' ? search.q : resultsSearchDefaults.q,
    confidence: parseConfidence(search.confidence),
    source: parseSource(search.source),
    view: parseView(search.view),
    page: parsePositiveInt(search.page, resultsSearchDefaults.page),
    pageSize: parsePageSize(search.pageSize),
  };
}
