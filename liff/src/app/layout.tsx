export const metadata = {
  title: 'LIFF App',
  description: 'LINE LIFF app',
};

import '../styles/globals.css';
import '../styles/main.css';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
