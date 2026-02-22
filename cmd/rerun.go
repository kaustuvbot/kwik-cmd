package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/spf13/cobra"
)

var (
	runLastDryRun bool
)

var runLastCmd = &cobra.Command{
	Use:   "rerun",
	Short: "Re-run the last tracked command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.Init(); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		commands, err := db.GetRecentCommands(1)
		if err != nil {
			return fmt.Errorf("failed to get commands: %w", err)
		}

		if len(commands) == 0 {
			fmt.Println("No commands to re-run.")
			return nil
		}

		lastCmd := commands[0].FullCommand
		fmt.Printf("Re-running: %s\n", lastCmd)

		if runLastDryRun {
			fmt.Println("(dry-run - not executing)")
			return nil
		}

		// Execute the command
		parts := strings.Fields(lastCmd)
		if len(parts) == 0 {
			return nil
		}

		execCmd := exec.Command(parts[0], parts[1:]...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		return execCmd.Run()
	},
}

func init() {
	runLastCmd.Flags().BoolVarP(&runLastDryRun, "dry-run", "n", false, "Show what would be run without executing")
	rootCmd.AddCommand(runLastCmd)
}
