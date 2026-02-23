#!/bin/bash
# kwik-cmd Installer for Ubuntu/Debian
# Usage: curl -sSL https://raw.githubusercontent.com/kaustuvbot/kwik-cmd/main/install.sh | bash

set -e

VERSION="0.1.0"
INSTALL_DIR="/usr/local/bin"

echo "=== kwik-cmd Installer ==="
echo "Version: $VERSION"
echo ""

# Detect shell
DETECTED_SHELL=""
if [ -n "$ZSH_VERSION" ]; then
    DETECTED_SHELL="zsh"
elif [ -n "$BASH_VERSION" ]; then
    DETECTED_SHELL="bash"
else
    # Fallback: check login shell
    DETECTED_SHELL=$(basename "$SHELL" 2>/dev/null || echo "bash")
fi

echo "Detected shell: $DETECTED_SHELL"

# Determine RC file and hook
case "$DETECTED_SHELL" in
    zsh)
        RC_FILE="$HOME/.zshrc"
        HOOK_NAME="zsh_hook.sh"
        ;;
    bash|sh|*)
        RC_FILE="$HOME/.bashrc"
        HOOK_NAME="bash_hook.sh"
        ;;
esac

echo "Using RC file: $RC_FILE"

# Install binary
echo "Installing kwik-cmd to $INSTALL_DIR..."

# Try to download first
DOWNLOADED=false
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH_NAME="amd64" ;;
    aarch64|arm64) ARCH_NAME="arm64" ;;
    *) ARCH_NAME="amd64" ;;
esac

TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

if command -v curl &> /dev/null; then
    curl -sSL -o kwik-cmd "https://github.com/kaustuvbot/kwik-cmd/releases/download/v${VERSION}/kwik-cmd_${VERSION}_linux_${ARCH_NAME}" 2>/dev/null && DOWNLOADED=true
fi

if [ "$DOWNLOADED" = true ] && [ -f kwik-cmd ]; then
    chmod +x kwik-cmd
    sudo mv kwik-cmd "$INSTALL_DIR/kwik-cmd"
else
    # Build from source
    echo "Building from source..."
    cd /tmp
    rm -rf kwik-cmd
    git clone --depth 1 https://github.com/kaustuvbot/kwik-cmd.git
    cd kwik-cmd
    go build -o kwik-cmd .
    sudo mv kwik-cmd "$INSTALL_DIR/kwik-cmd"
    rm -rf /tmp/kwik-cmd
fi

cd "$HOME"

# Copy shell hook and autocomplete
mkdir -p "$HOME/.kwik-cmd/shell"
cp "./kwik-cmd/shell/$HOOK_NAME" "$HOME/.kwik-cmd/shell/" 2>/dev/null || true
cp "./kwik-cmd/shell/zsh_autocomplete.sh" "$HOME/.kwik-cmd/shell/" 2>/dev/null || true
cp "/tmp/kwik-cmd/shell/$HOOK_NAME" "$HOME/.kwik-cmd/shell/" 2>/dev/null || true
cp "/tmp/kwik-cmd/shell/zsh_autocomplete.sh" "$HOME/.kwik-cmd/shell/" 2>/dev/null || true

# Add to RC file if not already present
HOOK_LINE="[ -f \"\$HOME/.kwik-cmd/shell/${HOOK_NAME}\" ] && source \"\$HOME/.kwik-cmd/shell/${HOOK_NAME}\""

# Add autocomplete for zsh
if [ "$DETECTED_SHELL" = "zsh" ]; then
    AUTOCOMPLETE_LINE="[ -f \"\$HOME/.kwik-cmd/shell/zsh_autocomplete.sh\" ] && source \"\$HOME/.kwik-cmd/shell/zsh_autocomplete.sh\""
fi

if [ -f "$RC_FILE" ]; then
    if ! grep -q "kwik-cmd/shell" "$RC_FILE" 2>/dev/null; then
        echo "" >> "$RC_FILE"
        echo "# kwik-cmd command tracking" >> "$RC_FILE"
        echo "$HOOK_LINE" >> "$RC_FILE"
        if [ "$DETECTED_SHELL" = "zsh" ]; then
            echo "$AUTOCOMPLETE_LINE" >> "$RC_FILE"
        fi
        echo "Added hook to $RC_FILE"
    else
        echo "Hook already present in $RC_FILE"
    fi
else
    # Create RC file with hook
    echo "# kwik-cmd command tracking" > "$RC_FILE"
    echo "$HOOK_LINE" >> "$RC_FILE"
    echo "Created $RC_FILE with hook"
fi

echo ""
echo "=== Installation Complete! ==="
echo ""
echo "Restart your terminal or run:"
echo "  source $RC_FILE"
echo ""
echo "Usage:"
echo "  kwik-cmd track \"git commit -m 'fix'\"  # Track a command"
echo "  kwik-cmd suggest \"git\"                # Get suggestions"
echo "  kwik-cmd stats                         # View statistics"
echo ""
echo "Test:"
echo "  kwik-cmd --version"
