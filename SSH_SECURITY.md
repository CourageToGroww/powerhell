# PowerHell SSH Security Guide 🔐

This guide explains how to secure your PowerHell SSH server so that:
1. Only authorized users can connect
2. Users can ONLY access PowerHell (no shell escape)
3. The server is protected against common attacks

## Quick Start

```bash
# Run the security setup script
sudo ./scripts/secure_ssh_setup.sh

# Add a user with password
echo "student:securePassword123" >> ~/.powerhell/ssh/passwd

# Start the secure server
./powerhell -ssh -secure -config ~/.powerhell/ssh/config.json
```

## Security Layers

### 1. SSH Authentication

The secure server requires authentication before allowing any connection:

**Password Authentication:**
```bash
# Add users to ~/.powerhell/ssh/passwd
echo "username:password" >> ~/.powerhell/ssh/passwd

# Connect with password
ssh username@server -p 2222
```

**Public Key Authentication:**
```bash
# Add user's public key
mkdir -p ~/.powerhell/ssh/authorized_keys/
cat user_key.pub >> ~/.powerhell/ssh/authorized_keys/username

# Connect with key
ssh -i private_key username@server -p 2222
```

### 2. Application Isolation

Users connecting via SSH:
- ✅ Can ONLY access the PowerHell training application
- ❌ Cannot escape to system shell
- ❌ Cannot execute system commands
- ❌ Cannot access server files
- ❌ Cannot open new SSH sessions

### 3. Network Security

**Firewall Rules (UFW):**
```bash
# Allow SSH on port 2222
sudo ufw allow 2222/tcp

# Rate limit connections
sudo ufw limit 2222/tcp

# Block all other ports
sudo ufw default deny incoming
sudo ufw enable
```

**IP Whitelisting (optional):**
```bash
# Only allow from specific IPs
sudo ufw allow from 192.168.1.0/24 to any port 2222
sudo ufw allow from 10.0.0.5 to any port 2222
```

### 4. Brute Force Protection

**Fail2ban Configuration:**
```bash
# Install fail2ban
sudo apt install fail2ban

# Configuration is auto-created by setup script
sudo systemctl restart fail2ban
```

This will:
- Ban IPs after 3 failed attempts
- Ban duration: 1 hour
- Monitor window: 10 minutes

### 5. User Restrictions

**Create Restricted User:**
```bash
# Create user with no shell access
sudo useradd -r -s /usr/sbin/nologin powerhell

# Run server as restricted user
sudo -u powerhell ./powerhell -ssh -secure
```

### 6. Process Isolation

**Systemd Hardening:**
```ini
# /etc/systemd/system/powerhell-ssh.service
[Service]
# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/home/powerhell/.powerhell
```

## Configuration File

Create `~/.powerhell/ssh/config.json`:
```json
{
  "host": "0.0.0.0",
  "port": 2222,
  "host_key": "/path/to/ssh_host_key",
  "auth": {
    "enable_password": true,
    "password_file": "/path/to/passwd",
    "enable_public_key": true,
    "authorized_keys_dir": "/path/to/authorized_keys/"
  },
  "security": {
    "max_auth_tries": 3,
    "login_grace_time": 120,
    "max_sessions": 10,
    "allowed_users": ["student", "teacher"],
    "denied_users": ["root"]
  }
}
```

## Running Securely

### Development:
```bash
# Basic secure mode
./powerhell -ssh -secure

# With custom config
./powerhell -ssh -secure -config ~/.powerhell/ssh/config.json
```

### Production:
```bash
# Using systemd
sudo systemctl start powerhell-ssh
sudo systemctl enable powerhell-ssh

# Check status
sudo systemctl status powerhell-ssh
```

### Docker:
```dockerfile
FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN adduser -D -s /usr/sbin/nologin powerhell

COPY powerhell /usr/local/bin/
COPY ssh_config.json /etc/powerhell/

USER powerhell
EXPOSE 2222

CMD ["powerhell", "-ssh", "-secure", "-config", "/etc/powerhell/ssh_config.json"]
```

## Monitoring & Logs

### View SSH Connections:
```bash
# Real-time logs
journalctl -u powerhell-ssh -f

# Failed login attempts
grep "Failed" /var/log/powerhell-ssh.log

# Successful logins
grep "Successful" /var/log/powerhell-ssh.log
```

### Monitor Active Sessions:
```bash
# Show active SSH connections
ss -tn sport = :2222

# Count active sessions
ss -tn sport = :2222 | grep ESTAB | wc -l
```

## Security Checklist

- [ ] Generated unique SSH host key
- [ ] Enabled authentication (password or public key)
- [ ] Created restricted user account
- [ ] Configured firewall rules
- [ ] Enabled fail2ban
- [ ] Set up systemd service with hardening
- [ ] Removed default/example credentials
- [ ] Configured logging and monitoring
- [ ] Tested that users cannot escape to shell
- [ ] Verified PowerHell-only access

## Testing Security

### Test Authentication:
```bash
# Should fail without credentials
ssh localhost -p 2222

# Should succeed with valid credentials
ssh student@localhost -p 2222
```

### Test Shell Escape (should all fail):
```bash
# Try to execute commands (won't work)
!ls
!bash
system("ls")

# Try to escape via SSH (won't work)
~C
~!
```

### Test Brute Force Protection:
```bash
# Make multiple failed attempts
for i in {1..5}; do
  ssh wronguser@localhost -p 2222
done

# Check if IP is banned
sudo fail2ban-client status powerhell-ssh
```

## Troubleshooting

### Permission Denied:
- Check user exists in passwd file
- Verify public key is in authorized_keys
- Check file permissions (600 for keys)

### Connection Refused:
- Verify server is running: `ps aux | grep powerhell`
- Check firewall: `sudo ufw status`
- Check port: `ss -tlnp | grep 2222`

### Users Can Access Shell:
- Ensure using secure mode: `-secure` flag
- Verify bubbletea middleware is active
- Check no shell escape in PowerHell app

## Important Notes

1. **No Shell Access**: The secure configuration ensures users can ONLY interact with PowerHell, not the underlying system
2. **Defense in Depth**: Multiple security layers protect against different attack vectors
3. **Regular Updates**: Keep the system and PowerHell updated for security patches
4. **Audit Logs**: Regularly review logs for suspicious activity
5. **Principle of Least Privilege**: Users get minimum necessary access