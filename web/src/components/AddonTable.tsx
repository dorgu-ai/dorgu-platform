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
import { AddonInfo } from '../lib/api';

interface AddonTableProps {
  addons?: AddonInfo[];
}

export function AddonTable({ addons }: AddonTableProps) {
  if (!addons || addons.length === 0) {
    return (
      <Alert>
        <InfoIcon className="h-4 w-4" />
        <AlertDescription>
          No addons installed. Run <code className="text-sm">dorgu cluster setup</code> to install components.
        </AlertDescription>
      </Alert>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Version</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Namespace</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {addons.map((addon) => (
          <TableRow key={addon.name}>
            <TableCell className="font-medium">{addon.name}</TableCell>
            <TableCell className="text-muted-foreground">
              {addon.version || '—'}
            </TableCell>
            <TableCell>
              <Badge variant={addon.healthy ? 'default' : 'destructive'}>
                {addon.healthy ? 'Healthy' : 'Unhealthy'}
              </Badge>
            </TableCell>
            <TableCell className="text-muted-foreground">
              {addon.namespace || '—'}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
