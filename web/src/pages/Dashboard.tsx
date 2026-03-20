import { ClusterList } from '../components/ClusterList';

export default function Dashboard() {
  return (
    <div className="container mx-auto py-8 space-y-6">
      <div className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight">Cluster Dashboard</h1>
        <p className="text-muted-foreground">
          View and manage your ClusterPersona resources
        </p>
      </div>

      <ClusterList />
    </div>
  );
}
