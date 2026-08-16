import type {ReactNode} from 'react';
import {
  useTable,
  type ColumnDef,
  type OnChangeFn,
  type PaginationState,
  type RowData,
} from '@tanstack/react-table';

import {features, type DataTableFeatures} from '@/components/ui/data-table-features';
import {
  Empty,
  EmptyHeader,
  EmptyTitle,
} from '@/components/ui/empty';
import {Field, FieldLabel} from '@/components/ui/field';
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {Skeleton} from '@/components/ui/skeleton';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {cn} from '@/lib/utils';

const PAGE_SIZE_ITEMS = [
  {label: '10', value: '10'},
  {label: '20', value: '20'},
  {label: '25', value: '25'},
  {label: '50', value: '50'},
];

type PageItem = number | 'ellipsis';

function getPageItems(pageCount: number, currentPage: number): PageItem[] {
  if (pageCount <= 5) {
    return Array.from({length: pageCount}, (_, i) => i + 1);
  }

  const items: PageItem[] = [1];
  const start = Math.max(2, currentPage - 1);
  const end = Math.min(pageCount - 1, currentPage + 1);

  if (start > 2) {
    items.push('ellipsis');
  }
  for (let page = start; page <= end; page++) {
    items.push(page);
  }
  if (end < pageCount - 1) {
    items.push('ellipsis');
  }
  items.push(pageCount);
  return items;
}

export type DataTableProps<TData extends RowData> = {
  columns: ColumnDef<DataTableFeatures, TData>[];
  data: TData[];
  rowCount: number;
  pagination: PaginationState;
  onPaginationChange: OnChangeFn<PaginationState>;
  getRowId?: (originalRow: TData, index: number) => string;
  isLoading?: boolean;
  isFetching?: boolean;
  toolbar?: ReactNode;
  empty?: ReactNode;
  /** Hide rows-per-page and page links (nested / complete pages). */
  showPagination?: boolean;
  /** Flush table without the sealed border shell — for accordion details. */
  variant?: 'default' | 'embedded';
};

export function DataTable<TData extends RowData>({
  columns,
  data,
  rowCount,
  pagination,
  onPaginationChange,
  getRowId,
  isLoading = false,
  isFetching = false,
  toolbar,
  empty,
  showPagination = true,
  variant = 'default',
}: DataTableProps<TData>) {
  const table = useTable(
    {
      features,
      data,
      columns,
      getRowId,
      rowCount,
      manualPagination: true,
      state: {pagination},
      onPaginationChange,
    },
    (state) => state,
  );

  const pageCount = Math.max(table.getPageCount(), 1);
  const colSpan = columns.length;
  const canPrevious = table.getCanPreviousPage();
  const canNext = table.getCanNextPage();
  const paginationVisible =
    showPagination && !isLoading && rowCount > 0;
  const skeletonRows =
    variant === 'embedded'
      ? Math.min(pagination.pageSize, 4)
      : pagination.pageSize;
  const embedded = variant === 'embedded';

  return (
    <div className={cn('flex flex-col', embedded ? 'gap-0' : 'gap-4')}>
      {toolbar}
      <div
        className={cn(
          'overflow-hidden',
          !embedded && 'rounded-2xl border',
          isFetching && !isLoading && 'opacity-70',
        )}
      >
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead
                    key={header.id}
                    className={header.column.columnDef.meta?.headerClassName}
                  >
                    {header.isPlaceholder ? null : (
                      <table.FlexRender header={header} />
                    )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({length: skeletonRows}, (_, rowIndex) => (
                <TableRow key={`skeleton-${rowIndex}`}>
                  {columns.map((_, cellIndex) => (
                    <TableCell key={`skeleton-${rowIndex}-${cellIndex}`}>
                      <Skeleton className="h-4 w-full" />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getAllCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className={cell.column.columnDef.meta?.cellClassName}
                    >
                      <table.FlexRender cell={cell} />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={colSpan} className="h-24 p-0">
                  {empty ?? (
                    <Empty className="border-0">
                      <EmptyHeader>
                        <EmptyTitle>No results</EmptyTitle>
                      </EmptyHeader>
                    </Empty>
                  )}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      {paginationVisible ? (
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <Field orientation="horizontal" className="w-fit">
            <FieldLabel htmlFor="rows-per-page">Rows per page</FieldLabel>
            <Select
              items={PAGE_SIZE_ITEMS}
              value={String(pagination.pageSize)}
              onValueChange={(value) => {
                if (!value) {
                  return;
                }
                onPaginationChange({
                  pageIndex: 0,
                  pageSize: Number(value),
                });
              }}
            >
              <SelectTrigger className="w-20" id="rows-per-page" size="sm">
                <SelectValue />
              </SelectTrigger>
              <SelectContent align="start" alignItemWithTrigger={false}>
                <SelectGroup>
                  {PAGE_SIZE_ITEMS.map((item) => (
                    <SelectItem key={item.value} value={item.value}>
                      {item.label}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </Field>
          <Pagination className="mx-0 w-auto sm:justify-end">
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  href="#"
                  aria-disabled={!canPrevious || undefined}
                  tabIndex={canPrevious ? undefined : -1}
                  className={cn(
                    !canPrevious && 'pointer-events-none opacity-50',
                  )}
                  onClick={(event) => {
                    event.preventDefault();
                    table.previousPage();
                  }}
                />
              </PaginationItem>
              {getPageItems(pageCount, pagination.pageIndex + 1).map(
                (item, index) =>
                  item === 'ellipsis' ? (
                    <PaginationItem key={`ellipsis-${index}`}>
                      <PaginationEllipsis />
                    </PaginationItem>
                  ) : (
                    <PaginationItem key={item}>
                      <PaginationLink
                        href="#"
                        isActive={item === pagination.pageIndex + 1}
                        onClick={(event) => {
                          event.preventDefault();
                          table.setPageIndex(item - 1);
                        }}
                      >
                        {item}
                      </PaginationLink>
                    </PaginationItem>
                  ),
              )}
              <PaginationItem>
                <PaginationNext
                  href="#"
                  aria-disabled={!canNext || undefined}
                  tabIndex={canNext ? undefined : -1}
                  className={cn(!canNext && 'pointer-events-none opacity-50')}
                  onClick={(event) => {
                    event.preventDefault();
                    table.nextPage();
                  }}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      ) : null}
    </div>
  );
}
