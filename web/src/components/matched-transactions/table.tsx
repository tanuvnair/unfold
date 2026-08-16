import {Fragment, useEffect, useState} from 'react';
import {keepPreviousData, useQuery} from '@tanstack/react-query';
import {useNavigate, useSearch} from '@tanstack/react-router';
import {HugeiconsIcon} from '@hugeicons/react';
import {
  ArrowDown01Icon,
  ArrowRight01Icon,
  Search01Icon,
  SearchRemoveIcon,
} from '@hugeicons/core-free-icons';

import {
  confidenceBadgeVariant,
  confidenceLabel,
  sourceBadgeVariant,
  sourceLabel,
} from '@/components/matched-transactions/badge-variants';
import {columns} from '@/components/matched-transactions/columns';
import {Alert, AlertDescription, AlertTitle} from '@/components/ui/alert';
import {Badge} from '@/components/ui/badge';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import {DataTable} from '@/components/ui/data-table';
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty';
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field';
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from '@/components/ui/input-group';
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemMedia,
  ItemSeparator,
  ItemTitle,
} from '@/components/ui/item';
import {Skeleton} from '@/components/ui/skeleton';
import {ToggleGroup, ToggleGroupItem} from '@/components/ui/toggle-group';
import {
  fetchReportSummary,
  fetchReportTransactions,
  type ConfidenceFilter,
  type Report,
  type SourceFilter,
  type SummaryGroup,
  type TransactionTypeFilter,
} from '@/lib/api';
import {useDebouncedValue} from '@/hooks/use-debounced-value';
import type {ResultsView} from '@/lib/results-search';

const TYPE_ITEMS: Array<{label: string; value: TransactionTypeFilter}> = [
  {label: 'All', value: 'all'},
  {label: 'DR', value: 'DR'},
  {label: 'CR', value: 'CR'},
];

const CONFIDENCE_ITEMS: Array<{label: string; value: ConfidenceFilter}> = [
  {label: 'All', value: 'all'},
  {label: 'High', value: 'high'},
  {label: 'Medium', value: 'medium'},
  {label: 'Low', value: 'low'},
];

const SOURCE_ITEMS: Array<{label: string; value: SourceFilter}> = [
  {label: 'All', value: 'all'},
  {label: 'Keyword', value: 'keyword'},
  {label: 'Pattern', value: 'recurrence'},
  {label: 'Both', value: 'both'},
];

const VIEW_ITEMS: Array<{label: string; value: ResultsView}> = [
  {label: 'Grouped', value: 'grouped'},
  {label: 'Transactions', value: 'transactions'},
];

function formatINR(amount: number): string {
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    maximumFractionDigits: 0,
  }).format(amount);
}

function ResultsFilters({
  searchInput,
  onSearchInput,
}: {
  searchInput: string;
  onSearchInput: (value: string) => void;
}) {
  const search = useSearch({from: '/results'});
  const navigate = useNavigate({from: '/results'});

  function setType(type: TransactionTypeFilter) {
    void navigate({search: (prev) => ({...prev, type, page: 1})});
  }
  function setConfidence(confidence: ConfidenceFilter) {
    void navigate({search: (prev) => ({...prev, confidence, page: 1})});
  }
  function setSource(source: SourceFilter) {
    void navigate({search: (prev) => ({...prev, source, page: 1})});
  }
  function setView(view: ResultsView) {
    void navigate({search: (prev) => ({...prev, view, page: 1})});
  }

  return (
    <FieldGroup>
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Field className="w-full max-w-sm">
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
              onChange={(event) => onSearchInput(event.target.value)}
              placeholder="Search descriptions..."
            />
          </InputGroup>
        </Field>
        <Field orientation="horizontal" className="w-fit shrink-0">
          <FieldLabel className="sr-only" id="txn-view-label">
            Results view
          </FieldLabel>
          <ToggleGroup
            variant="outline"
            spacing={0}
            value={[search.view]}
            aria-labelledby="txn-view-label"
            onValueChange={(value) => {
              const next = value[0] as ResultsView | undefined;
              if (next) {
                setView(next);
              }
            }}
          >
            {VIEW_ITEMS.map((item) => (
              <ToggleGroupItem key={item.value} value={item.value}>
                {item.label}
              </ToggleGroupItem>
            ))}
          </ToggleGroup>
        </Field>
      </div>
      <FieldGroup className="gap-3 sm:flex-row sm:flex-wrap sm:items-center">
        <Field orientation="horizontal" className="w-fit">
          <FieldLabel className="sr-only" id="txn-source-label">
            Detection source
          </FieldLabel>
          <ToggleGroup
            variant="outline"
            size="sm"
            spacing={0}
            value={[search.source]}
            aria-labelledby="txn-source-label"
            onValueChange={(value) => {
              const next = value[0] as SourceFilter | undefined;
              if (next) {
                setSource(next);
              }
            }}
          >
            {SOURCE_ITEMS.map((item) => (
              <ToggleGroupItem key={item.value} value={item.value}>
                {item.label}
              </ToggleGroupItem>
            ))}
          </ToggleGroup>
        </Field>
        <Field orientation="horizontal" className="w-fit">
          <FieldLabel className="sr-only" id="txn-confidence-label">
            Confidence
          </FieldLabel>
          <ToggleGroup
            variant="outline"
            size="sm"
            spacing={0}
            value={[search.confidence]}
            aria-labelledby="txn-confidence-label"
            onValueChange={(value) => {
              const next = value[0] as ConfidenceFilter | undefined;
              if (next) {
                setConfidence(next);
              }
            }}
          >
            {CONFIDENCE_ITEMS.map((item) => (
              <ToggleGroupItem key={item.value} value={item.value}>
                {item.label}
              </ToggleGroupItem>
            ))}
          </ToggleGroup>
        </Field>
        <Field orientation="horizontal" className="w-fit">
          <FieldLabel className="sr-only" id="txn-type-label">
            Debit or credit
          </FieldLabel>
          <ToggleGroup
            variant="outline"
            size="sm"
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
      </FieldGroup>
    </FieldGroup>
  );
}

function PayeeGroupRow({report, group}: {report: Report; group: SummaryGroup}) {
  const [open, setOpen] = useState(false);
  const search = useSearch({from: '/results'});

  const detail = useQuery({
    queryKey: ['report-transactions', report.id, 'payee', group.payee, search],
    queryFn: () =>
      fetchReportTransactions(report.id, {
        q: '',
        type: search.type,
        confidence: search.confidence,
        source: search.source,
        page: 0,
        pageSize: Math.max(group.occurrenceCount, 10),
        payee: group.payee,
      }),
    enabled: open,
  });

  const amountLabel = group.isMonthlyCadence
    ? `${formatINR(group.latestAmount)}/mo`
    : formatINR(group.totalAmount);
  const countLabel = `× ${group.occurrenceCount}`;

  return (
    <Item
      size="sm"
      className="w-full min-w-0 flex-col flex-nowrap items-stretch rounded-none border-0 px-0 py-0"
      role="listitem"
    >
      <Collapsible open={open} onOpenChange={setOpen} className="w-full min-w-0">
        <CollapsibleTrigger className="flex h-auto w-full max-w-full min-w-0 shrink cursor-pointer items-center justify-start gap-3 rounded-2xl px-3.5 py-3 text-left text-sm font-medium outline-none transition-all focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/30 disabled:pointer-events-none disabled:opacity-50">
          <ItemMedia variant="icon" className="shrink-0">
            <HugeiconsIcon
              icon={open ? ArrowDown01Icon : ArrowRight01Icon}
              strokeWidth={2}
            />
          </ItemMedia>
          <ItemContent className="min-w-0 flex-1 overflow-hidden">
            <ItemTitle className="block w-full min-w-0 max-w-full truncate">
              {group.payee}
            </ItemTitle>
            <ItemDescription className="truncate">
              <span className="tabular-nums">
                {amountLabel} · {countLabel}
              </span>
            </ItemDescription>
          </ItemContent>
          <ItemActions className="shrink-0">
            {group.confidence ? (
              <Badge variant={confidenceBadgeVariant(group.confidence)}>
                {confidenceLabel(group.confidence)}
              </Badge>
            ) : null}
            {group.source ? (
              <Badge variant={sourceBadgeVariant(group.source)}>
                {sourceLabel(group.source)}
              </Badge>
            ) : null}
          </ItemActions>
        </CollapsibleTrigger>
        <CollapsibleContent className="pb-3 pl-7">
          {detail.isError ? (
            <Alert variant="destructive">
              <AlertTitle>Could not load transactions</AlertTitle>
              <AlertDescription>
                {(detail.error as Error).message}
              </AlertDescription>
            </Alert>
          ) : (
            <DataTable
              columns={columns}
              data={detail.data?.rows ?? []}
              rowCount={detail.data?.rowCount ?? 0}
              pagination={{
                pageIndex: 0,
                pageSize: Math.max(group.occurrenceCount, 10),
              }}
              onPaginationChange={() => {}}
              getRowId={(row) => row.id}
              isLoading={detail.isPending}
              isFetching={detail.isFetching}
              showPagination={false}
              variant="embedded"
              empty={
                <Empty className="border-0 py-6">
                  <EmptyHeader>
                    <EmptyTitle>No rows</EmptyTitle>
                  </EmptyHeader>
                </Empty>
              }
            />
          )}
        </CollapsibleContent>
      </Collapsible>
    </Item>
  );
}

function GroupsLoading() {
  return (
    <div className="flex flex-col gap-3">
      <Skeleton className="h-4 w-48" />
      <ItemGroup className="gap-0">
        <Skeleton className="h-14 w-full" />
        <ItemSeparator className="my-0" />
        <Skeleton className="h-14 w-full" />
        <ItemSeparator className="my-0" />
        <Skeleton className="h-14 w-full" />
      </ItemGroup>
    </div>
  );
}

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

  const summaryQuery = useQuery({
    queryKey: [
      'report-summary',
      report.id,
      search.q,
      search.type,
      search.confidence,
      search.source,
    ],
    queryFn: () =>
      fetchReportSummary(report.id, {
        q: search.q,
        type: search.type,
        confidence: search.confidence,
        source: search.source,
      }),
    placeholderData: keepPreviousData,
    enabled: search.view === 'grouped',
  });

  const txQuery = useQuery({
    queryKey: ['report-transactions', report.id, search],
    queryFn: () =>
      fetchReportTransactions(report.id, {
        q: search.q,
        type: search.type,
        confidence: search.confidence,
        source: search.source,
        page: search.page - 1,
        pageSize: search.pageSize,
      }),
    placeholderData: keepPreviousData,
    enabled: search.view === 'transactions',
  });

  const rows = txQuery.data?.rows ?? [];
  const rowCount = txQuery.data?.rowCount ?? 0;

  useEffect(() => {
    if (search.view !== 'transactions' || txQuery.isPending) {
      return;
    }
    const maxPage = Math.max(1, Math.ceil(rowCount / search.pageSize) || 1);
    if (search.page > maxPage) {
      void navigate({
        search: (prev) => ({...prev, page: maxPage}),
        replace: true,
      });
    }
  }, [
    navigate,
    rowCount,
    search.page,
    search.pageSize,
    search.view,
    txQuery.isPending,
  ]);

  const pagination = {
    pageIndex: search.page - 1,
    pageSize: search.pageSize,
  };

  const activeError =
    search.view === 'grouped' ? summaryQuery.error : txQuery.error;
  const isError =
    search.view === 'grouped' ? summaryQuery.isError : txQuery.isError;
  const groups = summaryQuery.data?.groups ?? [];

  return (
    <div className="flex flex-col gap-6">
      <ResultsFilters
        searchInput={searchInput}
        onSearchInput={setSearchInput}
      />
      {isError ? (
        <Alert variant="destructive">
          <AlertTitle>Could not load results</AlertTitle>
          <AlertDescription>
            {(activeError as Error).message}. Try analyzing the statement again.
          </AlertDescription>
        </Alert>
      ) : null}

      {search.view === 'grouped' ? (
        <div className="flex flex-col gap-4">
          {summaryQuery.data && summaryQuery.data.estimatedMonthlyTotal > 0 ? (
            <FieldDescription>
              Estimated monthly{' '}
              <span className="font-medium tabular-nums text-foreground">
                {formatINR(summaryQuery.data.estimatedMonthlyTotal)}
              </span>
            </FieldDescription>
          ) : null}
          {summaryQuery.isPending && !summaryQuery.data ? (
            <GroupsLoading />
          ) : groups.length === 0 ? (
            <Empty className="border-0">
              <EmptyHeader>
                <EmptyMedia variant="icon">
                  <HugeiconsIcon icon={SearchRemoveIcon} strokeWidth={2} />
                </EmptyMedia>
                <EmptyTitle>No matching groups</EmptyTitle>
                <EmptyDescription>
                  Nothing in this statement matched that search or filter.
                </EmptyDescription>
              </EmptyHeader>
            </Empty>
          ) : (
            <ItemGroup className="min-w-0 gap-0">
              {groups.map((group, index) => (
                <Fragment key={group.payee}>
                  {index > 0 ? <ItemSeparator className="my-0" /> : null}
                  <PayeeGroupRow report={report} group={group} />
                </Fragment>
              ))}
            </ItemGroup>
          )}
        </div>
      ) : (
        <DataTable
          columns={columns}
          data={rows}
          rowCount={rowCount}
          pagination={pagination}
          onPaginationChange={(updater) => {
            const next =
              typeof updater === 'function' ? updater(pagination) : updater;
            void navigate({
              search: (prev) => ({
                ...prev,
                page:
                  next.pageSize !== search.pageSize ? 1 : next.pageIndex + 1,
                pageSize: next.pageSize,
              }),
            });
          }}
          getRowId={(row) => row.id}
          isLoading={txQuery.isPending}
          isFetching={txQuery.isFetching}
          empty={
            isError ? (
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
      )}
    </div>
  );
}
