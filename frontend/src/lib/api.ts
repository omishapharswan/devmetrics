const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export interface Scan {
  id: number;
  folderPath: string;
  scannedAt: string;
  totalFiles: number;
  avgComplexity: number;
  avgHealthScore: number;
}

export interface FileMetrics {
  id: number;
  scanId: number;
  path: string;
  extension: string;
  sizeBytes: number;
  linesOfCode: number;
  commentLines: number;
  commentRatio: number;
  complexity: number;
  isDuplicate: boolean;
  healthScore: number;
}

export interface ScanResult {
  scan: Scan;
  files: FileMetrics[];
}

export class ApiError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ApiError";
  }
}

export async function runScan(folderPath: string): Promise<ScanResult> {
  const res = await fetch(`${API_BASE_URL}/api/scan`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ folderPath }),
  });

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new ApiError(body?.error ?? `Scan failed with status ${res.status}`);
  }

  return res.json();
}
