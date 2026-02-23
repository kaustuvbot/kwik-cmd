# kwik-cmd

<p align="center">
  <img src="https://img.shields.io/badge/version-0.1.0-blue" alt="Version">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green" alt="License">
</p>

A high-performance CLI tool written in Go that tracks terminal commands automatically, learns usage patterns, and suggests intelligent command completions. Outperforms traditional shell history (Ctrl+R) with structured understanding and smart ranking.

## ✨ Features

### Core Features
- **Command Tracking** - Automatically tracks commands executed in yourIntelligent Suggestions** - terminal
- ** Ranks commands by recency + frequency + directory context
- **Keyword Search** - Search commands by keywords with fuzzy matching
- **Pattern Detection** - Detects command patterns (e.g., git subcommands)
- **Failure Analysis** - Tracks command success/failure rates
- **Alias Suggestions** - Suggests aliases based on usage patterns

### Additional Features
- **Colored Output** - Beautiful terminal output with colors
- **Export/Import** - Export history to JSON/CSV, import from backup
- **Interactive Mode** - Interactive command picker
- **Quick Pick** - List recent commands for quick copy
- **Rerun** - Re-run previous commands
- **Version Check** - Check for updates
- **Shell Completions** - Bash and Zsh autocomplete

## 🚀 Quick Start

### Installation

#### From Source
```bash
git clone https://github.com/kaustuvbot/kwik-cmd.git
cd kwik-cmd
go build -o kwik-cmd .
sudo mv kwik-cmd /usr/local/bin/
```

#### Using Install Script
```bash
curl -sSL https://raw.githubusercontent.com/kaustuvbot/kwik-cmd/main/install.sh | bash
```

### Shell Integration

#### Bash
Add to your `~/.bashrc`:
```bash
source /path/to/kwik-cmd/shell/bash_hook.sh
```

#### Zsh
Add to your `~/.zshrc`:
```bash
source /path/to/kwik-cmd/shell/zsh_hook.sh
```

Then restart your terminal or run `source ~/.bashrc` (or `~/.zshrc`).

## 📖 Commands

### Tracking
```bash
kwik-cmd track "git commit -m 'fix bug'"     # Track a command
kwik-cmd track "docker build" --exit-code 0   # Track with exit status
kwik-cmd track "npm test" --exit-code 1       # Track failed command
```

### Suggestions
```bash
kwik-cmd suggest "git"     # Get suggestions for partial command
kwik-cmd suggest           # Get all top suggestions
```

### Search
```bash
kwik-cmd search "commit message"  # Search by keywords
kwik-cmd search "docker run"        # Find docker commands
```

### Statistics & Analysis
```bash
kwik-cmd stats      # View usage statistics
kwik-cmd analyze   # Analyze patterns and get alias suggestions
kwik-cmd quick      # Quick pick from recent commands
```

### Export/Import
```bash
kwik-cmd export                 # Export to JSON (default)
kwik-cmd export data.csv -f csv # Export to CSV
kwik-cmd import backup.json     # Import from backup
```

### Utilities
```bash
kwik-cmd rerun             # Re-run the last command
kwik-cmd rerun --dry-run   # Show what would be run
kwik-cmd copy 1            # Copy command #1 to clipboard
kwik-cmd interactive      # Start interactive mode
kwik-cmd version           # Show version
kwik-cmd check-update     # Check for updates
kwik-cmd reset            # Clear all history (⚠️)
```

## 🔧 Configuration

Config file location: `~/.kwik-cmd/config.yaml`

```yaml
database_path: ~/.kwik-cmd/commands.db
max_suggestions: 10
recency_weight: 0.4
frequency_weight: 0.4
directory_weight: 0.2
enable_colors: true
shell_integration: auto
```

## 📊 Ranking Algorithm

Suggestions are ranked using a weighted scoring system:

| Factor | Weight | Description |
|--------|--------|-------------|
| Recency | 40% | More recent commands score higher |
| Frequency | 40% | Frequently used commands score higher |
| Directory Context | 20% | Commands used in current directory are boosted |

## 🗄️ Database

- Location: `~/.kwik-cmd/commands.db`
- Type: SQLite (embedded, no external dependencies)

### Schema
```sql
commands     - Store command history
flags        - Command flags and meanings
keywords     - Extracted keywords for search
usage_stats  - Success/failure tracking
```

## 🏗️ Architecture

```
kwik-cmd/
├── cmd/              # CLI commands
│   ├── root.go       # Root command
│   ├── track.go      # Track command
│   ├── suggest.go    # Suggest command
│   ├── search.go     # Search command
│   ├── stats.go      # Stats command
│   ├── analyze.go    # Analyze command
│   ├── export.go     # Export/import
│   └── ...
├── internal/
│   ├── db/           # Database layer
│   ├── parser/       # Command parsing
│   ├── tracker/      # Tracking logic
│   ├── suggester/    # Suggestion engine
│   ├── config/       # Configuration
│   └── export/       # Export/import
└── shell/            # Shell integration scripts
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

## 🔗 Links

- [GitHub Repository](https://github.com/kaustuvbot/kwik-cmd)
- [Releases](https://github.com/kaustuvbot/kwik-cmd/releases)
- [Report Issues](https://github.com/kaustuvbot/kwik-cmd/issues)
