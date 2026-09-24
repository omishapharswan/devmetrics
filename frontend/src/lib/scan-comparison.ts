import type { ScanResult } from "@/lib/api";

export interface MetricDelta {
  label: string;
  before: number;
  after: number;
  delta: number;
}

export interface FileChange {
  path: string;
  status: "added" | "removed" | "changed" | "unchanged";
  healthBefore: number | null;
  healthAfter: number | null;
  healthDelta: number;
}

export interface ScanComparison {
  metrics: MetricDelta[];
  fileChanges: FileChange[];
}

function metric(label: string, before: number, after: number): MetricDelta {
  return { label, before, after, delta: after - before };
}

/**
 * Compares two scans, matching files by their relative path (file IDs
 * differ across scans since each scan re-inserts its own rows).
 */
export function compareScans(before: ScanResult, after: ScanResult): ScanComparison {
  const metrics: MetricDelta[] = [
    metric("Total Files", before.scan.totalFiles, after.scan.totalFiles),
    metric("Avg Complexity", before.scan.avgComplexity, after.scan.avgComplexity),
    metric("Avg Health Score", before.scan.avgHealthScore, after.scan.avgHealthScore),
  ];

  const beforeByPath = new Map(before.files.map((f) => [f.path, f]));
  const afterByPath = new Map(after.files.map((f) => [f.path, f]));
  const allPaths = new Set([...beforeByPath.keys(), ...afterByPath.keys()]);

  const fileChanges: FileChange[] = Array.from(allPaths)
    .map((path): FileChange => {
      const b = beforeByPath.get(path);
      const a = afterByPath.get(path);
      if (!b) {
        return { path, status: "added", healthBefore: null, healthAfter: a!.healthScore, healthDelta: 0 };
      }
      if (!a) {
        return { path, status: "removed", healthBefore: b.healthScore, healthAfter: null, healthDelta: 0 };
      }
      const healthDelta = a.healthScore - b.healthScore;
      return {
        path,
        status: healthDelta === 0 && a.complexity === b.complexity ? "unchanged" : "changed",
        healthBefore: b.healthScore,
        healthAfter: a.healthScore,
        healthDelta,
      };
    })
    .sort((x, y) => x.path.localeCompare(y.path));

  return { metrics, fileChanges };
}
