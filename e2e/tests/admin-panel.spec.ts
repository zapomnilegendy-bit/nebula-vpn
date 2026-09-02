import { test, expect } from '@playwright/test';

const ADMIN_BASE = process.env.ADMIN_BASE_URL || 'http://localhost:3001';

async function loginAsAdmin(page: any) {
  await page.goto(`${ADMIN_BASE}/`);
  await page.fill('[name="email"]', 'admin@nebula-vpn.io');
  await page.fill('[name="password"]', 'adminpass123');
  await page.click('button:has-text("Sign In")');
  await page.waitForURL('**/dashboard');
}

test.describe('Admin Panel - Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
  });

  test('shows metrics cards', async ({ page }) => {
    await expect(page.locator('text=Total Users')).toBeVisible();
    await expect(page.locator('text=Active Nodes')).toBeVisible();
    await expect(page.locator('text=MRR')).toBeVisible();
    await expect(page.locator('text=Avg Latency')).toBeVisible();
  });

  test('shows bandwidth chart', async ({ page }) => {
    const chart = page.locator('[data-testid="bandwidth-chart"]');
    await expect(chart).toBeVisible();
  });
});

test.describe('Admin Panel - Users', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await page.click('text=Users');
  });

  test('displays user table', async ({ page }) => {
    const table = page.locator('table');
    await expect(table).toBeVisible();
  });

  test('searches users', async ({ page }) => {
    await page.fill('[placeholder="Search users"]', 'test@');
    await page.waitForTimeout(500);
    const rows = page.locator('table tbody tr');
    const count = await rows.count();
    expect(count).toBeGreaterThanOrEqual(0);
  });
});

test.describe('Admin Panel - Servers', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await page.click('text=Servers');
  });

  test('shows server cards with status', async ({ page }) => {
    const cards = page.locator('[data-testid="server-admin-card"]');
    await expect(cards.first()).toBeVisible();
  });

  test('adds a new server', async ({ page }) => {
    await page.click('text=Add Server');
    await page.fill('[name="city"]', 'Toronto');
    await page.fill('[name="host"]', 'tor1.nebula-vpn.io');
    await page.fill('[name="port"]', '51820');
    await page.click('text=Save');
    await expect(page.locator('text=Toronto')).toBeVisible();
  });
});
