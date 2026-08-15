import {createFileRoute, Link, redirect} from '@tanstack/react-router';
import {HugeiconsIcon} from '@hugeicons/react';
import {SearchRemoveIcon} from '@hugeicons/core-free-icons';

import {BrandLockup, Page} from '@/components/layout/page';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {getLatestReport} from '@/lib/report-store';

const DISPLAY_COLUMNS = [
  'Transaction Date',
  'Description',
  'Amount',
  'Dr / Cr',
] as const;

export const Route = createFileRoute('/results')({
  beforeLoad: () => {
    if (!getLatestReport()) {
      throw redirect({to: '/'});
    }
  },
  component: ResultsPage,
});

function ResultsPage() {
  const report = getLatestReport();
  if (!report) {
    return null;
  }

  return (
    <Page size="lg" className="flex flex-col gap-8">
      <BrandLockup />
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div className="flex flex-col gap-2">
          <div className="flex flex-wrap items-center gap-3">
            <h1 className="font-heading text-3xl font-semibold tracking-tight">
              Autopay matches
            </h1>
            <Badge variant="secondary">{report.transaction_count}</Badge>
          </div>
          <p className="text-muted-foreground">{report.bank_name}</p>
        </div>
        <Button
          variant="outline"
          render={<Link to="/" />}
          nativeButton={false}
        >
          Analyze another
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Matched transactions</CardTitle>
          <CardDescription>
            Rows that matched autopay / mandate keywords in your statement.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {report.transactions.length === 0 ? (
            <Empty className="border">
              <EmptyHeader>
                <EmptyMedia variant="icon">
                  <HugeiconsIcon icon={SearchRemoveIcon} strokeWidth={2} />
                </EmptyMedia>
                <EmptyTitle>No autopay matches</EmptyTitle>
                <EmptyDescription>
                  No autopay matches found for this statement.
                </EmptyDescription>
              </EmptyHeader>
            </Empty>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  {DISPLAY_COLUMNS.map((col) => (
                    <TableHead key={col}>{col}</TableHead>
                  ))}
                </TableRow>
              </TableHeader>
              <TableBody>
                {report.transactions.map((row, i) => (
                  <TableRow key={`${row['Sl. No.'] ?? i}-${row.Description}`}>
                    {DISPLAY_COLUMNS.map((col) => (
                      <TableCell
                        key={col}
                        className={
                          col === 'Description' ? 'max-w-md' : undefined
                        }
                      >
                        <span
                          className={
                            col === 'Description' ? 'break-all' : undefined
                          }
                        >
                          {row[col] ?? '—'}
                        </span>
                      </TableCell>
                    ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </Page>
  );
}
