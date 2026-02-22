package cmd

import (
	"github.com/kaustuvbot/kwik-cmd/internal/tracker"
	"github.com/spf13/cobra"
)

var (
	trackSuccess bool
	trackExitCode int
)

var trackCmd = &cobra.Command{
	Use:   "track \"<command>\"",
	Short: "Track a command execution",
	Long: `Track a command execution. Use --exit-code to record success/failure.
Examples:
  kwik-cmd track "git commit -m 'fix bug'"
  kwik-cmd track "docker build" --exit-code 0
  kwik-cmd track "npm test" --exit-code 1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return tracker.TrackCommandWithStatus(args[0], trackSuccess, trackExitCode)
	},
}

func init() {
	trackCmd.Flags().BoolVarP(&trackSuccess, "success", "s", true, "Command succeeded")
	trackCmd.Flags().IntVarP(&trackExitCode, "exit-code", "e", 0, "Command exit code")
}
