package cmd

import (
	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset command history",
	RunE: func(cmd *cobra.Command, args []string) error {
		return db.Reset()
	},
}
