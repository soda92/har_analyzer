package printer

import (
	"fmt"
	"strings"

	"har_analyzer/pkg/har"
	"github.com/fatih/color"
)

// PrintTable outputs the entries in a beautifully formatted table
func PrintTable(entries []har.Entry) {
	if len(entries) == 0 {
		fmt.Println("No matching entries found.")
		return
	}

	Bold.Printf("%-4s %-6s %-10s %-10s %-8s %s\n", "#", "METHOD", "STATUS", "DURATION", "TYPE", "URL")
	fmt.Println(strings.Repeat("-", 100))

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
	Bold.Printf("Total: %d entries\n", len(entries))
}

func colorizeMethod(padded, raw string) string {
	switch strings.ToUpper(raw) {
	case "GET":
		return color.GreenString(padded)
	case "POST":
		return color.BlueString(padded)
	case "PUT", "PATCH":
		return color.YellowString(padded)
	case "DELETE":
		return color.RedString(padded)
	default:
		return color.HiBlackString(padded)
	}
}

func colorizeStatus(padded string, status int) string {
	if status >= 200 && status < 300 {
		return color.GreenString(padded)
	} else if status >= 300 && status < 400 {
		return color.YellowString(padded)
	} else if status >= 400 && status < 600 {
		return color.RedString(padded)
	} else {
		return color.HiBlackString(padded)
	}
}

func colorizeDuration(padded string, val float64) string {
	if val <= 0 {
		return color.HiBlackString(padded)
	}
	if val < 200 {
		return color.GreenString(padded)
	} else if val < 1000 {
		return color.YellowString(padded)
	} else {
		return color.RedString(padded)
	}
}

func colorizeType(padded, raw string) string {
	rawLower := strings.ToLower(raw)
	if rawLower == "xhr" || rawLower == "fetch" || strings.HasPrefix(rawLower, "xhr") {
		return color.CyanString(padded)
	}
	return color.HiBlackString(padded)
}
