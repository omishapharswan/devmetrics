import { describe, expect, it } from "vitest";
import type { FileMetrics } from "@/lib/api";
import { filterFiles, sortFiles, uniqueExtensions } from "@/lib/file-table";

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

describe("sortFiles", () => {
  const files = [
    makeFile({ id: 1, path: "b.ts", complexity: 5 }),
    makeFile({ id: 2, path: "a.ts", complexity: 20 }),
    makeFile({ id: 3, path: "c.ts", complexity: 1 }),
  ];

  it("sorts numeric columns ascending and descending", () => {
    expect(sortFiles(files, "complexity", "asc").map((f) => f.id)).toEqual([3, 1, 2]);
    expect(sortFiles(files, "complexity", "desc").map((f) => f.id)).toEqual([2, 1, 3]);
  });

  it("sorts string columns alphabetically", () => {
    expect(sortFiles(files, "path", "asc").map((f) => f.path)).toEqual(["a.ts", "b.ts", "c.ts"]);
  });

  it("does not mutate the input array", () => {
    const original = [...files];
    sortFiles(files, "complexity", "asc");
    expect(files).toEqual(original);
  });
});

describe("filterFiles", () => {
  const files = [
    makeFile({ id: 1, path: "src/index.ts", extension: ".ts", isDuplicate: false }),
    makeFile({ id: 2, path: "src/utils.go", extension: ".go", isDuplicate: true }),
    makeFile({ id: 3, path: "test/index.ts", extension: ".ts", isDuplicate: false }),
  ];

  it("filters by case-insensitive path substring", () => {
    const result = filterFiles(files, { query: "INDEX", extension: "all", duplicatesOnly: false });
    expect(result.map((f) => f.id)).toEqual([1, 3]);
  });

  it("filters by extension", () => {
    const result = filterFiles(files, { query: "", extension: ".go", duplicatesOnly: false });
    expect(result.map((f) => f.id)).toEqual([2]);
  });

  it("filters to duplicates only", () => {
    const result = filterFiles(files, { query: "", extension: "all", duplicatesOnly: true });
    expect(result.map((f) => f.id)).toEqual([2]);
  });

  it("combines all filters", () => {
    const result = filterFiles(files, { query: "src", extension: ".ts", duplicatesOnly: false });
    expect(result.map((f) => f.id)).toEqual([1]);
  });
});

describe("uniqueExtensions", () => {
  it("returns sorted, deduplicated extensions", () => {
    const files = [makeFile({ extension: ".ts" }), makeFile({ extension: ".go" }), makeFile({ extension: ".ts" })];
    expect(uniqueExtensions(files)).toEqual([".go", ".ts"]);
  });
});
