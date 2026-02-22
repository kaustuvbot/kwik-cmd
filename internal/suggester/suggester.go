package suggester

import (
	"fmt"
	"os"
	"strings"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
)

// Suggest provides command suggestions using ranking engine
func Suggest(partial string) error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get current directory for context awareness
	currentDir, _ := os.Getwd()

	partial = strings.TrimSpace(partial)

	// Use ranking engine for intelligent suggestions
	ranked, err := db.GetRankedCommands(partial, currentDir, 10)
	if err != nil {
		return fmt.Errorf("failed to get ranked commands: %w", err)
	}

	if len(ranked) == 0 {
		fmt.Println("No commands found. Start tracking commands with 'kwik-cmd track <command>'")
		return nil
	}

	fmt.Printf("=== Suggestions for '%s' ===\n", partial)
	fmt.Println("(Ranked by: recency + frequency + directory context)\n")

	for i, rc := range ranked {
		dirTag := ""
		if currentDir != "" && rc.Directory == currentDir {
			dirTag = " [current dir]"
		}
		fmt.Printf("  %d. %s (score: %.2f, used: %d times)%s\n",
			i+1, rc.FullCommand, rc.Score, rc.Frequency, dirTag)
	}

	return nil
}

// Search searches commands by keywords
func Search(keywords string) error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	keywords = strings.TrimSpace(keywords)
	if keywords == "" {
		return fmt.Errorf("please provide search keywords")
	}

	// Get current directory for context
	currentDir, _ := os.Getwd()

	// Split keywords
	words := strings.Fields(strings.ToLower(keywords))

	// Get all recent commands and filter
	all, err := db.GetRecentCommands(100)
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	fmt.Printf("=== Search results for '%s' ===\n\n", keywords)

	var matches []db.RankedCommand
	for _, c := range all {
		cmdLower := strings.ToLower(c.FullCommand)
		match := true
		for _, w := range words {
			if !strings.Contains(cmdLower, w) {
				match = false
				break
			}
		}
		if match {
			// Calculate simple score for search results
			score := float64(c.Frequency) * 0.5
			matches = append(matches, db.RankedCommand{
				Command: c,
				Score:   score,
			})
		}
	}

	if len(matches) == 0 {
		// Try keyword search from DB
		keywordResults, err := db.SearchByKeyword(keywords)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}
		for _, c := range keywordResults {
			matches = append(matches, db.RankedCommand{
				Command: c,
				Score:   float64(c.Frequency),
			})
		}
	}

	if len(matches) == 0 {
		fmt.Println("No matching commands found.")
		return nil
	}

	// Sort by score
	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].Score > matches[i].Score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	for i, rc := range matches {
		if i >= 20 {
			break
		}
		dirTag := ""
		if currentDir != "" && rc.Directory == currentDir {
			dirTag = " [current dir]"
		}
		fmt.Printf("%d. %s\n   Used %d times, last: %s%s\n\n",
			i+1, rc.FullCommand, rc.Frequency, rc.LastUsed.Format("2006-01-02 15:04"), dirTag)
	}

	return nil
}
