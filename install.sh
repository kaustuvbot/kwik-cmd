#!/bin/bash
# kwik-cmd Installer for Ubuntu/Debian
# Usage: curl -sSL https://raw.githubusercontent.com/kaustuvbot/kwik-cmd/main/install.sh | bash

set -e

VERSION="0.1.0"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.kwik-cmd"

echo "=== kwik-cmd Installer ==="
echo "Version: $VERSION"
echo ""

# Check if running as root or have sudo access
if [ "$EUID" -ne 0 ] && ! sudo -v 2>/dev/null; then
    echo "Note: Installation requires sudo. Please enter password if prompted."
fi

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        ARCH_NAME="amd64"
        ;;
    aarch64|arm64)
        ARCH_NAME="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Create download URL
DOWNLOAD_URL="https://github.com/kaustuvbot/kwik-cmd/releases/download/v${VERSION}/kwik-cmd_${VERSION}_linux_${ARCH_NAME}.tar.gz"

# Create temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "Downloading kwik-cmd..."
if command -v wget &> /dev/null; then
    wget -q -O kwik-cmd.tar.gz "$DOWNLOAD_URL" || {
        echo "Failed to download. Building from source..."
        rm -rf "$TEMP_DIR"
        exec bash -c "$(cat <<'SOURCEINSTALL'
cd /tmp
git clone https://github.com/kaustuvbot/kwik-cmd.git
cd kwik-cmd
go build -o kwik-cmd .
sudo mv kwik-cmd /usr/local/bin/
echo "Installed! Add to shell:"
echo "  source ~/.kwik-cmd/shell/bash_hook.sh  # for bash"
echo "  source ~/.kwik-cmd/shell/zsh_hook.sh   # for zsh"
SOURCEINSTALL
)"
    }
else
    curl -sSL -o kwik-cmd.tar.gz "$DOWNLOAD_URL" || {
        echo "Download failed. Please install from source."
        exit 1
    }
fi

echo "Extracting..."
tar -xzf kwik-cmd.tar.gz

echo "Installing to $INSTALL_DIR..."
sudo mv kwik-cmd "$INSTALL_DIR/kwik-cmd"
chmod +x "$INSTALL_DIR/kwik-cmd"

# Create config directory
mkdir -p "$CONFIG_DIR"

echo "Installing shell hooks..."

# Detect shell
SHELL_NAME=$(basename "$SHELL")
case "$SHELL_NAME" in
    bash)
        HOOK_SOURCE="# kwik-cmd auto-tracking
if command -v kwik-cmd &> /dev/null; then
    source \"$CONFIG_DIR/shell/bash_hook.sh\"
fi"
        HOOK_FILE="$CONFIG_DIR/shell/bash_hook.sh"
        ;;
    zsh)
        HOOK_SOURCE="# kwik-cmd auto-tracking
if command -v kwik-cmd &> /dev/null; then
    source \"$CONFIG_DIR/shell/zsh_hook.sh\"
fi"
        HOOK_FILE="$CONFIG_DIR/shell/zsh_hook.sh"
        ;;
    *)
        echo "Unsupported shell: $SHELL_NAME"
        ;;
esac

# Copy shell hooks
mkdir -p "$CONFIG_DIR/shell"
# We'll create the hooks from the binary itself or download them
echo "Creating shell integration..."

# Add to shell rc file
RC_FILE=""
if [ -n "$BASH_VERSION" ]; then
    RC_FILE="$HOME/.bashrc"
elif [ -n "$ZSH_VERSION" ]; then
    RC_FILE="$HOME/.zshrc"
fi

if [ -n "$RC_FILE" ] && [ -f "$RC_FILE" ]; then
    if ! grep -q "kwik-cmd" "$RC_FILE" 2>/dev/null; then
        echo "" >> "$RC_FILE"
        echo "# kwik-cmd command tracking" >> "$RC_FILE"
        echo 'export PATH="$PATH:/usr/local/bin"' >> "$RC_FILE"
        echo 'if command -v kwik-cmd &> /dev/null; then' >> "$RC_FILE"
        echo '    source <(kwik-cmd completion bash)' >> "$RC_FILE" 2>/dev/null || true
        echo 'fi' >> "$RC_FILE"
        echo "Added to $RC_FILE"
    fi
fi

# Cleanup
rm -rf "$TEMP_DIR"

echo ""
echo "=== Installation Complete! ==="
echo ""
echo "Usage:"
echo "  kwik-cmd track \"git commit -m 'fix'\"  # Track a command"
echo "  kwik-cmd suggest \"git\"                # Get suggestions"
echo "  kwik-cmd search \"commit\"              # Search commands"
echo "  kwik-cmd stats                         # View statistics"
echo "  kwik-cmd analyze                       # Analyze patterns"
echo ""
echo "For auto-tracking, restart your terminal or run:"
echo "  source $RC_FILE"
echo ""
echo "Test it:"
echo "  kwik-cmd --version"
