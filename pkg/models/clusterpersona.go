package models

import (
	"time"
)

// ClusterPersona represents the serialized ClusterPersona CRD for the API.
type ClusterPersona struct {
	Name   string               `json:"name"`
	Spec   ClusterPersonaSpec   `json:"spec"`
	Status ClusterPersonaStatus `json:"status"`
}

// ClusterPersonaSpec matches the CRD spec fields we care about.
type ClusterPersonaSpec struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Environment string `json:"environment,omitempty"`
}

// ClusterPersonaStatus matches the CRD status fields.
type ClusterPersonaStatus struct {
	Phase             string           `json:"phase,omitempty"`
	LastDiscovery     *time.Time       `json:"lastDiscovery,omitempty"`
	Nodes             []NodeInfo       `json:"nodes,omitempty"`
	ResourceSummary   *ResourceSummary `json:"resourceSummary,omitempty"`
	Addons            []AddonInfo      `json:"addons,omitempty"`
	KubernetesVersion string           `json:"kubernetesVersion,omitempty"`
	Platform          string           `json:"platform,omitempty"`
	ApplicationCount  int32            `json:"applicationCount,omitempty"`
}

// NodeInfo represents a cluster node.
type NodeInfo struct {
	Name             string         `json:"name"`
	Role             string         `json:"role,omitempty"`
	Ready            bool           `json:"ready"`
	Capacity         *NodeResources `json:"capacity,omitempty"`
	Allocatable      *NodeResources `json:"allocatable,omitempty"`
	KubeletVersion   string         `json:"kubeletVersion,omitempty"`
	ContainerRuntime string         `json:"containerRuntime,omitempty"`
}

// NodeResources represents node resource quantities.
type NodeResources struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	Pods   string `json:"pods,omitempty"`
}

// ResourceSummary represents cluster-wide resource information.
type ResourceSummary struct {
	TotalCPU          string `json:"totalCPU,omitempty"`
	TotalMemory       string `json:"totalMemory,omitempty"`
	AllocatableCPU    string `json:"allocatableCPU,omitempty"`
	AllocatableMemory string `json:"allocatableMemory,omitempty"`
	RunningPods       int32  `json:"runningPods,omitempty"`
}

// AddonInfo represents an installed cluster addon/component.
type AddonInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Healthy   bool   `json:"healthy"`
	Version   string `json:"version,omitempty"`
}

// ClusterList is a wrapper for the list of clusters.
type ClusterList struct {
	Clusters []ClusterPersona `json:"clusters"`
}
