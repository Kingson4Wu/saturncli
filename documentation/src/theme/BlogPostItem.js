import React from 'react';
import BlogPostItem from '@theme-original/BlogPostItem';
import StructuredData from '@site/src/components/StructuredData';

export default function BlogPostItemWithStructuredData(props) {
  const { frontMatter, metadata } = props;

  // Prepare structured data for blog post
  const structuredData = {
    "@context": "https://schema.org",
    "@type": "BlogPosting",
    "mainEntityOfPage": {
      "@type": "WebPage",
      "@id": metadata?.permalink || ""
    },
    "headline": frontMatter?.title || metadata?.title,
    "description": frontMatter?.description || metadata?.description,
    "image": frontMatter?.image || "/img/docusaurus-social-card.jpg",
    "author": {
      "@type": "Organization",
      "name": "Saturn CLI Team"
    },
    "publisher": {
      "@type": "Organization",
      "name": "Saturn CLI",
      "logo": {
        "@type": "ImageObject",
        "url": "/img/logo.svg"
      }
    },
    "datePublished": frontMatter?.date,
    "dateModified": frontMatter?.date
  };

  return (
    <>
      <StructuredData data={structuredData} />
      <BlogPostItem {...props} />
    </>
  );
}