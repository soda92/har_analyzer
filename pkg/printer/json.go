package printer

import (
	"encoding/json"
	"fmt"
	"log"

	"har_analyzer/pkg/har"
)

// PrintJSON outputs the entries in JSON format
func PrintJSON(entries []har.Entry) {
	summaries := make([]EntrySummary, len(entries))
	for i, entry := range entries {
		resType := entry.GetResourceType()
		if resType == "" {
			if entry.HasXHRHeader() {
				resType = "xhr*"
			} else {
				resType = "unknown"
			}
		}
		summaries[i] = EntrySummary{
			Index:        i,
			Method:       entry.GetMethod(),
			URL:          entry.GetURL(),
			Status:       entry.GetStatus(),
			StatusText:   entry.GetStatusText(),
			Time:         entry.GetTime(),
			ResourceType: resType,
		}
	}

	jsonData, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}
