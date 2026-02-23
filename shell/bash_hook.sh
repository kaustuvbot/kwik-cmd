#!/bin/bash
# kwik-cmd Bash Integration - Auto-tracking
# Add to ~/.bashrc: source /path/to/kwik-cmd/shell/bash_hook.sh

# Commands to ignore
KWIk_IGNORE="cd ls lla ll la pwd echo exit export declare typeset unset shift
local readonly help which what time fg bg jobs kill builtin test [ 
true false logout shopt umask setx setenv printenv eval exec
source alias unalias"

# Check if kwik-cmd exists
command -v kwik-cmd >/dev/null 2>&1 || return

# Auto-track commands
kwik_preexec() {
    local cmd="$BASH_COMMAND"
    
    # Skip empty
    [ -z "$cmd" ] && return
    
    # Skip kwik-cmd itself
    [[ "$cmd" == kwik-cmd* ]] && return
    
    # Skip ignored commands
    local base="${cmd%% *}"
    for ignore in $KWIk_IGNORE; do
        [ "$base" = "$ignore" ] && return
    done
    
    # Track in background
    kwik-cmd track "$cmd" 2>/dev/null &
}

PROMPT_COMMAND="kwik_preexec"
