import { test, expect } from '@playwright/test';
import { chromium } from 'playwright';

const BASE = process.env.BASE_URL || 'http://localhost:3000';

// Helper to simulate auth
async function login(page: any) {
  await page.goto(`${BASE}/`);
  await page.fill('[name="email"]', 'test@nebula-vpn.io');
  await page.fill('[name="password"]', 'testpass123');
  await page.click('button:has-text("Sign In")');
  await page.waitForURL('**/dashboard');
}

test.describe('User Panel - Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('shows connection ring', async ({ page }) => {
    await expect(page.locator('svg')).toBeVisible();
  });

  test('shows data usage stats', async ({ page }) => {
    await expect(page.locator('text=Download')).toBeVisible();
    await expect(page.locator('text=Upload')).toBeVisible();
  });

  test('shows recent sessions', async ({ page }) => {
    const sessions = page.locator('table');
    await expect(sessions).toBeVisible();
  });
});

test.describe('User Panel - Servers', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.click('text=Servers');
  });

  test('displays server list', async ({ page }) => {
    const servers = page.locator('[data-testid="server-card"]');
    await expect(servers.first()).toBeVisible();
  });

  test('filters by region', async ({ page }) => {
    await page.fill('[placeholder="Search"]', 'Frankfurt');
    const results = page.locator('[data-testid="server-card"]');
    await expect(results).toHaveCount(1);
  });

  test('shows latency bars', async ({ page }) => {
    const latencyBar = page.locator('[data-testid="latency-bar"]').first();
    await expect(latencyBar).toBeVisible();
  });
});

test.describe('User Panel - Devices', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.click('text=Devices');
  });

  test('adds new device', async ({ page }) => {
    await page.click('text=Add Device');
    await page.fill('[name="device-name"]', 'My MacBook');
    await page.selectOption('[name="platform"]', 'macOS');
    await page.click('text=Save');
    await expect(page.locator('text=My MacBook')).toBeVisible();
  });

  test('downloads config', async ({ page }) => {
    const downloadBtn = page.locator('[data-testid="download-config"]').first();
    await expect(downloadBtn).toBeVisible();
  });
});

test.describe('User Panel - Subscription', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.click('text=Subscription');
  });

  test('shows plan comparison', async ({ page }) => {
    await expect(page.locator('text=Free')).toBeVisible();
    await expect(page.locator('text=Pro')).toBeVisible();
    await expect(page.locator('text=Ultimate')).toBeVisible();
  });

  test('highlights current plan', async ({ page }) => {
    const currentPlan = page.locator('[data-testid="current-plan"]');
    await expect(currentPlan).toBeVisible();
    await expect(currentPlan).toContainText('Pro');
  });
});
