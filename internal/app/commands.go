package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func runSummary(args []string) int {
	fs := flag.NewFlagSet("summary", flag.ContinueOnError)
	input := fs.String("input", "", "path to a JSONL log file")
	jsonOutput := fs.Bool("json", false, "emit machine-readable JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *input == "" {
		fmt.Fprintln(os.Stderr, "-input is required")
		return 2
	}

	events, err := loadJSONL(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	summary := summarize(events)

	if *jsonOutput {
		body, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Println(string(body))
		return 0
	}

	printSummary(summary)
	return 0
}

func runDiagnose(args []string) int {
	fs := flag.NewFlagSet("diagnose", flag.ContinueOnError)
	input := fs.String("input", "", "path to a JSONL log file")
	jsonOutput := fs.Bool("json", false, "emit machine-readable JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *input == "" {
		fmt.Fprintln(os.Stderr, "-input is required")
		return 2
	}

	events, err := loadJSONL(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	report := diagnose(events)

	if *jsonOutput {
		body, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Println(string(body))
		return 0
	}

	printDiagnose(report)
	return 0
}
