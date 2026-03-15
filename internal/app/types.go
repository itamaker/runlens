package app

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

type ErrorCluster struct {
	Signature string   `json:"signature"`
	Count     int      `json:"count"`
	Tools     []string `json:"tools,omitempty"`
}

type FlakyTool struct {
	Name         string  `json:"name"`
	Runs         int     `json:"runs"`
	Failures     int     `json:"failures"`
	FailureRate  float64 `json:"failure_rate"`
	AvgLatencyMS float64 `json:"avg_latency_ms"`
	P95LatencyMS int     `json:"p95_latency_ms"`
	AvgTokens    float64 `json:"avg_tokens"`
}

type Transition struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Count int    `json:"count"`
}

type OutlierEvent struct {
	Index     int    `json:"index"`
	Tool      string `json:"tool"`
	Agent     string `json:"agent,omitempty"`
	LatencyMS int    `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
	Reason    string `json:"reason"`
}

type DiagnoseReport struct {
	Summary       Summary        `json:"summary"`
	ErrorClusters []ErrorCluster `json:"error_clusters,omitempty"`
	FlakyTools    []FlakyTool    `json:"flaky_tools,omitempty"`
	Transitions   []Transition   `json:"transitions,omitempty"`
	Outliers      []OutlierEvent `json:"outliers,omitempty"`
}

type toolStats struct {
	runs      int
	failures  int
	latencies []int
	tokens    int
}
