export function normalize(text: string) {
  return text.normalize().replace(/\r\n|\r|\n/g, '').trim();
}
