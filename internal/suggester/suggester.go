package suggester

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kaustuvbot/kwik-cmd/internal/db"
)

var (
	bold    = color.New(color.Bold)
	cyan    = color.New(color.FgCyan)
	green   = color.New(color.FgGreen)
	yellow  = color.New(color.FgYellow)
	magenta = color.New(color.FgMagenta)
	white   = color.New(color.FgWhite)
	dim     = color.New(color.FgBlack)
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
		yellow.Println("No commands found. Start tracking commands with 'kwik-cmd track <command>'")
		return nil
	}

	cyan.Print("=== Suggestions for '")
	white.Print(partial)
	cyan.Println("' ===")
	dim.Print("(Ranked by: recency + frequency + directory context)")

	for i, rc := range ranked {
		// Number in bold cyan
		bold.Printf("  %d. ", i+1)
		
		// Command in green
		green.Print(rc.FullCommand)
		
		// Score in dim
		dim.Printf(" (score: %.2f, used: %d times)", rc.Score, rc.Frequency)
		
		// Current dir tag in magenta
		if currentDir != "" && rc.Directory == currentDir {
			magenta.Print(" [current dir]")
		}
		fmt.Println()
	}

	return nil
}

// SuggestPlain returns command suggestions as plain text (no colors/headers)
// Used by zsh shell integration for inline suggestions
func SuggestPlain(partial string, limit int) ([]string, error) {
	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	currentDir, _ := os.Getwd()
	partial = strings.TrimSpace(partial)

	ranked, err := db.GetRankedCommands(partial, currentDir, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get ranked commands: %w", err)
	}

	commands := make([]string, 0, len(ranked))
	for _, rc := range ranked {
		commands = append(commands, rc.FullCommand)
	}

	return commands, nil
}

// SuggestPlainSplit returns command suggestions split by recent and frequent
// Returns a map with "recent" and "frequent" keys
// Used by zsh shell integration for categorized suggestions
func SuggestPlainSplit(partial string, limit int) (map[string][]string, error) {
	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	partial = strings.TrimSpace(partial)

	// Get recent commands
	recentCmds, err := db.GetRecentCommands(limit * 2)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent commands: %w", err)
	}

	// Get frequent commands
	frequentCmds, err := db.GetTopCommands(limit * 2)
	if err != nil {
		return nil, fmt.Errorf("failed to get frequent commands: %w", err)
	}

	// Filter by partial match if provided
	var recentFiltered []string
	var frequentFiltered []string

	for _, c := range recentCmds {
		if partial == "" || strings.HasPrefix(strings.ToLower(c.Base), strings.ToLower(partial)) ||
			strings.HasPrefix(strings.ToLower(c.FullCommand), strings.ToLower(partial)) {
			recentFiltered = append(recentFiltered, c.FullCommand)
			if len(recentFiltered) >= limit {
				break
			}
		}
	}

	for _, c := range frequentCmds {
		if partial == "" || strings.HasPrefix(strings.ToLower(c.Base), strings.ToLower(partial)) ||
			strings.HasPrefix(strings.ToLower(c.FullCommand), strings.ToLower(partial)) {
			frequentFiltered = append(frequentFiltered, c.FullCommand)
			if len(frequentFiltered) >= limit {
				break
			}
		}
	}

	return map[string][]string{
		"recent":   recentFiltered,
		"frequent": frequentFiltered,
	}, nil
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

	cyan.Print("=== Search results for '")
	white.Print(keywords)
	cyan.Print("' ===")

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
		yellow.Println("No matching commands found.")
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
		bold.Printf("%d. ", i+1)
		green.Println(rc.FullCommand)
		dim.Printf("   Used %d times, last: %s", rc.Frequency, rc.LastUsed.Format("2006-01-02 15:04"))
		if currentDir != "" && rc.Directory == currentDir {
			magenta.Print(" [current dir]")
		}
		fmt.Println()
	}

	return nil
}

// SearchPlain returns search results as plain text (no colors/headers)
// Used by zsh shell integration for keyword search
func SearchPlain(keywords string, limit int) ([]string, error) {
	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	keywords = strings.TrimSpace(keywords)
	if keywords == "" {
		return nil, fmt.Errorf("please provide search keywords")
	}

	words := strings.Fields(strings.ToLower(keywords))
	all, err := db.GetRecentCommands(100)
	if err != nil {
		return nil, fmt.Errorf("failed to get commands: %w", err)
	}

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
			score := float64(c.Frequency) * 0.5
			matches = append(matches, db.RankedCommand{
				Command: c,
				Score:   score,
			})
		}
	}

	if len(matches) == 0 {
		keywordResults, err := db.SearchByKeyword(keywords)
		if err != nil {
			return nil, fmt.Errorf("search failed: %w", err)
		}
		for _, c := range keywordResults {
			matches = append(matches, db.RankedCommand{
				Command: c,
				Score:   float64(c.Frequency),
			})
		}
	}

	if len(matches) == 0 {
		return []string{}, nil
	}

	// Sort by score
	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].Score > matches[i].Score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	commands := make([]string, 0, len(matches))
	for i, rc := range matches {
		if i >= limit {
			break
		}
		commands = append(commands, rc.FullCommand)
	}

	return commands, nil
}
