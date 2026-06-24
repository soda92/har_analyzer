package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ANSI color escape sequences
var (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
)

// HAR represents the structure of the HTTP Archive file
type HAR struct {
	Log Log `json:"log"`
}

// Log represents the log field inside HAR
type Log struct {
	Version string                   `json:"version"`
	Creator interface{}              `json:"creator,omitempty"`
	Browser interface{}              `json:"browser,omitempty"`
	Pages   interface{}              `json:"pages,omitempty"`
	Entries []map[string]interface{} `json:"entries"`
	Comment string                   `json:"comment,omitempty"`
}

// EntrySummary is used for printing structured JSON summaries
type EntrySummary struct {
	Index        int     `json:"index"`
	Method       string  `json:"method"`
	URL          string  `json:"url"`
	Status       int     `json:"status"`
	StatusText   string  `json:"statusText"`
	Time         float64 `json:"time_ms"`
	ResourceType string  `json:"resourceType"`
}

// FilterOpts holds all user-specified filtering parameters
type FilterOpts struct {
	Methods    []string
	URLPattern string
	Status     string
	Types      []string
}

func initColors(noColor bool) {
	if noColor {
		colorReset = ""
		colorBold = ""
		colorRed = ""
		colorGreen = ""
		colorYellow = ""
		colorBlue = ""
		colorMagenta = ""
		colorCyan = ""
		colorGray = ""
	}
}

func main() {
	var (
		fileOpt    string
		iOpt       string
		outOpt     string
		oOpt       string
		methodOpt  string
		mOpt       string
		urlOpt     string
		uOpt       string
		statusOpt  string
		sOpt       string
		typeOpt    string
		tOpt       string
		formatOpt  string
		fOpt       string
		showOpt    int
		dOpt       int
		noColorOpt bool
	)

	flag.StringVar(&fileOpt, "file", "", "Input HAR file path")
	flag.StringVar(&iOpt, "i", "", "Input HAR file path (alias)")
	flag.StringVar(&outOpt, "out", "", "Output filtered HAR file path")
	flag.StringVar(&oOpt, "o", "", "Output filtered HAR file path (alias)")
	flag.StringVar(&methodOpt, "method", "", "Filter by HTTP method (comma-separated, e.g. GET,POST)")
	flag.StringVar(&mOpt, "m", "", "Filter by HTTP method (alias)")
	flag.StringVar(&urlOpt, "url", "", "Filter by URL (regex or substring match)")
	flag.StringVar(&uOpt, "u", "", "Filter by URL (alias)")
	flag.StringVar(&statusOpt, "status", "", "Filter by HTTP status code (e.g. 200, 4xx, 5xx)")
	flag.StringVar(&sOpt, "s", "", "Filter by HTTP status (alias)")
	flag.StringVar(&typeOpt, "type", "", "Filter by resource type (comma-separated, default: xhr,fetch)")
	flag.StringVar(&tOpt, "t", "", "Filter by resource type (alias)")
	flag.StringVar(&formatOpt, "format", "", "Output format: table, json, csv, har (default: table)")
	flag.StringVar(&fOpt, "f", "", "Output format (alias)")
	flag.IntVar(&showOpt, "show", -1, "Inspect detailed headers and body of the entry at index")
	flag.IntVar(&dOpt, "d", -1, "Inspect detailed headers and body (alias)")
	flag.BoolVar(&noColorOpt, "no-color", false, "Disable colorized output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%sHAR Analyzer - A fast Go tool to filter and inspect XHR/fetch requests in HAR files%s\n\n", colorBold, colorReset)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  har_analyzer [options] [path/to/file.har]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -i, -file <path>       Path to the input HAR file (can also be passed as positional argument or '-' for stdin)\n")
		fmt.Fprintf(os.Stderr, "  -o, -out <path>        Path to write the filtered HAR file\n")
		fmt.Fprintf(os.Stderr, "  -m, -method <methods>  Filter by HTTP method(s) (comma-separated: GET,POST,etc.)\n")
		fmt.Fprintf(os.Stderr, "  -u, -url <pattern>     Filter by URL pattern (regex or substring match)\n")
		fmt.Fprintf(os.Stderr, "  -s, -status <pattern>  Filter by HTTP status code (e.g. 200, 4xx, 5xx)\n")
		fmt.Fprintf(os.Stderr, "  -t, -type <types>      Filter by resource type(s) (comma-separated, default: xhr,fetch; use 'all' for all)\n")
		fmt.Fprintf(os.Stderr, "  -f, -format <format>   Output format: table, json, csv, har (default: table)\n")
		fmt.Fprintf(os.Stderr, "  -d, -show <index>      Inspect detailed headers and body of the entry at 0-based index\n")
		fmt.Fprintf(os.Stderr, "  -no-color              Disable ANSI colors in terminal output\n")
	}

	flag.Parse()

	// Initialize color options
	initColors(noColorOpt)

	// Resolve aliases
	filePath := getFlagString(fileOpt, iOpt)
	if filePath == "" && len(flag.Args()) > 0 {
		filePath = flag.Arg(0)
	}

	outPath := getFlagString(outOpt, oOpt)

	var methods []string
	methodVal := getFlagString(methodOpt, mOpt)
	if methodVal != "" {
		parts := strings.Split(methodVal, ",")
		for _, p := range parts {
			if p = strings.TrimSpace(p); p != "" {
				methods = append(methods, p)
			}
		}
	}

	urlPattern := getFlagString(urlOpt, uOpt)
	statusPattern := getFlagString(statusOpt, sOpt)

	typeVal := getFlagString(typeOpt, tOpt)
	if typeVal == "" {
		typeVal = "xhr,fetch" // default resource types
	}
	var resourceTypes []string
	parts := strings.Split(typeVal, ",")
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			resourceTypes = append(resourceTypes, p)
		}
	}

	formatVal := getFlagString(formatOpt, fOpt)
	if formatVal == "" {
		formatVal = "table"
	}

	inspectIndex := getFlagInt(showOpt, dOpt)

	// Read and parse HAR file
	har, err := readHAR(filePath)
	if err != nil {
		log.Fatalf("Error reading HAR file: %v", err)
	}

	// Filter entries
	opts := FilterOpts{
		Methods:    methods,
		URLPattern: urlPattern,
		Status:     statusPattern,
		Types:      resourceTypes,
	}
	filtered := filterEntries(har.Log.Entries, opts)

	// Inspect single request detail if requested
	if inspectIndex >= 0 {
		if inspectIndex >= len(filtered) {
			log.Fatalf("Index %d is out of range. There are only %d filtered entries (0-%d).", 
				inspectIndex, len(filtered), len(filtered)-1)
		}
		printDetailedEntry(filtered[inspectIndex], inspectIndex)
		return
	}

	// Write to file if specified
	if outPath != "" {
		err := writeHAR(har, filtered, outPath)
		if err != nil {
			log.Fatalf("Error saving filtered HAR to %s: %v", outPath, err)
		}
		fmt.Printf("%sSaved %d filtered entries to %s%s\n\n", colorGreen, len(filtered), outPath, colorReset)
	}

	// Print to stdout in requested format
	if outPath == "" || strings.ToLower(formatVal) != "har" {
		switch strings.ToLower(formatVal) {
		case "table":
			printTable(filtered)
		case "json":
			printJSON(filtered)
		case "csv":
			printCSV(filtered)
		case "har":
			if err := writeHAR(har, filtered, ""); err != nil {
				log.Fatalf("Error printing HAR: %v", err)
			}
		default:
			log.Fatalf("Unknown format: %s. Supported formats: table, json, csv, har", formatVal)
		}
	}
}

func getFlagString(long, short string) string {
	if long != "" {
		return long
	}
	return short
}

func getFlagInt(long, short int) int {
	if long != -1 {
		return long
	}
	return short
}

func readHAR(filePath string) (*HAR, error) {
	var reader io.Reader
	if filePath == "" || filePath == "-" {
		// Read from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return nil, fmt.Errorf("no input HAR file specified and no stdin pipe detected (run with -h for help)")
		}
		reader = os.Stdin
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var har HAR
	if err := json.Unmarshal(data, &har); err != nil {
		return nil, err
	}

	if har.Log.Entries == nil {
		return nil, fmt.Errorf("invalid HAR structure: log.entries is missing")
	}

	return &har, nil
}

func writeHAR(har *HAR, filteredEntries []map[string]interface{}, outputPath string) error {
	outputHar := HAR{
		Log: Log{
			Version: har.Log.Version,
			Creator: har.Log.Creator,
			Browser: har.Log.Browser,
			Pages:   har.Log.Pages,
			Entries: filteredEntries,
			Comment: har.Log.Comment,
		},
	}

	jsonData, err := json.MarshalIndent(outputHar, "", "  ")
	if err != nil {
		return err
	}

	if outputPath == "" {
		_, err = os.Stdout.Write(jsonData)
		return err
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

func filterEntries(entries []map[string]interface{}, opts FilterOpts) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, entry := range entries {
		if matchEntry(entry, opts) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func matchEntry(entry map[string]interface{}, opts FilterOpts) bool {
	// 1. Method check
	if len(opts.Methods) > 0 {
		method := strings.ToUpper(getMethod(entry))
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
		urlStr := getURL(entry)
		if !matchURL(urlStr, opts.URLPattern) {
			return false
		}
	}

	// 3. Status check
	if opts.Status != "" {
		status := getStatus(entry)
		if !matchStatus(status, opts.Status) {
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
			resType := strings.ToLower(getResourceType(entry))
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
				if hasXhrOrFetch && hasXHRHeader(entry) {
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

func matchURL(urlStr, pattern string) bool {
	if pattern == "" {
		return true
	}
	re, err := regexp.Compile(pattern)
	if err == nil {
		return re.MatchString(urlStr)
	}
	return strings.Contains(strings.ToLower(urlStr), strings.ToLower(pattern))
}

func matchStatus(status int, pattern string) bool {
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

// Helpers for accessing fields in the raw JSON maps

func getRequest(entry map[string]interface{}) map[string]interface{} {
	if req, ok := entry["request"].(map[string]interface{}); ok {
		return req
	}
	return nil
}

func getResponse(entry map[string]interface{}) map[string]interface{} {
	if resp, ok := entry["response"].(map[string]interface{}); ok {
		return resp
	}
	return nil
}

func getResourceType(entry map[string]interface{}) string {
	if t, ok := entry["_resourceType"].(string); ok {
		return t
	}
	return ""
}

func getURL(entry map[string]interface{}) string {
	req := getRequest(entry)
	if req == nil {
		return ""
	}
	if url, ok := req["url"].(string); ok {
		return url
	}
	return ""
}

func getMethod(entry map[string]interface{}) string {
	req := getRequest(entry)
	if req == nil {
		return ""
	}
	if method, ok := req["method"].(string); ok {
		return method
	}
	return ""
}

func getStatus(entry map[string]interface{}) int {
	resp := getResponse(entry)
	if resp == nil {
		return 0
	}
	if statusVal, ok := resp["status"].(float64); ok {
		return int(statusVal)
	}
	return 0
}

func getStatusText(entry map[string]interface{}) string {
	resp := getResponse(entry)
	if resp == nil {
		return ""
	}
	if text, ok := resp["statusText"].(string); ok {
		return text
	}
	return ""
}

func getTime(entry map[string]interface{}) float64 {
	if tVal, ok := entry["time"].(float64); ok {
		return tVal
	}
	return 0.0
}

func getRequestHeaders(entry map[string]interface{}) []map[string]interface{} {
	req := getRequest(entry)
	if req == nil {
		return nil
	}
	if headers, ok := req["headers"].([]interface{}); ok {
		res := make([]map[string]interface{}, 0, len(headers))
		for _, h := range headers {
			if hMap, ok := h.(map[string]interface{}); ok {
				res = append(res, hMap)
			}
		}
		return res
	}
	return nil
}

func hasXHRHeader(entry map[string]interface{}) bool {
	headers := getRequestHeaders(entry)
	for _, h := range headers {
		name, _ := h["name"].(string)
		val, _ := h["value"].(string)
		if strings.ToLower(name) == "x-requested-with" && strings.ToLower(val) == "xmlhttprequest" {
			return true
		}
	}
	return false
}

// Formatting and Printing functions

func colorizeMethod(padded, raw string) string {
	switch strings.ToUpper(raw) {
	case "GET":
		return colorGreen + padded + colorReset
	case "POST":
		return colorBlue + padded + colorReset
	case "PUT", "PATCH":
		return colorYellow + padded + colorReset
	case "DELETE":
		return colorRed + padded + colorReset
	default:
		return colorGray + padded + colorReset
	}
}

func colorizeStatus(padded string, status int) string {
	if status >= 200 && status < 300 {
		return colorGreen + padded + colorReset
	} else if status >= 300 && status < 400 {
		return colorYellow + padded + colorReset
	} else if status >= 400 && status < 600 {
		return colorRed + padded + colorReset
	} else {
		return colorGray + padded + colorReset
	}
}

func colorizeDuration(padded string, val float64) string {
	if val <= 0 {
		return colorGray + padded + colorReset
	}
	if val < 200 {
		return colorGreen + padded + colorReset
	} else if val < 1000 {
		return colorYellow + padded + colorReset
	} else {
		return colorRed + padded + colorReset
	}
}

func colorizeType(padded, raw string) string {
	rawLower := strings.ToLower(raw)
	if rawLower == "xhr" || rawLower == "fetch" || strings.HasPrefix(rawLower, "xhr") {
		return colorCyan + padded + colorReset
	}
	return colorGray + padded + colorReset
}

func printTable(entries []map[string]interface{}) {
	if len(entries) == 0 {
		fmt.Println("No matching entries found.")
		return
	}

	fmt.Printf("%s%-4s %-6s %-10s %-10s %-8s %s%s\n", 
		colorBold, "#", "METHOD", "STATUS", "DURATION", "TYPE", "URL", colorReset)
	fmt.Println(strings.Repeat("-", 100))

	for i, entry := range entries {
		method := getMethod(entry)
		status := getStatus(entry)
		statusText := getStatusText(entry)
		timeVal := getTime(entry)
		resType := getResourceType(entry)
		urlStr := getURL(entry)

		if resType == "" {
			if hasXHRHeader(entry) {
				resType = "xhr*"
			} else {
				resType = "unknown"
			}
		}

		statusStr := fmt.Sprintf("%d", status)
		if statusText != "" {
			statusStr = fmt.Sprintf("%d %s", status, statusText)
		}
		if len(statusStr) > 10 {
			statusStr = statusStr[:9] + "…"
		}

		durStr := ""
		if timeVal > 0 {
			if timeVal < 1000 {
				durStr = fmt.Sprintf("%.1f ms", timeVal)
			} else {
				durStr = fmt.Sprintf("%.2f s", timeVal/1000.0)
			}
		} else {
			durStr = "-"
		}

		displayURL := urlStr
		if len(displayURL) > 55 {
			displayURL = displayURL[:52] + "..."
		}

		methodPadded := fmt.Sprintf("%-6s", method)
		statusPadded := fmt.Sprintf("%-10s", statusStr)
		durPadded := fmt.Sprintf("%-10s", durStr)
		typePadded := fmt.Sprintf("%-8s", resType)

		methodColored := colorizeMethod(methodPadded, method)
		statusColored := colorizeStatus(statusPadded, status)
		durColored := colorizeDuration(durPadded, timeVal)
		typeColored := colorizeType(typePadded, resType)

		fmt.Printf("%-4d %s %s %s %s %s\n", 
			i, methodColored, statusColored, durColored, typeColored, displayURL)
	}

	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%sTotal: %d entries%s\n", colorBold, len(entries), colorReset)
}

func printJSON(entries []map[string]interface{}) {
	summaries := make([]EntrySummary, len(entries))
	for i, entry := range entries {
		resType := getResourceType(entry)
		if resType == "" {
			if hasXHRHeader(entry) {
				resType = "xhr*"
			} else {
				resType = "unknown"
			}
		}
		summaries[i] = EntrySummary{
			Index:        i,
			Method:       getMethod(entry),
			URL:          getURL(entry),
			Status:       getStatus(entry),
			StatusText:   getStatusText(entry),
			Time:         getTime(entry),
			ResourceType: resType,
		}
	}

	jsonData, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

func printCSV(entries []map[string]interface{}) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	writer.Write([]string{"Index", "Method", "Status", "StatusText", "DurationMs", "ResourceType", "URL"})

	for i, entry := range entries {
		method := getMethod(entry)
		status := getStatus(entry)
		statusText := getStatusText(entry)
		timeVal := getTime(entry)
		resType := getResourceType(entry)
		urlStr := getURL(entry)

		if resType == "" {
			if hasXHRHeader(entry) {
				resType = "xhr*"
			} else {
				resType = "unknown"
			}
		}

		writer.Write([]string{
			strconv.Itoa(i),
			method,
			strconv.Itoa(status),
			statusText,
			fmt.Sprintf("%.2f", timeVal),
			resType,
			urlStr,
		})
	}
}

func printDetailedEntry(entry map[string]interface{}, idx int) {
	method := getMethod(entry)
	urlStr := getURL(entry)
	status := getStatus(entry)
	statusText := getStatusText(entry)
	timeVal := getTime(entry)
	resType := getResourceType(entry)

	fmt.Printf("%s=== Entry #%d Detailed View ===%s\n", colorBold, idx, colorReset)
	fmt.Printf("%sGeneral:%s\n", colorBold, colorReset)
	fmt.Printf("  URL:           %s%s%s\n", colorCyan, urlStr, colorReset)
	fmt.Printf("  Method:        %s\n", colorizeMethod(method, method))
	fmt.Printf("  Status:        %s\n", colorizeStatus(fmt.Sprintf("%d %s", status, statusText), status))
	fmt.Printf("  Duration:      %s\n", colorizeDuration(fmt.Sprintf("%.2f ms", timeVal), timeVal))
	fmt.Printf("  Resource Type: %s\n", colorizeType(resType, resType))

	// Request Headers
	req := getRequest(entry)
	if req != nil {
		fmt.Printf("\n%sRequest Headers:%s\n", colorBold, colorReset)
		headers := getRequestHeaders(entry)
		for _, h := range headers {
			name, _ := h["name"].(string)
			val, _ := h["value"].(string)
			fmt.Printf("  %s%s:%s %s\n", colorGray, name, colorReset, val)
		}

		// Request Body
		if postData, ok := req["postData"].(map[string]interface{}); ok {
			mimeType, _ := postData["mimeType"].(string)
			text, _ := postData["text"].(string)
			fmt.Printf("\n%sRequest Body (%s):%s\n", colorBold, mimeType, colorReset)
			if text != "" {
				if strings.Contains(strings.ToLower(mimeType), "json") {
					var raw map[string]interface{}
					if err := json.Unmarshal([]byte(text), &raw); err == nil {
						indented, _ := json.MarshalIndent(raw, "  ", "  ")
						fmt.Printf("  %s\n", string(indented))
					} else {
						var rawArr []interface{}
						if err := json.Unmarshal([]byte(text), &rawArr); err == nil {
							indented, _ := json.MarshalIndent(rawArr, "  ", "  ")
							fmt.Printf("  %s\n", string(indented))
						} else {
							fmt.Printf("  %s\n", text)
						}
					}
				} else {
					fmt.Printf("  %s\n", text)
				}
			} else if params, ok := postData["params"].([]interface{}); ok && len(params) > 0 {
				for _, p := range params {
					if pMap, ok := p.(map[string]interface{}); ok {
						name, _ := pMap["name"].(string)
						val, _ := pMap["value"].(string)
						fmt.Printf("  %s = %s\n", name, val)
					}
				}
			}
		}
	}

	// Response Headers
	resp := getResponse(entry)
	if resp != nil {
		fmt.Printf("\n%sResponse Headers:%s\n", colorBold, colorReset)
		if headers, ok := resp["headers"].([]interface{}); ok {
			for _, h := range headers {
				if hMap, ok := h.(map[string]interface{}); ok {
					name, _ := hMap["name"].(string)
					val, _ := hMap["value"].(string)
					fmt.Printf("  %s%s:%s %s\n", colorGray, name, colorReset, val)
				}
			}
		}

		// Response Body
		if content, ok := resp["content"].(map[string]interface{}); ok {
			sizeVal, _ := content["size"].(float64)
			mimeType, _ := content["mimeType"].(string)
			text, _ := content["text"].(string)
			encoding, _ := content["encoding"].(string)

			fmt.Printf("\n%sResponse Content (Size: %.0f bytes, MimeType: %s):%s\n", 
				colorBold, sizeVal, mimeType, colorReset)

			if text != "" {
				decodedText := text
				if strings.ToLower(encoding) == "base64" {
					decodedBytes, err := base64.StdEncoding.DecodeString(text)
					if err == nil {
						decodedText = string(decodedBytes)
					}
				}

				if strings.Contains(strings.ToLower(mimeType), "json") {
					var raw map[string]interface{}
					if err := json.Unmarshal([]byte(decodedText), &raw); err == nil {
						indented, _ := json.MarshalIndent(raw, "  ", "  ")
						fmt.Printf("  %s\n", string(indented))
					} else {
						var rawArr []interface{}
						if err := json.Unmarshal([]byte(decodedText), &rawArr); err == nil {
							indented, _ := json.MarshalIndent(rawArr, "  ", "  ")
							fmt.Printf("  %s\n", string(indented))
						} else {
							printTruncatedText(decodedText)
						}
					}
				} else {
					printTruncatedText(decodedText)
				}
			} else {
				fmt.Println("  (No content text or body is empty)")
			}
		}
	}
}

func printTruncatedText(t string) {
	const maxLen = 2000
	if len(t) > maxLen {
		fmt.Printf("  %s\n  ... (truncated %d bytes) ...\n", t[:maxLen], len(t)-maxLen)
	} else {
		fmt.Printf("  %s\n", t)
	}
}
