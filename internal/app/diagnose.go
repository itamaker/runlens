package app

import (
	"fmt"
	"sort"
	"strings"
)

func diagnose(events []LogEvent) DiagnoseReport {
	report := DiagnoseReport{
		Summary: summarize(events),
	}
	if len(events) == 0 {
		return report
	}

	errorClusters := map[string]*ErrorCluster{}
	toolStatsIndex := map[string]*toolStats{}
	transitionCounts := map[string]int{}
	globalP95 := report.Summary.P95LatencyMS

	for i, event := range events {
		stats := toolStatsIndex[event.Tool]
		if stats == nil {
			stats = &toolStats{}
			toolStatsIndex[event.Tool] = stats
		}
		stats.runs++
		stats.latencies = append(stats.latencies, event.LatencyMS)
		stats.tokens += event.TokensIn + event.TokensOut
		if !event.Success {
			stats.failures++
		}

		if i > 0 && events[i-1].Tool != "" && event.Tool != "" {
			key := events[i-1].Tool + "->" + event.Tool
			transitionCounts[key]++
		}

		if !event.Success || strings.TrimSpace(event.Error) != "" {
			signature := normalizeError(event.Error)
			cluster := errorClusters[signature]
			if cluster == nil {
				cluster = &ErrorCluster{Signature: signature}
				errorClusters[signature] = cluster
			}
			cluster.Count++
			cluster.Tools = appendIfMissing(cluster.Tools, event.Tool)
		}
	}

	for tool, stats := range toolStatsIndex {
		sort.Ints(stats.latencies)
		p95Latency := stats.latencies[percentileIndex(len(stats.latencies), 0.95)]
		failureRate := float64(stats.failures) / float64(stats.runs)
		avgLatency := meanInts(stats.latencies)
		avgTokens := float64(stats.tokens) / float64(stats.runs)

		if failureRate >= 0.2 || p95Latency > int(float64(globalP95)*1.25) {
			report.FlakyTools = append(report.FlakyTools, FlakyTool{
				Name:         tool,
				Runs:         stats.runs,
				Failures:     stats.failures,
				FailureRate:  failureRate,
				AvgLatencyMS: avgLatency,
				P95LatencyMS: p95Latency,
				AvgTokens:    avgTokens,
			})
		}
	}
	sort.Slice(report.FlakyTools, func(i, j int) bool {
		if report.FlakyTools[i].FailureRate == report.FlakyTools[j].FailureRate {
			return report.FlakyTools[i].P95LatencyMS > report.FlakyTools[j].P95LatencyMS
		}
		return report.FlakyTools[i].FailureRate > report.FlakyTools[j].FailureRate
	})

	for _, cluster := range errorClusters {
		sort.Strings(cluster.Tools)
		report.ErrorClusters = append(report.ErrorClusters, *cluster)
	}
	sort.Slice(report.ErrorClusters, func(i, j int) bool {
		if report.ErrorClusters[i].Count == report.ErrorClusters[j].Count {
			return report.ErrorClusters[i].Signature < report.ErrorClusters[j].Signature
		}
		return report.ErrorClusters[i].Count > report.ErrorClusters[j].Count
	})

	for key, count := range transitionCounts {
		parts := strings.SplitN(key, "->", 2)
		report.Transitions = append(report.Transitions, Transition{
			From:  parts[0],
			To:    parts[1],
			Count: count,
		})
	}
	sort.Slice(report.Transitions, func(i, j int) bool {
		if report.Transitions[i].Count == report.Transitions[j].Count {
			if report.Transitions[i].From == report.Transitions[j].From {
				return report.Transitions[i].To < report.Transitions[j].To
			}
			return report.Transitions[i].From < report.Transitions[j].From
		}
		return report.Transitions[i].Count > report.Transitions[j].Count
	})
	if len(report.Transitions) > 5 {
		report.Transitions = report.Transitions[:5]
	}

	for i, event := range events {
		stats := toolStatsIndex[event.Tool]
		reason := outlierReason(event, stats, globalP95)
		if reason == "" {
			continue
		}
		report.Outliers = append(report.Outliers, OutlierEvent{
			Index:     i + 1,
			Tool:      event.Tool,
			Agent:     event.Agent,
			LatencyMS: event.LatencyMS,
			Error:     event.Error,
			Reason:    reason,
		})
	}
	sort.Slice(report.Outliers, func(i, j int) bool {
		if report.Outliers[i].LatencyMS == report.Outliers[j].LatencyMS {
			return report.Outliers[i].Index < report.Outliers[j].Index
		}
		return report.Outliers[i].LatencyMS > report.Outliers[j].LatencyMS
	})
	if len(report.Outliers) > 5 {
		report.Outliers = report.Outliers[:5]
	}

	return report
}

func printDiagnose(report DiagnoseReport) {
	printSummary(report.Summary)

	if len(report.FlakyTools) > 0 {
		fmt.Println("Flaky tools:")
		for _, tool := range report.FlakyTools {
			fmt.Printf("- %s: failure %.0f%%, avg %.1fms, p95 %dms\n",
				tool.Name, tool.FailureRate*100, tool.AvgLatencyMS, tool.P95LatencyMS)
		}
	}
	if len(report.ErrorClusters) > 0 {
		fmt.Println("Error clusters:")
		for _, cluster := range report.ErrorClusters {
			fmt.Printf("- %s: %d (%s)\n", cluster.Signature, cluster.Count, strings.Join(cluster.Tools, ", "))
		}
	}
	if len(report.Transitions) > 0 {
		fmt.Println("Top transitions:")
		for _, edge := range report.Transitions {
			fmt.Printf("- %s -> %s: %d\n", edge.From, edge.To, edge.Count)
		}
	}
	if len(report.Outliers) > 0 {
		fmt.Println("Outliers:")
		for _, outlier := range report.Outliers {
			fmt.Printf("- #%d %s: %dms (%s)\n", outlier.Index, outlier.Tool, outlier.LatencyMS, outlier.Reason)
		}
	}
}

func normalizeError(text string) string {
	text = strings.TrimSpace(strings.ToLower(text))
	if text == "" {
		return "unspecified failure"
	}
	return strings.Join(strings.Fields(text), " ")
}

func appendIfMissing(items []string, item string) []string {
	for _, existing := range items {
		if existing == item {
			return items
		}
	}
	return append(items, item)
}

func outlierReason(event LogEvent, stats *toolStats, globalP95 int) string {
	var reasons []string
	if !event.Success {
		reasons = append(reasons, "failed")
	}
	if stats != nil {
		threshold := meanInts(stats.latencies) + 2*stddevInts(stats.latencies)
		if float64(event.LatencyMS) >= threshold && event.LatencyMS > globalP95 {
			reasons = append(reasons, "tool-latency-outlier")
		}
	}
	if event.LatencyMS > globalP95 && globalP95 > 0 && len(reasons) == 0 {
		reasons = append(reasons, "global-latency-outlier")
	}
	return strings.Join(reasons, ", ")
}
