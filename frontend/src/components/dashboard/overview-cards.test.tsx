import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { OverviewCards } from "@/components/dashboard/overview-cards";
import type { FileMetrics, Scan } from "@/lib/api";

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

describe("OverviewCards", () => {
  it("renders scan-level aggregates and the derived duplicate count", () => {
    const scan: Scan = {
      id: 1,
      folderPath: "/proj",
      scannedAt: "2026-01-01T00:00:00Z",
      totalFiles: 3,
      avgComplexity: 12.345,
      avgHealthScore: 87.6,
    };
    const files = [
      makeFile({ id: 1, isDuplicate: true }),
      makeFile({ id: 2, isDuplicate: false }),
      makeFile({ id: 3, isDuplicate: true }),
    ];

    render(<OverviewCards scan={scan} files={files} />);

    expect(screen.getByText("Total Files")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
    expect(screen.getByText("12.3")).toBeInTheDocument();
    expect(screen.getByText("88")).toBeInTheDocument();
    expect(screen.getByText("Duplicate Files")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
  });
});
