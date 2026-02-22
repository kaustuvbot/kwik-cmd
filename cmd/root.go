package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kwik-cmd",
	Short: "A high-performance CLI tool for command tracking and intelligent suggestions",
	Long: `kwik-cmd tracks terminal commands automatically, learns usage patterns,
and suggests intelligent command completions. Works natively with Bash and Zsh.`,
	Version: "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(trackCmd)
	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(resetCmd)
}
