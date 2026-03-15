package app

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

func TestDiagnoseFindsClustersAndTransitions(t *testing.T) {
	t.Parallel()

	events := []LogEvent{
		{Agent: "planner", Tool: "search", LatencyMS: 80, TokensIn: 50, TokensOut: 20, Success: true},
		{Agent: "planner", Tool: "retriever", LatencyMS: 110, TokensIn: 80, TokensOut: 25, Success: true},
		{Agent: "writer", Tool: "code", LatencyMS: 420, TokensIn: 120, TokensOut: 30, Success: false, Error: "Timeout"},
		{Agent: "writer", Tool: "code", LatencyMS: 390, TokensIn: 110, TokensOut: 25, Success: false, Error: "timeout"},
	}

	report := diagnose(events)
	if len(report.ErrorClusters) == 0 || report.ErrorClusters[0].Signature != "timeout" {
		t.Fatalf("error clusters = %+v, want timeout cluster", report.ErrorClusters)
	}
	if len(report.FlakyTools) == 0 || report.FlakyTools[0].Name != "code" {
		t.Fatalf("flaky tools = %+v, want code", report.FlakyTools)
	}
	if len(report.Transitions) == 0 {
		t.Fatalf("expected transitions in diagnose report")
	}
	if len(report.Outliers) == 0 {
		t.Fatalf("expected outliers in diagnose report")
	}
}
