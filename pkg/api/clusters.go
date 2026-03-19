package api

import (
	"encoding/json"
	"net/http"

	"github.com/dorgu-ai/dorgu-platform/pkg/models"
	"github.com/gorilla/mux"
)

// ClustersHandler handles API requests for ClusterPersona resources.
type ClustersHandler struct {
	// Agent 3 will provide a data source (watcher)
	// For now, we'll return mock data
}

// NewClustersHandler creates a new clusters API handler.
func NewClustersHandler() *ClustersHandler {
	return &ClustersHandler{}
}

// ListClusters handles GET /api/clusters
func (h *ClustersHandler) ListClusters(w http.ResponseWriter, r *http.Request) {
	// TODO: Agent 3 will provide real data from K8s watcher
	// For now, return empty list or mock data for testing

	clusters := models.ClusterList{
		Clusters: []models.ClusterPersona{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusters)
}

// GetCluster handles GET /api/clusters/:name
func (h *ClustersHandler) GetCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// TODO: Agent 3 will provide real data from K8s watcher
	// For now, return 404 for any cluster

	http.Error(w, "Cluster not found: "+name, http.StatusNotFound)
}

// RegisterRoutes registers the cluster API routes.
func (h *ClustersHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/clusters", h.ListClusters).Methods("GET")
	router.HandleFunc("/api/clusters/{name}", h.GetCluster).Methods("GET")
}
