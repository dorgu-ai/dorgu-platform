import { useQuery, UseQueryResult } from '@tanstack/react-query';
import { api, ClusterPersona } from '../lib/api';

export function useClusters(): UseQueryResult<ClusterPersona[], Error> {
  return useQuery({
    queryKey: ['clusters'],
    queryFn: api.getClusters,
    staleTime: 10000,
  });
}
