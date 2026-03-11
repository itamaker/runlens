package main

import "testing"

func TestSummarize(t *testing.T) {
	t.Parallel()

	events := []LogEvent{
		{Tool: "search", LatencyMS: 80, TokensIn: 120, TokensOut: 30, Success: true},
		{Tool: "search", LatencyMS: 120, TokensIn: 100, TokensOut: 20, Success: true},
		{Tool: "code", LatencyMS: 240, TokensIn: 200, TokensOut: 50, Success: false},
	}

	summary := summarize(events)
	if summary.TotalRuns != 3 {
		t.Fatalf("TotalRuns = %d, want 3", summary.TotalRuns)
	}
	if summary.P95LatencyMS != 240 {
		t.Fatalf("P95LatencyMS = %d, want 240", summary.P95LatencyMS)
	}
	if len(summary.Tools) != 2 {
		t.Fatalf("tool summary count = %d, want 2", len(summary.Tools))
	}
}
