import { test, expect } from '@playwright/test'

const BASE_URL = process.env.BASE_URL || 'http://localhost:5173'

test.describe('docs', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(BASE_URL)
  })

  test('homepage loads', async ({ page }) => {
    await expect(page.locator('body')).toBeVisible()
  })
  test('route /guide/index loads', async ({ page }) => {
    await page.goto(BASE_URL + '/guide/index')
    await expect(page.locator('body')).toBeVisible()
  })
  test('route /index loads', async ({ page }) => {
    await page.goto(BASE_URL + '/index')
    await expect(page.locator('body')).toBeVisible()
  })
  test('route /zh-CN loads', async ({ page }) => {
    await page.goto(BASE_URL + '/zh-CN')
    await expect(page.locator('body')).toBeVisible()
  })
  test('route /zh-TW loads', async ({ page }) => {
    await page.goto(BASE_URL + '/zh-TW')
    await expect(page.locator('body')).toBeVisible()
  })
  test('route /fa loads', async ({ page }) => {
    await page.goto(BASE_URL + '/fa')
    await expect(page.locator('body')).toBeVisible()
  })
  test('route /fa-Latn loads', async ({ page }) => {
    await page.goto(BASE_URL + '/fa-Latn')
    await expect(page.locator('body')).toBeVisible()
  })
})
