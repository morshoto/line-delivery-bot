export async function loadZxing() {
  const worker = new Worker(
    new URL("../workers/zxing-worker.ts", import.meta.url),
  );
  return worker;
}
