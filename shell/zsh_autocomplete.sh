#!/bin/zsh
# kwik-cmd - Zsh auto-suggestions integration
# Mode 1: Inline ghost text (automatic as you type)
# Mode 2: Tab expands to full list picker (fzf or numbered menu)
# Global: Ctrl+R for keyword search

# ============================================================
# Setup
# ============================================================

# Source zsh-autosuggestions first
if [ -f "$ZSH_CUSTOM/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" ]; then
    source "$ZSH_CUSTOM/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
elif [ -f "$ZSH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" ]; then
    source "$ZSH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
fi

# ============================================================
# Mode 1: Inline ghost text using custom strategy
# ============================================================

# Custom strategy: receives typed buffer, sets global suggestion
_zsh_autosuggest_strategy_kwik() {
    local prefix="$1"
    [[ -z "$prefix" ]] && return

    local suggestion
    suggestion=$(kwik-cmd suggest --plain --limit 1 "$prefix" 2>/dev/null | head -1)

    # If we have a suggestion, append hint about Tab
    if [[ -n "$suggestion" ]]; then
        typeset -g suggestion="$suggestion"
    fi
}

# Use our custom strategy first, fall back to history
ZSH_AUTOSUGGEST_STRATEGY=(kwik history)
ZSH_AUTOSUGGEST_USE_ASYNC=true
ZSH_AUTOSUGGEST_HISTORY_IGNORE="cd |ls |ll |la |pwd |echo |exit |sudo |rm |clear"

# Custom key bindings for accepting suggestions with hints
_zsh_autosuggest_accept() {
    local suggestion="$suggestion"

    # Accept the suggestion
    zle autosuggest-execute

    # If there are more options available (user typed partial), show hint
    local prefix="$BUFFER"
    if [[ -n "$prefix" && "$suggestion" != "$prefix" ]]; then
        # Show notification that Tab can show more options
        zle -M "Press Tab for more options"
    fi
}

zle -N autosuggest-execute _zsh_autosuggest_accept

# ============================================================
# Mode 2: Tab expands to full list picker (improved)
# ============================================================

_kwik_tab_expand() {
    local prefix="$BUFFER"

    # Don't expand for empty buffer
    [[ -z "$prefix" ]] && zle expand-or-complete && return

    # Get suggestions with split mode (recent + frequent)
    local suggestions
    suggestions=$(kwik-cmd suggest --split --limit 10 "$prefix" 2>/dev/null)

    [[ -z "$suggestions" ]] && zle expand-or-complete && return

    # Parse suggestions into arrays
    local -a recent=()
    local -a frequent=()
    local section="recent"

    while IFS= read -r line; do
        case "$line" in
            "---RECENT---") section="recent" ;;
            "---FREQUENT---") section="frequent" ;;
            *)
                if [[ -n "$line" ]]; then
                    if [[ "$section" == "recent" ]]; then
                        recent+=("$line")
                    else
                        frequent+=("$line")
                    fi
                fi
                ;;
        esac
    done <<< "$suggestions"

    # If only one suggestion, accept it directly
    if [[ ${#recent[@]} -eq 0 && ${#frequent[@]} -eq 0 ]]; then
        zle expand-or-complete
        return
    fi

    local total=0
    for cmd in "$recent[@]"; do ((total++)); done
    for cmd in "$frequent[@]"; do ((total++)); done

    if [[ $total -eq 1 ]]; then
        local first_cmd
        first_cmd="${recent[1]:-${frequent[1]}}"
        if [[ "$first_cmd" == "$prefix"* ]]; then
            BUFFER="$first_cmd"
            CURSOR=${#BUFFER}
            zle reset-prompt
            return
        fi
    fi

    # Try fzf if available
    if command -v fzf &>/dev/null; then
        local all_suggestions=()
        for cmd in "$recent[@]"; do all_suggestions+=("$cmd"); done
        for cmd in "$frequent[@]"; do all_suggestions+=("$cmd"); done

        local selected
        selected=$(printf '%s\n' "${all_suggestions[@]}" | \
            FZF_DEFAULT_OPTS="--height=50% --reverse --nth=1 --prompt='kwik> '" \
            fzf --expect=enter 2>/dev/null)

        local key cmd
        key=$(echo "$selected" | head -1)
        cmd=$(echo "$selected" | tail -1)

        if [[ -n "$cmd" ]]; then
            BUFFER="$cmd"
            CURSOR=${#BUFFER}
            zle reset-prompt
        fi
    else
        # Fallback: interactive numbered menu
        _kwik_interactive_menu prefix recent frequent
    fi
}

# Interactive menu with proper key handling
_kwik_interactive_menu() {
    local prefix="$1"
    shift
    local -a recent=("$@")

    local num_recent=${#recent[@]}
    local num_frequent=$(($# - num_recent))

    # Save/restore terminal state
    local saved_tty
    saved_tty=$(stty -g 2>/dev/null)

    # Show menu
    echo ""
    echo "\033[1;36m  Recently Used:\033[0m"
    local i=1
    for cmd in "$recent[@]"; do
        if [[ $i -eq 1 ]]; then
            echo "\033[32m  → [$i] $cmd\033[0m"
        else
            echo "    [$i] $cmd"
        fi
        ((i++))
    done

    if [[ $# -gt $num_recent ]]; then
        echo ""
        echo "\033[1;35m  Most Frequently Used:\033[0m"
        for ((j=1; j<=num_frequent; j++)); do
            local idx=$((num_recent + j))
            echo "    [$idx] ${@[$idx]}"
        done
    fi

    echo ""
    echo "\033[90m  [↑↓] navigate  [1-$((num_recent + num_frequent))] select  [Enter] accept  [Esc] cancel\033[0m"

    # Read single key with timeout
    local key
    local char
    local timeout=0

    # Use a simple read approach
    read -k 1 char

    # Restore terminal
    stty "$saved_tty" 2>/dev/null

    case "$char" in
        $'\e')  # Escape - cancel
            echo ""
            ;;
        $'\r'|$'\n')  # Enter - accept first
            BUFFER="$recent[1]"
            CURSOR=${#BUFFER}
            ;;
        [0-9])
            local num="$char"
            # Check for multi-digit
            if [[ $(($num_recent + num_frequent)) -gt 9 ]]; then
                local next_char
                read -k 1 next_char 2>/dev/null
                if [[ "$next_char" =~ [0-9] ]]; then
                    num="$num$next_char"
                fi
            fi

            if [[ $num -le $(($num_recent + num_frequent)) && $num -ge 1 ]]; then
                if [[ $num -le $num_recent ]]; then
                    BUFFER="$recent[$num]"
                else
                    BUFFER="${@[$num]}"
                fi
                CURSOR=${#BUFFER}
            fi
            ;;
        $'\x1b')  # Arrow keys
            read -k 1 char 2>/dev/null
            read -k 1 char 2>/dev/null
            # For now, just refresh and show menu again on arrow
            ;;
    esac

    zle reset-prompt
}

# Register the widget
zle -N kwik-tab-expand _kwik_tab_expand

# Bind Tab to our custom widget
bindkey "^I" kwik-tab-expand

# ============================================================
# Global: Ctrl+R for keyword search
# ============================================================

_kwik_keyword_search() {
    # Read search keywords from command line
    echo ""
    echo "\033[1;33m  kwik-cmd keyword search:\033[0m"
    echo -n "  > "
    local keywords
    read keywords

    [[ -z "$keywords" ]] && zle reset-prompt && return

    local results
    results=$(kwik-cmd search --plain --limit 15 "$keywords" 2>/dev/null)

    [[ -z "$results" ]] && echo "  No results found" && zle reset-prompt && return

    local count
    count=$(echo "$results" | wc -l)

    if [[ $count -eq 1 ]]; then
        BUFFER="$results"
        CURSOR=${#BUFFER}
        zle reset-prompt
        return
    fi

    # Use fzf if available for search results
    if command -v fzf &>/dev/null; then
        local selected
        selected=$(echo "$results" | \
            FZF_DEFAULT_OPTS="--height=50% --reverse --prompt='search> '" \
            fzf --expect=enter 2>/dev/null)

        local cmd
        cmd=$(echo "$selected" | tail -1)

        if [[ -n "$cmd" ]]; then
            BUFFER="$cmd"
            CURSOR=${#BUFFER}
        fi
    else
        # Numbered menu for search results
        echo ""
        local i=1
        while IFS= read -r line; do
            if [[ -n "$line" ]]; then
                if [[ $i -eq 1 ]]; then
                    echo "\033[32m  → [$i] $line\033[0m"
                else
                    echo "    [$i] $line"
                fi
                ((i++))
            fi
        done <<< "$results"

        echo ""
        echo "\033[90m  [1-9] select  [Enter] accept first  [Esc] cancel\033[0m"

        local char
        read -k 1 char

        case "$char" in
            [1-9])
                local num=$char
                local line
                line=$(echo "$results" | sed -n "${num}p")
                if [[ -n "$line" ]]; then
                    BUFFER="$line"
                    CURSOR=${#BUFFER}
                fi
                ;;
        esac
    fi

    zle reset-prompt
}

zle -N kwik-keyword-search _kwik_keyword_search
bindkey "^R" kwik-keyword-search

# ============================================================
# Configuration
# ============================================================

echo "kwik-cmd suggestions loaded!"
echo "  Type command prefix → inline ghost suggestion appears"
echo "  Press Tab → full list picker (recent + frequent)"
echo "  Press Right Arrow → accept suggestion"
echo "  Ctrl+R → keyword search"
