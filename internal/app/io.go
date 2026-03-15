package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

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
