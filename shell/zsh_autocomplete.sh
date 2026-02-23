#!/bin/zsh
# kwik-cmd Zsh Auto-suggestions
# Shows suggestions as you type, like zsh-autosuggestions

# Check if kwik-cmd exists
(( $+commands[kwik-cmd] )) || return

# Suggestion state
typeset -g KWIK_SUGGESTION=""
typeset -g KWIK_LAST_BUFFER=""

# Get suggestion for current input
kwik-fetch-suggestion() {
    local buffer="$BUFFER"
    local cursor=$CURSOR
    
    # Don't suggest for empty or very short input
    (( ${#buffer} < 1 )) && return
    
    # Skip ignored commands
    local base="${buffer%% *}"
    case "$base" in
        cd|ls|ll|la|l|ls-|pwd|echo|exit|export|declare|typeset|unset|shift|\
        local|readonly|help|which|what|time|fg|bg|jobs|kill|builtin|test|\
        true|false|logout|shopt|umask|set|setenv|printenv|eval|exec|source|alias|unalias|\
        sudo|rm|mkdir|touch|cat|grep|sed|awk|find|xargs|sort|uniq|head|tail|wc)
            return
            ;;
    esac
    
    # Get suggestion
    KWIK_SUGGESTION=$(kwik-cmd suggest "$buffer" 2>/dev/null | sed -n '2p' | sed 's/^[ ]*[0-9]*\. //')
    KWIK_LAST_BUFFER="$buffer"
}

# Widget: Accept suggestion
kwik-accept-suggestion() {
    if [ -n "$KWIK_SUGGESTION" ] && [ "$BUFFER" != "$KWIK_SUGGESTION" ]; then
        BUFFER="$KWIK_SUGGESTION"
        CURSOR=${#BUFFER}
        KWIK_SUGGESTION=""
    fi
}

# Widget: Accept and execute
kwik-accept-and-run() {
    if [ -n "$KWIK_SUGGESTION" ]; then
        BUFFER="$KWIK_SUGGESTION"
        zle accept-line
    fi
}

# Widget: Clear suggestion
kwik-clear-suggestion() {
    KWIK_SUGGESTION=""
}

# Widget: Show suggestion (bind to cursor movement)
kwik-suggest-show() {
    # Only fetch if buffer changed
    if [ "$BUFFER" != "$KWIK_LAST_BUFFER" ]; then
        kwik-fetch-suggestion
    fi
    
    # Show in zle message area
    if [ -n "$KWIK_SUGGESTION" ]; then
        zle -M "==> $KWIK_SUGGESTION"
    else
        zle -M ""
    fi
}

# Create widgets
zle -N kwik-accept-suggestion
zle -N kwik-accept-and-run
zle -N kwik-clear-suggestion
zle -N kwik-suggest-show

# Bind keys
bindkey "^F" kwik-accept-suggestion         # Ctrl+F to accept
bindkey "^E" kwik-accept-suggestion         # Ctrl+E to accept  
bindkey "^[[A" kwik-accept-suggestion      # Up arrow to accept
bindkey "^[f" kwik-accept-suggestion       # Alt+f to accept word

# Hook into zle - show suggestion after cursor moves
zle -N zle-keymap-select kwik-suggest-show
zle -N zle-line-init kwik-suggest-show

# Clear on line finish
zle -N zle-line-finish kwik-clear-suggestion

# Initial load
kwik-fetch-suggestion

echo "kwik-cmd auto-suggestions loaded!"
echo "Type a command and wait - suggestions appear below"
echo "Press Ctrl+F or Up Arrow to accept suggestion"
