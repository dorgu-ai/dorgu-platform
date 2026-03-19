package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dorgu-ai/dorgu-platform/pkg/models"
	"github.com/gorilla/mux"
)

func TestListClusters(t *testing.T) {
	handler := NewClustersHandler()
	req := httptest.NewRequest("GET", "/api/clusters", nil)
	w := httptest.NewRecorder()

	handler.ListClusters(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.ClusterList
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Clusters) != 0 {
		t.Errorf("Expected empty cluster list, got %d clusters", len(result.Clusters))
	}
}

func TestGetCluster_NotFound(t *testing.T) {
	handler := NewClustersHandler()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest("GET", "/api/clusters/test-cluster", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
