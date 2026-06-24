---
title: Architecture
weight: 3
---

# Architecture

This document describes the internal modules and design patterns of the **HAR Analyzer**.

---

## Package Layout

The application codebase is structured to separate concerns between command execution, core logic, and presentation printers:

```text
├── main.go               # Project entry point
├── cmd/
│   └── root.go           # CLI flags registration and Cobra command setup
└── pkg/
    ├── har/
    │   ├── har.go        # Types and getters (HAR, Log, Entry)
    │   ├── filter.go     # Filter evaluation (URL, method, status matchers)
    │   ├── parser.go     # Stream reading and file saving functions
    │   └── *_test.go     # Package unit tests
    └── printer/
        ├── printer.go    # Output color configurations
        ├── table.go      # CLI table format rendering
        ├── json.go       # JSON summary rendering
        ├── csv.go        # CSV layout rendering
        └── detail.go     # Detailed header/body inspector rendering
```

---

## Custom Metadata Preservation

A typical challenge when parsing JSON in Go is that mapping keys to rigid structs forces undeclared fields to be discarded. 

Browser devtools frequently inject custom, browser-specific attributes into HAR records (such as `_initiator` for scripts call stacks, `_priority`, etc.). To prevent losing this metadata when filtering logs and writing them back, the HAR Analyzer represents entries as a dynamic map:

```go
type Entry map[string]interface{}
```

Then, idiomatic helper methods are attached to `Entry` to extract nested details safely:

```go
func (e Entry) GetMethod() string {
	req := e.GetRequest()
	if req == nil {
		return ""
	}
	if method, ok := req["method"].(string); ok {
		return method
	}
	return ""
}
```

This guarantees 100% data preservation of custom debugging metadata.

---

## Testing Strategy

Core filtering engines are unit tested in `pkg/har`:
- **Getter Testing**: Asserts parsing behaviors on simulated, mock HAR JSON payloads.
- **Filter Matrix Testing**: Asserts pattern matching rules for methods, wildcard status ranges (such as `2xx` or `5xx`), and URL regex expressions.
