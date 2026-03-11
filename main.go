package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

type LogEvent struct {
	Timestamp string `json:"timestamp"`
	Agent     string `json:"agent"`
	Tool      string `json:"tool"`
	LatencyMS int    `json:"latency_ms"`
	TokensIn  int    `json:"tokens_in"`
	TokensOut int    `json:"tokens_out"`
	Success   bool   `json:"success"`
	Error     string `json:"error"`
}

type ToolSummary struct {
	Name         string  `json:"name"`
	Runs         int     `json:"runs"`
	Failures     int     `json:"failures"`
	AvgLatencyMS float64 `json:"avg_latency_ms"`
}

type Summary struct {
	TotalRuns      int           `json:"total_runs"`
	SuccessRate    float64       `json:"success_rate"`
	AvgLatencyMS   float64       `json:"avg_latency_ms"`
	P95LatencyMS   int           `json:"p95_latency_ms"`
	TotalTokensIn  int           `json:"total_tokens_in"`
	TotalTokensOut int           `json:"total_tokens_out"`
	Tools          []ToolSummary `json:"tools"`
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}

	switch args[0] {
	case "summary":
		return runSummary(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", args[0])
		usage()
		return 2
	}
}

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

func usage() {
	fmt.Println("runlens summarizes agent run logs.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  runlens summary -input examples/run.jsonl")
}

func loadJSONL(path string) ([]LogEvent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	defer file.Close()

	var events []LogEvent
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++
		if scanner.Text() == "" {
			continue
		}

		var event LogEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf("decode JSONL line %d: %w", line, err)
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan log file: %w", err)
	}
	return events, nil
}

func summarize(events []LogEvent) Summary {
	if len(events) == 0 {
		return Summary{}
	}

	var successCount int
	var totalLatency int
	var totalTokensIn int
	var totalTokensOut int
	latencies := make([]int, 0, len(events))
	toolIndex := map[string]*ToolSummary{}

	for _, event := range events {
		if event.Success {
			successCount++
		}
		totalLatency += event.LatencyMS
		totalTokensIn += event.TokensIn
		totalTokensOut += event.TokensOut
		latencies = append(latencies, event.LatencyMS)

		summary, ok := toolIndex[event.Tool]
		if !ok {
			summary = &ToolSummary{Name: event.Tool}
			toolIndex[event.Tool] = summary
		}
		summary.Runs++
		summary.AvgLatencyMS += float64(event.LatencyMS)
		if !event.Success {
			summary.Failures++
		}
	}

	toolSummaries := make([]ToolSummary, 0, len(toolIndex))
	for _, summary := range toolIndex {
		summary.AvgLatencyMS = summary.AvgLatencyMS / float64(summary.Runs)
		toolSummaries = append(toolSummaries, *summary)
	}
	sort.Slice(toolSummaries, func(i, j int) bool {
		if toolSummaries[i].Runs == toolSummaries[j].Runs {
			return toolSummaries[i].Name < toolSummaries[j].Name
		}
		return toolSummaries[i].Runs > toolSummaries[j].Runs
	})

	sort.Ints(latencies)
	p95Index := int(0.95*float64(len(latencies)-1) + 0.5)

	return Summary{
		TotalRuns:      len(events),
		SuccessRate:    float64(successCount) / float64(len(events)),
		AvgLatencyMS:   float64(totalLatency) / float64(len(events)),
		P95LatencyMS:   latencies[p95Index],
		TotalTokensIn:  totalTokensIn,
		TotalTokensOut: totalTokensOut,
		Tools:          toolSummaries,
	}
}

func printSummary(summary Summary) {
	fmt.Printf("Runs: %d\n", summary.TotalRuns)
	fmt.Printf("Success rate: %.2f%%\n", summary.SuccessRate*100)
	fmt.Printf("Latency avg/p95: %.1fms / %dms\n", summary.AvgLatencyMS, summary.P95LatencyMS)
	fmt.Printf("Tokens in/out: %d / %d\n", summary.TotalTokensIn, summary.TotalTokensOut)
	fmt.Println("Tool usage:")
	for _, tool := range summary.Tools {
		fmt.Printf("- %s: %d runs, %d failures, %.1fms avg\n", tool.Name, tool.Runs, tool.Failures, tool.AvgLatencyMS)
	}
}
