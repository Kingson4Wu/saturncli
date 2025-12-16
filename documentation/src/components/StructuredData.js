import React from 'react';

// Component to add JSON-LD structured data to pages
const StructuredData = ({ data }) => {
  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(data) }}
    />
  );
};

// Default structured data for Saturn CLI
export const getDefaultStructuredData = (pageData = {}) => {
  const defaultData = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    "name": "Saturn CLI",
    "headline": pageData.title || "Saturn CLI - Go Job Execution Toolkit",
    "description": pageData.description || "A lightweight client/server toolkit for job execution in Go applications. Execute shell-style jobs, manage background processes, and enable efficient inter-process communication with Unix domain sockets.",
    "applicationCategory": "DeveloperApplication",
    "operatingSystem": ["Linux", "MacOS", "Windows"],
    "offers": {
      "@type": "Offer",
      "price": "0",
      "priceCurrency": "USD"
    },
    "author": {
      "@type": "Person",
      "name": "Kingson Wu",
      "url": "https://github.com/Kingson4Wu"
    },
    "softwareVersion": "1.0.0",
    "isAccessibleForFree": true,
    "softwareHelp": {
      "@type": "CreativeWork",
      "name": "Saturn CLI Documentation",
      "url": "https://kingson4wu.github.io/saturncli/"
    },
    "featureList": [
      "Job execution",
      "Background processing",
      "Unix domain sockets communication",
      "Inter-process communication",
      "Process management",
      "CLI automation",
      "Client-server architecture"
    ]
  };

  return { ...defaultData, ...pageData };
};

export default StructuredData;