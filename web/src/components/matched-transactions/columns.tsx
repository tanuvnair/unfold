import {createColumnHelper} from '@tanstack/react-table';

import {
  confidenceBadgeVariant,
  confidenceLabel,
  sourceBadgeVariant,
  sourceLabel,
  sourceTitle,
} from '@/components/matched-transactions/badge-variants';
import {Badge} from '@/components/ui/badge';
import {type DataTableFeatures} from '@/components/ui/data-table-features';
import {type MatchedTransaction} from '@/lib/api';

const columnHelper = createColumnHelper<DataTableFeatures, MatchedTransaction>();

export const columns = columnHelper.columns([
  columnHelper.accessor('transactionDate', {
    header: 'Transaction Date',
  }),
  columnHelper.accessor('description', {
    header: 'Description',
    meta: {cellClassName: 'max-w-md whitespace-normal'},
    cell: ({getValue}) => (
      <span className="break-all">{getValue() || '—'}</span>
    ),
  }),
  columnHelper.accessor('amount', {
    header: () => <div className="text-right">Amount</div>,
    meta: {headerClassName: 'text-right', cellClassName: 'text-right'},
    cell: ({getValue}) => (
      <span className="font-medium tabular-nums">{getValue() || '—'}</span>
    ),
  }),
  columnHelper.accessor('confidence', {
    header: 'Confidence',
    cell: ({row, getValue}) => {
      const confidence = getValue();
      const source = row.original.source;
      if (!confidence && !source) {
        return '—';
      }
      return (
        <div className="flex flex-wrap items-center gap-1">
          {confidence ? (
            <Badge variant={confidenceBadgeVariant(confidence)}>
              {confidenceLabel(confidence)}
            </Badge>
          ) : null}
          {source ? (
            <Badge
              variant={sourceBadgeVariant(source)}
              title={sourceTitle(source)}
            >
              {sourceLabel(source)}
            </Badge>
          ) : null}
        </div>
      );
    },
  }),
  columnHelper.accessor('type', {
    header: 'Dr / Cr',
    cell: ({getValue}) => {
      const type = getValue();
      if (!type) {
        return '—';
      }
      return (
        <Badge variant={type === 'CR' ? 'success' : 'destructive'}>{type}</Badge>
      );
    },
  }),
]);
