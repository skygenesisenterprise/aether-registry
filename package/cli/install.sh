#!/bin/bash
set -e

INSTALL_URL="${INSTALL_URL:-https://bank.skygenesisenterprise.com/cli/install}"
VERSION="${VERSION:-v1.0.0}"

echo "Installing Aether Bank CLI ${VERSION}..."

if [ ! -d "$HOME/.bank" ]; then
    mkdir -p "$HOME/.bank"
    echo "Created config directory: $HOME/.bank"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_DIR="$SCRIPT_DIR"

if [ -f "$CLI_DIR/dist/bank-${VERSION}-linux-amd64" ]; then
    BINARY_PATH="$CLI_DIR/dist/bank-${VERSION}-linux-amd64"
else
    echo "Binary not found in dist folder. Building..."
    cd "$CLI_DIR"
    go build -o bank .
    BINARY_PATH="$CLI_DIR/bank"
fi

cp "$BINARY_PATH" "$HOME/.bank/bank"
chmod +x "$HOME/.bank/bank"

if [ -f "$CLI_DIR/config.sample.yaml" ]; then
    if [ ! -f "$HOME/.bank/config.yaml" ]; then
        cp "$CLI_DIR/config.sample.yaml" "$HOME/.bank/config.yaml"
        echo "Created config file: $HOME/.bank/config.yaml"
    fi
fi

echo ""
echo "Installation complete!"
echo ""
echo "Add to your PATH:"
echo "  export PATH=\"\$HOME/.bank:\$PATH\""
echo ""
echo "Or use with full path:"
echo "  $HOME/.bank/bank --help"
echo ""
echo "Quick start:"
echo "  $HOME/.bank/bank auth login --email your@email.com --password yourpassword"
echo "  $HOME/.bank/bank user list"