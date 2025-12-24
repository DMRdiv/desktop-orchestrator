package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information (will be set via ldflags in future)
var (
	Version   = "0.1.0-dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display version, git commit, and build date for dspo.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "dspo version %s\n", Version)
		fmt.Fprintf(cmd.OutOrStdout(), "  commit: %s\n", GitCommit)
		fmt.Fprintf(cmd.OutOrStdout(), "  built:  %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
