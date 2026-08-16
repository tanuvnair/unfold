import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {defineConfig, loadEnv} from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import tanstackRouter from '@tanstack/router-plugin/vite';

const rootDir = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(rootDir, '..');

export default defineConfig(({mode}) => {
  const rootEnv = loadEnv(mode, repoRoot, '');
  const webEnv = loadEnv(mode, rootDir, '');
  const env = {...rootEnv, ...webEnv};

  const apiHost = env.UNFOLD_API_HOST || '127.0.0.1';
  const apiPort = env.UNFOLD_API_PORT || '8080';
  const apiProxyTarget =
    env.UNFOLD_API_PROXY_TARGET ||
    env.UNFOLD_API_ADDR ||
    `http://${apiHost}:${apiPort}`;

  const webHost = env.UNFOLD_WEB_HOST || 'localhost';
  const webPort = Number(env.UNFOLD_WEB_PORT || env.PORT || 5173);

  // Empty = same-origin /api (dev proxy or production behind one host).
  const webAPIBase = (
    env.UNFOLD_WEB_API_BASE ||
    env.VITE_API_BASE ||
    ''
  ).replace(/\/$/, '');

  return {
    // Load UNFOLD_* / VITE_* from the monorepo root .env into the client bundle.
    envDir: repoRoot,
    envPrefix: ['VITE_', 'UNFOLD_WEB_'],
    plugins: [
      tanstackRouter({
        target: 'react',
        autoCodeSplitting: true,
      }),
      react(),
      tailwindcss(),
    ],
    resolve: {
      alias: {
        '@': path.resolve(rootDir, './src'),
      },
    },
    define: {
      // Ensure build picks up UNFOLD_WEB_API_BASE even if only set in shell/make.
      'import.meta.env.UNFOLD_WEB_API_BASE': JSON.stringify(webAPIBase),
    },
    server: {
      host: webHost,
      port: webPort,
      // Only used when the browser calls relative /api (UNFOLD_WEB_API_BASE empty).
      proxy: {
        '/api': {
          target: apiProxyTarget.startsWith('http')
            ? apiProxyTarget
            : `http://${apiProxyTarget}`,
          changeOrigin: true,
        },
      },
    },
  };
});
