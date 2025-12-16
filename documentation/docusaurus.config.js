// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Saturn CLI - Go Job Execution Toolkit | Unix Domain Sockets & Background Jobs',
  tagline: 'A Go CLI client implementation based on the VipShop Saturn distributed task scheduling system. Execute shell-style jobs, manage background processes, and enable efficient inter-process communication with Unix domain sockets.',
  favicon: 'img/favicon.svg',

  // Set the production url of your site here
  url: 'https://kingson4wu.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/saturncli/',

  // GitHub pages deployment config.
  organizationName: 'Kingson4Wu', // Usually your GitHub org/user name.
  projectName: 'saturncli', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Update deprecated markdown options
  markdown: {
    format: 'mdx',
    mermaid: false,
    // Move deprecated options here for Docusaurus v3
  },

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/Kingson4Wu/saturncli/tree/main/documentation/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/Kingson4Wu/saturncli/tree/main/documentation/',
          onUntruncatedBlogPosts: 'ignore', // Add this to suppress the truncation warning
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        sitemap: {
          changefreq: 'weekly',
          priority: 0.5,
          filename: 'sitemap.xml',
        },
      }),
    ],
  ],

  plugins: [
    [
      require.resolve('@docusaurus/plugin-google-gtag'),
      {
        trackingID: 'G-XXXXXXXXXX',
        anonymizeIP: true,
      },
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/docusaurus-social-card.png',
      metadata: [
        {name: 'keywords', content: 'saturn cli, go cli, job execution, background jobs, unix domain sockets, ipc, inter-process communication, golang, cli toolkit, server client, task runner, job scheduler, automation, embedded cli, go library, unix sockets, process management, background processing'},
        {name: 'description', content: 'Saturn CLI - A lightweight client/server toolkit for job execution in Go applications. Execute shell-style jobs, manage background processes, and enable inter-process communication with Unix domain sockets.'},
        {name: 'author', content: 'Kingson Wu'},
        {name: 'robots', content: 'index, follow'},
        {name: 'theme-color', content: '#25c2a0'},
        {property: 'og:title', content: 'Saturn CLI Documentation - Go Job Execution Toolkit'},
        {property: 'og:description', content: 'A lightweight client/server toolkit for job execution in Go. Execute background jobs, manage processes, and enable efficient inter-process communication.'},
        {property: 'og:image', content: 'img/docusaurus-social-card.jpg'},
        {property: 'og:url', content: 'https://kingson4wu.github.io/saturncli/'},
        {property: 'og:type', content: 'website'},
        {name: 'twitter:card', content: 'summary_large_image'},
        {name: 'twitter:title', content: 'Saturn CLI Documentation - Go Job Execution Toolkit'},
        {name: 'twitter:description', content: 'A lightweight client/server toolkit for job execution in Go applications. Execute shell-style jobs and manage background processes.'},
        {name: 'twitter:image', content: 'img/docusaurus-social-card.jpg'},
      ],
      navbar: {
        title: 'Saturn CLI',
        logo: {
          alt: 'Saturn CLI Logo',
          src: 'img/logo.svg',
          width: 32,
          height: 32,
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'tutorialSidebar',
            position: 'left',
            label: 'Documentation',
          },
          {to: '/blog', label: 'Blog', position: 'left'},
          {
            href: 'https://github.com/Kingson4Wu/saturncli',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Getting Started',
                to: '/docs/getting-started',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub Discussions',
                href: 'https://github.com/Kingson4Wu/saturncli/discussions',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/Kingson4Wu/saturncli',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Saturn CLI, Inc. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;