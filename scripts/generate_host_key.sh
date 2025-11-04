#!/bin/bash

# Generate SSH host key for PowerHell SSH server

KEY_PATH="../ssh_host_key"

echo "🔑 Generating SSH host key for PowerHell..."

# Generate ED25519 key (recommended for modern SSH)
ssh-keygen -t ed25519 -f "$KEY_PATH" -N "" -C "powerhell-ssh-server"

echo "✅ Host key generated at: $KEY_PATH"
echo ""
echo "To use this key with your SSH server, run:"
echo "  ./powerhell -ssh -hostkey $KEY_PATH"
echo ""
echo "⚠️  Keep this key secure and do not share it!"