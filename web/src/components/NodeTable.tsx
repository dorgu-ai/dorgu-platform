import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { InfoIcon } from 'lucide-react';
import { NodeInfo } from '../lib/api';

interface NodeTableProps {
  nodes?: NodeInfo[];
}

export function NodeTable({ nodes }: NodeTableProps) {
  if (!nodes || nodes.length === 0) {
    return (
      <Alert>
        <InfoIcon className="h-4 w-4" />
        <AlertDescription>No node information available</AlertDescription>
      </Alert>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Role</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Version</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {nodes.map((node) => (
          <TableRow key={node.name}>
            <TableCell className="font-medium">{node.name}</TableCell>
            <TableCell>
              <Badge variant="outline" className="capitalize">
                {node.role || 'unknown'}
              </Badge>
            </TableCell>
            <TableCell>
              <Badge variant={node.ready ? 'default' : 'destructive'}>
                {node.ready ? 'Ready' : 'NotReady'}
              </Badge>
            </TableCell>
            <TableCell className="text-muted-foreground">
              {node.kubeletVersion || '—'}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
