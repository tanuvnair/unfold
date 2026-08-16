import {Fragment, useEffect, useState} from 'react';
import {keepPreviousData, useQuery} from '@tanstack/react-query';
import {useNavigate, useSearch} from '@tanstack/react-router';
import {format, parse} from 'date-fns';
import {HugeiconsIcon} from '@hugeicons/react';
import {
  ArrowDown01Icon,
  ArrowRight01Icon,
  Calendar01Icon,
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
import {Button} from '@/components/ui/button';
import {Calendar} from '@/components/ui/calendar';
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
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {Separator} from '@/components/ui/separator';
import {Skeleton} from '@/components/ui/skeleton';
import {ToggleGroup, ToggleGroupItem} from '@/components/ui/toggle-group';
import {
  fetchReportSummary,
  fetchReportTransactions,
  type ConfidenceFilter,
  type Report,
  type SourceFilter,
  type SummaryGroup,
} from '@/lib/api';
import {useDebouncedValue} from '@/hooks/use-debounced-value';
import {
  resultsSearchDefaults,
  type ResultsView,
} from '@/lib/results-search';
import {cn} from '@/lib/utils';

function parseISODateLocal(value: string): Date | undefined {
  if (!value) {
    return undefined;
  }
  return parse(value, 'yyyy-MM-dd', new Date());
}

function formatISODateLocal(date: Date): string {
  return format(date, 'yyyy-MM-dd');
}

function DatePickerField({
  id,
  label,
  value,
  onChange,
  timeZone,
  disabled,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (next: string) => void;
  timeZone?: string;
  disabled?: (date: Date) => boolean;
}) {
  const [open, setOpen] = useState(false);
  const selected = parseISODateLocal(value);

  return (
    <Field className="w-fit">
      <FieldLabel htmlFor={id}>{label}</FieldLabel>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger
          render={
            <Button
              id={id}
              variant="outline"
              size="sm"
              data-empty={!selected}
              className={cn(
                'min-w-40 justify-start text-left font-normal',
                'data-[empty=true]:text-muted-foreground',
              )}
            />
          }
        >
          <HugeiconsIcon
            icon={Calendar01Icon}
            strokeWidth={2}
            data-icon="inline-start"
          />
          {selected ? (
            format(selected, 'LLL dd, y')
          ) : (
            <span>Pick a date</span>
          )}
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="start">
          <Calendar
            mode="single"
            selected={selected}
            defaultMonth={selected}
            onSelect={(date) => {
              onChange(date ? formatISODateLocal(date) : '');
              setOpen(false);
            }}
            disabled={disabled}
            timeZone={timeZone}
          />
        </PopoverContent>
      </Popover>
    </Field>
  );
}

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
  const [timeZone, setTimeZone] = useState<string | undefined>(undefined);

  useEffect(() => {
    setTimeZone(Intl.DateTimeFormat().resolvedOptions().timeZone);
  }, []);

  const fromDate = parseISODateLocal(search.from);
  const toDate = parseISODateLocal(search.to);

  function setConfidence(confidence: ConfidenceFilter) {
    void navigate({search: (prev) => ({...prev, confidence, page: 1})});
  }
  function setSource(source: SourceFilter) {
    void navigate({search: (prev) => ({...prev, source, page: 1})});
  }
  function setFrom(from: string) {
    void navigate({search: (prev) => ({...prev, from, page: 1})});
  }
  function setTo(to: string) {
    void navigate({search: (prev) => ({...prev, to, page: 1})});
  }
  function clearFilters() {
    onSearchInput('');
    void navigate({
      search: (prev) => ({
        ...prev,
        q: resultsSearchDefaults.q,
        from: resultsSearchDefaults.from,
        to: resultsSearchDefaults.to,
        confidence: resultsSearchDefaults.confidence,
        source: resultsSearchDefaults.source,
        page: 1,
      }),
    });
  }

  const filtersActive =
    Boolean(searchInput.trim()) ||
    Boolean(search.q.trim()) ||
    Boolean(search.from) ||
    Boolean(search.to) ||
    search.confidence !== resultsSearchDefaults.confidence ||
    search.source !== resultsSearchDefaults.source;

  return (
    <FieldGroup>
      <FieldGroup className="gap-3 sm:flex-row sm:items-center">
        <Field className="min-w-0 flex-1">
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
        <Field className="w-fit shrink-0">
          <FieldLabel className="sr-only" htmlFor="txn-clear-filters">
            Clear filters
          </FieldLabel>
          <Button
            id="txn-clear-filters"
            type="button"
            variant="ghost"
            size="sm"
            disabled={!filtersActive}
            onClick={clearFilters}
          >
            Clear filters
          </Button>
        </Field>
      </FieldGroup>

      <FieldGroup className="gap-4 sm:flex-row sm:flex-wrap sm:items-end">
        <DatePickerField
          id="txn-from"
          label="From"
          value={search.from}
          onChange={setFrom}
          timeZone={timeZone}
          disabled={toDate ? (date) => date > toDate : undefined}
        />
        <DatePickerField
          id="txn-to"
          label="To"
          value={search.to}
          onChange={setTo}
          timeZone={timeZone}
          disabled={fromDate ? (date) => date < fromDate : undefined}
        />
        <Field className="w-fit">
          <FieldLabel id="txn-source-label">Source</FieldLabel>
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
        <Field className="w-fit">
          <FieldLabel id="txn-confidence-label">Confidence</FieldLabel>
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
      </FieldGroup>
    </FieldGroup>
  );
}

export function ResultsViewToggle() {
  const search = useSearch({from: '/results'});
  const navigate = useNavigate({from: '/results'});

  return (
    <ToggleGroup
      variant="outline"
      size="sm"
      spacing={0}
      value={[search.view]}
      aria-label="Results view"
      onValueChange={(value) => {
        const next = value[0] as ResultsView | undefined;
        if (next) {
          void navigate({search: (prev) => ({...prev, view: next, page: 1})});
        }
      }}
    >
      {VIEW_ITEMS.map((item) => (
        <ToggleGroupItem key={item.value} value={item.value}>
          {item.label}
        </ToggleGroupItem>
      ))}
    </ToggleGroup>
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
        type: 'DR',
        confidence: search.confidence,
        source: search.source,
        page: 0,
        pageSize: Math.max(group.occurrenceCount, 10),
        payee: group.payee,
        from: search.from,
        to: search.to,
      }),
    enabled: open,
  });

  const amountLabel = group.isMonthlyCadence
    ? `${formatINR(group.latestAmount)}/mo`
    : formatINR(group.totalAmount);
  const countLabel = `× ${group.occurrenceCount}`;

  return (
    <Collapsible
      open={open}
      onOpenChange={setOpen}
      role="listitem"
      className="w-full min-w-0"
    >
      <Item
        size="sm"
        className="w-full min-w-0"
        render={<CollapsibleTrigger />}
      >
        <ItemMedia variant="icon">
          <HugeiconsIcon
            icon={open ? ArrowDown01Icon : ArrowRight01Icon}
            strokeWidth={2}
          />
        </ItemMedia>
        <ItemContent className="min-w-0 flex-1 overflow-hidden">
          <ItemTitle className="min-w-0 max-w-full">
            {group.payee}
          </ItemTitle>
          <ItemDescription className="truncate">
            <span className="tabular-nums">
              {amountLabel} · {countLabel}
            </span>
          </ItemDescription>
        </ItemContent>
        <ItemActions>
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
      </Item>
      <CollapsibleContent className="pb-3 pl-11">
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
      search.confidence,
      search.source,
      search.from,
      search.to,
    ],
    queryFn: () =>
      fetchReportSummary(report.id, {
        q: search.q,
        type: 'DR',
        confidence: search.confidence,
        source: search.source,
        from: search.from,
        to: search.to,
      }),
    placeholderData: keepPreviousData,
    enabled: search.view === 'grouped',
  });

  const txQuery = useQuery({
    queryKey: ['report-transactions', report.id, search],
    queryFn: () =>
      fetchReportTransactions(report.id, {
        q: search.q,
        type: 'DR',
        confidence: search.confidence,
        source: search.source,
        page: search.page - 1,
        pageSize: search.pageSize,
        from: search.from,
        to: search.to,
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
      <Separator />
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
              <span className="tabular-nums">
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
