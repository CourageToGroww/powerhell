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

// AuthConfig holds authentication configuration
type AuthConfig struct {
	// Enable password authentication
	EnablePassword bool
	// Map of username to password (use bcrypt in production)
	Passwords map[string]string
	
	// Enable public key authentication
	EnablePublicKey bool
	// Map of username to authorized public keys
	AuthorizedKeys map[string][]ssh.PublicKey
}

// SecureConfig holds SSH server configuration with auth
type SecureConfig struct {
	Host       string
	Port       int
	HostKeyPEM []byte
	Auth       AuthConfig
}

// SecureSSHServer represents the SSH server instance with authentication
type SecureSSHServer struct {
	config SecureConfig
	server *ssh.Server
}

// NewSecureSSHServer creates a new SSH server instance with authentication
func NewSecureSSHServer(config SecureConfig) *SecureSSHServer {
	return &SecureSSHServer{
		config: config,
	}
}

// passwordHandler handles password authentication
func (s *SecureSSHServer) passwordHandler(ctx ssh.Context, password string) bool {
	if !s.config.Auth.EnablePassword {
		return false
	}
	
	username := ctx.User()
	expectedPassword, ok := s.config.Auth.Passwords[username]
	if !ok {
		log.Printf("Failed login attempt for unknown user: %s from %s", username, ctx.RemoteAddr())
		return false
	}
	
	// In production, use bcrypt.CompareHashAndPassword
	if password != expectedPassword {
		log.Printf("Failed login attempt for user: %s from %s", username, ctx.RemoteAddr())
		return false
	}
	
	log.Printf("Successful password authentication for user: %s from %s", username, ctx.RemoteAddr())
	return true
}

// publicKeyHandler handles public key authentication
func (s *SecureSSHServer) publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	if !s.config.Auth.EnablePublicKey {
		return false
	}
	
	username := ctx.User()
	authorizedKeys, ok := s.config.Auth.AuthorizedKeys[username]
	if !ok {
		log.Printf("Failed public key auth for unknown user: %s from %s", username, ctx.RemoteAddr())
		return false
	}
	
	for _, authorizedKey := range authorizedKeys {
		if ssh.KeysEqual(key, authorizedKey) {
			log.Printf("Successful public key authentication for user: %s from %s", username, ctx.RemoteAddr())
			return true
		}
	}
	
	log.Printf("Failed public key auth for user: %s from %s (key not authorized)", username, ctx.RemoteAddr())
	return false
}

// Start starts the SSH server with authentication
func (s *SecureSSHServer) Start(teaHandler func(ssh.Session) (tea.Model, []tea.ProgramOption)) error {
	// Create server with wish middleware and authentication
	server, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)),
		wish.WithHostKeyPEM(s.config.HostKeyPEM),
		wish.WithPasswordAuth(s.passwordHandler),
		wish.WithPublicKeyAuth(s.publicKeyHandler),
		wish.WithMiddleware(
			s.authLoggingMiddleware(), // Custom auth logging
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create SSH server: %w", err)
	}

	s.server = server

	// Start server in goroutine
	go func() {
		log.Printf("🔐 Secure PowerHell SSH server starting on %s:%d", s.config.Host, s.config.Port)
		log.Printf("🔑 Authentication required: Password=%v, PublicKey=%v", 
			s.config.Auth.EnablePassword, s.config.Auth.EnablePublicKey)
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

// authLoggingMiddleware adds authentication context to logs
func (s *SecureSSHServer) authLoggingMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			// Log authenticated session
			log.Printf("Authenticated session started: user=%s, addr=%s", 
				sess.User(), sess.RemoteAddr())
			next(sess)
			log.Printf("Session ended: user=%s", sess.User())
		}
	}
}

// LoadAuthorizedKeysFile loads authorized keys from a file
func LoadAuthorizedKeysFile(filename string) ([]ssh.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var keys []ssh.PublicKey
	for len(data) > 0 {
		key, _, _, rest, err := ssh.ParseAuthorizedKey(data)
		if err != nil {
			// Skip invalid lines
			data = rest
			continue
		}
		keys = append(keys, key)
		data = rest
	}
	
	return keys, nil
}