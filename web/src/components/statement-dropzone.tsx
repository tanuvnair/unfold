import {useRef, useState, type ChangeEvent, type DragEvent} from 'react';
import {HugeiconsIcon} from '@hugeicons/react';
import {Cancel01Icon, File01Icon, Upload01Icon} from '@hugeicons/core-free-icons';

import {Button} from '@/components/ui/button';
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty';
import {cn} from '@/lib/utils';

function formatBytes(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function StatementDropzone({
  id,
  file,
  disabled = false,
  invalid = false,
  onFile,
}: {
  id: string;
  file: File | null;
  disabled?: boolean;
  invalid?: boolean;
  onFile: (file: File | null) => void;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const dragDepth = useRef(0);
  const [dragging, setDragging] = useState(false);

  function takeFile(next: File | null) {
    if (inputRef.current && next === null) {
      inputRef.current.value = '';
    }
    onFile(next);
  }

  function onInputChange(event: ChangeEvent<HTMLInputElement>) {
    takeFile(event.target.files?.[0] ?? null);
  }

  function onDragEnter(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    if (disabled) {
      return;
    }
    dragDepth.current += 1;
    setDragging(true);
  }

  function onDragOver(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
  }

  function onDragLeave(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    dragDepth.current = Math.max(0, dragDepth.current - 1);
    if (dragDepth.current === 0) {
      setDragging(false);
    }
  }

  function onDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    dragDepth.current = 0;
    setDragging(false);
    if (disabled) {
      return;
    }
    takeFile(event.dataTransfer.files?.[0] ?? null);
  }

  return (
    <div
      onDragEnter={onDragEnter}
      onDragOver={onDragOver}
      onDragLeave={onDragLeave}
      onDrop={onDrop}
      className={cn(
        'has-[:focus-visible]:rounded-3xl has-[:focus-visible]:ring-3 has-[:focus-visible]:ring-ring/30',
        disabled && 'pointer-events-none cursor-not-allowed opacity-50',
      )}
    >
      <input
        ref={inputRef}
        id={id}
        type="file"
        accept=".csv,text/csv"
        disabled={disabled}
        className="sr-only"
        aria-invalid={invalid || undefined}
        onChange={onInputChange}
      />
      {file ? (
        <div className="flex items-center gap-3 rounded-3xl border border-dashed px-4 py-3">
          <label
            htmlFor={id}
            className="flex min-w-0 flex-1 cursor-pointer items-center gap-3"
          >
            <EmptyMedia variant="icon" className="mb-0">
              <HugeiconsIcon icon={File01Icon} strokeWidth={2} />
            </EmptyMedia>
            <span className="min-w-0 flex-1">
              <span className="block truncate text-sm font-medium">
                {file.name}
              </span>
              <span className="block text-sm text-muted-foreground tabular-nums">
                {formatBytes(file.size)}
              </span>
            </span>
          </label>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            aria-label="Remove statement"
            onClick={() => takeFile(null)}
          >
            <HugeiconsIcon
              icon={Cancel01Icon}
              strokeWidth={2}
              data-icon="inline-start"
            />
          </Button>
        </div>
      ) : (
        <label htmlFor={id} className="block cursor-pointer">
          <Empty
            className={cn(
              'min-h-32 border p-8',
              dragging && 'bg-muted',
              invalid && 'border-destructive ring-3 ring-destructive/20',
            )}
          >
            <EmptyHeader>
              <EmptyMedia variant="icon">
                <HugeiconsIcon icon={Upload01Icon} strokeWidth={2} />
              </EmptyMedia>
              <EmptyTitle>
                {dragging ? 'Drop to attach' : 'Drop a CSV here, or choose a file'}
              </EmptyTitle>
              <EmptyDescription>.csv, max 10 MB</EmptyDescription>
            </EmptyHeader>
          </Empty>
        </label>
      )}
    </div>
  );
}
