package db

import (
	"fmt"
	"strings"
	"time"
)

// FailureStats represents failure statistics for a command
type FailureStats struct {
	CommandID   int64
	FullCommand string
	TotalRuns   int
	Failures    int
	SuccessRate float64
	LastFailure *time.Time
}

// GetFailureStats returns failure statistics for all commands
func GetFailureStats() ([]FailureStats, error) {
	rows, err := db.Query(`
		SELECT 
			c.id,
			c.full_command,
			COUNT(us.id) as total_runs,
			SUM(CASE WHEN us.success = 0 THEN 1 ELSE 0 END) as failures,
			MAX(CASE WHEN us.success = 0 THEN us.used_at END) as last_failure
		FROM commands c
		LEFT JOIN usage_stats us ON c.id = us.command_id
		GROUP BY c.id
		HAVING failures > 0
		ORDER BY failures DESC
		LIMIT 20
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []FailureStats
	for rows.Next() {
		var s FailureStats
		if err := rows.Scan(&s.CommandID, &s.FullCommand, &s.TotalRuns, &s.Failures, &s.LastFailure); err != nil {
			return nil, err
		}
		if s.TotalRuns > 0 {
			s.SuccessRate = float64(s.TotalRuns-s.Failures) / float64(s.TotalRuns) * 100
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// PatternGroup represents a group of related commands
type PatternGroup struct {
	BaseCommand string
	Commands    []string
	RunCount    int
}

// DetectPatterns detects command patterns (e.g., git subcommands, docker commands)
func DetectPatterns() ([]PatternGroup, error) {
	rows, err := db.Query(`
		SELECT base, COUNT(DISTINCT full_command) as cmd_count, SUM(frequency) as total_runs
		FROM commands
		GROUP BY base
		HAVING cmd_count > 1
		ORDER BY total_runs DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patterns []PatternGroup
	for rows.Next() {
		var p PatternGroup
		var count int
		if err := rows.Scan(&p.BaseCommand, &count, &p.RunCount); err != nil {
			return nil, err
		}

		// Get all commands for this base
		cmdRows, err := db.Query(`
			SELECT full_command FROM commands 
			WHERE base = ? ORDER BY frequency DESC LIMIT 10
		`, p.BaseCommand)
		if err != nil {
			continue
		}
		for cmdRows.Next() {
			var cmd string
			if err := cmdRows.Scan(&cmd); err != nil {
				break
			}
			p.Commands = append(p.Commands, cmd)
		}
		cmdRows.Close()

		patterns = append(patterns, p)
	}
	return patterns, nil
}

// SuggestAlias suggests aliases based on command patterns
func SuggestAliases() ([]string, error) {
	patterns, err := DetectPatterns()
	if err != nil {
		return nil, err
	}

	var aliases []string
	commonCommands := map[string]string{
		"git":           "g",
		"docker":        "d",
		"kubectl":       "k",
		"npm":           "n",
		"yarn":          "y",
		"terraform":     "tf",
		"docker-compose": "dc",
	}

	for _, p := range patterns {
		if alias, ok := commonCommands[p.BaseCommand]; ok {
			aliases = append(aliases, fmt.Sprintf("alias %s='%s'", alias, p.BaseCommand))
		}
	}

	// Suggest based on long commands
	longCmds, err := db.Query(`
		SELECT full_command, frequency FROM commands 
		WHERE LENGTH(full_command) > 20 AND frequency > 2
		ORDER BY frequency DESC LIMIT 5
	`)
	if err != nil {
		return aliases, nil
	}
	defer longCmds.Close()

	for longCmds.Next() {
		var cmd string
		var freq int
		if err := longCmds.Scan(&cmd, &freq); err != nil {
			break
		}
		// Create alias name from command words
		words := strings.Fields(cmd)
		if len(words) >= 2 {
			aliasName := strings.ToLower(words[0][:1] + words[1][:2])
			aliases = append(aliases, fmt.Sprintf("alias %s='%s'", aliasName, cmd))
		}
	}

	return aliases, nil
}
