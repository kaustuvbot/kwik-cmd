package cmd

import (
	"fmt"

	"github.com/kaustuvbot/kwik-cmd/internal/suggester"
	"github.com/spf13/cobra"
)

var (
	plainFlag  bool
	splitFlag  bool
	limitFlag  int
)

var suggestCmd = &cobra.Command{
	Use:   "suggest \"<partial>\"",
	Short: "Get command suggestions for partial input",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if splitFlag {
			result, err := suggester.SuggestPlainSplit(args[0], limitFlag)
			if err != nil {
				return err
			}
			fmt.Println("---RECENT---")
			for _, c := range result["recent"] {
				fmt.Println(c)
			}
			fmt.Println("---FREQUENT---")
			for _, c := range result["frequent"] {
				fmt.Println(c)
			}
			return nil
		}
		if plainFlag {
			commands, err := suggester.SuggestPlain(args[0], limitFlag)
			if err != nil {
				return err
			}
			for _, c := range commands {
				fmt.Println(c)
			}
			return nil
		}
		return suggester.Suggest(args[0])
	},
}

func init() {
	suggestCmd.Flags().BoolVarP(&plainFlag, "plain", "p", false, "Output plain text (one command per line, no colors/headers)")
	suggestCmd.Flags().BoolVarP(&splitFlag, "split", "s", false, "Output split by recent and frequent (for shell integration)")
	suggestCmd.Flags().IntVarP(&limitFlag, "limit", "l", 10, "Maximum number of suggestions to return")
}
