# kwik-cmd Zsh Autocomplete
# Place in ~/.zshrc or source this file

# Initialize completion system
autoload -Uz compinit

# kwik-cmd completion handler
_kwik-cmd-complete() {
    local -a cmd
    cmd=(
        'track:Track a command'
        'suggest:Get suggestions'
        'search:Search commands'
        'stats:Show statistics'
        'analyze:Analyze patterns'
        'export:Export history'
        'import:Import history'
        'quick:Quick pick'
        'rerun:Re-run last command'
        'copy:Copy to clipboard'
        'version:Show version'
        'reset:Reset history'
    )
    
    _describe 'command' cmd
}

# Register completions
compdef _kwik-cmd-complete kwik-cmd

# For command arguments - provide suggestions based on history
_kwik-cmd-arg-complete() {
    local -a suggestions
    local current_word="$WORDS[CURRENT]"
    
    # Get suggestions from kwik-cmd
    if (( $+commands[kwik-cmd] )); then
        suggestions=($(kwik-cmd suggest "$current_word" 2>/dev/null | sed -n 's/^[ ]*[0-9]*\. //p'))
    fi
    
    if (( ${#suggestions} > 0 )); then
        _describe 'command' suggestions
    fi
}

# Setup aliases completion
aliascompletion() {
    compdef _kwik-cmd-arg-complete $@
}

# Bind auto-suggestion widget
if (( $+widgets[kwik-cmd-suggest] )); then
    zle -N kwik-cmd-suggest
fi

kwik-cmd-suggest() {
    local buffer="$BUFFER"
    local -a suggestions
    
    # Get suggestions for current buffer
    if (( $+commands[kwik-cmd] )) && [ -n "$buffer" ]; then
        suggestions=($(kwik-cmd suggest "$buffer" 2>/dev/null | sed -n 's/^[ ]*[0-9]*\. //p'))
        
        if (( ${#suggestions} > 0 )); then
            # Show first suggestion
            zle -M "${suggestions[1]}"
        fi
    fi
}

# Accept suggestion with Ctrl+F
kwik-cmd-accept-suggestion() {
    local buffer="$BUFFER"
    local -a suggestions
    
    if (( $+commands[kwik-cmd] )) && [ -n "$buffer" ]; then
        suggestions=($(kwik-cmd suggest "$buffer" 2>/dev/null | sed -n 's/^[ ]*[0-9]*\. //p'))
        
        if (( ${#suggestions} > 0 )); then
            BUFFER="$suggestions[1]"
            CURSOR=${#BUFFER}
        fi
    fi
}

zle -N kwik-cmd-accept-suggestion

# Key bindings
bindkey "^F" kwik-cmd-accept-suggestion

echo "kwik-cmd autocomplete loaded. Use Ctrl+F to accept suggestions."
