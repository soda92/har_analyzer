package har

import (
	"regexp"
	"strconv"
	"strings"
)

// FilterOpts holds all user-specified filtering parameters
type FilterOpts struct {
	Methods    []string
	URLPattern string
	Status     string
	Types      []string
}

// FilterEntries filters a slice of entries based on the given FilterOpts
func FilterEntries(entries []Entry, opts FilterOpts) []Entry {
	var filtered []Entry
	for _, entry := range entries {
		if MatchEntry(entry, opts) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// MatchEntry checks if an entry matches the filtering criteria
func MatchEntry(entry Entry, opts FilterOpts) bool {
	// 1. Method check
	if len(opts.Methods) > 0 {
		method := strings.ToUpper(entry.GetMethod())
		matched := false
		for _, m := range opts.Methods {
			if method == strings.ToUpper(m) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 2. URL check
	if opts.URLPattern != "" {
		urlStr := entry.GetURL()
		if !MatchURL(urlStr, opts.URLPattern) {
			return false
		}
	}

	// 3. Status check
	if opts.Status != "" {
		status := entry.GetStatus()
		if !MatchStatus(status, opts.Status) {
			return false
		}
	}

	// 4. Resource Type check
	if len(opts.Types) > 0 {
		hasAll := false
		for _, t := range opts.Types {
			if strings.ToLower(t) == "all" {
				hasAll = true
				break
			}
		}

		if !hasAll {
			resType := strings.ToLower(entry.GetResourceType())
			matched := false
			for _, t := range opts.Types {
				if resType == strings.ToLower(t) {
					matched = true
					break
				}
			}

			// Fallback checks
			if !matched {
				hasXhrOrFetch := false
				for _, t := range opts.Types {
					lt := strings.ToLower(t)
					if lt == "xhr" || lt == "fetch" {
						hasXhrOrFetch = true
						break
					}
				}
				if hasXhrOrFetch && entry.HasXHRHeader() {
					matched = true
				}
			}

			if !matched {
				return false
			}
		}
	}

	return true
}

// MatchURL checks if a URL matches the pattern (regex or substring)
func MatchURL(urlStr, pattern string) bool {
	if pattern == "" {
		return true
	}
	re, err := regexp.Compile(pattern)
	if err == nil {
		return re.MatchString(urlStr)
	}
	return strings.Contains(strings.ToLower(urlStr), strings.ToLower(pattern))
}

// MatchStatus checks if a status code matches the pattern (specific or class e.g., 2xx)
func MatchStatus(status int, pattern string) bool {
	if pattern == "" {
		return true
	}
	statusStr := strconv.Itoa(status)
	pattern = strings.ToLower(pattern)

	if len(pattern) == 3 && strings.HasSuffix(pattern, "xx") {
		return statusStr[0] == pattern[0]
	}

	return statusStr == pattern
}
