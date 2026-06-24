---
title: HAR Analyzer Docs
---

# HAR Analyzer

Welcome to the documentation site for the **HAR Analyzer**!

HAR Analyzer is a fast, modular, and dependency-free command-line tool written in Go to filter, inspect, and analyze XHR and Fetch requests in HTTP Archive (HAR) files. 

It keeps all original custom browser metadata (such as custom debug keys, initiator stacks, and priority parameters) completely intact when exporting filtered logs back to a new HAR file.

<div style="margin: 2rem 0;">
  <a href="./docs/" style="display: inline-block; padding: 0.75rem 1.5rem; background-color: #2563eb; color: #ffffff; border-radius: 0.375rem; font-weight: 600; text-decoration: none; box-shadow: 0 4px 6px -1px rgba(37, 99, 235, 0.4);">
    Read the Docs →
  </a>
</div>

## Key Features

- 🏎️ **Fast & Modular**: Built with a clean Go standard library core, structured via `cobra`, and styled using `fatih/color`.
- 🔍 **Granular Filtering**: Filter by HTTP method, URL pattern (regex or substring), status code (with wildcard support e.g. `200` or `5xx`), and resource type.
- 💾 **Data Integrity**: Dynamically parses HAR entries to guarantee no browser-specific debugging extensions or metadata are discarded when saving.
- 🔬 **Rich Inline Inspector**: Print headers, cookies, and pretty-printed query/response payloads (with base64 decoding and JSON indentation built-in).
