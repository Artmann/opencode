package main

import (
	"github.com/sst/opencode/cmd"
	"github.com/sst/opencode/internal/logging"
	"github.com/sst/opencode/internal/status"
)

func main() {
	defer handleErrors()

	cmd.Execute()
}

func handleErrors() {
	logging.RecoverPanic("main", func() {
		status.Error("Application terminated due to unhandled panic.")
	})
}
