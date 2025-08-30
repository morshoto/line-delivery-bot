export type AppConfig = {
  liffId: string;
  apiBase: string;
  useSharedToken: boolean;
  sharedToken: string;
  oidcEnabled: boolean;
  env?: 'prod' | 'stg' | 'dev';
};

let cache: AppConfig | null = null;

function toBool(v: unknown): boolean {
  if (typeof v === 'boolean') return v;
  if (typeof v === 'string') return /^(true|1|yes|on)$/i.test(v);
  return false;
}

export async function loadConfig(): Promise<AppConfig> {
  if (cache) return cache;
  const env = import.meta.env as any;
  const cfg: AppConfig = {
    liffId: env.VITE_LIFF_ID ?? '',
    apiBase: env.VITE_API_BASE ?? '',
    useSharedToken: toBool(env.VITE_USE_SHARED_TOKEN),
    sharedToken: env.VITE_SHARED_TOKEN ?? '',
    oidcEnabled: toBool(env.VITE_OIDC_ENABLED),
    env: (env.VITE_APP_ENV ?? env.MODE ?? 'dev') as 'prod' | 'stg' | 'dev',
  };
  if (!cfg.liffId || cfg.liffId === 'YOUR_LIFF_ID') {
    throw new Error('VITE_LIFF_ID が未設定です。liff/.env などに正しい LIFF ID を設定してください');
  }
  cache = cfg;
  return cache;
}
