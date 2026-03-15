package app

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runTUI()
	}

	switch args[0] {
	case "summary":
		return runSummary(args[1:])
	case "diagnose":
		return runDiagnose(args[1:])
	case "tui", "interactive":
		return runTUI()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", args[0])
		usage()
		return 2
	}
}

func usage() {
	fmt.Println("runlens summarizes and diagnoses agent run logs.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  runlens                    # launch Bubble Tea TUI")
	fmt.Println("  runlens summary -input examples/run.jsonl")
	fmt.Println("  runlens diagnose -input examples/run.jsonl")
}
