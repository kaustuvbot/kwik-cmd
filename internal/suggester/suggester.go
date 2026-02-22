package suggester

import (
	"fmt"
	"strings"

	"github.com/kaustuvbot/kwik-cmd/internal/db"
)

// Suggest provides command suggestions for partial input
func Suggest(partial string) error {
	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	partial = strings.TrimSpace(partial)
	if partial == "" {
		// Show recent and top if no input
		return showSuggestions("", 3, 3)
	}

	// Get commands that match the partial input
	recent, err := db.GetRecentCommands(10)
	if err != nil {
		return fmt.Errorf("failed to get recent commands: %w", err)
	}

	top, err := db.GetTopCommands(10)
	if err != nil {
		return fmt.Errorf("failed to get top commands: %w", err)
	}

	// Filter by partial match
	var recentMatched []db.Command
	var topMatched []db.Command

	for _, c := range recent {
		if strings.HasPrefix(strings.ToLower(c.Base), strings.ToLower(partial)) ||
			strings.HasPrefix(strings.ToLower(c.FullCommand), strings.ToLower(partial)) {
			recentMatched = append(recentMatched, c)
		}
	}

	for _, c := range top {
		if strings.HasPrefix(strings.ToLower(c.Base), strings.ToLower(partial)) ||
			strings.HasPrefix(strings.ToLower(c.FullCommand), strings.ToLower(partial)) {
			topMatched = append(topMatched, c)
		}
	}

	// Remove duplicates
	seen := make(map[int64]bool)
	var uniqueRecent []db.Command
	for _, c := range recentMatched {
		if !seen[c.ID] {
			seen[c.ID] = true
			uniqueRecent = append(uniqueRecent, c)
		}
	}

	seen = make(map[int64]bool)
	var uniqueTop []db.Command
	for _, c := range topMatched {
		if !seen[c.ID] {
			seen[c.ID] = true
			uniqueTop = append(uniqueTop, c)
		}
	}

	// Display results
	fmt.Printf("=== Suggestions for '%s' ===\n\n", partial)

	if len(uniqueRecent) > 0 {
		fmt.Println("Recent:")
		for i, c := range uniqueRecent {
			if i >= 5 {
				break
			}
			fmt.Printf("  %d. %s (used %d times)\n", i+1, c.FullCommand, c.Frequency)
		}
		fmt.Println()
	}

	if len(uniqueTop) > 0 {
		fmt.Println("Most Used:")
		for i, c := range uniqueTop {
			if i >= 5 {
				break
			}
			fmt.Printf("  %d. %s (used %d times)\n", i+1, c.FullCommand, c.Frequency)
		}
	}

	if len(uniqueRecent) == 0 && len(uniqueTop) == 0 {
		fmt.Println("No matching commands found.")
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

	// Split keywords
	words := strings.Fields(strings.ToLower(keywords))

	// Get all recent commands and filter
	all, err := db.GetRecentCommands(100)
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	fmt.Printf("=== Search results for '%s' ===\n\n", keywords)

	var matches []db.Command
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
			matches = append(matches, c)
		}
	}

	if len(matches) == 0 {
		// Try keyword search
		matches, err = db.SearchByKeyword(keywords)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}
	}

	if len(matches) == 0 {
		fmt.Println("No matching commands found.")
		return nil
	}

	for i, c := range matches {
		if i >= 20 {
			break
		}
		fmt.Printf("%d. %s\n   Used %d times, last: %s\n\n",
			i+1, c.FullCommand, c.Frequency, c.LastUsed.Format("2006-01-02 15:04"))
	}

	return nil
}

// showSuggestions shows recent and top commands
func showSuggestions(prefix string, recentCount, topCount int) error {
	recent, err := db.GetRecentCommands(recentCount)
	if err != nil {
		return err
	}

	top, err := db.GetTopCommands(topCount)
	if err != nil {
		return err
	}

	fmt.Println("=== Recent Commands ===")
	for i, c := range recent {
		fmt.Printf("%d. %s\n", i+1, c.FullCommand)
	}

	fmt.Println("\n=== Most Used Commands ===")
	for i, c := range top {
		fmt.Printf("%d. %s (used %d times)\n", i+1, c.FullCommand, c.Frequency)
	}

	return nil
}
