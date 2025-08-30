declare global {
  interface Window {
    liff: any;
  }
}

import type { AppConfig } from './env';

export async function initLiff(cfg: AppConfig) {
  await window.liff?.init({ liffId: cfg.liffId });
  if (!window.liff?.isInClient()) {
    // ブラウザ起動時でも動くが、今回用途は in-app を想定
  }
}

export async function getGroupIdOrThrow(): Promise<string> {
  const ctx = window.liff!.getContext();
  if (ctx.type !== 'group' || !ctx.groupId) {
    throw new Error('グループから開いてください');
  }
  return ctx.groupId;
}

export async function getProfileSafe() {
  try {
    const p = await window.liff!.getProfile();
    return { displayName: p.displayName, userId: p.userId };
  } catch {
    return { displayName: '', userId: '' };
  }
}

export async function getAuthHeader(cfg: AppConfig): Promise<Record<string,string>> {
  if (cfg.oidcEnabled) {
    const idt = await window.liff!.getIDToken();
    return idt ? { Authorization: `Bearer ${idt}` } : {};
  }
  if (cfg.useSharedToken && cfg.sharedToken) {
    return { 'X-Shared-Token': cfg.sharedToken };
  }
  return {};
}
