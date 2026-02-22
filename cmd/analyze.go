package cmd

import (
	"fmt"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze command patterns and suggest improvements",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.Init(); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		fmt.Println("=== Command Pattern Analysis ===\n")

		// Show patterns
		patterns, err := db.DetectPatterns()
		if err != nil {
			return fmt.Errorf("failed to detect patterns: %w", err)
		}

		if len(patterns) > 0 {
			fmt.Println("Detected Patterns:")
			for _, p := range patterns {
				fmt.Printf("\n  %s (used %d times):\n", p.BaseCommand, p.RunCount)
				for i, c := range p.Commands {
					if i >= 5 {
						fmt.Println("    ...")
						break
					}
					fmt.Printf("    - %s\n", c)
				}
			}
		} else {
			fmt.Println("No patterns detected yet. Track more commands!")
		}

		// Show failure stats
		fmt.Println("\n=== Failure Analysis ===")
		failures, err := db.GetFailureStats()
		if err != nil {
			return fmt.Errorf("failed to get failure stats: %w", err)
		}

		if len(failures) > 0 {
			for _, f := range failures {
				fmt.Printf("\n  %s\n", f.FullCommand)
				fmt.Printf("    Runs: %d, Failures: %d, Success Rate: %.1f%%\n",
					f.TotalRuns, f.Failures, f.SuccessRate)
			}
		} else {
			fmt.Println("No failures recorded yet.")
		}

		// Show alias suggestions
		fmt.Println("\n=== Alias Suggestions ===")
		aliases, err := db.SuggestAliases()
		if err != nil {
			return fmt.Errorf("failed to suggest aliases: %w", err)
		}

		if len(aliases) > 0 {
			fmt.Println("Add these to your ~/.bashrc or ~/.zshrc:")
			for _, a := range aliases {
				fmt.Printf("  %s\n", a)
			}
		} else {
			fmt.Println("No alias suggestions yet. Track more commands!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
