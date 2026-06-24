---
title: Getting Started
weight: 1
---

# Getting Started

Learn how to build, run, and verify the HAR Analyzer CLI tool on your local environment.

## Prerequisites

Before building the tool, ensure you have the following prerequisites installed:
- **Go Compiler**: Version `1.18` or higher (we developed it using Go `1.26`).
- **Git** (optional, for cloning the source code).

---

## Build from Source

To compile the self-contained executable binary from source, run:

1. Clone or navigate to the project directory:
   ```bash
   cd har_analyzer
   ```

2. Download Go module dependencies (such as Cobra CLI and fatih/color):
   ```bash
   go mod tidy
   ```

3. Build the binary:
   ```bash
   go build -o har_analyzer main.go
   ```

This creates an executable binary named `har_analyzer` in your current directory.

---

## Verifying the Installation

To verify that the CLI binary works correctly, print the version or help description:

```bash
./har_analyzer -h
```

You should see output similar to this:

```text
HAR Analyzer is a quick and modular CLI tool written in Go to filter
and analyze network requests (like XHR and fetch) exported in HTTP Archive (HAR) format.
It supports printing to tables, JSON summaries, CSV, or writing back to filtered HAR files.

Usage:
  har_analyzer [path/to/file.har] [flags]

Flags:
  -i, --file string     Input HAR file path (or '-' for stdin)
  -f, --format string   Output format: table, json, csv, har (default "table")
  -h, --help            help for har_analyzer
  -m, --method string   Filter by HTTP method(s) (comma-separated, e.g. GET,POST)
      --no-color        Disable colorized terminal output
  -o, --out string      Output filtered HAR file path
  -d, --show int        Inspect detailed headers and body of entry at index (default -1)
  -s, --status string   Filter by HTTP status code (e.g. 200, 4xx, 5xx)
  -t, --type string     Filter by resource type(s) (comma-separated, or 'all') (default "xhr,fetch")
  -u, --url string      Filter by URL (regex or substring match)
```
