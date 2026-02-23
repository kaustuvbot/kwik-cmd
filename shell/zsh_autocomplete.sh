#!/bin/zsh
# kwik-cmd - Show suggestions in RPROMPT (right side of prompt)

# Check if kwik-cmd exists
(( $+commands[kwik-cmd] )) || return

# Update RPROMPT with suggestion
kwik-update-rprompt() {
    local buffer="$BUFFER"
    
    # Don't for empty or very short
    (( ${#buffer} < 1 )) && return
    
    # Skip ignored commands
    local base="${buffer%% *}"
    case "$base" in
        cd|ls|ll|la|l|pwd|echo|exit|export|declare|typeset|unset|shift|\
        local|readonly|help|which|what|time|fg|bg|jobs|kill|builtin|test|\
        true|false|logout|shopt|umask|set|setenv|printenv|eval|exec|source|alias|unalias|\
        sudo|rm|mkdir|touch|cat|grep|sed|awk|find|xargs|sort|uniq|head|tail|wc|clear)
            RPROMPT=""
            return
            ;;
    esac
    
    # Get suggestion
    local suggestion
    suggestion=$(kwik-cmd suggest "$buffer" 2>/dev/null | sed -n '2p' | sed 's/^[ ]*[0-9]*\. //')
    
    if [ -n "$suggestion" ]; then
        RPROMPT="%F{243}$suggestion%f"
    else
        RPROMPT=""
    fi
}

# Accept suggestion widget
kwik-accept-rprompt() {
    if [ -n "$RPROMPT" ]; then
        local suggestion
        suggestion=$(kwik-cmd suggest "$BUFFER" 2>/dev/null | sed -n '2p' | sed 's/^[ ]*[0-9]*\. //')
        if [ -n "$suggestion" ]; then
            BUFFER="$suggestion"
            CURSOR=${#BUFFER}
            RPROMPT=""
        fi
    fi
}

# Create widgets
zle -N kwik-update-rprompt
zle -N kwik-accept-rprompt

# Hook into zle - update on every keypress
zle -N zle-keymap-select kwik-update-rprompt
zle -N zle-line-init kwik-update-rprompt

# Accept with Ctrl+F
bindkey "^F" kwik-accept-rprompt

echo "kwik-cmd RPROMPT suggestions loaded!"
echo "Suggestions appear on the right side of your prompt"
echo "Press Ctrl+F to accept"
