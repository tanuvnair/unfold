/** Matches the API multipart cap in core/cmd/api. */
export const MAX_STATEMENT_BYTES = 10 << 20;

export function validateStatementFile(file: File): string | null {
  const name = file.name.toLowerCase();
  if (!name.endsWith('.csv')) {
    return 'File must be a .csv statement export.';
  }
  if (file.size > MAX_STATEMENT_BYTES) {
    return 'Upload too large (max 10 MB).';
  }
  return null;
}
