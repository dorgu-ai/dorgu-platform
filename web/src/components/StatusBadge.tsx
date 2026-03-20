import { Badge } from '@/components/ui/badge';

interface StatusBadgeProps {
  phase?: string;
  health?: string;
}

export function StatusBadge({ phase, health }: StatusBadgeProps) {
  // Determine badge variant based on phase
  const getPhaseVariant = (phase?: string): 'default' | 'secondary' | 'destructive' | 'outline' => {
    if (!phase) return 'outline';

    const lowerPhase = phase.toLowerCase();
    if (lowerPhase === 'ready') return 'default';
    if (lowerPhase === 'pending' || lowerPhase === 'initializing') return 'secondary';
    if (lowerPhase === 'failed' || lowerPhase === 'error') return 'destructive';
    return 'outline';
  };

  // Determine color class based on health (used as dot indicator)
  const getHealthColor = (health?: string): string => {
    if (!health) return 'bg-gray-400';

    const lowerHealth = health.toLowerCase();
    if (lowerHealth === 'healthy') return 'bg-green-500';
    if (lowerHealth === 'degraded') return 'bg-yellow-500';
    if (lowerHealth === 'unhealthy') return 'bg-red-500';
    return 'bg-gray-400';
  };

  const variant = getPhaseVariant(phase);
  const healthColor = getHealthColor(health);

  return (
    <div className="flex items-center gap-2">
      <Badge variant={variant}>{phase || 'Unknown'}</Badge>
      {health && (
        <div className="flex items-center gap-1.5">
          <div className={`h-2 w-2 rounded-full ${healthColor}`} />
          <span className="text-xs text-muted-foreground">{health}</span>
        </div>
      )}
    </div>
  );
}
