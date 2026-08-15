/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly UNFOLD_WEB_API_BASE?: string;
  readonly UNFOLD_WEB_HOST?: string;
  readonly UNFOLD_WEB_PORT?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
