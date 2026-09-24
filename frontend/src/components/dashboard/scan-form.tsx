"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ApiError, runScan, type ScanResult } from "@/lib/api";

interface ScanFormProps {
  onScanComplete: (result: ScanResult) => void;
}

export function ScanForm({ onScanComplete }: ScanFormProps) {
  const [folderPath, setFolderPath] = useState("");
  const [isScanning, setIsScanning] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!folderPath.trim()) return;

    setIsScanning(true);
    setError(null);
    try {
      const result = await runScan(folderPath.trim());
      onScanComplete(result);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to reach the DevMetrics API.");
    } finally {
      setIsScanning(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-2 sm:flex-row sm:items-start">
      <div className="flex-1">
        <Input
          value={folderPath}
          onChange={(e) => setFolderPath(e.target.value)}
          placeholder="/path/to/your/project"
          aria-label="Folder path to scan"
          disabled={isScanning}
        />
        {error && <p className="mt-1 text-sm text-destructive">{error}</p>}
      </div>
      <Button type="submit" disabled={isScanning || !folderPath.trim()}>
        {isScanning ? "Scanning..." : "Run Scan"}
      </Button>
    </form>
  );
}
