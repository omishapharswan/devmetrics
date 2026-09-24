import { describe, expect, it } from "vitest";
import type { FileMetrics } from "@/lib/api";
import {
  complexityDistribution,
  duplicateCount,
  extensionBreakdown,
  healthDistribution,
} from "@/lib/dashboard-metrics";

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

describe("complexityDistribution", () => {
  it("buckets files into the correct complexity ranges", () => {
    const files = [
      makeFile({ complexity: 3 }),
      makeFile({ complexity: 8 }),
      makeFile({ complexity: 15 }),
      makeFile({ complexity: 35 }),
      makeFile({ complexity: 100 }),
    ];
    const buckets = complexityDistribution(files);
    expect(buckets).toEqual([
      { bucket: "1-5", count: 1 },
      { bucket: "6-10", count: 1 },
      { bucket: "11-20", count: 1 },
      { bucket: "21-40", count: 1 },
      { bucket: "41+", count: 1 },
    ]);
  });

  it("returns zero counts for an empty file list", () => {
    const buckets = complexityDistribution([]);
    expect(buckets.every((b) => b.count === 0)).toBe(true);
  });
});

describe("healthDistribution", () => {
  it("buckets a boundary value into the higher bucket (inclusive upper bound)", () => {
    const buckets = healthDistribution([makeFile({ healthScore: 75 })]);
    const fair = buckets.find((b) => b.bucket.startsWith("Fair"));
    const good = buckets.find((b) => b.bucket.startsWith("Good"));
    expect(fair?.count).toBe(1);
    expect(good?.count).toBe(0);
  });
});

describe("extensionBreakdown", () => {
  it("counts and sorts extensions by frequency, descending", () => {
    const files = [
      makeFile({ extension: ".go" }),
      makeFile({ extension: ".ts" }),
      makeFile({ extension: ".ts" }),
      makeFile({ extension: ".ts" }),
    ];
    expect(extensionBreakdown(files)).toEqual([
      { bucket: ".ts", count: 3 },
      { bucket: ".go", count: 1 },
    ]);
  });
});

describe("duplicateCount", () => {
  it("counts only files flagged as duplicates", () => {
    const files = [
      makeFile({ isDuplicate: true }),
      makeFile({ isDuplicate: false }),
      makeFile({ isDuplicate: true }),
    ];
    expect(duplicateCount(files)).toBe(2);
  });
});
