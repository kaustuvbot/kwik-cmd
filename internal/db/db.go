package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, ".kwik-cmd")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "commands.db")
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		base TEXT NOT NULL,
		subcommand TEXT,
		full_command TEXT NOT NULL,
		frequency INTEGER DEFAULT 1,
		last_used DATETIME DEFAULT CURRENT_TIMESTAMP,
		directory TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS flags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command_id INTEGER NOT NULL,
		flag TEXT NOT NULL,
		meaning TEXT,
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS keywords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command_id INTEGER NOT NULL,
		keyword TEXT NOT NULL,
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS usage_stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command_id INTEGER NOT NULL,
		success BOOLEAN DEFAULT TRUE,
		exit_code INTEGER,
		used_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_commands_base ON commands(base);
	CREATE INDEX IF NOT EXISTS idx_commands_directory ON commands(directory);
	CREATE INDEX IF NOT EXISTS idx_keywords_keyword ON keywords(keyword);
	CREATE INDEX IF NOT EXISTS idx_usage_stats_used_at ON usage_stats(used_at);
	`

	_, err := db.Exec(schema)
	return err
}

func Close() {
	if db != nil {
		db.Close()
	}
}

type Command struct {
	ID          int64
	Base        string
	Subcommand  string
	FullCommand string
	Frequency   int
	LastUsed    time.Time
	Directory   string
}

func AddCommand(base, subcommand, fullCommand, directory string) (int64, error) {
	// Check if command already exists
	var existingID int64
	err := db.QueryRow(`
		SELECT id FROM commands 
		WHERE base = ? AND full_command = ? AND directory = ?
	`, base, fullCommand, directory).Scan(&existingID)

	if err == nil {
		// Update frequency and last_used
		_, err = db.Exec(`
			UPDATE commands 
			SET frequency = frequency + 1, last_used = CURRENT_TIMESTAMP 
			WHERE id = ?
		`, existingID)
		return existingID, err
	}

	// Insert new command
	result, err := db.Exec(`
		INSERT INTO commands (base, subcommand, full_command, directory)
		VALUES (?, ?, ?, ?)
	`, base, subcommand, fullCommand, directory)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func RecordUsage(commandID int64, success bool, exitCode int) error {
	_, err := db.Exec(`
		INSERT INTO usage_stats (command_id, success, exit_code)
		VALUES (?, ?, ?)
	`, commandID, success, exitCode)
	return err
}

func GetRecentCommands(limit int) ([]Command, error) {
	rows, err := db.Query(`
		SELECT id, base, subcommand, full_command, frequency, last_used, directory
		FROM commands
		ORDER BY last_used DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var c Command
		if err := rows.Scan(&c.ID, &c.Base, &c.Subcommand, &c.FullCommand, &c.Frequency, &c.LastUsed, &c.Directory); err != nil {
			return nil, err
		}
		commands = append(commands, c)
	}
	return commands, nil
}

func GetTopCommands(limit int) ([]Command, error) {
	rows, err := db.Query(`
		SELECT id, base, subcommand, full_command, frequency, last_used, directory
		FROM commands
		ORDER BY frequency DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var c Command
		if err := rows.Scan(&c.ID, &c.Base, &c.Subcommand, &c.FullCommand, &c.Frequency, &c.LastUsed, &c.Directory); err != nil {
			return nil, err
		}
		commands = append(commands, c)
	}
	return commands, nil
}

func SearchByKeyword(keyword string) ([]Command, error) {
	rows, err := db.Query(`
		SELECT DISTINCT c.id, c.base, c.subcommand, c.full_command, c.frequency, c.last_used, c.directory
		FROM commands c
		JOIN keywords k ON c.id = k.command_id
		WHERE k.keyword LIKE ?
		ORDER BY c.frequency DESC
		LIMIT 20
	`, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var c Command
		if err := rows.Scan(&c.ID, &c.Base, &c.Subcommand, &c.FullCommand, &c.Frequency, &c.LastUsed, &c.Directory); err != nil {
			return nil, err
		}
		commands = append(commands, c)
	}
	return commands, nil
}

func GetStats() (totalCommands int, totalExecutions int, err error) {
	err = db.QueryRow("SELECT COUNT(*), COALESCE(SUM(frequency), 0) FROM commands").Scan(&totalCommands, &totalExecutions)
	return
}

func Reset() error {
	_, err := db.Exec("DELETE FROM usage_stats; DELETE FROM keywords; DELETE FROM flags; DELETE FROM commands;")
	return err
}

// GetDB returns the database connection for external use
func GetDB() *sql.DB {
	return db
}

// WithContext returns a context with the database connection
func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "db", db)
}
