import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import type { FileMetrics, Scan } from "@/lib/api";
import { duplicateCount } from "@/lib/dashboard-metrics";

interface OverviewCardsProps {
  scan: Scan;
  files: FileMetrics[];
}

export function OverviewCards({ scan, files }: OverviewCardsProps) {
  const stats = [
    { label: "Total Files", value: scan.totalFiles.toLocaleString() },
    { label: "Avg Complexity", value: scan.avgComplexity.toFixed(1) },
    { label: "Avg Health Score", value: scan.avgHealthScore.toFixed(0) },
    { label: "Duplicate Files", value: duplicateCount(files).toLocaleString() },
  ];

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat) => (
        <Card key={stat.label}>
          <CardHeader>
            <CardDescription>{stat.label}</CardDescription>
            <CardTitle className="text-3xl font-semibold tabular-nums">{stat.value}</CardTitle>
          </CardHeader>
        </Card>
      ))}
    </div>
  );
}
