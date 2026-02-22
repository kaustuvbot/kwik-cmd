package cmd

import (
	"fmt"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/spf13/cobra"
)

var copyIndex int

var copyCmd = &cobra.Command{
	Use:   "copy [index]",
	Short: "Copy a command to clipboard",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.Init(); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		index := 0
		if len(args) > 0 {
			fmt.Sscanf(args[0], "%d", &index)
			index-- // Convert to 0-indexed
		}

		commands, err := db.GetRecentCommands(10)
		if err != nil {
			return fmt.Errorf("failed to get commands: %w", err)
		}

		if len(commands) == 0 {
			fmt.Println("No commands to copy.")
			return nil
		}

		if index < 0 || index >= len(commands) {
			fmt.Printf("Invalid index. Choose 1-%d\n", len(commands))
			return nil
		}

		selectedCmd := commands[index].FullCommand
		fmt.Printf("Command: %s\n", selectedCmd)
		fmt.Println("(Use xclip or pbcopy to copy to clipboard)")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
