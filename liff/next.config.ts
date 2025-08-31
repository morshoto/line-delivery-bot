import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  reactStrictMode: true,
  // Ensure we only render LIFF-dependent code on the client
  // (Components using LIFF should be marked with 'use client')
  typedRoutes: false,
};

export default nextConfig;
