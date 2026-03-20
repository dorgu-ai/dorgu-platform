package watcher

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ClusterPersona represents the CRD as defined in dorgu-operator.
// This is a simplified version focusing on fields we need for the dashboard.
type ClusterPersona struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterPersonaSpec   `json:"spec,omitempty"`
	Status            ClusterPersonaStatus `json:"status,omitempty"`
}

// ClusterPersonaList contains a list of ClusterPersona.
type ClusterPersonaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterPersona `json:"items"`
}

// ClusterPersonaSpec defines the desired state.
type ClusterPersonaSpec struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Environment string `json:"environment,omitempty"`
}

// ClusterPersonaStatus defines the observed state.
type ClusterPersonaStatus struct {
	Phase             string           `json:"phase,omitempty"`
	LastDiscovery     *metav1.Time     `json:"lastDiscovery,omitempty"`
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

// AddonInfo represents an installed cluster addon.
type AddonInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Healthy   bool   `json:"healthy"`
	Version   string `json:"version,omitempty"`
}

// DeepCopyInto is required for runtime.Object interface.
func (in *ClusterPersona) DeepCopyInto(out *ClusterPersona) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is required for runtime.Object interface.
func (in *ClusterPersona) DeepCopy() *ClusterPersona {
	if in == nil {
		return nil
	}
	out := new(ClusterPersona)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is required for runtime.Object interface.
func (in *ClusterPersona) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}
