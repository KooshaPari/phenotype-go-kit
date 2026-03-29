import assert from 'node:assert/strict'
import test from 'node:test'
import { createSiteMeta } from '../../.vitepress/site-meta.mjs'

test('site meta exports createSiteMeta', () => {
  const meta = createSiteMeta({ base: '/' })
  assert.ok(meta.title)
  assert.ok(meta.themeConfig)
  assert.ok(meta.themeConfig.nav)
  assert.ok(Array.isArray(meta.themeConfig.nav))
})

test('site meta nav has expected items', () => {
  const meta = createSiteMeta({ base: '/' })
  const links = meta.themeConfig.nav.map(n => n.link)
  assert.ok(links.length > 0, 'nav should not be empty')
})
