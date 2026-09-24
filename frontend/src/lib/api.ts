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

export interface Dependency {
  id: number;
  scanId: number;
  fromFileId: number;
  toFileId: number;
}

export interface ScanResult {
  scan: Scan;
  files: FileMetrics[];
  dependencies: Dependency[];
}

export class ApiError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, init);
  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new ApiError(body?.error ?? `Request failed with status ${res.status}`);
  }
  return res.json();
}

export function runScan(folderPath: string): Promise<ScanResult> {
  return request<ScanResult>("/api/scan", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ folderPath }),
  });
}

export async function listScans(): Promise<Scan[]> {
  const { scans } = await request<{ scans: Scan[] }>("/api/scans");
  return scans;
}

export function getScan(id: number): Promise<ScanResult> {
  return request<ScanResult>(`/api/scans/${id}`);
}
