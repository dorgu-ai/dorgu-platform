package watcher

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Watcher watches ClusterPersona CRDs and maintains an in-memory cache.
type Watcher struct {
	client   dynamic.Interface
	informer cache.SharedIndexInformer
	stopCh   chan struct{}
	clusters map[string]*ClusterPersona
	mu       sync.RWMutex
	eventCh  chan Event
}

// Event represents a CRD event.
type Event struct {
	Type    string           // "added", "modified", "deleted"
	Cluster *ClusterPersona
}

// NewWatcher creates a new ClusterPersona watcher.
func NewWatcher(kubeconfig, kubeContext string) (*Watcher, error) {
	// Build kubeconfig path
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	// Build config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Define ClusterPersona GVR (Group-Version-Resource)
	gvr := schema.GroupVersionResource{
		Group:    "dorgu.io",
		Version:  "v1",
		Resource: "clusterpersonas",
	}

	// Create dynamic informer factory
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, time.Minute)
	informer := factory.ForResource(gvr).Informer()

	w := &Watcher{
		client:   dynamicClient,
		informer: informer,
		stopCh:   make(chan struct{}),
		clusters: make(map[string]*ClusterPersona),
		eventCh:  make(chan Event, 100),
	}

	// Add event handlers
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    w.onAdd,
		UpdateFunc: w.onUpdate,
		DeleteFunc: w.onDelete,
	})

	return w, nil
}

// Start starts the watcher.
func (w *Watcher) Start(ctx context.Context) error {
	log.Println("Starting ClusterPersona watcher...")
	go w.informer.Run(w.stopCh)

	// Wait for cache to sync
	if !cache.WaitForCacheSync(w.stopCh, w.informer.HasSynced) {
		return fmt.Errorf("failed to sync ClusterPersona cache")
	}

	log.Println("ClusterPersona watcher started and synced")
	return nil
}

// Stop stops the watcher.
func (w *Watcher) Stop() {
	log.Println("Stopping ClusterPersona watcher...")
	close(w.stopCh)
}

// GetClusters returns all cached ClusterPersonas.
func (w *Watcher) GetClusters() []*ClusterPersona {
	w.mu.RLock()
	defer w.mu.RUnlock()

	clusters := make([]*ClusterPersona, 0, len(w.clusters))
	for _, c := range w.clusters {
		clusters = append(clusters, c)
	}
	return clusters
}

// GetCluster returns a single ClusterPersona by name.
func (w *Watcher) GetCluster(name string) (*ClusterPersona, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cluster, found := w.clusters[name]
	return cluster, found
}

// Events returns the event channel for WebSocket broadcasting (Agent 7).
func (w *Watcher) Events() <-chan Event {
	return w.eventCh
}

// onAdd handles ClusterPersona add events.
func (w *Watcher) onAdd(obj interface{}) {
	cluster := w.convertToClusterPersona(obj)
	if cluster == nil {
		return
	}

	w.mu.Lock()
	w.clusters[cluster.Name] = cluster
	w.mu.Unlock()

	log.Printf("ClusterPersona added: %s", cluster.Name)

	// Send event to WebSocket broadcast channel (Agent 7 will consume)
	select {
	case w.eventCh <- Event{Type: "added", Cluster: cluster}:
	default:
		log.Println("Event channel full, dropping event")
	}
}

// onUpdate handles ClusterPersona update events.
func (w *Watcher) onUpdate(oldObj, newObj interface{}) {
	cluster := w.convertToClusterPersona(newObj)
	if cluster == nil {
		return
	}

	w.mu.Lock()
	w.clusters[cluster.Name] = cluster
	w.mu.Unlock()

	log.Printf("ClusterPersona updated: %s", cluster.Name)

	select {
	case w.eventCh <- Event{Type: "modified", Cluster: cluster}:
	default:
		log.Println("Event channel full, dropping event")
	}
}

// onDelete handles ClusterPersona delete events.
func (w *Watcher) onDelete(obj interface{}) {
	cluster := w.convertToClusterPersona(obj)
	if cluster == nil {
		return
	}

	w.mu.Lock()
	delete(w.clusters, cluster.Name)
	w.mu.Unlock()

	log.Printf("ClusterPersona deleted: %s", cluster.Name)

	select {
	case w.eventCh <- Event{Type: "deleted", Cluster: cluster}:
	default:
		log.Println("Event channel full, dropping event")
	}
}

// convertToClusterPersona converts unstructured object to ClusterPersona.
func (w *Watcher) convertToClusterPersona(obj interface{}) *ClusterPersona {
	unstructuredObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Println("Failed to convert object to unstructured")
		return nil
	}

	// Basic conversion (simplified for MVP)
	cluster := &ClusterPersona{
		ObjectMeta: metav1.ObjectMeta{
			Name: unstructuredObj.GetName(),
		},
	}

	// TODO: Properly unmarshal spec and status from unstructured.Object
	// For MVP, this provides basic name tracking
	// Agent 7 will enhance with full field mapping

	return cluster
}
