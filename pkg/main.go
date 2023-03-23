package main

import (
	"github.com/spf13/cobra"

	"github.com/aiven/go-api-schemas/pkg/cmd"
	"github.com/aiven/go-api-schemas/pkg/util"
)

// logger is the logger of the application.
var logger = &util.Logger{}

// rootCmd is the root command for the application.
var rootCmd *cobra.Command

// setup sets up the application.
func setup() {
	util.SetupLogger(logger)

	rootCmd = cmd.NewCmdRoot(logger)
}

// main is the entrypoint for the application.
func main() {
	setup()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
