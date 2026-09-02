import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  timeout: 30000,
  retries: 1,
  use: {
    headless: true,
    screenshot: 'only-on-failure',
    trace: 'retain-on-failure',
  },
  projects: [
    {
      name: 'User Panel',
      use: { baseURL: process.env.BASE_URL || 'http://localhost:3000' },
    },
    {
      name: 'Admin Panel',
      use: { baseURL: process.env.ADMIN_BASE_URL || 'http://localhost:3001' },
    },
  ],
});
