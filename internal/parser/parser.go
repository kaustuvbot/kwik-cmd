package parser

import (
	"strings"

	"github.com/samber/lo"
)

// ParsedCommand represents a parsed shell command
type ParsedCommand struct {
	Base       string
	Subcommand string
	FullCmd    string
	Flags      []string
	Args       []string
	Directory  string
}

// ParseCommand parses a shell command into its components
func ParseCommand(cmd string) *ParsedCommand {
	// Remove leading/trailing whitespace
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return nil
	}

	// Split by whitespace
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	base := parts[0]

	// Handle paths (./script, /path/to/command)
	if strings.Contains(base, "/") {
		base = strings.TrimPrefix(base, "./")
		base = strings.TrimPrefix(base, "/")
		if idx := strings.LastIndex(base, "/"); idx >= 0 {
			base = base[idx+1:]
		}
	}

	// Remove common aliases prefixes
	base = strings.TrimPrefix(base, "sudo ")
	base = strings.TrimPrefix(base, "doas ")

	var subcommand string
	var flags []string
	var args []string

	if len(parts) > 1 {
		// Look for subcommand (usually not starting with -)
		for i := 1; i < len(parts); i++ {
			part := parts[i]
			if strings.HasPrefix(part, "-") {
				// It's a flag
				flags = append(flags, part)
				// Handle flags with values (--flag=value, -f value)
				if strings.Contains(part, "=") {
					// Already has value attached
				} else if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "-") {
					// Next part is value
					flags = append(flags, parts[i+1])
					i++
				}
			} else if subcommand == "" {
				// First non-flag is subcommand
				subcommand = part
			} else {
				// Rest are arguments
				args = append(args, part)
			}
		}
	}

	return &ParsedCommand{
		Base:       base,
		Subcommand: subcommand,
		FullCmd:    cmd,
		Flags:      lo.Uniq(flags),
		Args:       args,
	}
}

// ExtractKeywords extracts searchable keywords from a command
func ExtractKeywords(cmd *ParsedCommand) []string {
	keywords := []string{cmd.Base}

	if cmd.Subcommand != "" {
		keywords = append(keywords, cmd.Subcommand)
	}

	// Add common command keywords
	keywords = append(keywords, extractCommonKeywords(cmd.Base, cmd.Subcommand)...)

	return lo.Uniq(keywords)
}

// extractCommonKeywords returns keywords for common commands
func extractCommonKeywords(base, subcommand string) []string {
	keywordsMap := map[string][]string{
		"git": {
			"commit", "push", "pull", "fetch", "merge", "rebase", "branch",
			"checkout", "status", "log", "diff", "add", "stash", "reset",
		},
		"docker": {
			"run", "build", "pull", "push", "ps", "logs", "exec", "stop",
			"start", "rm", "images", "-compose",
		},
		"npm": {
			"install", "run", "start", "test", "build", "publish", "init",
		},
		"yarn": {
			"add", "run", "start", "build", "publish",
		},
		"go": {
			"run", "build", "test", "get", "install", "mod", "fmt", "vet",
		},
		"kubectl": {
			"get", "apply", "delete", "describe", "logs", "exec", "port-forward",
		},
		"terraform": {
			"init", "plan", "apply", "destroy", "validate", "fmt",
		},
		"make": {
			"install", "build", "clean", "test",
		},
		"python": {
			"python", "pip", "virtualenv", "venv",
		},
		"ls": {
			"list", "directory", "files",
		},
		"cd": {
			"change", "directory", "navigate",
		},
		"vim": {
			"editor", "edit", "vim", "vi",
		},
	}

	if subcommand != "" {
		if kws, ok := keywordsMap[base]; ok {
			for _, kw := range kws {
				if strings.Contains(subcommand, kw) {
					return []string{kw}
				}
			}
		}
	}

	return keywordsMap[base]
}

// FlagMeaning returns the meaning of common flags
func FlagMeaning(flag string) string {
	meanings := map[string]string{
		"-m":    "message",
		"--message": "message",
		"-f":    "force",
		"--force":   "force",
		"-r":    "recursive",
		"--recursive": "recursive",
		"-v":    "verbose",
		"--verbose":  "verbose",
		"-d":    "directory",
		"--directory": "directory",
		"-n":    "dry-run",
		"--dry-run":   "dry-run",
		"-y":    "yes",
		"--yes":      "yes",
		"-p":    "port",
		"--port":     "port",
		"-t":    "tag",
		"--tag":      "tag",
		"--help":     "help",
		"--version":  "version",
	}

	return meanings[flag]
}
