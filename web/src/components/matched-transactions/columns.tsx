import {createColumnHelper} from '@tanstack/react-table';

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
