package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/dorgu-ai/dorgu-platform/pkg/api"
	"github.com/dorgu-ai/dorgu-platform/pkg/watcher"
	ws "github.com/dorgu-ai/dorgu-platform/pkg/websocket"
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
	watcher    *watcher.Watcher
	wsHub      *ws.Hub
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

// setupRoutes configures all HTTP routes.
func (s *Server) setupRoutes() {
	// API routes
	clustersHandler := api.NewClustersHandler(s.watcher)
	clustersHandler.RegisterRoutes(s.router)

	// WebSocket endpoint
	s.router.HandleFunc("/ws", s.wsHub.ServeWS)

	// Static file serving (frontend)
	s.setupStaticRoutes()
}

// setupStaticRoutes serves static files for the frontend.
func (s *Server) setupStaticRoutes() {
	staticFiles, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Println("Frontend dist not found, serving placeholder")
		s.router.HandleFunc("/", s.servePlaceholder).Methods("GET")
		return
	}

	// Check if this is the React app (has index.html with React content)
	// or just the placeholder
	if _, err := fs.Stat(staticFiles, "assets"); err != nil {
		// No assets directory means this is the placeholder, not the React build
		fileServer := http.FileServer(http.FS(staticFiles))
		s.router.PathPrefix("/").Handler(fileServer)
		return
	}

	// Serve the React SPA with proper fallback to index.html
	s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try to serve the file directly
		if _, err := fs.Stat(staticFiles, path); err == nil {
			http.FileServer(http.FS(staticFiles)).ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for client-side routing
		indexFile, err := fs.ReadFile(staticFiles, "index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexFile)
	})
}

// servePlaceholder serves a simple HTML placeholder.
func (s *Server) servePlaceholder(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Dorgu Platform</title>
    <style>
        body { font-family: system-ui; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .status { color: #28a745; }
        a { color: #007bff; }
    </style>
</head>
<body>
    <h1>Dorgu Platform</h1>
    <p class="status">✓ Backend is running</p>
    <p>The frontend will be available after running: <code>make frontend && make build</code></p>
    <h2>API Endpoints</h2>
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

	s.watcher = w

	// Initialize WebSocket hub
	s.wsHub = ws.NewHub()
	go s.wsHub.Run()

	// Bridge watcher events to WebSocket hub
	go func() {
		eventTypeMap := map[string]string{
			"added":    "cluster.added",
			"modified": "cluster.updated",
			"deleted":  "cluster.deleted",
		}
		for event := range w.Events() {
			wsType, ok := eventTypeMap[event.Type]
			if !ok {
				wsType = event.Type
			}
			name := ""
			if event.Cluster != nil {
				name = event.Cluster.Name
			}
			log.Printf("Broadcasting WebSocket event: %s - %s", wsType, name)
			s.wsHub.Broadcast(wsType, map[string]string{
				"name": name,
			})
		}
	}()

	// Setup routes with watcher
	s.setupRoutes()

	addr := fmt.Sprintf(":%d", s.Port)
	log.Printf("Dorgu Platform starting on http://localhost%s", addr)
	log.Printf("API available at http://localhost%s/api/clusters", addr)
	log.Printf("WebSocket available at ws://localhost%s/ws", addr)
	log.Printf("Frontend available at http://localhost%s/", addr)

	return http.ListenAndServe(addr, s.router)
}
