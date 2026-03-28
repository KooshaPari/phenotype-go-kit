import { withMermaid } from 'vitepress-plugin-mermaid'

export default withMermaid({
  title: 'phenotype-go-kit',
  description: 'Shared Go toolkit for the Phenotype ecosystem — utilities, middleware, and infrastructure primitives.',
  appearance: 'dark',
  lastUpdated: true,
  themeConfig: {
    nav: [{ text: 'Home', link: '/' }],
    sidebar: [],
    search: { provider: 'local' },
  },
  mermaid: { theme: 'dark' },
})
