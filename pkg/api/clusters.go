package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dorgu-ai/dorgu-platform/pkg/models"
	"github.com/dorgu-ai/dorgu-platform/pkg/watcher"
	"github.com/gorilla/mux"
)

// ClustersHandler handles API requests for ClusterPersona resources.
type ClustersHandler struct {
	watcher *watcher.Watcher
}

// NewClustersHandler creates a new clusters API handler.
func NewClustersHandler(w *watcher.Watcher) *ClustersHandler {
	return &ClustersHandler{
		watcher: w,
	}
}

// ListClusters handles GET /api/clusters
func (h *ClustersHandler) ListClusters(w http.ResponseWriter, r *http.Request) {
	if h.watcher == nil {
		response := models.ClusterList{
			Clusters: []models.ClusterPersona{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	clusters := h.watcher.GetClusters()

	// Convert watcher.ClusterPersona to models.ClusterPersona
	modelClusters := make([]models.ClusterPersona, len(clusters))
	for i, c := range clusters {
		modelClusters[i] = h.convertToModel(c)
	}

	response := models.ClusterList{
		Clusters: modelClusters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCluster handles GET /api/clusters/:name
func (h *ClustersHandler) GetCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if h.watcher == nil {
		http.Error(w, "Cluster not found", http.StatusNotFound)
		return
	}

	cluster, found := h.watcher.GetCluster(name)
	if !found {
		http.Error(w, "Cluster not found", http.StatusNotFound)
		return
	}

	modelCluster := h.convertToModel(cluster)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modelCluster)
}

// convertToModel converts watcher.ClusterPersona to models.ClusterPersona.
func (h *ClustersHandler) convertToModel(c *watcher.ClusterPersona) models.ClusterPersona {
	return models.ClusterPersona{
		Name: c.Name,
		Spec: models.ClusterPersonaSpec{
			Name:        c.Spec.Name,
			Description: c.Spec.Description,
			Environment: c.Spec.Environment,
		},
		Status: h.convertStatus(c.Status),
	}
}

// convertStatus converts status fields.
func (h *ClustersHandler) convertStatus(s watcher.ClusterPersonaStatus) models.ClusterPersonaStatus {
	// Convert time pointer
	var lastDiscovery *time.Time
	if s.LastDiscovery != nil {
		t := s.LastDiscovery.Time
		lastDiscovery = &t
	}

	// Convert nodes
	nodes := make([]models.NodeInfo, len(s.Nodes))
	for i, n := range s.Nodes {
		var capacity *models.NodeResources
		if n.Capacity != nil {
			capacity = &models.NodeResources{
				CPU:    n.Capacity.CPU,
				Memory: n.Capacity.Memory,
				Pods:   n.Capacity.Pods,
			}
		}
		var allocatable *models.NodeResources
		if n.Allocatable != nil {
			allocatable = &models.NodeResources{
				CPU:    n.Allocatable.CPU,
				Memory: n.Allocatable.Memory,
				Pods:   n.Allocatable.Pods,
			}
		}
		nodes[i] = models.NodeInfo{
			Name:             n.Name,
			Role:             n.Role,
			Ready:            n.Ready,
			Capacity:         capacity,
			Allocatable:      allocatable,
			KubeletVersion:   n.KubeletVersion,
			ContainerRuntime: n.ContainerRuntime,
		}
	}

	// Convert addons
	addons := make([]models.AddonInfo, len(s.Addons))
	for i, a := range s.Addons {
		addons[i] = models.AddonInfo{
			Name:      a.Name,
			Namespace: a.Namespace,
			Healthy:   a.Healthy,
			Version:   a.Version,
		}
	}

	// Convert resource summary
	var resourceSummary *models.ResourceSummary
	if s.ResourceSummary != nil {
		resourceSummary = &models.ResourceSummary{
			TotalCPU:          s.ResourceSummary.TotalCPU,
			TotalMemory:       s.ResourceSummary.TotalMemory,
			AllocatableCPU:    s.ResourceSummary.AllocatableCPU,
			AllocatableMemory: s.ResourceSummary.AllocatableMemory,
			RunningPods:       s.ResourceSummary.RunningPods,
		}
	}

	return models.ClusterPersonaStatus{
		Phase:             s.Phase,
		LastDiscovery:     lastDiscovery,
		Nodes:             nodes,
		ResourceSummary:   resourceSummary,
		Addons:            addons,
		KubernetesVersion: s.KubernetesVersion,
		Platform:          s.Platform,
		ApplicationCount:  s.ApplicationCount,
	}
}

// RegisterRoutes registers the cluster API routes.
func (h *ClustersHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/clusters", h.ListClusters).Methods("GET")
	router.HandleFunc("/api/clusters/{name}", h.GetCluster).Methods("GET")
}
