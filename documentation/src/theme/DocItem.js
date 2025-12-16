import React from 'react';
import DocItem from '@theme-original/DocItem';
import StructuredData, { getDefaultStructuredData } from '@site/src/components/StructuredData';

export default function DocItemWithStructuredData(props) {
  const { content: DocContent } = props;
  const { metadata } = DocContent;
  
  // Prepare structured data based on page metadata
  const structuredData = getDefaultStructuredData({
    name: metadata.title,
    headline: metadata.title,
    description: metadata.description || metadata.frontMatter.description,
  });

  return (
    <>
      <StructuredData data={structuredData} />
      <DocItem {...props} />
    </>
  );
}