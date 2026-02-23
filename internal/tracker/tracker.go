package tracker

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kaustuvbot/kwik-cmd/internal/db"
	"github.com/kaustuvbot/kwik-cmd/internal/parser"
)

var (
	bold   = color.New(color.Bold)
	cyan   = color.New(color.FgCyan)
	green  = color.New(color.FgGreen)
	yellow = color.New(color.FgYellow)
	white  = color.New(color.FgWhite)
	dim    = color.New(color.FgBlack)
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

	// Colored output
	green.Print("✓ Tracked: ")
	white.Print(parsed.FullCmd)
	
	if success && exitCode == 0 {
		cyan.Print(" [success]")
	} else {
		yellow.Printf(" [failed, exit=%d]", exitCode)
	}
	fmt.Println()
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

	bold.Println("=== kwik-cmd Statistics ===")
	fmt.Print("Total unique commands: ")
	cyan.Printf("%d\n", total)
	fmt.Print("Total executions: ")
	green.Printf("%d\n", executions)

	// Show recent commands
	recent, err := db.GetRecentCommands(5)
	if err != nil {
		return fmt.Errorf("failed to get recent commands: %w", err)
	}

	if len(recent) > 0 {
		bold.Print("\n=== Recent Commands ===\n")
		for i, c := range recent {
			bold.Printf("%d. ", i+1)
			green.Print(c.FullCommand)
			dim.Printf(" (used %d times, last: %s)\n",
				c.Frequency, c.LastUsed.Format("2006-01-02 15:04"))
		}
	}

	// Show top commands
	top, err := db.GetTopCommands(5)
	if err != nil {
		return fmt.Errorf("failed to get top commands: %w", err)
	}

	if len(top) > 0 {
		bold.Print("\n=== Most Used Commands ===\n")
		for i, c := range top {
			bold.Printf("%d. ", i+1)
			green.Print(c.FullCommand)
			dim.Printf(" (used %d times)\n", c.Frequency)
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
