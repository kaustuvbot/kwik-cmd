package cmd

import (
	"fmt"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/spf13/cobra"
)

var quickPickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Quick pick from recent commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.Init(); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		commands, err := db.GetRecentCommands(5)
		if err != nil {
			return fmt.Errorf("failed to get commands: %w", err)
		}

		if len(commands) == 0 {
			fmt.Println("No commands tracked yet.")
			return nil
		}

		fmt.Println("Recent commands:")
		for i, c := range commands {
			fmt.Printf("  %d. %s\n", i+1, c.FullCommand)
		}
		fmt.Println("\nCopy the command you want to use.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(quickPickCmd)
}
