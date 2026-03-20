import { useParams, useNavigate } from 'react-router-dom';
import { useCluster } from '../hooks/useCluster';
import { StatusBadge } from '../components/StatusBadge';
import { NodeTable } from '../components/NodeTable';
import { AddonTable } from '../components/AddonTable';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertCircle, ArrowLeft } from 'lucide-react';

export default function ClusterView() {
  const { name } = useParams<{ name: string }>();
  const navigate = useNavigate();
  const { data: cluster, isLoading, isError, error } = useCluster(name || '');

  if (isLoading) {
    return (
      <div className="container mx-auto py-8 space-y-6">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-32 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="container mx-auto py-8">
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            Failed to load cluster: {error?.message || 'Unknown error'}
          </AlertDescription>
        </Alert>
        <Button
          variant="outline"
          className="mt-4"
          onClick={() => navigate('/')}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Dashboard
        </Button>
      </div>
    );
  }

  if (!cluster) {
    return (
      <div className="container mx-auto py-8">
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Cluster Not Found</AlertTitle>
          <AlertDescription>
            ClusterPersona "{name}" does not exist.
          </AlertDescription>
        </Alert>
        <Button
          variant="outline"
          className="mt-4"
          onClick={() => navigate('/')}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Dashboard
        </Button>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 space-y-6">
      {/* Header with breadcrumb */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate('/')}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>
        <div className="flex-1 space-y-1">
          <h1 className="text-3xl font-bold tracking-tight">{cluster.name}</h1>
          <p className="text-muted-foreground">
            {cluster.spec.description || 'No description'}
          </p>
        </div>
      </div>

      {/* Status Overview Card */}
      <Card>
        <CardHeader>
          <CardTitle>Status Overview</CardTitle>
          <CardDescription>Current cluster state and health</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <div className="text-sm font-medium text-muted-foreground mb-1">
                Environment
              </div>
              <div className="text-lg font-semibold capitalize">
                {cluster.spec.environment || 'N/A'}
              </div>
            </div>
            <div>
              <div className="text-sm font-medium text-muted-foreground mb-1">
                Phase
              </div>
              <StatusBadge
                phase={cluster.status?.phase}
              />
            </div>
            <div>
              <div className="text-sm font-medium text-muted-foreground mb-1">
                Nodes
              </div>
              <div className="text-lg font-semibold">
                {cluster.status?.nodes ? cluster.status.nodes.length : 0}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Nodes Card */}
      <Card>
        <CardHeader>
          <CardTitle>Nodes</CardTitle>
          <CardDescription>Kubernetes nodes in this cluster</CardDescription>
        </CardHeader>
        <CardContent>
          <NodeTable nodes={cluster.status?.nodes} />
        </CardContent>
      </Card>

      {/* Addons Card */}
      <Card>
        <CardHeader>
          <CardTitle>Installed Addons</CardTitle>
          <CardDescription>
            Components installed via <code className="text-sm">dorgu cluster setup</code>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <AddonTable addons={cluster.status?.addons} />
        </CardContent>
      </Card>
    </div>
  );
}
