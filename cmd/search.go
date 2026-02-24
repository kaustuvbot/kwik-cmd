package cmd

import (
	"fmt"

	"github.com/kaustuvbot/kwik-cmd/internal/suggester"
	"github.com/spf13/cobra"
)

var (
	searchPlainFlag bool
	searchLimitFlag int
)

var searchCmd = &cobra.Command{
	Use:   "search \"<keywords>\"",
	Short: "Search commands by keywords",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if searchPlainFlag {
			commands, err := suggester.SearchPlain(args[0], searchLimitFlag)
			if err != nil {
				return err
			}
			for _, c := range commands {
				fmt.Println(c)
			}
			return nil
		}
		return suggester.Search(args[0])
	},
}

func init() {
	searchCmd.Flags().BoolVarP(&searchPlainFlag, "plain", "p", false, "Output plain text (one command per line, no colors/headers)")
	searchCmd.Flags().IntVarP(&searchLimitFlag, "limit", "l", 10, "Maximum number of results to return")
}
