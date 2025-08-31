import { cleanse, validate } from "./sanitize";
import { getAuthHeader } from "../services/liff";
import type { AppConfig } from "../services/env";

export type ScanPayload = {
  groupId: string;
  qrText: string;
  displayName: string;
  userId: string;
  meta?: Record<string, unknown>;
};

export async function postScan(cfg: AppConfig, p: ScanPayload) {
  const qrText = cleanse(p.qrText);
  validate(p.groupId, qrText);
  const headers = {
    "Content-Type": "application/json",
    ...(await getAuthHeader(cfg)),
  };

  const ctrl = new AbortController();
  const t = setTimeout(() => ctrl.abort(), 8000);

  try {
    const res = await fetch(`${cfg.apiBase}/api/scan`, {
      method: "POST",
      headers,
      body: JSON.stringify({ ...p, qrText }),
      signal: ctrl.signal,
    });
    if (!res.ok) throw new Error(`POST失敗: ${res.status}`);
    return await res.json();
  } finally {
    clearTimeout(t);
  }
}
