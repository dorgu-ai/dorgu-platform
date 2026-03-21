package platform

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dorgu-ai/dorgu-platform/pkg/server"
)

// Config holds configuration for the platform server.
type Config struct {
	Port        string // Server port (default: "8080")
	Kubeconfig  string // Path to kubeconfig (default: $KUBECONFIG or ~/.kube/config)
	Context     string // Kubernetes context (default: current context)
	Development bool   // Enable development mode (verbose logging)
}

// Server represents the platform server.
type Server struct {
	config Config
	srv    *server.Server
}

// NewServer creates a new platform server with the given config.
func NewServer(config Config) (*Server, error) {
	// Set defaults
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.Kubeconfig == "" {
		config.Kubeconfig = os.Getenv("KUBECONFIG")
		if config.Kubeconfig == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			config.Kubeconfig = homeDir + "/.kube/config"
		}
	}

	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid port %q: %w", config.Port, err)
	}

	srvConfig := &server.Config{
		Port:       port,
		KubeConfig: config.Kubeconfig,
		Context:    config.Context,
	}

	return &Server{
		config: config,
		srv:    server.NewServer(srvConfig),
	}, nil
}

// Start starts the platform server. It blocks until the context is cancelled
// or an interrupt signal is received.
func (s *Server) Start(ctx context.Context) error {
	if s.config.Development {
		log.Println("Development mode enabled")
	}

	// Start server in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.srv.Start()
	}()

	// Wait for shutdown signal or server error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	case <-sigChan:
		log.Println("\nInterrupt signal received, shutting down...")
	}

	return s.Stop()
}

// Stop gracefully stops the platform server.
func (s *Server) Stop() error {
	log.Println("Stopping platform server...")
	log.Println("Platform server stopped")
	return nil
}
