package har

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ReadHAR parses a HAR structure from file path or stdin
func ReadHAR(filePath string) (*HAR, error) {
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

// WriteHAR serializes filtered entries back to HAR format
func WriteHAR(originalHar *HAR, filteredEntries []Entry, outputPath string) error {
	outputHar := HAR{
		Log: Log{
			Version: originalHar.Log.Version,
			Creator: originalHar.Log.Creator,
			Browser: originalHar.Log.Browser,
			Pages:   originalHar.Log.Pages,
			Entries: filteredEntries,
			Comment: originalHar.Log.Comment,
		},
	}

	jsonData, err := json.MarshalIndent(outputHar, "", "  ")
	if err != nil {
		return err
	}

	if outputPath == "" || outputPath == "-" {
		_, err = os.Stdout.Write(jsonData)
		return err
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}
