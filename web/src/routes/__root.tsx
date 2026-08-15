import {Outlet, createRootRoute} from '@tanstack/react-router';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';

import {AppHeader} from '@/components/layout/page';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60_000,
      retry: 1,
    },
  },
});

export const Route = createRootRoute({
  component: RootLayout,
});

function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="flex min-h-svh flex-col bg-background">
        <AppHeader />
        <div className="flex flex-1 flex-col">
          <Outlet />
        </div>
      </div>
    </QueryClientProvider>
  );
}
