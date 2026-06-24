package printer

import (
	"github.com/fatih/color"
)

// InitColors initializes the color library setting
func InitColors(noColor bool) {
	color.NoColor = noColor
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

// Color helpers using fatih/color
var (
	Bold      = color.New(color.Bold)
	Red       = color.New(color.FgRed)
	Green     = color.New(color.FgGreen)
	Yellow    = color.New(color.FgYellow)
	Blue      = color.New(color.FgBlue)
	Magenta   = color.New(color.FgMagenta)
	Cyan      = color.New(color.FgCyan)
	Gray      = color.New(color.FgHiBlack)
)
