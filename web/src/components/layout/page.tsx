import {Link} from '@tanstack/react-router';

import {cn} from '@/lib/utils';

type PageSize = 'md' | 'lg';

const pageWidth: Record<PageSize, string> = {
  md: 'max-w-3xl',
  lg: 'max-w-5xl',
};

export function AppHeader() {
  return (
    <header className="border-b">
      <div className="mx-auto flex h-14 w-full max-w-5xl items-center px-6">
        <Link
          to="/"
          className="font-heading text-lg font-semibold tracking-tight text-foreground"
        >
          unfold
        </Link>
      </div>
    </header>
  );
}

export function Page({
  children,
  size = 'md',
  className,
}: {
  children: React.ReactNode;
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
