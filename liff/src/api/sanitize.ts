export function cleanse(s: string, max = 2048) {
  const noCtl = s.replace(/[\u0000-\u001F\u007F]/g, "");
  return noCtl.trim().slice(0, max);
}

export function validate(groupId: string, qr: string) {
  if (groupId.length < 1 || groupId.length > 64) throw new Error("groupId不正");
  if (qr.length < 1 || qr.length > 2048) throw new Error("qrText不正");
}
