import {useEffect, useState} from 'react';
import {keepPreviousData, useQuery} from '@tanstack/react-query';
import {useNavigate, useSearch} from '@tanstack/react-router';
import {HugeiconsIcon} from '@hugeicons/react';
import {Search01Icon, SearchRemoveIcon} from '@hugeicons/core-free-icons';
import type {OnChangeFn, PaginationState} from '@tanstack/react-table';

import {columns} from '@/components/matched-transactions/columns';
import {Alert, AlertDescription, AlertTitle} from '@/components/ui/alert';
import {DataTable} from '@/components/ui/data-table';
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty';
import {Field, FieldLabel} from '@/components/ui/field';
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from '@/components/ui/input-group';
import {ToggleGroup, ToggleGroupItem} from '@/components/ui/toggle-group';
import {
  fetchReportTransactions,
  type Report,
  type TransactionTypeFilter,
} from '@/lib/api';
import {useDebouncedValue} from '@/hooks/use-debounced-value';

const TYPE_ITEMS: Array<{label: string; value: TransactionTypeFilter}> = [
  {label: 'All', value: 'all'},
  {label: 'DR', value: 'DR'},
  {label: 'CR', value: 'CR'},
];

export function MatchedTransactionsTable({report}: {report: Report}) {
  const search = useSearch({from: '/results'});
  const navigate = useNavigate({from: '/results'});
  const [searchInput, setSearchInput] = useState(search.q);
  const debouncedQ = useDebouncedValue(searchInput, 300);

  useEffect(() => {
    setSearchInput(search.q);
  }, [search.q]);

  useEffect(() => {
    if (debouncedQ === search.q) {
      return;
    }
    void navigate({
      search: (prev) => ({...prev, q: debouncedQ, page: 1}),
      replace: true,
    });
  }, [debouncedQ, navigate, search.q]);

  const query = useQuery({
    queryKey: ['report-transactions', report.id, search],
    queryFn: () =>
      fetchReportTransactions(report.id, {
        q: search.q,
        type: search.type,
        page: search.page - 1,
        pageSize: search.pageSize,
      }),
    placeholderData: keepPreviousData,
  });

  const rows = query.data?.rows ?? [];
  const rowCount = query.data?.rowCount ?? 0;

  useEffect(() => {
    if (query.isPending) {
      return;
    }
    const maxPage = Math.max(1, Math.ceil(rowCount / search.pageSize) || 1);
    if (search.page > maxPage) {
      void navigate({
        search: (prev) => ({...prev, page: maxPage}),
        replace: true,
      });
    }
  }, [navigate, query.isPending, rowCount, search.page, search.pageSize]);

  const pagination: PaginationState = {
    pageIndex: search.page - 1,
    pageSize: search.pageSize,
  };

  const onPaginationChange: OnChangeFn<PaginationState> = (updater) => {
    const next = typeof updater === 'function' ? updater(pagination) : updater;
    void navigate({
      search: (prev) => ({
        ...prev,
        page: next.pageSize !== search.pageSize ? 1 : next.pageIndex + 1,
        pageSize: next.pageSize,
      }),
    });
  };

  function setType(type: TransactionTypeFilter) {
    void navigate({
      search: (prev) => ({...prev, type, page: 1}),
    });
  }

  return (
    <DataTable
      columns={columns}
      data={rows}
      rowCount={rowCount}
      pagination={pagination}
      onPaginationChange={onPaginationChange}
      getRowId={(row) => row.id}
      isLoading={query.isPending}
      isFetching={query.isFetching}
      toolbar={
        <div className="flex flex-col gap-3">
          {query.isError ? (
            <Alert variant="destructive">
              <AlertTitle>Could not load transactions</AlertTitle>
              <AlertDescription>
                {(query.error as Error).message}. Try analyzing the statement
                again.
              </AlertDescription>
            </Alert>
          ) : null}
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <Field className="max-w-sm">
              <FieldLabel htmlFor="txn-search" className="sr-only">
                Search descriptions
              </FieldLabel>
              <InputGroup>
                <InputGroupAddon>
                  <HugeiconsIcon icon={Search01Icon} strokeWidth={2} />
                </InputGroupAddon>
                <InputGroupInput
                  id="txn-search"
                  value={searchInput}
                  onChange={(event) => setSearchInput(event.target.value)}
                  placeholder="Search descriptions..."
                />
              </InputGroup>
            </Field>
            <Field orientation="horizontal" className="w-fit">
              <FieldLabel className="sr-only" id="txn-type-label">
                Debit or credit
              </FieldLabel>
              <ToggleGroup
                variant="outline"
                spacing={0}
                value={[search.type]}
                aria-labelledby="txn-type-label"
                onValueChange={(value) => {
                  const next = value[0] as TransactionTypeFilter | undefined;
                  if (next) {
                    setType(next);
                  }
                }}
              >
                {TYPE_ITEMS.map((item) => (
                  <ToggleGroupItem key={item.value} value={item.value}>
                    {item.label}
                  </ToggleGroupItem>
                ))}
              </ToggleGroup>
            </Field>
          </div>
        </div>
      }
      empty={
        query.isError ? (
          <div className="h-24" />
        ) : (
          <Empty className="border-0">
            <EmptyHeader>
              <EmptyMedia variant="icon">
                <HugeiconsIcon icon={SearchRemoveIcon} strokeWidth={2} />
              </EmptyMedia>
              <EmptyTitle>No matching rows</EmptyTitle>
              <EmptyDescription>
                Nothing in this statement matched that search or filter.
              </EmptyDescription>
            </EmptyHeader>
          </Empty>
        )
      }
    />
  );
}
