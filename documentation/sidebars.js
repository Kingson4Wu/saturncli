// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  tutorialSidebar: [
    'intro',
    {
      type: 'category',
      label: 'Getting Started',
      items: [
        'getting-started',
        'getting-started/setup',
        'getting-started/quick-start',
        'getting-started/examples'
      ],
    },
    {
      type: 'category',
      label: 'Core Concepts',
      items: [
        'architecture',
        'client-api',
        'server-api',
        'cli-reference'
      ],
    },
    {
      type: 'category',
      label: 'Guides',
      items: [
        'embedding',
        'examples',
        'troubleshooting',
        'best-practices'
      ],
    },
    {
      type: 'category',
      label: 'Development',
      items: [
        'contributing',
        'development-setup',
        'testing'
      ],
    }
  ],
};

module.exports = sidebars;