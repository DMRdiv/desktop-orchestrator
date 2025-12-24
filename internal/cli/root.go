package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dspo",
	Short: "Desktop Profile Orchestrator",
	Long: `dspo is a local-first tool for capturing and replaying
Linux desktop configuration intent across machines.

It focuses on safe, additive reproduction of a workstation
environment using existing system tools.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Future: Add persistent flags here (e.g., --verbose, --config)
}
