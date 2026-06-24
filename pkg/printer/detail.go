package printer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"har_analyzer/pkg/har"
	"github.com/fatih/color"
)

// PrintDetailedEntry prints request and response headers & content
func PrintDetailedEntry(entry har.Entry, idx int) {
	method := entry.GetMethod()
	urlStr := entry.GetURL()
	status := entry.GetStatus()
	statusText := entry.GetStatusText()
	timeVal := entry.GetTime()
	resType := entry.GetResourceType()

	Bold.Printf("=== Entry #%d Detailed View ===\n", idx)
	Bold.Println("General:")
	fmt.Printf("  URL:           %s\n", color.CyanString(urlStr))
	fmt.Printf("  Method:        %s\n", colorizeMethod(method, method))
	fmt.Printf("  Status:        %s\n", colorizeStatus(fmt.Sprintf("%d %s", status, statusText), status))
	fmt.Printf("  Duration:      %s\n", colorizeDuration(fmt.Sprintf("%.2f ms", timeVal), timeVal))
	fmt.Printf("  Resource Type: %s\n", colorizeType(resType, resType))

	// Request Details
	req := entry.GetRequest()
	if req != nil {
		Bold.Println("\nRequest Headers:")
		headers := entry.GetRequestHeaders()
		for _, h := range headers {
			name, _ := h["name"].(string)
			val, _ := h["value"].(string)
			fmt.Printf("  %s %s\n", color.HiBlackString(name+":"), val)
		}

		// Request Body
		if postData, ok := req["postData"].(map[string]interface{}); ok {
			mimeType, _ := postData["mimeType"].(string)
			text, _ := postData["text"].(string)
			Bold.Printf("\nRequest Body (%s):\n", mimeType)
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

	// Response Details
	resp := entry.GetResponse()
	if resp != nil {
		Bold.Println("\nResponse Headers:")
		if headers, ok := resp["headers"].([]interface{}); ok {
			for _, h := range headers {
				if hMap, ok := h.(map[string]interface{}); ok {
					name, _ := hMap["name"].(string)
					val, _ := hMap["value"].(string)
					fmt.Printf("  %s %s\n", color.HiBlackString(name+":"), val)
				}
			}
		}

		// Response Body
		if content, ok := resp["content"].(map[string]interface{}); ok {
			sizeVal, _ := content["size"].(float64)
			mimeType, _ := content["mimeType"].(string)
			text, _ := content["text"].(string)
			encoding, _ := content["encoding"].(string)

			Bold.Printf("\nResponse Content (Size: %.0f bytes, MimeType: %s):\n", sizeVal, mimeType)

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
