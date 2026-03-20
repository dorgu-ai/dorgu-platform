import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface ClusterPersona {
  name: string;
  spec: {
    name: string;
    description?: string;
    environment?: string;
  };
  status: {
    phase?: string;
    kubernetesVersion?: string;
    platform?: string;
    nodes?: NodeInfo[];
    addons?: AddonInfo[];
    applicationCount?: number;
    resourceSummary?: ResourceSummary;
    lastDiscovery?: string;
  };
}

export interface NodeInfo {
  name: string;
  role?: string;
  ready: boolean;
  kubeletVersion?: string;
  containerRuntime?: string;
  capacity?: NodeResources;
  allocatable?: NodeResources;
}

export interface NodeResources {
  cpu?: string;
  memory?: string;
  pods?: string;
}

export interface AddonInfo {
  name: string;
  namespace?: string;
  healthy: boolean;
  version?: string;
}

export interface ResourceSummary {
  totalCPU?: string;
  totalMemory?: string;
  allocatableCPU?: string;
  allocatableMemory?: string;
  runningPods?: number;
}

export const api = {
  getClusters: async (): Promise<ClusterPersona[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/clusters`);
    return response.data.clusters || [];
  },

  getCluster: async (name: string): Promise<ClusterPersona> => {
    const response = await axios.get(`${API_BASE_URL}/api/clusters/${name}`);
    return response.data;
  },
};
