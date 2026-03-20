import { useQuery, UseQueryResult } from '@tanstack/react-query';
import { api, ClusterPersona } from '../lib/api';

export function useClusters(): UseQueryResult<ClusterPersona[], Error> {
  return useQuery({
    queryKey: ['clusters'],
    queryFn: api.getClusters,
    refetchInterval: 30000, // Refetch every 30 seconds (will be replaced by WebSocket in Agent 7)
    staleTime: 10000, // Consider data stale after 10 seconds
  });
}
