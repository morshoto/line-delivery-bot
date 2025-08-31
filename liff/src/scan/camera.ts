export async function startCamera() {
  return navigator.mediaDevices.getUserMedia({
    video: { facingMode: "environment" },
  });
}

export function stopCamera(stream: MediaStream) {
  stream.getTracks().forEach((t) => t.stop());
}
