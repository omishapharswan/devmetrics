import { describe, expect, it } from "vitest";
import type { FileMetrics, Scan, ScanResult } from "@/lib/api";
import { compareScans } from "@/lib/scan-comparison";

function makeScan(overrides: Partial<Scan>): Scan {
  return {
    id: 1,
    folderPath: "/proj",
    scannedAt: "2026-01-01T00:00:00Z",
    totalFiles: 1,
    avgComplexity: 1,
    avgHealthScore: 100,
    ...overrides,
  };
}

function makeFile(overrides: Partial<FileMetrics>): FileMetrics {
  return {
    id: 1,
    scanId: 1,
    path: "a.ts",
    extension: ".ts",
    sizeBytes: 100,
    linesOfCode: 10,
    commentLines: 1,
    commentRatio: 0.1,
    complexity: 1,
    isDuplicate: false,
    healthScore: 100,
    ...overrides,
  };
}

describe("compareScans", () => {
  it("computes metric deltas between the two scans", () => {
    const before: ScanResult = {
      scan: makeScan({ totalFiles: 2, avgComplexity: 5, avgHealthScore: 80 }),
      files: [],
      dependencies: [],
    };
    const after: ScanResult = {
      scan: makeScan({ totalFiles: 3, avgComplexity: 8, avgHealthScore: 70 }),
      files: [],
      dependencies: [],
    };

    const { metrics } = compareScans(before, after);
    expect(metrics).toEqual([
      { label: "Total Files", before: 2, after: 3, delta: 1 },
      { label: "Avg Complexity", before: 5, after: 8, delta: 3 },
      { label: "Avg Health Score", before: 80, after: 70, delta: -10 },
    ]);
  });

  it("classifies files as added, removed, changed, or unchanged by path", () => {
    const before: ScanResult = {
      scan: makeScan({}),
      files: [
        makeFile({ path: "stable.ts", healthScore: 90, complexity: 2 }),
        makeFile({ path: "removed.ts", healthScore: 50, complexity: 1 }),
        makeFile({ path: "regressed.ts", healthScore: 90, complexity: 2 }),
      ],
      dependencies: [],
    };
    const after: ScanResult = {
      scan: makeScan({}),
      files: [
        makeFile({ path: "stable.ts", healthScore: 90, complexity: 2 }),
        makeFile({ path: "added.ts", healthScore: 60, complexity: 3 }),
        makeFile({ path: "regressed.ts", healthScore: 40, complexity: 5 }),
      ],
      dependencies: [],
    };

    const { fileChanges } = compareScans(before, after);
    const byPath = Object.fromEntries(fileChanges.map((f) => [f.path, f]));

    expect(byPath["stable.ts"].status).toBe("unchanged");
    expect(byPath["removed.ts"].status).toBe("removed");
    expect(byPath["removed.ts"].healthAfter).toBeNull();
    expect(byPath["added.ts"].status).toBe("added");
    expect(byPath["added.ts"].healthBefore).toBeNull();
    expect(byPath["regressed.ts"].status).toBe("changed");
    expect(byPath["regressed.ts"].healthDelta).toBe(-50);
  });
});
