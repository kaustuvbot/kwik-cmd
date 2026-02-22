#!/bin/bash
# kwik-cmd Bash Integration
# Add to ~/.bashrc: source /path/to/kwik-cmd/shell/bash_hook.sh

# Track commands automatically
kwik_cmd_track() {
    local cmd="$1"
    if [ -n "$cmd" ] && [ -z "${cmd// }" ]; then
        # Don't track empty commands or common builtins
        case "$cmd" in
            cd|ls|pwd|echo|exit|export|readonly|local|declare|typeset|unset|shift)
                return
                ;;
        esac
        
        # Run in background to not slow down the shell
        kwik-cmd track "$cmd" 2>/dev/null &
    fi
}

# For bash, use PROMPT_COMMAND
kwik_preexec() {
    kwik_cmd_track "$BASH_COMMAND"
}

# Only enable if kwik-cmd is installed
if command -v kwik-cmd &> /dev/null; then
    PROMPT_COMMAND="kwik_preexec"
fi
