# HAR Analyzer

A fast, lightweight, and dependency-free command-line tool written in Go to filter and inspect XHR/fetch requests in HTTP Archive (HAR) files.

It preserves all original custom browser metadata (such as custom keys, initiator stacks, and priority parameters) when exporting filtered logs back to a new HAR file.

## Features

- **Fast & Dependency-Free**: Built purely using the Go standard library.
- **Advanced Filtering**:
  - Filter by HTTP Method (e.g., `GET`, `POST`, `PUT`, `DELETE`).
  - Filter by Response Status Code (e.g., specific status `200`, or class wildcards like `4xx` and `5xx`).
  - Filter by URL (substring matching or regular expressions).
  - Filter by Resource Type (defaults to `xhr` and `fetch`, but customizable to compile stylesheets, scripts, images, etc., or `all`).
- **Flexible Outputs**:
  - **Table Format**: Visual colorized CLI table with response speed highlights.
  - **JSON Format**: Export a clean JSON array of entry summaries.
  - **CSV Format**: Save summaries to a CSV file for analytical tooling.
  - **HAR Format**: Generate a valid, clean HTTP Archive file containing only the filtered entries.
- **Detailed Entry Inspection**: Inspect specific requests to view detailed headers, query parameters, cookies, and pretty-printed request/response bodies (supports auto-decoding of Base64 responses and pretty-printing of JSON payloads).

---

## Installation & Building

To build the executable from source, ensure you have Go installed (v1.16+), then run:

```bash
go build -o har_analyzer main.go
```

This will generate a self-contained executable binary named `har_analyzer` in your current directory.

---

## Usage Examples

### 1. View all XHR and Fetch requests (default)
Prints a beautiful colorized table of all `xhr` and `fetch` requests found in the HAR file:
```bash
./har_analyzer resources/gitee.com.har
```

### 2. Filter by status code
Get only requests that returned a status of `200`:
```bash
./har_analyzer -s 200 resources/gitee.com.har
```

Get only server error requests (`500`, `502`, `503`, etc.):
```bash
./har_analyzer -s 5xx resources/gitee.com.har
```

### 3. Filter by HTTP method
Get only `POST` requests:
```bash
./har_analyzer -m POST resources/gitee.com.har
```

Get both `POST` and `PUT` requests:
```bash
./har_analyzer -m POST,PUT resources/gitee.com.har
```

### 4. Filter by URL pattern
Filter for requests whose URL contains the string `/graphql`:
```bash
./har_analyzer -u "/graphql" resources/gitee.com.har
```

Filter using a regular expression:
```bash
./har_analyzer -u "api/v[0-9]/users" resources/gitee.com.har
```

### 5. Combine filters and save to a new HAR file
Filter all `POST` requests with status `200 OK` from the input file, save the output as a new HAR file called `filtered.har`, and show a table summary of what was saved:
```bash
./har_analyzer -m POST -s 200 -o filtered.har resources/gitee.com.har
```

### 6. Export to other formats
Output a summary of filtered requests in JSON format (ideal for piping to tools like `jq`):
```bash
./har_analyzer -f json resources/gitee.com.har
```

Output as a CSV list:
```bash
./har_analyzer -f csv resources/gitee.com.har > output.csv
```

### 7. Detailed inspection of a specific request
To view all request and response headers, query params, cookies, and pretty-printed bodies for a particular entry, use the `-d` (or `-show`) flag with the index of the request from the table output:
```bash
./har_analyzer -d 2 resources/gitee.com.har
```

---

## Command Line Arguments Reference

| Option | Alias | Description | Default |
| :--- | :--- | :--- | :--- |
| `-file <path>` | `-i` | Path to the input HAR file (can also be passed as a positional arg, or `-` for stdin) | - |
| `-out <path>` | `-o` | Path to save the filtered HAR file | - |
| `-method <list>` | `-m` | Filter by HTTP method(s) (comma-separated, case-insensitive) | - |
| `-url <pattern>` | `-u` | Filter by URL pattern (regex or substring match) | - |
| `-status <pattern>`| `-s` | Filter by HTTP status code (e.g. `200`, `4xx`, `5xx`) | - |
| `-type <list>` | `-t` | Filter by resource type(s) (comma-separated, or `all` for all types) | `xhr,fetch` |
| `-format <format>` | `-f` | Output format: `table`, `json`, `csv`, `har` | `table` |
| `-show <index>` | `-d` | Inspect detailed headers and body of the entry at 0-based index | -1 (disabled) |
| `-no-color` | - | Disable colorized terminal outputs | `false` |
