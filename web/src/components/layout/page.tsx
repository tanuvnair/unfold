import {Link} from '@tanstack/react-router';
import type {ReactNode} from 'react';

import {cn} from '@/lib/utils';

type PageSize = 'md' | 'lg';

const pageWidth: Record<PageSize, string> = {
  md: 'max-w-3xl',
  lg: 'max-w-5xl',
};

export function BrandLockup() {
  return (
    <Link
      to="/"
      aria-label="unfold home"
      className="flex w-fit items-center gap-3 text-foreground"
    >
      <img src="/favicon.svg" alt="" className="size-12" />
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
