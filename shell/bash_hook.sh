#!/bin/bash
# kwik-cmd Bash Integration
# Add to ~/.bashrc: source /path/to/kwik-cmd/shell/bash_hook.sh

# Commands to ignore (space-separated)
KWIk_CMD_IGNORE="cd ls lla ll la pwd echo exit export declare typeset unset shift
local readonly help which what time fg bg jobs kill builtin test [ 
true false logout shopt umask setx"

# Track commands automatically
kwik_cmd_track() {
    local cmd="$1"
    
    # Don't track empty commands
    if [ -z "$cmd" ] || [ -z "${cmd// }" ]; then
        return
    fi
    
    # Don't track kwik-cmd itself
    if [[ "$cmd" == kwik-cmd* ]]; then
        return
    fi
    
    # Check if command is in ignore list
    local base_cmd="${cmd%% *}"
    for ignore in $KWIk_CMD_IGNORE; do
        if [ "$base_cmd" = "$ignore" ]; then
            return
        done
    done
    
    # Don't track aliases expansion
    if [[ "$cmd" == *" "* ]] && [[ "$cmd" != *[![:space:]]* ]]; then
        return
    fi
    
    # Run in background to not slow down the shell
    kwik-cmd track "$cmd" 2>/dev/null &
}

# For bash, use PROMPT_COMMAND
kwik_preexec() {
    kwik_cmd_track "$BASH_COMMAND"
}

# Only enable if kwik-cmd is installed
if command -v kwik-cmd &> /dev/null; then
    PROMPT_COMMAND="kwik_preexec"
fi
