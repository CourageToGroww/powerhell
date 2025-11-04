# PowerHell SSH Server Setup 🔥

PowerHell can now be accessed remotely via SSH! Users can connect to your server and experience the interactive PowerShell training environment through their terminal.

## Quick Start

### 1. Build the Application
```bash
make build
# or
go build -o powerhell cmd/powerhell/main.go
```

### 2. Generate SSH Host Key (First Time Only)
```bash
make generate-key
# or
cd scripts && ./generate_host_key.sh
```

### 3. Run the SSH Server
```bash
# Default port (2222)
make ssh-server

# Or with custom options
./powerhell -ssh -port 3333 -host 0.0.0.0
```

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-ssh` | Run as SSH server instead of local app | false |
| `-host` | SSH server host address | 0.0.0.0 |
| `-port` | SSH server port | 2222 |
| `-hostkey` | Path to SSH host key file | (uses example key) |

## Connecting to the Server

Once the server is running, users can connect with:

```bash
# Local connection
ssh -p 2222 localhost

# Remote connection (replace YOUR_SERVER_IP)
ssh -p 2222 YOUR_SERVER_IP

# If you get a PTY error, use -t flag
ssh -t -p 2222 localhost
```

## Security Considerations

### Host Keys
- **Always generate a proper host key** for production use
- The example key in the code is for development only
- Keep your host key file secure and never commit it to version control

### Firewall
- Only expose the SSH port to trusted networks
- Consider using a VPN for remote access
- Use fail2ban or similar tools to prevent brute force attacks

### Authentication
- Currently, the server accepts any connection (no authentication)
- For production, implement proper authentication:
  - Password authentication
  - Public key authentication
  - Integration with existing user systems

## Advanced Configuration

### Running Behind a Proxy
```bash
# Using nginx as a TCP proxy
stream {
    server {
        listen 22;
        proxy_pass localhost:2222;
    }
}
```

### Systemd Service
Create `/etc/systemd/system/powerhell-ssh.service`:
```ini
[Unit]
Description=PowerHell SSH Server
After=network.target

[Service]
Type=simple
User=powerhell
ExecStart=/usr/local/bin/powerhell -ssh -hostkey /etc/powerhell/ssh_host_key
Restart=always

[Install]
WantedBy=multi-user.target
```

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o powerhell cmd/powerhell/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/powerhell /usr/local/bin/
EXPOSE 2222
CMD ["powerhell", "-ssh"]
```

## Features

- ✅ Full PowerHell training experience over SSH
- ✅ Animated intro screen works perfectly
- ✅ All navigation and interactions supported
- ✅ Automatic terminal size detection
- ✅ Clean disconnection handling

## Troubleshooting

### "No active terminal" Error
Use the `-t` flag when connecting:
```bash
ssh -t -p 2222 localhost
```

### Connection Refused
- Check if the server is running: `ps aux | grep powerhell`
- Check if the port is open: `netstat -tlnp | grep 2222`
- Check firewall rules: `sudo ufw status`

### Host Key Warnings
If you regenerate the host key, users will see a warning. They need to remove the old key:
```bash
ssh-keygen -R "[localhost]:2222"
```

## Future Enhancements

1. **Authentication System**
   - User accounts with progress tracking
   - OAuth integration
   - API key support

2. **Multi-User Support**
   - Isolated sessions
   - Shared learning spaces
   - Instructor mode

3. **Session Management**
   - Session persistence
   - Resume learning where left off
   - Session replay

4. **Monitoring**
   - Connection logs
   - Usage statistics
   - Learning analytics