#!/bin/zsh
# kwik-cmd - Integrates with zsh-autosuggestions
# Shows suggestions at bottom, navigate with Tab/arrows

# Source zsh-autosuggestions first
if [ -f "$ZSH_CUSTOM/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" ]; then
    source "$ZSH_CUSTOM/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
elif [ -f "$ZSH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" ]; then
    source "$ZSH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
fi

# Override to use kwik-cmd
_zsh_autosuggest_suggest() {
    local buffer="$1"
    [ -z "$buffer" ] && return
    
    # Get first suggestion from kwik-cmd
    kwik-cmd suggest "$buffer" 2>/dev/null | sed -n '2p' | sed 's/^[ ]*[0-9]*\. //'
}

# Enable
ZSH_AUTOSUGGEST_USE_ASYNC=false
ZSH_AUTOSUGGEST_HISTORY_IGNORE="cd |ls |ll |la |pwd |echo |exit |sudo |rm |clear"

echo "kwik-cmd + zsh-autosuggestions loaded!"
echo "Type command - suggestions appear at bottom"
echo "Press Tab or Right Arrow to accept"
