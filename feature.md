# Command Suggestions Feature

## Overview

kwik-cmd provides intelligent, context-aware command suggestions in Zsh. It learns from your command history (recency, frequency, directory context) and suggests relevant commands as you type.

## Two-Mode UX

### Mode 1: Inline Ghost Text (Automatic)

As you type a command prefix, kwik-cmd automatically shows the top suggestion in gray/dim text after your cursor.

```
➜  ~ g[it commit -m 'fix']
```

- **Trigger**: Typing any command prefix
- **Accept**: Press `→` (Right Arrow) or `End` key
- **More Options**: Press `Tab` to see all suggestions

### Mode 2: Full List Picker (Tab)

Press `Tab` to open a full list of suggestions organized by recent and frequent usage.

**With fzf (if installed):**
```
➜  ~ g[Tab pressed]
╭──────────────────────────────╮
│ > git commit -m 'fix'       │
│   git push origin main       │
│   git status                 │
│   git log --oneline          │
╰──────────────────────────────╯
```

**Without fzf (numbered menu):**
```
➜  ~ g
  Recently Used:
  → [1] git commit -m 'fix'
    [2] git push origin main
    [3] git status

  Most Frequently Used:
    [4] git add .
    [5] git pull

  [↑↓] navigate  [1-5] select  [Enter] accept  [Esc] cancel
```

- **Navigate**: Arrow keys or type number
- **Accept**: `Enter` or number key
- **Cancel**: `Esc`

## Global Keyword Search (Ctrl+R)

Press `Ctrl+R` anywhere to search your command history by keywords.

```
➜  ~ [Ctrl+R pressed]
  kwik-cmd keyword search:
  > commit fix

  → [1] git commit -m 'fix bug'
    [2] git commit -m 'init'
    [3] npm commit

  [1-9] select  [Enter] accept first  [Esc] cancel
```

Type keywords, see matching commands, select with number or arrows.

## Shell Integration Files

| File | Purpose |
|------|---------|
| `shell/zsh_autocomplete.sh` | Zsh suggestion integration |
| `shell/zsh_hook.sh` | Command tracking hook |

### Installation

Add to your `.zshrc`:

```zsh
# Track commands automatically
source /path/to/kwik-cmd/shell/zsh_hook.sh

# Enable suggestions
source /path/to/kwik-cmd/shell/zsh_autocomplete.sh
```

### Requirements

- Zsh 5.0+
- zsh-autosuggestions plugin (for ghost text rendering)
- Optional: fzf (for enhanced Tab picker and search)

## CLI Commands

### Suggest

```bash
# Colored output with rankings
kwik-cmd suggest "git"

# Plain output (for shell integration)
kwik-cmd suggest --plain "git"

# Split output (recent + frequent)
kwik-cmd suggest --split "git"

# With limit
kwik-cmd suggest --plain --limit 5 "git"
```

### Search

```bash
# Colored search results
kwik-cmd search "commit fix"

# Plain output (for shell integration)
kwik-cmd search --plain "commit fix"

# With limit
kwik-cmd search --plain --limit 20 "commit"
```

## Technical Details

### Strategy: `_zsh_autosuggest_strategy_kwik`

The inline ghost text uses zsh-autosuggestions' strategy API:

```zsh
_zsh_autosuggest_strategy_kwik() {
    local prefix="$1"
    typeset -g suggestion="$(kwik-cmd suggest --plain --limit 1 "$prefix")"
}
```

### Widget: `kwik-tab-expand`

The Tab picker is a custom ZLE widget that:
1. Gets current buffer content
2. Queries kwik-cmd for suggestions (using `--split` mode)
3. Launches fzf (or numbered menu fallback)
4. Replaces buffer with selected command

### Widget: `kwik-keyword-search`

The Ctrl+R search is a custom ZLE widget that:
1. Prompts for keywords
2. Queries kwik-cmd search
3. Shows results in fzf or numbered menu
4. Fills buffer with selected command

## Future Work

- [ ] Arrow key navigation in numbered menu (currently only fzf supports this)
- [ ] Preview window showing command details
- [ ] Bash integration parity
- [ ] Fish shell support
- [ ] Add commands to history without executing
