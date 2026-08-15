import type {TransactionTypeFilter} from '@/lib/api';

export const RESULTS_PAGE_SIZES = [10, 20, 25, 50] as const;

export type ResultsSearch = {
  q: string;
  type: TransactionTypeFilter;
  page: number;
  pageSize: number;
};

export const resultsSearchDefaults: ResultsSearch = {
  q: '',
  type: 'all',
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

function parseType(value: unknown): TransactionTypeFilter {
  if (value === 'DR' || value === 'CR' || value === 'all') {
    return value;
  }
  return resultsSearchDefaults.type;
}

export function parseResultsSearch(
  search: Record<string, unknown>,
): ResultsSearch {
  return {
    q: typeof search.q === 'string' ? search.q : resultsSearchDefaults.q,
    type: parseType(search.type),
    page: parsePositiveInt(search.page, resultsSearchDefaults.page),
    pageSize: parsePageSize(search.pageSize),
  };
}
