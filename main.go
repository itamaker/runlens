package main

import (
	"os"

	"github.com/itamaker/runlens/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
