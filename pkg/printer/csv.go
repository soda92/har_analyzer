package printer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"har_analyzer/pkg/har"
)

// PrintCSV outputs the entries in CSV format
func PrintCSV(entries []har.Entry) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	writer.Write([]string{"Index", "Method", "Status", "StatusText", "DurationMs", "ResourceType", "URL"})

	for i, entry := range entries {
		method := entry.GetMethod()
		status := entry.GetStatus()
		statusText := entry.GetStatusText()
		timeVal := entry.GetTime()
		resType := entry.GetResourceType()
		urlStr := entry.GetURL()

		if resType == "" {
			if entry.HasXHRHeader() {
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
