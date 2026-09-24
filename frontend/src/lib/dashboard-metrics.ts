import type { FileMetrics } from "@/lib/api";

export interface BucketCount {
  bucket: string;
  count: number;
}

const COMPLEXITY_BUCKETS: { label: string; max: number }[] = [
  { label: "1-5", max: 5 },
  { label: "6-10", max: 10 },
  { label: "11-20", max: 20 },
  { label: "21-40", max: 40 },
  { label: "41+", max: Infinity },
];

const HEALTH_BUCKETS: { label: string; max: number }[] = [
  { label: "Critical (0-25)", max: 25 },
  { label: "Poor (25-50)", max: 50 },
  { label: "Fair (50-75)", max: 75 },
  { label: "Good (75-100)", max: Infinity },
];

function bucketCounts(
  values: number[],
  buckets: { label: string; max: number }[],
): BucketCount[] {
  const counts = new Map(buckets.map((b) => [b.label, 0]));
  for (const value of values) {
    const bucket = buckets.find((b) => value <= b.max);
    if (bucket) counts.set(bucket.label, (counts.get(bucket.label) ?? 0) + 1);
  }
  return buckets.map((b) => ({ bucket: b.label, count: counts.get(b.label) ?? 0 }));
}

export function complexityDistribution(files: FileMetrics[]): BucketCount[] {
  return bucketCounts(
    files.map((f) => f.complexity),
    COMPLEXITY_BUCKETS,
  );
}

export function healthDistribution(files: FileMetrics[]): BucketCount[] {
  return bucketCounts(
    files.map((f) => f.healthScore),
    HEALTH_BUCKETS,
  );
}

export function extensionBreakdown(files: FileMetrics[]): BucketCount[] {
  const counts = new Map<string, number>();
  for (const f of files) {
    counts.set(f.extension, (counts.get(f.extension) ?? 0) + 1);
  }
  return Array.from(counts.entries())
    .map(([bucket, count]) => ({ bucket, count }))
    .sort((a, b) => b.count - a.count);
}

export function duplicateCount(files: FileMetrics[]): number {
  return files.filter((f) => f.isDuplicate).length;
}
