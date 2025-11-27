package main

import (
	"github.com/sst/opencode/cmd"
	"github.com/sst/opencode/internal/logging"
	"github.com/sst/opencode/internal/status"
)

func main() {
	defer handleErrors2()

	cmd.Execute()
}

func handleErrors2() {
	logging.RecoverPanic("main", func() {
		status.Error("Application terminated due to unhandled panic.")
	})
}
