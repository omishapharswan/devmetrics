"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { ApiError, getScan, listScans, type Scan, type ScanResult } from "@/lib/api";
import { compareScans, type ScanComparison } from "@/lib/scan-comparison";
import { ScanComparisonView } from "@/components/dashboard/scan-comparison-view";

interface ScanHistoryProps {
  onViewScan: (result: ScanResult) => void;
}

export function ScanHistory({ onViewScan }: ScanHistoryProps) {
  const [scans, setScans] = useState<Scan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selected, setSelected] = useState<number[]>([]);
  const [comparison, setComparison] = useState<ScanComparison | null>(null);
  const [comparing, setComparing] = useState(false);

  useEffect(() => {
    listScans()
      .then(setScans)
      .catch((err) => setError(err instanceof ApiError ? err.message : "Failed to load scan history."))
      .finally(() => setLoading(false));
  }, []);

  function toggleSelected(id: number, checked: boolean) {
    setComparison(null);
    setSelected((prev) => {
      if (checked) {
        const next = [...prev, id];
        return next.length > 2 ? next.slice(1) : next;
      }
      return prev.filter((s) => s !== id);
    });
  }

  async function handleCompare() {
    if (selected.length !== 2) return;
    setComparing(true);
    setError(null);
    try {
      // selected[0] is the older selection, selected[1] the newer, since
      // toggleSelected appends in click order.
      const [before, after] = await Promise.all([getScan(selected[0]), getScan(selected[1])]);
      setComparison(compareScans(before, after));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to load scans for comparison.");
    } finally {
      setComparing(false);
    }
  }

  async function handleView(id: number) {
    try {
      onViewScan(await getScan(id));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to load scan.");
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="text-base">Scan History</CardTitle>
          <Button size="sm" disabled={selected.length !== 2 || comparing} onClick={handleCompare}>
            {comparing ? "Comparing..." : "Compare Selected"}
          </Button>
        </CardHeader>
        <CardContent>
          {error && <p className="mb-2 text-sm text-destructive">{error}</p>}
          {loading ? (
            <p className="text-sm text-muted-foreground">Loading...</p>
          ) : scans.length === 0 ? (
            <p className="text-sm text-muted-foreground">No scans yet. Run one from the Scan tab.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-8" />
                  <TableHead>Folder</TableHead>
                  <TableHead>Scanned At</TableHead>
                  <TableHead>Files</TableHead>
                  <TableHead>Avg Complexity</TableHead>
                  <TableHead>Avg Health</TableHead>
                  <TableHead />
                </TableRow>
              </TableHeader>
              <TableBody>
                {scans.map((scan) => (
                  <TableRow key={scan.id}>
                    <TableCell>
                      <Checkbox
                        checked={selected.includes(scan.id)}
                        onCheckedChange={(checked) => toggleSelected(scan.id, checked === true)}
                        aria-label={`Select scan ${scan.id} for comparison`}
                      />
                    </TableCell>
                    <TableCell className="max-w-xs truncate font-mono text-xs" title={scan.folderPath}>
                      {scan.folderPath}
                    </TableCell>
                    <TableCell>{new Date(scan.scannedAt).toLocaleString()}</TableCell>
                    <TableCell className="tabular-nums">{scan.totalFiles}</TableCell>
                    <TableCell className="tabular-nums">{scan.avgComplexity.toFixed(1)}</TableCell>
                    <TableCell className="tabular-nums">{scan.avgHealthScore.toFixed(0)}</TableCell>
                    <TableCell>
                      <Button size="sm" variant="outline" onClick={() => handleView(scan.id)}>
                        View
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {comparison && <ScanComparisonView comparison={comparison} />}
    </div>
  );
}
