import type { FileMetrics } from "@/lib/api";

export type SortColumn =
  | "path"
  | "extension"
  | "linesOfCode"
  | "complexity"
  | "commentRatio"
  | "healthScore";

export type SortDirection = "asc" | "desc";

export function sortFiles(
  files: FileMetrics[],
  column: SortColumn,
  direction: SortDirection,
): FileMetrics[] {
  const factor = direction === "asc" ? 1 : -1;
  return [...files].sort((a, b) => {
    const av = a[column];
    const bv = b[column];
    if (typeof av === "string" && typeof bv === "string") {
      return av.localeCompare(bv) * factor;
    }
    return (Number(av) - Number(bv)) * factor;
  });
}

export interface FileFilters {
  query: string;
  extension: string | "all";
  duplicatesOnly: boolean;
}

export function filterFiles(files: FileMetrics[], filters: FileFilters): FileMetrics[] {
  const query = filters.query.trim().toLowerCase();
  return files.filter((f) => {
    if (query && !f.path.toLowerCase().includes(query)) return false;
    if (filters.extension !== "all" && f.extension !== filters.extension) return false;
    if (filters.duplicatesOnly && !f.isDuplicate) return false;
    return true;
  });
}

export function uniqueExtensions(files: FileMetrics[]): string[] {
  return Array.from(new Set(files.map((f) => f.extension))).sort();
}
