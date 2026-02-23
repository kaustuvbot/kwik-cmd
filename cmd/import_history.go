package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaustuvbot/kwik-cmd/internal/tracker"
	"github.com/spf13/cobra"
)

var importHistoryFile string

var importHistoryCmd = &cobra.Command{
	Use:   "import-history",
	Short: "Import commands from shell history file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return importHistory()
	},
}

func importHistory() error {
	// Default zsh history locations
	historyFiles := []string{
		importHistoryFile,
		os.Getenv("HISTFILE"),
		filepath.Join(os.Getenv("HOME"), ".zsh_history"),
		filepath.Join(os.Getenv("HOME"), ".bash_history"),
		filepath.Join(os.Getenv("HOME"), ".history", "zsh"),
		filepath.Join(os.Getenv("HOME"), ".history", "bash"),
	}

	var historyFile string
	for _, f := range historyFiles {
		if f != "" {
			info, err := os.Stat(f)
			if err == nil && !info.IsDir() {
				historyFile = f
				break
			}
		}
	}

	if historyFile == "" {
		return fmt.Errorf("Could not find history file. Use --file flag")
	}

	fmt.Printf("Importing from: %s\n", historyFile)

	file, err := os.Open(historyFile)
	if err != nil {
		return fmt.Errorf("Failed to open history file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	imported := 0
	skipped := 0

	ignored := map[string]bool{
		"cd": true, "ls": true, "ll": true, "la": true, "lla": true,
		"pwd": true, "echo": true, "exit": true, "export": true,
		"declare": true, "typeset": true, "unset": true, "shift": true,
		"local": true, "readonly": true, "help": true, "which": true,
		"time": true, "fg": true, "bg": true, "jobs": true, "kill": true,
		"builtin": true, "test": true, "[": true, "true": true,
		"false": true, "logout": true, "shopt": true, "umask": true,
		"set": true, "setenv": true, "printenv": true, "eval": true,
		"exec": true, "source": true, "alias": true, "unalias": true,
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines
		if line == "" {
			continue
		}

		// Remove zsh history prefix (timestamp/flags)
		if strings.HasPrefix(line, ": ") {
			parts := strings.SplitN(line, ";", 2)
			if len(parts) > 1 {
				line = strings.TrimSpace(parts[1])
			}
		}

		// Skip comments and special commands
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			skipped++
			continue
		}

		// Get base command
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		base := fields[0]

		// Skip ignored
		if ignored[base] {
			skipped++
			continue
		}

		// Track command
		if err := tracker.TrackCommand(line); err != nil {
			// Silently continue on error
			continue
		}
		imported++
	}

	fmt.Printf("\nImport complete!\n")
	fmt.Printf("Imported: %d commands\n", imported)
	fmt.Printf("Skipped: %d commands\n", skipped)

	return nil
}

func init() {
	importHistoryCmd.Flags().StringVarP(&importHistoryFile, "file", "f", "", "History file path")
	rootCmd.AddCommand(importHistoryCmd)
}
