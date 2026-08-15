import {useEffect, useMemo, useState} from 'react';
import {createFileRoute, useNavigate} from '@tanstack/react-router';
import {useMutation, useQuery} from '@tanstack/react-query';
import {HugeiconsIcon} from '@hugeicons/react';
import {Upload01Icon} from '@hugeicons/core-free-icons';

import {Page} from '@/components/layout/page';
import {Alert, AlertDescription, AlertTitle} from '@/components/ui/alert';
import {Button} from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
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
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {Skeleton} from '@/components/ui/skeleton';
import {Spinner} from '@/components/ui/spinner';
import {analyzeStatement, fetchBanks} from '@/lib/api';
import {setLatestReport} from '@/lib/report-store';

export const Route = createFileRoute('/')({
  component: HomePage,
});

function HomePage() {
  const navigate = useNavigate();
  const [bank, setBank] = useState('');
  const [file, setFile] = useState<File | null>(null);

  const banksQuery = useQuery({
    queryKey: ['banks'],
    queryFn: fetchBanks,
  });

  const selectableBanks = useMemo(
    () => banksQuery.data?.filter((b) => b.has_parser) ?? [],
    [banksQuery.data],
  );

  useEffect(() => {
    if (!bank && selectableBanks.length === 1) {
      setBank(selectableBanks[0].key);
    }
  }, [bank, selectableBanks]);

  const analyzeMutation = useMutation({
    mutationFn: async () => {
      if (!file) {
        throw new Error('Choose a bank statement CSV first.');
      }
      return analyzeStatement(file, bank);
    },
    onSuccess: (report) => {
      setLatestReport(report);
      void navigate({to: '/results'});
    },
  });

  const bankItems = useMemo(
    () => [
      {
        label: 'Select a bank',
        value: null as string | null,
      },
      ...selectableBanks.map((b) => ({
        label: b.bank_name,
        value: b.key,
      })),
    ],
    [selectableBanks],
  );

  const pending = analyzeMutation.isPending;
  const canSubmit =
    Boolean(file && bank) && selectableBanks.length > 0 && !pending;

  return (
    <Page size="md" className="flex flex-col gap-8">
      <div className="flex flex-col gap-2">
        <h1 className="font-heading text-3xl font-semibold tracking-tight">
          Find hidden subscriptions
        </h1>
        <p className="text-muted-foreground">
          Upload a bank statement CSV to surface autopay and mandate charges.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Analyze statement</CardTitle>
          <CardDescription>
            Choose a bank profile, then upload the statement export.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            id="analyze-form"
            onSubmit={(e) => {
              e.preventDefault();
              analyzeMutation.mutate();
            }}
          >
            <FieldGroup>
              {banksQuery.isLoading ? (
                <Field>
                  <FieldLabel>Bank</FieldLabel>
                  <Skeleton className="h-8 w-full" />
                </Field>
              ) : (
                <Field data-disabled={selectableBanks.length === 0 || undefined}>
                  <FieldLabel htmlFor="bank">Bank</FieldLabel>
                  <Select
                    items={bankItems}
                    value={bank || null}
                    onValueChange={(value) => setBank(value ?? '')}
                    disabled={selectableBanks.length === 0}
                  >
                    <SelectTrigger id="bank" className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        {bankItems.map((item) => (
                          <SelectItem
                            key={item.value ?? 'placeholder'}
                            value={item.value}
                          >
                            {item.label}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </Field>
              )}

              <Field>
                <FieldLabel htmlFor="statement">Statement CSV</FieldLabel>
                <InputGroup>
                  <InputGroupAddon>
                    <HugeiconsIcon icon={Upload01Icon} strokeWidth={2} />
                  </InputGroupAddon>
                  <InputGroupInput
                    id="statement"
                    type="file"
                    accept=".csv,text/csv"
                    onChange={(e) => {
                      setFile(e.target.files?.[0] ?? null);
                    }}
                  />
                </InputGroup>
                {file ? (
                  <FieldDescription>Selected: {file.name}</FieldDescription>
                ) : (
                  <FieldDescription>
                    Supported bank statement exports (.csv).
                  </FieldDescription>
                )}
              </Field>

              {banksQuery.isError ? (
                <Alert variant="destructive">
                  <AlertTitle>Could not load banks</AlertTitle>
                  <AlertDescription>
                    {(banksQuery.error as Error).message}. Is the API running?
                  </AlertDescription>
                </Alert>
              ) : null}

              {analyzeMutation.isError ? (
                <Alert variant="destructive">
                  <AlertTitle>Analyze failed</AlertTitle>
                  <AlertDescription>
                    {(analyzeMutation.error as Error).message}
                  </AlertDescription>
                </Alert>
              ) : null}
            </FieldGroup>
          </form>
        </CardContent>
        <CardFooter className="justify-end border-t">
          <Button
            type="submit"
            form="analyze-form"
            size="lg"
            disabled={!canSubmit}
          >
            {pending ? <Spinner data-icon="inline-start" /> : null}
            {pending ? 'Analyzing…' : 'Analyze statement'}
          </Button>
        </CardFooter>
      </Card>
    </Page>
  );
}
