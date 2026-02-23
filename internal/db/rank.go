package db

import (
	"math"
	"time"
)

// Weighted ranking configuration
const (
	RecencyWeight     = 0.4
	FrequencyWeight   = 0.4
	DirectoryWeight  = 0.2
	DaysSinceMaxScore = 30.0 // Commands older than 30 days get 0 recency score
)

// RankedCommand includes ranking score
type RankedCommand struct {
	Command
	Score float64
}

// GetRankedCommands returns commands sorted by weighted ranking score
func GetRankedCommands(partial, currentDir string, limit int) ([]RankedCommand, error) {
	commands, err := GetRecentCommands(100)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var ranked []RankedCommand

	for _, c := range commands {
		// Calculate recency score (0-1, newer = higher)
		daysSince := now.Sub(c.LastUsed).Hours() / 24
		recencyScore := math.Max(0, 1-(daysSince/DaysSinceMaxScore))

		// Frequency score (0-1, normalize based on max frequency)
		maxFreq := 10 // Assume max frequency of 10 for normalization
		frequencyScore := math.Min(1, float64(c.Frequency)/float64(maxFreq))

		// Directory match score
		directoryScore := 0.0
		if currentDir != "" && c.Directory != "" {
			if c.Directory == currentDir {
				directoryScore = 1.0
			} else if containsParent(currentDir, c.Directory) || containsParent(c.Directory, currentDir) {
				directoryScore = 0.5
			}
		}

		// Partial match bonus
		partialScore := 0.0
		if partial != "" {
			lowerPartial := partial
			lowerBase := c.Base
			lowerFull := c.FullCommand
			if len(lowerPartial) > 0 && (len(lowerBase) >= len(lowerPartial) && lowerBase[:len(lowerPartial)] == lowerPartial ||
				len(lowerFull) >= len(lowerPartial) && lowerFull[:len(lowerPartial)] == lowerPartial) {
				partialScore = 0.3 // Boost exact prefix matches
			}
		}

		// Calculate total score
		score := (RecencyWeight*recencyScore +
			FrequencyWeight*frequencyScore +
			DirectoryWeight*directoryScore +
			partialScore)

		ranked = append(ranked, RankedCommand{
			Command: c,
			Score:   score,
		})
	}

	// Sort by score descending
	for i := 0; i < len(ranked)-1; i++ {
		for j := i + 1; j < len(ranked); j++ {
			if ranked[j].Score > ranked[i].Score {
				ranked[i], ranked[j] = ranked[j], ranked[i]
			}
		}
	}

	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}

	return ranked, nil
}

// containsParent checks if path1 is a parent of path2 or vice versa
func containsParent(path1, path2 string) bool {
	if len(path1) > len(path2) {
		return path2+"/" == path1[:len(path2)+1]
	}
	if len(path2) > len(path1) {
		return path1+"/" == path2[:len(path1)+1]
	}
	return false
}

// AddKeyword adds a keyword to a command
func AddKeyword(commandID int64, keyword string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO keywords (command_id, keyword) VALUES (?, ?)
	`, commandID, keyword)
	return err
}

// AddFlag adds a flag to a command
func AddFlag(commandID int64, flag, meaning string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO flags (command_id, flag, meaning) VALUES (?, ?, ?)
	`, commandID, flag, meaning)
	return err
}

// GetFlags returns all flags for a command
func GetFlags(commandID int64) ([]string, error) {
	rows, err := db.Query("SELECT flag FROM flags WHERE command_id = ?", commandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []string
	for rows.Next() {
		var flag string
		if err := rows.Scan(&flag); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}
	return flags, nil
}
