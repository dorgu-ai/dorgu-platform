import { useQuery, UseQueryResult } from '@tanstack/react-query';
import { api, ClusterPersona } from '../lib/api';

export function useCluster(name: string): UseQueryResult<ClusterPersona, Error> {
  return useQuery({
    queryKey: ['cluster', name],
    queryFn: () => api.getCluster(name),
    staleTime: 10000,
    enabled: !!name,
  });
}
