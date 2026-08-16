import {createFileRoute, Link, redirect, stripSearchParams} from '@tanstack/react-router';
import {HugeiconsIcon} from '@hugeicons/react';
import {SearchRemoveIcon} from '@hugeicons/core-free-icons';

import {MatchedTransactionsTable} from '@/components/matched-transactions/table';
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
import {Separator} from '@/components/ui/separator';
import {getLatestReport} from '@/lib/report-store';
import {
  parseResultsSearch,
  resultsSearchDefaults,
} from '@/lib/results-search';

export const Route = createFileRoute('/results')({
  validateSearch: parseResultsSearch,
  search: {
    middlewares: [stripSearchParams(resultsSearchDefaults)],
  },
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
          <CardTitle>Matched charges</CardTitle>
          <CardDescription>
            Expand a payee for the underlying rows, or switch to the flat
            transaction list.
          </CardDescription>
        </CardHeader>
        <Separator />
        <CardContent>
          {report.transaction_count === 0 ? (
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
            <MatchedTransactionsTable report={report} />
          )}
        </CardContent>
      </Card>
    </Page>
  );
}
