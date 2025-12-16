import React from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  
  return (
    <Layout
      title={`Welcome to ${siteConfig.title}`}
      description="Documentation for Saturn CLI - A lightweight client/server toolkit for job execution">
      <main className="container margin-vert--lg">
        <div className="row">
          <div className="col col--8 col--offset--2">
            <div className="text--center padding-horiz--md">
              <h1>{siteConfig.title}</h1>
              <p className="hero__subtitle">{siteConfig.tagline}</p>
              
              <div className="margin-vert--lg">
                <Link className="button button--primary button--lg margin-horiz--sm" to="/saturncli/docs/getting-started">
                  Get Started
                </Link>
                <Link className="button button--secondary button--lg margin-horiz--sm" to="/saturncli/docs/intro">
                  Learn More
                </Link>
              </div>
              
              <div className="margin-vert--lg">
                <h2>What is Saturn CLI?</h2>
                <p>
                  Saturn CLI is a lightweight client/server toolkit that lets you trigger and monitor shell-style jobs
                  from your Go services or the command line. The client communicates with a long-running daemon over
                  Unix domain sockets (macOS/Linux) or HTTP (Windows), providing a fast, secure channel for
                  orchestrating background work.
                </p>
                <p>
                  <strong>About this project:</strong> Saturn CLI is a Go CLI client implementation based on the original
                  <a href="https://github.com/vipshop/Saturn" target="_blank" rel="noopener noreferrer"> VipShop Saturn distributed task scheduling system</a>.
                  This project provides a lightweight alternative for executing jobs from Go applications using Unix domain sockets for efficient local communication.
                </p>
              </div>
              
              <div className="margin-vert--lg">
                <h3>Quick Navigation</h3>
                <div className="row">
                  <div className="col">
                    <Link to="/saturncli/docs/getting-started">
                      <h4>Getting Started</h4>
                      <p>New to Saturn CLI? Start here</p>
                    </Link>
                  </div>
                  <div className="col">
                    <Link to="/saturncli/docs/architecture">
                      <h4>Architecture</h4>
                      <p>Understand the system design</p>
                    </Link>
                  </div>
                  <div className="col">
                    <Link to="/saturncli/docs/examples">
                      <h4>Examples</h4>
                      <p>Practical usage scenarios</p>
                    </Link>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </Layout>
  );
}