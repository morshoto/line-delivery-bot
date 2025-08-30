export async function detectBarcode(video: HTMLVideoElement): Promise<string | null> {
  const anyWindow = window as any;
  if ('BarcodeDetector' in anyWindow) {
    const detector = new anyWindow.BarcodeDetector({ formats: ['qr_code'] });
    const codes = await detector.detect(video);
    return codes[0]?.rawValue ?? null;
  }
  return null;
}
