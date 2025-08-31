const marks: Record<string, number> = {};

export function mark(label: "t0" | "t1" | "t2") {
  marks[label] = performance.now();
}

export function getTimings() {
  return marks;
}
