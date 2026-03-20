import { useQuery, UseQueryResult } from '@tanstack/react-query';
import { api, ClusterPersona } from '../lib/api';

export function useCluster(name: string): UseQueryResult<ClusterPersona, Error> {
  return useQuery({
    queryKey: ['cluster', name],
    queryFn: () => api.getCluster(name),
    refetchInterval: 30000, // Refetch every 30 seconds (will be replaced by WebSocket in Agent 7)
    staleTime: 10000, // Consider data stale after 10 seconds
    enabled: !!name, // Only run query if name is provided
  });
}
