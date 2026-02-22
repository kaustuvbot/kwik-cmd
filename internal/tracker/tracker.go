package tracker

import (
	"fmt"
	"os"
	"time"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/kaustuvbot/kwik-cmd/internal/parser"
)

// TrackCommand tracks a command execution
func TrackCommand(cmd string) error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	parsed := parser.ParseCommand(cmd)
	if parsed == nil {
		return fmt.Errorf("failed to parse command")
	}

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		dir = ""
	}

	// Add to database
	commandID, err := db.AddCommand(parsed.Base, parsed.Subcommand, parsed.FullCmd, dir)
	if err != nil {
		return fmt.Errorf("failed to add command: %w", err)
	}

	// Record usage
	if err := db.RecordUsage(commandID, true, 0); err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}

	fmt.Printf("Tracked: %s\n", parsed.FullCmd)
	return nil
}

// ShowStats displays command usage statistics
func ShowStats() error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	total, executions, err := db.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	fmt.Printf("=== kwik-cmd Statistics ===\n")
	fmt.Printf("Total unique commands: %d\n", total)
	fmt.Printf("Total executions: %d\n", executions)

	// Show recent commands
	recent, err := db.GetRecentCommands(5)
	if err != nil {
		return fmt.Errorf("failed to get recent commands: %w", err)
	}

	if len(recent) > 0 {
		fmt.Printf("\n=== Recent Commands ===\n")
		for i, c := range recent {
			fmt.Printf("%d. %s (used %d times, last: %s)\n",
				i+1, c.FullCommand, c.Frequency, c.LastUsed.Format("2006-01-02 15:04"))
		}
	}

	// Show top commands
	top, err := db.GetTopCommands(5)
	if err != nil {
		return fmt.Errorf("failed to get top commands: %w", err)
	}

	if len(top) > 0 {
		fmt.Printf("\n=== Most Used Commands ===\n")
		for i, c := range top {
			fmt.Printf("%d. %s (used %d times)\n", i+1, c.FullCommand, c.Frequency)
		}
	}

	return nil
}

// GetCurrentDirectory returns the current working directory
func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// GetTimestamp returns the current timestamp
func GetTimestamp() time.Time {
	return time.Now()
}
