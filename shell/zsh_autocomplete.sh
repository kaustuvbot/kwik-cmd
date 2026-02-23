#!/bin/zsh
# kwik-cmd - Integrates with zsh-autosuggestions
# Shows suggestions at bottom, navigate with Tab/arrows

# Source zsh-autosuggestions first
if [ -f "$ZSH_CUSTOM/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" ]; then
    source "$ZSH_CUSTOM/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
elif [ -f "$ZSH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" ]; then
    source "$ZSH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
fi

# Override to use kwik-cmd - extract just the command
_zsh_autosuggest_suggest() {
    local buffer="$1"
    [ -z "$buffer" ] && return
    
    # Get suggestion - skip first 3 lines (header), get first match
    local suggestion
    suggestion=$(kwik-cmd suggest "$buffer" 2>/dev/null | grep -v "===" | grep -v "^$" | grep -v "Ranked by" | head -1 | sed 's/^[ ]*[0-9]*\. //' | sed 's/ (score:.*//')
    
    [ -n "$suggestion" ] && echo "$suggestion"
}

# Enable
ZSH_AUTOSUGGEST_USE_ASYNC=false
ZSH_AUTOSUGGEST_HISTORY_IGNORE="cd |ls |ll |la |pwd |echo |exit |sudo |rm |clear"

echo "kwik-cmd suggestions loaded!"
echo "Type command - suggestions appear at bottom"
echo "Press Tab or Right Arrow to accept"
