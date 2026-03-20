import { useNavigate } from 'react-router-dom';
import { useClusters } from '../hooks/useClusters';
import { StatusBadge } from './StatusBadge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Skeleton } from '@/components/ui/skeleton';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertCircle } from 'lucide-react';

export function ClusterList() {
  const navigate = useNavigate();
  const { data: clusters, isLoading, isError, error } = useClusters();

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-2">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-10 w-full" />
      </div>
    );
  }

  // Error state
  if (isError) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>
          Failed to load clusters: {error?.message || 'Unknown error'}
        </AlertDescription>
      </Alert>
    );
  }

  // Empty state
  if (!clusters || clusters.length === 0) {
    return (
      <Alert>
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>No clusters found</AlertTitle>
        <AlertDescription>
          No ClusterPersona resources found in your Kubernetes cluster. Create one using{' '}
          <code className="text-sm">dorgu cluster init</code>.
        </AlertDescription>
      </Alert>
    );
  }

  // Data state
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Environment</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Nodes</TableHead>
          <TableHead>Description</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {clusters.map((cluster) => (
          <TableRow
            key={cluster.name}
            className="cursor-pointer hover:bg-muted/50"
            onClick={() => navigate(`/cluster/${cluster.name}`)}
          >
            <TableCell className="font-medium">{cluster.name}</TableCell>
            <TableCell>
              <span className="capitalize">{cluster.spec.environment || 'N/A'}</span>
            </TableCell>
            <TableCell>
              <StatusBadge
                phase={cluster.status?.phase}
              />
            </TableCell>
            <TableCell>
              {cluster.status?.nodes ? cluster.status.nodes.length : 0}
            </TableCell>
            <TableCell className="text-muted-foreground max-w-md truncate">
              {cluster.spec.description || '—'}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
