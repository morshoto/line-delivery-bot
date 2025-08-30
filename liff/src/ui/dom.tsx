import { createRoot } from 'react-dom/client';
import type { ReactNode } from 'react';

export function showToast(message: string, timeout = 2000) {
  const div = document.createElement('div');
  div.textContent = message;
  div.style.position = 'fixed';
  div.style.bottom = '1rem';
  div.style.left = '50%';
  div.style.transform = 'translateX(-50%)';
  div.style.background = '#333';
  div.style.color = '#fff';
  div.style.padding = '0.5rem 1rem';
  div.style.borderRadius = '4px';
  document.body.appendChild(div);
  setTimeout(() => div.remove(), timeout);
}

export function renderDialog(node: ReactNode) {
  const container =
    document.getElementById('modal') || document.body.appendChild(document.createElement('div'));
  const root = createRoot(container);
  root.render(node);
  return () => root.unmount();
}
