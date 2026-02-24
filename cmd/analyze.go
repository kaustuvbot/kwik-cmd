package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/spf13/cobra"
)

var (
	bold    = color.New(color.Bold)
	cyan    = color.New(color.FgCyan)
	green   = color.New(color.FgGreen)
	yellow  = color.New(color.FgYellow)
	magenta = color.New(color.FgMagenta)
	dim     = color.New(color.FgBlack)
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze command patterns and suggest improvements",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.Init(); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		bold.Print("=== Command Pattern Analysis ===")

		// Show patterns
		patterns, err := db.DetectPatterns()
		if err != nil {
			return fmt.Errorf("failed to detect patterns: %w", err)
		}

		if len(patterns) > 0 {
			bold.Println("Detected Patterns:")
			for _, p := range patterns {
				cyan.Printf("\n  %s ", p.BaseCommand)
				dim.Printf("(used %d times):\n", p.RunCount)
				for i, c := range p.Commands {
					if i >= 5 {
						dim.Println("    ...")
						break
					}
					green.Printf("    - %s\n", c)
				}
			}
		} else {
			yellow.Println("No patterns detected yet. Track more commands!")
		}

		// Show failure stats
		bold.Print("\n=== Failure Analysis ===\n")
		failures, err := db.GetFailureStats()
		if err != nil {
			return fmt.Errorf("failed to get failure stats: %w", err)
		}

		if len(failures) > 0 {
			for _, f := range failures {
				green.Printf("\n  %s\n", f.FullCommand)
				dim.Printf("    Runs: %d, Failures: %d, Success Rate: %.1f%%\n",
					f.TotalRuns, f.Failures, f.SuccessRate)
			}
		} else {
			dim.Println("No failures recorded yet.")
		}

		// Show alias suggestions
		bold.Print("\n=== Alias Suggestions ===\n")
		aliases, err := db.SuggestAliases()
		if err != nil {
			return fmt.Errorf("failed to suggest aliases: %w", err)
		}

		if len(aliases) > 0 {
			magenta.Println("Add these to your ~/.bashrc or ~/.zshrc:")
			for _, a := range aliases {
				green.Printf("  %s\n", a)
			}
		} else {
			dim.Println("No alias suggestions yet. Track more commands!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
