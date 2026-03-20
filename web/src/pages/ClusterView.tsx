import { useParams } from 'react-router-dom';

export default function ClusterView() {
  const { name } = useParams<{ name: string }>();

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">Cluster: {name}</h2>
      <p className="text-gray-600">
        Agent 6 will implement the cluster detail view here.
      </p>
    </div>
  );
}
