package cmd

import (
	"github.com/kaustuvbot/kwik-cmd/internal/tracker"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show command usage statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tracker.ShowStats()
	},
}
