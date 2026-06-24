---
title: CLI Usage & Flags
weight: 2
---

# CLI Usage & Flags

HAR Analyzer provides powerful command-line flags to filter, format, and inspect network logs.

---

## Command Syntax

```bash
./har_analyzer [options] [path/to/file.har]
```

If no file path is specified, the CLI will look for input piped from standard input (`stdin`), allowing integration with other CLI tools.

---

## Filtering Flags

### Method Filter (`-m`, `--method`)
Filter requests by case-insensitive HTTP method name. You can supply multiple comma-separated methods:
```bash
./har_analyzer -m GET resources/gitee.com.har
./har_analyzer -m POST,PUT resources/gitee.com.har
```

### URL Filter (`-u`, `--url`)
Filters by matching the URL. If the query compiles as a valid regular expression, it will execute regex matching. Otherwise, it defaults to case-insensitive substring matching:
```bash
# Substring match
./har_analyzer -u "/api/v1" resources/gitee.com.har

# Regex match
./har_analyzer -u "graphql$" resources/gitee.com.har
```

### Status Code Filter (`-s`, `--status`)
Filters by response HTTP status code. Supports exact status numbers or wildcard classes (like `2xx`, `4xx`, `5xx`):
```bash
# Exact status code
./har_analyzer -s 200 resources/gitee.com.har

# Error classes
./har_analyzer -s 5xx resources/gitee.com.har
```

### Resource Type Filter (`-t`, `--type`)
Filters requests by their developer tools resource type. Defaults to `xhr,fetch`. You can override it to other standard types or set it to `all` to disable resource-type filtering:
```bash
# Get stylesheets and scripts instead of XHR
./har_analyzer -t stylesheet,script resources/gitee.com.har

# Include all types
./har_analyzer -t all resources/gitee.com.har
```

---

## Output Formats (`-f`, `--format`)

Use the `-f` / `--format` flag to control stdout representation:

| Format Name | Command Example | Output Description |
| :--- | :--- | :--- |
| `table` (Default) | `./har_analyzer -f table file.har` | Beautiful colorized table showing indexing, method, status, duration, type, and URL. |
| `json` | `./har_analyzer -f json file.har` | Array of parsed entry summaries containing basic metrics. |
| `csv` | `./har_analyzer -f csv file.har` | Comma-separated sheet representation. |
| `har` | `./har_analyzer -f har file.har` | Valid HAR JSON containing only the filtered entries. |

---

## Exporting & Inspecting

### Saving Filtered Logs (`-o`, `--out`)
Save the filtered HAR results directly to a new file. It prints a table summary of what was saved to stdout:
```bash
./har_analyzer -m POST -s 200 -o filtered.har resources/gitee.com.har
```

### Detailed Inspection (`-d`, `--show`)
Inspect full request headers, request query parameters, cookies, request body, response headers, and response body for a specific index from the filtered list:
```bash
./har_analyzer -d 2 resources/gitee.com.har
```

Features of the inspector:
*   **Pretty Printing**: JSON bodies (request and response) are parsed and indented automatically.
*   **Base64 Decoding**: If the response body is encoded in base64, it is decoded inline automatically.
*   **Buffer Truncation**: Massive payloads are truncated to prevent bloating the terminal screen.
