#!/bin/zsh
# kwik-cmd Zsh Integration - Auto-tracking
# Add to ~/.zshrc: source /path/to/kwik-cmd/shell/zsh_hook.sh

# Commands to ignore
KWIk_IGNORE="cd ls lla ll la pwd echo exit export declare typeset unset shift local readonly help which what time fg bg jobs kill builtin test [ true false logout shopt umask setx setenv printenv eval exec"

# Check if kwik-cmd is available
[ -z $commands[(i)kwik-cmd] ] && return

# Auto-track commands
preexec() {
    local cmd="$1"
    
    # Skip empty
    [ -z "$cmd" ] && return
    
    # Skip kwik-cmd itself
    [[ "$cmd" == kwik-cmd* ]] && return
    
    # Skip ignored commands
    local base="${cmd%% *}"
    for ignore in $=KWIk_IGNORE; do
        [ "$base" = "$ignore" ] && return
    done
    
    # Track in background
    kwik-cmd track "$cmd" 2>/dev/null &
}

# Source this file to enable auto-tracking
