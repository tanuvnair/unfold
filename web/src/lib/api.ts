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

export type MatchedTransaction = {
  id: string;
  transactionDate: string;
  description: string;
  amount: string;
  type: string;
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
    page: number;
    pageSize: number;
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
    })),
    rowCount: data.row_count,
  };
}
