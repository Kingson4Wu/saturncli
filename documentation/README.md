# Saturn CLI Documentation

This directory contains the documentation for the Saturn CLI project, built with Docusaurus and set up for GitHub Pages deployment.

## About Saturn CLI

Saturn CLI is a Go CLI client implementation based on the [VipShop Saturn](https://github.com/vipshop/Saturn) distributed task scheduling system. This project provides a lightweight client/server toolkit that enables you to trigger and monitor shell-style jobs from your Go services or the command line, with Unix domain sockets for efficient local communication.

## Overview

The documentation provides comprehensive information about Saturn CLI, including:

- Getting started guides
- API references
- Architecture documentation
- Usage examples
- Best practices
- Contributing guidelines

## Local Development

To run the documentation locally:

1. Install Node.js and npm
2. Install dependencies:
   ```bash
   cd documentation
   npm install
   ```

3. Run the local server:
   ```bash
   npm start
   ```

4. View the documentation at `http://localhost:3000`

## Structure

- `docs/` - Documentation markdown files
- `blog/` - Blog posts (if any)
- `src/` - Custom React components and CSS
- `static/` - Static assets
- `docusaurus.config.js` - Main Docusaurus configuration
- `sidebars.js` - Navigation sidebar configuration
- `package.json` - NPM dependencies and scripts

## Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The site is published at:

`https://kingson4wu.github.io/saturncli/`

GitHub Pages is configured to use GitHub Actions or similar to build and deploy the Docusaurus site.

## Contributing to Documentation

We welcome contributions to improve the documentation! To contribute:

1. Fork the repository
2. Create a feature branch
3. Make your changes in the appropriate markdown files in the `docs/` directory
4. Test locally before submitting
5. Submit a pull request

For more information about contributing, see the [contributing guidelines](./docs/contributing.md).

## License

The documentation is licensed under the same license as the Saturn CLI project. See the [LICENSE](../LICENSE) file for details.