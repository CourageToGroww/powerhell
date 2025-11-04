package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

// Config holds SSH server configuration
type Config struct {
	Host       string
	Port       int
	HostKeyPEM []byte
}

// SSHServer represents the SSH server instance
type SSHServer struct {
	config Config
	server *ssh.Server
}

// NewSSHServer creates a new SSH server instance
func NewSSHServer(config Config) *SSHServer {
	return &SSHServer{
		config: config,
	}
}

// Start starts the SSH server
func (s *SSHServer) Start(teaHandler func(ssh.Session) (tea.Model, []tea.ProgramOption)) error {
	// Create server with wish middleware
	server, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)),
		wish.WithHostKeyPEM(s.config.HostKeyPEM),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Allows for better terminal support
			logging.Middleware(),    // Logs SSH connections
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create SSH server: %w", err)
	}

	s.server = server

	// Start server in goroutine
	go func() {
		log.Printf("🔥 PowerHell SSH server starting on %s:%d", s.config.Host, s.config.Port)
		log.Printf("🚀 Users can connect with: ssh -p %d %s", s.config.Port, s.config.Host)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("SSH server error: %s", err)
		}
	}()

	// Wait for interrupt signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	// Graceful shutdown
	log.Println("Shutting down SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	
	return nil
}

// GenerateHostKey generates a new ED25519 host key
func GenerateHostKey() []byte {
	// This is a simple example host key. In production, you should:
	// 1. Generate a proper key using ssh-keygen
	// 2. Store it securely
	// 3. Load it from a file
	
	// Example ED25519 private key (DO NOT USE IN PRODUCTION)
	return []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBrH1MUrUXGhGqL1PeOr0fTayW4Ej6n+q4rAiJFHJCqkwAAAJgGPK5SBjyu
UgAAAAtzc2gtZWQyNTUxOQAAACBrH1MUrUXGhGqL1PeOr0fTayW4Ej6n+q4rAiJFHJCqkw
AAAECh8EzpGzW+M2DdNVMz+LMJr9v6oCPsu0K3pBKPr8T4nGsfUxStRcaEaovU946vR9Nr
JbgSPqf6risCIkUckKqTAAAAEGZvb0BleGFtcGxlLmNvbQECAw==
-----END OPENSSH PRIVATE KEY-----`)
}