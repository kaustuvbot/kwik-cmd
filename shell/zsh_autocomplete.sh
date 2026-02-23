#!/bin/zsh
# kwik-cmd Zsh Auto-suggestions
# Replaces zsh-autosuggestions with kwik-cmd powered suggestions

# Check if kwik-cmd exists
(( $+commands[kwik-cmd] )) || return

# Disable zsh-autosuggestions if running
if (( $+functions[_zsh_autosuggest_disable] )); then
    _zsh_autosuggest_disable
fi

# State
typeset -g KWIK_SUGGESTION=""
typeset -g KWIK_BUFFER=""

# Fetch suggestion
kwik-fetch() {
    local buffer="$BUFFER"
    
    # Only fetch if buffer changed
    [ "$buffer" = "$KWIK_BUFFER" ] && return
    KWIK_BUFFER="$buffer"
    
    # Don't suggest for empty or very short
    (( ${#buffer} < 1 )) && return
    
    # Skip ignored commands
    local base="${buffer%% *}"
    case "$base" in
        cd|ls|ll|la|l|pwd|echo|exit|export|declare|typeset|unset|shift|\
        local|readonly|help|which|what|time|fg|bg|jobs|kill|builtin|test|\
        true|false|logout|shopt|umask|set|setenv|printenv|eval|exec|source|alias|unalias|\
        sudo|rm|mkdir|touch|cat|grep|sed|awk|find|xargs|sort|uniq|head|tail|wc|clear)
            KWIK_SUGGESTION=""
            return
            ;;
    esac
    
    # Get suggestion
    KWIK_SUGGESTION=$(kwik-cmd suggest "$buffer" 2>/dev/null | sed -n '2p' | sed 's/^[ ]*[0-9]*\. //')
    
    # Display suggestion
    if [ -n "$KWIK_SUGGESTION" ]; then
        zle -M "$KWIK_SUGGESTION"
    else
        zle -M ""
    fi
}

# Accept suggestion
kwik-accept() {
    if [ -n "$KWIK_SUGGESTION" ] && [ "$BUFFER" != "$KWIK_SUGGESTION" ]; then
        BUFFER="$KWIK_SUGGESTION"
        CURSOR=${#BUFFER}
        KWIK_SUGGESTION=""
        zle -M ""
    fi
}

# Clear
kwik-clear() {
    KWIK_SUGGESTION=""
    KWIK_BUFFER=""
}

# Create widgets
zle -N kwik-fetch
zle -N kwik-accept
zle -N kwik-clear

# Bind keys
bindkey "^F" kwik-accept
bindkey "^E" kwik-accept

# Hook into zle
zle -N zle-line-init kwik-fetch
zle -N zle-keymap-select kwik-fetch
zle -N zle-line-finish kwik-clear

# Initial fetch
kwik-fetch

echo "kwik-cmd suggestions loaded!"
echo "Type command and wait - suggestions appear below"
echo "Press Ctrl+F or Enter to accept"
