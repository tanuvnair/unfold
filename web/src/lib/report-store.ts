import type {Report} from '@/lib/api';

let latestReport: Report | null = null;

export function setLatestReport(report: Report | null) {
  latestReport = report;
}

export function getLatestReport(): Report | null {
  return latestReport;
}
