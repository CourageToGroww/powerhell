#!/bin/bash

# PowerHell SSH Security Setup Script
# This script helps secure your PowerHell SSH server

set -e

echo "🔐 PowerHell SSH Security Setup"
echo "================================"

# Check if running as root (for firewall rules)
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  Note: Run as root/sudo for firewall configuration"
fi

# 1. Create SSH keys directory
SSH_DIR="$HOME/.powerhell/ssh"
mkdir -p "$SSH_DIR"
chmod 700 "$SSH_DIR"

# 2. Generate host key if it doesn't exist
HOST_KEY="$SSH_DIR/ssh_host_ed25519_key"
if [ ! -f "$HOST_KEY" ]; then
    echo "🔑 Generating SSH host key..."
    ssh-keygen -t ed25519 -f "$HOST_KEY" -N "" -C "powerhell@$(hostname)"
    chmod 600 "$HOST_KEY"
    echo "✅ Host key generated: $HOST_KEY"
else
    echo "✅ Host key already exists: $HOST_KEY"
fi

# 3. Create authorized_keys directory
AUTH_KEYS_DIR="$SSH_DIR/authorized_keys"
mkdir -p "$AUTH_KEYS_DIR"
chmod 700 "$AUTH_KEYS_DIR"

# 4. Create password file (for demo - use proper auth in production)
PASSWD_FILE="$SSH_DIR/passwd"
if [ ! -f "$PASSWD_FILE" ]; then
    echo "📝 Creating password file..."
    cat > "$PASSWD_FILE" <<EOF
# PowerHell SSH Users
# Format: username:password (use bcrypt hashes in production)
# Example:
# student:learningPowerHell123
# admin:superSecurePass456
EOF
    chmod 600 "$PASSWD_FILE"
    echo "✅ Password file created: $PASSWD_FILE"
    echo "⚠️  Add users to this file in format: username:password"
else
    echo "✅ Password file exists: $PASSWD_FILE"
fi

# 5. Create systemd service file
if [ "$EUID" -eq 0 ]; then
    echo "📄 Creating systemd service..."
    cat > /etc/systemd/system/powerhell-ssh.service <<EOF
[Unit]
Description=PowerHell Secure SSH Server
After=network.target

[Service]
Type=simple
User=powerhell
Group=powerhell
ExecStart=/usr/local/bin/powerhell -ssh -secure -hostkey $HOST_KEY -port 2222
Restart=always
RestartSec=10

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/home/powerhell/.powerhell

[Install]
WantedBy=multi-user.target
EOF
    echo "✅ Systemd service created"
fi

# 6. Create restricted user for running the service
if [ "$EUID" -eq 0 ]; then
    if ! id -u powerhell >/dev/null 2>&1; then
        echo "👤 Creating powerhell user..."
        useradd -r -s /usr/sbin/nologin -d /home/powerhell -m powerhell
        echo "✅ User 'powerhell' created"
    else
        echo "✅ User 'powerhell' already exists"
    fi
fi

# 7. Configure firewall (UFW)
if [ "$EUID" -eq 0 ] && command -v ufw >/dev/null 2>&1; then
    echo "🔥 Configuring firewall..."
    
    # Allow SSH on custom port
    ufw allow 2222/tcp comment "PowerHell SSH"
    
    # Limit connections to prevent brute force
    ufw limit 2222/tcp
    
    echo "✅ Firewall rules added"
    echo "ℹ️  Don't forget to enable UFW: sudo ufw enable"
fi

# 8. Create fail2ban configuration
if [ "$EUID" -eq 0 ] && command -v fail2ban-server >/dev/null 2>&1; then
    echo "🛡️  Creating fail2ban configuration..."
    cat > /etc/fail2ban/jail.d/powerhell-ssh.conf <<EOF
[powerhell-ssh]
enabled = true
port = 2222
filter = sshd
logpath = /var/log/powerhell-ssh.log
maxretry = 3
bantime = 3600
findtime = 600
EOF
    echo "✅ Fail2ban configuration created"
fi

# 9. Create configuration file
CONFIG_FILE="$SSH_DIR/config.json"
cat > "$CONFIG_FILE" <<EOF
{
  "host": "0.0.0.0",
  "port": 2222,
  "host_key": "$HOST_KEY",
  "auth": {
    "enable_password": true,
    "password_file": "$PASSWD_FILE",
    "enable_public_key": true,
    "authorized_keys_dir": "$AUTH_KEYS_DIR"
  },
  "security": {
    "max_auth_tries": 3,
    "login_grace_time": 120,
    "max_sessions": 10,
    "allowed_users": ["student", "admin"],
    "denied_users": ["root"]
  }
}
EOF
chmod 600 "$CONFIG_FILE"
echo "✅ Configuration file created: $CONFIG_FILE"

# 10. Create example authorized_keys file
EXAMPLE_USER="student"
EXAMPLE_AUTH_KEYS="$AUTH_KEYS_DIR/$EXAMPLE_USER"
if [ ! -f "$EXAMPLE_AUTH_KEYS" ]; then
    echo "📝 Creating example authorized_keys for user: $EXAMPLE_USER"
    cat > "$EXAMPLE_AUTH_KEYS" <<EOF
# Add SSH public keys here, one per line
# Format: ssh-rsa AAAAB3... user@host
# or: ssh-ed25519 AAAAC3... user@host
EOF
    chmod 600 "$EXAMPLE_AUTH_KEYS"
fi

echo ""
echo "🎉 Security setup complete!"
echo ""
echo "📋 Next steps:"
echo "1. Add users to: $PASSWD_FILE"
echo "2. Add SSH public keys to: $AUTH_KEYS_DIR/username"
echo "3. Review configuration: $CONFIG_FILE"
echo "4. Start the secure server:"
echo "   ./powerhell -ssh -secure -config $CONFIG_FILE"
echo ""
echo "🔒 Security features enabled:"
echo "✅ SSH authentication required"
echo "✅ Host key verification"
echo "✅ Firewall rules (if root)"
echo "✅ Fail2ban protection (if available)"
echo "✅ Restricted user account"
echo "✅ No shell access - PowerHell only"