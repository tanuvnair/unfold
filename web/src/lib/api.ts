export type Bank = {
  key: string;
  bank_name: string;
  has_parser: boolean;
};

export type Report = {
  id: string;
  bank_name: string;
  transaction_count: number;
};

export type TransactionTypeFilter = 'all' | 'DR' | 'CR';

export type ConfidenceFilter = 'all' | 'high' | 'medium' | 'low';

export type SourceFilter = 'all' | 'keyword' | 'recurrence' | 'both';

export type MatchedTransaction = {
  id: string;
  transactionDate: string;
  description: string;
  amount: string;
  type: string;
  confidence: string;
  source: string;
};

export type ApiError = {
  error: string;
};

type TransactionPageResponse = {
  rows: Array<{
    id: string;
    transaction_date: string;
    description: string;
    amount: string;
    type: string;
    confidence: string;
    source: string;
  }>;
  row_count: number;
  page: number;
  page_size: number;
};

/** Absolute API origin for split deployments, or "" for same-origin `/api`. */
function apiBase(): string {
  const raw = import.meta.env.UNFOLD_WEB_API_BASE as string | undefined;
  if (!raw) {
    return '';
  }
  return raw.replace(/\/$/, '');
}

function apiURL(path: string): string {
  return `${apiBase()}${path}`;
}

async function readError(res: Response): Promise<string> {
  try {
    const body = (await res.json()) as ApiError;
    if (body.error) {
      return body.error;
    }
  } catch {
    // fall through
  }
  return res.statusText || `Request failed (${res.status})`;
}

export async function fetchBanks(): Promise<Bank[]> {
  const res = await fetch(apiURL('/api/banks'));
  if (!res.ok) {
    throw new Error(await readError(res));
  }
  const data = (await res.json()) as {banks: Bank[]};
  return data.banks ?? [];
}

export async function analyzeStatement(
  file: File,
  bank: string,
): Promise<Report> {
  const form = new FormData();
  form.append('file', file);
  if (bank) {
    form.append('bank', bank);
  }
  const res = await fetch(apiURL('/api/analyze'), {
    method: 'POST',
    body: form,
  });
  if (!res.ok) {
    throw new Error(await readError(res));
  }
  return (await res.json()) as Report;
}

export async function fetchReportTransactions(
  id: string,
  params: {
    q: string;
    type: TransactionTypeFilter;
    confidence: ConfidenceFilter;
    source: SourceFilter;
    page: number;
    pageSize: number;
    payee?: string;
  },
): Promise<{rows: MatchedTransaction[]; rowCount: number}> {
  const search = new URLSearchParams({
    page: String(params.page),
    page_size: String(params.pageSize),
  });
  if (params.q.trim()) {
    search.set('q', params.q.trim());
  }
  if (params.type !== 'all') {
    search.set('type', params.type);
  }
  if (params.confidence !== 'all') {
    search.set('confidence', params.confidence);
  }
  if (params.source !== 'all') {
    search.set('source', params.source);
  }
  if (params.payee?.trim()) {
    search.set('payee', params.payee.trim());
  }

  const res = await fetch(
    apiURL(`/api/reports/${encodeURIComponent(id)}/transactions?${search}`),
  );
  if (!res.ok) {
    throw new Error(await readError(res));
  }
  const data = (await res.json()) as TransactionPageResponse;
  return {
    rows: (data.rows ?? []).map((row) => ({
      id: row.id,
      transactionDate: row.transaction_date,
      description: row.description,
      amount: row.amount,
      type: row.type,
      confidence: row.confidence,
      source: row.source,
    })),
    rowCount: data.row_count,
  };
}

export type SummaryGroup = {
  payee: string;
  occurrenceCount: number;
  totalAmount: number;
  avgAmount: number;
  latestAmount: number;
  firstSeen: string;
  lastSeen: string;
  confidence: string;
  source: string;
  monthlyEstimate: number;
  isMonthlyCadence: boolean;
};

export type ReportSummary = {
  groups: SummaryGroup[];
  estimatedMonthlyTotal: number;
  groupCount: number;
};

type SummaryResponse = {
  groups: Array<{
    payee: string;
    occurrence_count: number;
    total_amount: number;
    avg_amount: number;
    latest_amount: number;
    first_seen: string;
    last_seen: string;
    confidence: string;
    source: string;
    monthly_estimate: number;
    is_monthly_cadence: boolean;
  }>;
  estimated_monthly_total: number;
  group_count: number;
};

export async function fetchReportSummary(
  id: string,
  params: {
    q: string;
    type: TransactionTypeFilter;
    confidence: ConfidenceFilter;
    source: SourceFilter;
  },
): Promise<ReportSummary> {
  const search = new URLSearchParams();
  if (params.q.trim()) {
    search.set('q', params.q.trim());
  }
  if (params.type !== 'all') {
    search.set('type', params.type);
  }
  if (params.confidence !== 'all') {
    search.set('confidence', params.confidence);
  }
  if (params.source !== 'all') {
    search.set('source', params.source);
  }
  const qs = search.toString();
  const res = await fetch(
    apiURL(
      `/api/reports/${encodeURIComponent(id)}/summary${qs ? `?${qs}` : ''}`,
    ),
  );
  if (!res.ok) {
    throw new Error(await readError(res));
  }
  const data = (await res.json()) as SummaryResponse;
  return {
    groups: (data.groups ?? []).map((g) => ({
      payee: g.payee,
      occurrenceCount: g.occurrence_count,
      totalAmount: g.total_amount,
      avgAmount: g.avg_amount,
      latestAmount: g.latest_amount,
      firstSeen: g.first_seen,
      lastSeen: g.last_seen,
      confidence: g.confidence,
      source: g.source,
      monthlyEstimate: g.monthly_estimate,
      isMonthlyCadence: g.is_monthly_cadence,
    })),
    estimatedMonthlyTotal: data.estimated_monthly_total,
    groupCount: data.group_count,
  };
}
