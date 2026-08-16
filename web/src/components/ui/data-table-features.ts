import {
  rowPaginationFeature,
  tableFeatures,
} from '@tanstack/react-table';

declare module '@tanstack/react-table' {
  interface ColumnMeta<TFeatures, TData, TValue> {
    cellClassName?: string;
    headerClassName?: string;
  }
}

export const features = tableFeatures({
  rowPaginationFeature,
});

export type DataTableFeatures = typeof features;
