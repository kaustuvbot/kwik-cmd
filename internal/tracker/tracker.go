package tracker

import (
	"fmt"
	"os"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/kaustuvbot/kwik-cmd/internal/parser"
)

// TrackCommand tracks a command execution (legacy, assumes success)
func TrackCommand(cmd string) error {
	return TrackCommandWithStatus(cmd, true, 0)
}

// TrackCommandWithStatus tracks a command with its exit status
func TrackCommandWithStatus(cmd string, success bool, exitCode int) error {
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

	// Add keywords
	keywords := parser.ExtractKeywords(parsed)
	for _, kw := range keywords {
		if err := db.AddKeyword(commandID, kw); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to add keyword %s: %v\n", kw, err)
		}
	}

	// Add flags with meanings
	for _, flag := range parsed.Flags {
		meaning := parser.FlagMeaning(flag)
		if err := db.AddFlag(commandID, flag, meaning); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to add flag %s: %v\n", flag, err)
		}
	}

	// Record usage with success/failure status
	if err := db.RecordUsage(commandID, success, exitCode); err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}

	status := "success"
	if !success || exitCode != 0 {
		status = "failed"
	}
	fmt.Printf("Tracked: %s [%s, exit=%d]\n", parsed.FullCmd, status, exitCode)
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
