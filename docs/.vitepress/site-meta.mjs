export function createSiteMeta({ base = '/' } = {}) {
  return {
    base,
    title: 'phenotype-go-kit',
    description: 'Documentation',
    themeConfig: {
      nav: [
        { text: 'Home', link: base || '/' },
      ],
    },
  }
}
