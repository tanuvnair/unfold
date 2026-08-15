import {Link} from '@tanstack/react-router';
import type {ReactNode} from 'react';

import {cn} from '@/lib/utils';

type PageSize = 'md' | 'lg';

const pageWidth: Record<PageSize, string> = {
  md: 'max-w-3xl',
  lg: 'max-w-5xl',
};

function UnfoldMark({className}: {className?: string}) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
      className={cn('size-6', className)}
    >
      <rect
        x="4.5"
        y="3.5"
        width="13"
        height="13"
        rx="2.5"
        stroke="currentColor"
        strokeWidth="1.5"
      />
      <rect
        x="6.5"
        y="7.5"
        width="13"
        height="13"
        rx="2.5"
        fill="currentColor"
        fillOpacity="0.12"
        stroke="currentColor"
        strokeWidth="1.5"
      />
    </svg>
  );
}

export function BrandLockup() {
  return (
    <Link
      to="/"
      aria-label="unfold home"
      className="flex w-fit items-center gap-3 text-foreground"
    >
      <span className="flex size-12 items-center justify-center rounded-xl bg-primary/10 text-primary ring-1 ring-primary/20">
        <UnfoldMark />
      </span>
      <span className="font-heading text-2xl font-semibold tracking-tight">
        unfold
      </span>
    </Link>
  );
}

export function Page({
  children,
  size = 'md',
  className,
}: {
  children: ReactNode;
  size?: PageSize;
  className?: string;
}) {
  return (
    <div
      className={cn(
        'mx-auto w-full px-6 py-10',
        'animate-in fade-in-0 slide-in-from-bottom-2 fill-mode-both duration-300',
        pageWidth[size],
        className,
      )}
    >
      {children}
    </div>
  );
}
