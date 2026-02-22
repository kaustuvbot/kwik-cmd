package cmd

import (
	"github.com/kaustuvbot/kwik-cmd/internal/suggester"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search \"<keywords>\"",
	Short: "Search commands by keywords",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return suggester.Search(args[0])
	},
}
