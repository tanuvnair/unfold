export type Bank = {
  key: string;
  bank_name: string;
  has_parser: boolean;
};

export type TransactionRow = Record<string, string>;

export type Report = {
  bank_name: string;
  transaction_count: number;
  transactions: TransactionRow[];
};

export type ApiError = {
  error: string;
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
