import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import type { ScanComparison } from "@/lib/scan-comparison";

function formatDelta(delta: number): string {
  const rounded = Math.round(delta * 10) / 10;
  if (rounded === 0) return "±0";
  return rounded > 0 ? `+${rounded}` : `${rounded}`;
}

function deltaVariant(delta: number, higherIsBetter: boolean): "default" | "secondary" | "destructive" {
  if (delta === 0) return "secondary";
  const improved = higherIsBetter ? delta > 0 : delta < 0;
  return improved ? "default" : "destructive";
}

export function ScanComparisonView({ comparison }: { comparison: ScanComparison }) {
  const changed = comparison.fileChanges.filter((f) => f.status !== "unchanged");

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Comparison</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-6">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          {comparison.metrics.map((m) => (
            <div key={m.label} className="flex flex-col gap-1 rounded-lg border p-3">
              <span className="text-sm text-muted-foreground">{m.label}</span>
              <div className="flex items-baseline gap-2">
                <span className="text-lg font-semibold tabular-nums">{m.after.toFixed(1)}</span>
                <Badge variant={deltaVariant(m.delta, m.label !== "Avg Complexity")}>
                  {formatDelta(m.delta)}
                </Badge>
              </div>
              <span className="text-xs text-muted-foreground">was {m.before.toFixed(1)}</span>
            </div>
          ))}
        </div>

        <div>
          <p className="mb-2 text-sm font-medium">File changes ({changed.length})</p>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Path</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Health Before</TableHead>
                <TableHead>Health After</TableHead>
                <TableHead>Delta</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {changed.map((f) => (
                <TableRow key={f.path}>
                  <TableCell className="max-w-xs truncate font-mono text-xs" title={f.path}>
                    {f.path}
                  </TableCell>
                  <TableCell className="capitalize">{f.status}</TableCell>
                  <TableCell className="tabular-nums">{f.healthBefore?.toFixed(0) ?? "—"}</TableCell>
                  <TableCell className="tabular-nums">{f.healthAfter?.toFixed(0) ?? "—"}</TableCell>
                  <TableCell>
                    {f.status === "changed" && (
                      <Badge variant={deltaVariant(f.healthDelta, true)}>{formatDelta(f.healthDelta)}</Badge>
                    )}
                  </TableCell>
                </TableRow>
              ))}
              {changed.length === 0 && (
                <TableRow>
                  <TableCell colSpan={5} className="text-center text-muted-foreground">
                    No file-level changes between these scans.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}
