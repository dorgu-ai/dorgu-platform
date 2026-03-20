package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/dorgu-ai/dorgu-platform/pkg/api"
	"github.com/dorgu-ai/dorgu-platform/pkg/watcher"
	"github.com/gorilla/mux"
)

//go:embed static
var staticFS embed.FS

// Server represents the Dorgu Platform HTTP server.
type Server struct {
	Port       int
	KubeConfig string
	Context    string
	router     *mux.Router
}

// NewServer creates a new server instance.
func NewServer(config *Config) *Server {
	return &Server{
		Port:       config.Port,
		KubeConfig: config.KubeConfig,
		Context:    config.Context,
		router:     mux.NewRouter(),
	}
}

// setupStaticRoutes serves static files for the frontend.
func (s *Server) setupStaticRoutes() {
	// Try to serve from embedded FS (after frontend build)
	// Fallback to a simple message if dist doesn't exist yet

	staticFiles, err := fs.Sub(staticFS, "static")
	if err != nil {
		// Static files not embedded yet, serve placeholder
		s.router.HandleFunc("/", s.servePlaceholder).Methods("GET")
		return
	}

	fileServer := http.FileServer(http.FS(staticFiles))
	s.router.PathPrefix("/").Handler(fileServer)
}

// servePlaceholder serves a simple HTML placeholder.
func (s *Server) servePlaceholder(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Dorgu Platform</title>
</head>
<body>
    <h1>Dorgu Platform</h1>
    <p>Backend is running. Frontend will be available after Agent 4 completes.</p>
    <p>API Endpoints:</p>
    <ul>
        <li><a href="/api/clusters">GET /api/clusters</a></li>
    </ul>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	// Initialize K8s watcher
	w, err := watcher.NewWatcher(s.KubeConfig, s.Context)
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	// Start watcher in background
	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}

	// Setup routes with watcher
	clustersHandler := api.NewClustersHandler(w)
	clustersHandler.RegisterRoutes(s.router)

	s.setupStaticRoutes()

	addr := fmt.Sprintf(":%d", s.Port)
	log.Printf("Dorgu Platform starting on http://localhost%s", addr)
	log.Printf("Watching ClusterPersona CRDs in current kube-context")
	log.Printf("API available at http://localhost%s/api/clusters", addr)

	return http.ListenAndServe(addr, s.router)
}
