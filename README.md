# kwik-cmd

A high-performance CLI tool written in Go that tracks terminal commands automatically, learns usage patterns, and suggests intelligent command completions.

## Features

- Command Tracking - Automatically tracks commands executed in your terminal
- Intelligent Suggestions - Ranks commands by recency + frequency + directory context
- Keyword Search - Search commands by keywords with fuzzy matching
- Pattern Detection - Detects command patterns (e.g., git subcommands)
- Failure Analysis - Tracks command success/failure rates
- Alias Suggestions - Suggests aliases based on usage patterns
- Shell Integration - Works with Bash and Zsh

## Installation

### From Source

```bash
git clone https://github.com/kaustuvbot/kwik-cmd.git
cd kwik-cmd
go build -o kwik-cmd .
sudo mv kwik-cmd /usr/local/bin/
```

### Using Install Script

```bash
curl -sSL https://raw.githubusercontent.com/kaustuvbot/kwik-cmd/main/install.sh | bash
```

### Pre-built Binaries

Download from the Releases page.

## Shell Integration

### Bash

Add to your ~/.bashrc:

```bash
source /path/to/kwik-cmd/shell/bash_hook.sh
```

### Zsh

Add to your ~/.zshrc:

```bash
source /path/to/kwik-cmd/shell/zsh_hook.sh
```

## Usage

### Track a command

```bash
kwik-cmd track "git commit -m 'fix bug'"
kwik-cmd track "docker build" --exit-code 0
```

### Get suggestions

```bash
kwik-cmd suggest "git"
kwik-cmd suggest
```

### Search commands

```bash
kwik-cmd search "commit message"
```

### View statistics

```bash
kwik-cmd stats
```

### Analyze patterns

```bash
kwik-cmd analyze
```

### Reset history

```bash
kwik-cmd reset
```

## Additional Commands

### Export/Import

```bash
kwik-cmd export
kwik-cmd export data.csv -f csv
kwik-cmd import backup.json
```

### Quick pick

```bash
kwik-cmd quick
```

### Rerun

```bash
kwik-cmd rerun
kwik-cmd rerun --dry-run
```

### Copy to clipboard

```bash
kwik-cmd copy 1
```

### Interactive mode

```bash
kwik-cmd interactive
```

### Version

```bash
kwik-cmd version
kwik-cmd check-update
```

## Configuration

Config file location: ~/.kwik-cmd/config.yaml

```yaml
database_path: ~/.kwik-cmd/commands.db
max_suggestions: 10
recency_weight: 0.4
frequency_weight: 0.4
directory_weight: 0.2
enable_colors: true
shell_integration: auto
```

## Ranking Algorithm

Suggestions are ranked using a weighted scoring system:

- Recency (40%): More recent commands score higher
- Frequency (40% score higher
-): Frequently used commands Directory Context (20%): Commands used in current directory are boosted

## Database

Commands are stored in ~/.kwik-cmd/commands.db (SQLite)

## License

MIT
