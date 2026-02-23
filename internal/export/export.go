package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
)

// ExportJSON exports commands to JSON
func ExportJSON(filename string) error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	commands, err := db.GetRecentCommands(1000)
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	data, err := json.MarshalIndent(commands, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// ExportCSV exports commands to CSV
func ExportCSV(filename string) error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	commands, err := db.GetRecentCommands(1000)
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	writer.Write([]string{"ID", "Base", "Subcommand", "Full Command", "Frequency", "Last Used", "Directory"})

	// Data
	for _, c := range commands {
		writer.Write([]string{
			fmt.Sprintf("%d", c.ID),
			c.Base,
			c.Subcommand,
			c.FullCommand,
			fmt.Sprintf("%d", c.Frequency),
			c.LastUsed.Format("2006-01-02 15:04:05"),
			c.Directory,
		})
	}

	return nil
}

// ImportJSON imports commands from JSON
func ImportJSON(filename string) error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var commands []db.Command
	if err := json.Unmarshal(data, &commands); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	// Re-import each command
	for _, c := range commands {
		_, err := db.AddCommand(c.Base, c.Subcommand, c.FullCommand, c.Directory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to import %s: %v\n", c.FullCommand, err)
		}
	}

	return nil
}
