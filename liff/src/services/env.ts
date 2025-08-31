export type AppConfig = {
  liffId: string;
  apiBase: string;
  useSharedToken: boolean;
  sharedToken: string;
  oidcEnabled: boolean;
  env?: "prod" | "stg" | "dev";
};

let cache: AppConfig | null = null;

function toBool(v: unknown): boolean {
  if (typeof v === "boolean") return v;
  if (typeof v === "string") return /^(true|1|yes|on)$/i.test(v);
  return false;
}

export async function loadConfig(): Promise<AppConfig> {
  if (cache) return cache;
  // In Next.js client code, only NEXT_PUBLIC_* variables are inlined at build time.
  const cfg: AppConfig = {
    liffId: process.env.NEXT_PUBLIC_LIFF_ID ?? "",
    apiBase: process.env.NEXT_PUBLIC_API_BASE ?? "",
    useSharedToken: toBool(process.env.NEXT_PUBLIC_USE_SHARED_TOKEN),
    sharedToken: process.env.NEXT_PUBLIC_SHARED_TOKEN ?? "",
    oidcEnabled: toBool(process.env.NEXT_PUBLIC_OIDC_ENABLED),
    env: (process.env.NEXT_PUBLIC_APP_ENV ?? "dev") as "prod" | "stg" | "dev",
  };
  if (!cfg.liffId || cfg.liffId === "YOUR_LIFF_ID") {
    throw new Error(
      "LIFF ID が未設定です。liff/.env などに正しい LIFF ID を設定してください",
    );
  }
  cache = cfg;
  return cache;
}
