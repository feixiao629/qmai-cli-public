package main

import (
	"fmt"
	"os"

	"github.com/madaima/qmai-cli/cmd"
)

// version info set by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := cmd.NewRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
