package cmd

import (
	"log"
	"strings"

	"har_analyzer/pkg/har"
	"har_analyzer/pkg/printer"
	"github.com/spf13/cobra"
)

var (
	fileOpt    string
	outOpt     string
	methodOpt  string
	urlOpt     string
	statusOpt  string
	typeOpt    string
	formatOpt  string
	showOpt    int
	noColorOpt bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "har_analyzer [path/to/file.har]",
	Short: "A fast Go tool to filter and inspect XHR/fetch requests in HAR files",
	Long: `HAR Analyzer is a quick and modular CLI tool written in Go to filter
and analyze network requests (like XHR and fetch) exported in HTTP Archive (HAR) format.
It supports printing to tables, JSON summaries, CSV, or writing back to filtered HAR files.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize colors
		printer.InitColors(noColorOpt)

		// Resolve input path
		filePath := fileOpt
		if filePath == "" && len(args) > 0 {
			filePath = args[0]
		}

		// Parse filter flags
		var methods []string
		if methodOpt != "" {
			parts := strings.Split(methodOpt, ",")
			for _, p := range parts {
				if p = strings.TrimSpace(p); p != "" {
					methods = append(methods, p)
				}
			}
		}

		var resourceTypes []string
		if typeOpt == "" {
			typeOpt = "xhr,fetch"
		}
		parts := strings.Split(typeOpt, ",")
		for _, p := range parts {
			if p = strings.TrimSpace(p); p != "" {
				resourceTypes = append(resourceTypes, p)
			}
		}

		// Read original HAR
		originalHar, err := har.ReadHAR(filePath)
		if err != nil {
			log.Fatalf("Error reading HAR: %v", err)
		}

		// Filter
		opts := har.FilterOpts{
			Methods:    methods,
			URLPattern: urlOpt,
			Status:     statusOpt,
			Types:      resourceTypes,
		}
		filtered := har.FilterEntries(originalHar.Log.Entries, opts)

		// Detailed inspector
		if showOpt >= 0 {
			if showOpt >= len(filtered) {
				log.Fatalf("Index %d is out of range. There are only %d filtered entries (0-%d).", 
					showOpt, len(filtered), len(filtered)-1)
			}
			printer.PrintDetailedEntry(filtered[showOpt], showOpt)
			return
		}

		// Write to file if specified
		if outOpt != "" {
			err := har.WriteHAR(originalHar, filtered, outOpt)
			if err != nil {
				log.Fatalf("Error saving filtered HAR to %s: %v", outOpt, err)
			}
			printer.Green.Printf("Saved %d filtered entries to %s\n\n", len(filtered), outOpt)
		}

		// Print stdout in requested format
		if outOpt == "" || strings.ToLower(formatOpt) != "har" {
			switch strings.ToLower(formatOpt) {
			case "table", "":
				printer.PrintTable(filtered)
			case "json":
				printer.PrintJSON(filtered)
			case "csv":
				printer.PrintCSV(filtered)
			case "har":
				if err := har.WriteHAR(originalHar, filtered, ""); err != nil {
					log.Fatalf("Error printing HAR: %v", err)
				}
			default:
				log.Fatalf("Unknown format: %s. Supported formats: table, json, csv, har", formatOpt)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Execution error: %v", err)
	}
}

func init() {
	RootCmd.Flags().StringVarP(&fileOpt, "file", "i", "", "Input HAR file path (or '-' for stdin)")
	RootCmd.Flags().StringVarP(&outOpt, "out", "o", "", "Output filtered HAR file path")
	RootCmd.Flags().StringVarP(&methodOpt, "method", "m", "", "Filter by HTTP method(s) (comma-separated, e.g. GET,POST)")
	RootCmd.Flags().StringVarP(&urlOpt, "url", "u", "", "Filter by URL (regex or substring match)")
	RootCmd.Flags().StringVarP(&statusOpt, "status", "s", "", "Filter by HTTP status code (e.g. 200, 4xx, 5xx)")
	RootCmd.Flags().StringVarP(&typeOpt, "type", "t", "xhr,fetch", "Filter by resource type(s) (comma-separated, or 'all')")
	RootCmd.Flags().StringVarP(&formatOpt, "format", "f", "table", "Output format: table, json, csv, har")
	RootCmd.Flags().IntVarP(&showOpt, "show", "d", -1, "Inspect detailed headers and body of entry at index")
	RootCmd.Flags().BoolVar(&noColorOpt, "no-color", false, "Disable colorized terminal output")
}
