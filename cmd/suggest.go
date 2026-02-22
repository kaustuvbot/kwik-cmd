package cmd

import (
	"github.com/kaustuvbot/kwik-cmd/internal/suggester"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:   "suggest \"<partial>\"",
	Short: "Get command suggestions for partial input",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return suggester.Suggest(args[0])
	},
}
