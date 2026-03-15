package app

import (
	"fmt"
	"math"
	"sort"
)

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
	p95Index := percentileIndex(len(latencies), 0.95)

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

func percentileIndex(length int, percentile float64) int {
	if length <= 1 {
		return 0
	}
	index := int(percentile*float64(length-1) + 0.5)
	if index < 0 {
		return 0
	}
	if index >= length {
		return length - 1
	}
	return index
}

func meanInts(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0
	for _, value := range values {
		total += value
	}
	return float64(total) / float64(len(values))
}

func stddevInts(values []int) float64 {
	if len(values) <= 1 {
		return 0
	}
	mean := meanInts(values)
	var sum float64
	for _, value := range values {
		delta := float64(value) - mean
		sum += delta * delta
	}
	return math.Sqrt(sum / float64(len(values)))
}
