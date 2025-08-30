export type AppConfig = {
  liffId: string;
  apiBase: string;
  useSharedToken: boolean;
  sharedToken: string;
  oidcEnabled: boolean;
  env?: 'prod'|'stg'|'dev';
};

let cache: AppConfig | null = null;

export async function loadConfig(): Promise<AppConfig> {
  if (cache) return cache;
  const res = await fetch('/config/config.json', { cache: 'no-store' });
  if (!res.ok) throw new Error('設定の取得に失敗しました');
  cache = await res.json();
  return cache!;
}
