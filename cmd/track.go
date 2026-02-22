package cmd

import (
	"github.com/kaustuvbot/kwik-cmd/internal/tracker"
	"github.com/spf13/cobra"
)

var trackCmd = &cobra.Command{
	Use:   "track \"<command>\"",
	Short: "Track a command execution",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return tracker.TrackCommand(args[0])
	},
}
