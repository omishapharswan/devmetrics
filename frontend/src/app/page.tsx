"use client";

import { useState } from "react";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ScanForm } from "@/components/dashboard/scan-form";
import { OverviewCards } from "@/components/dashboard/overview-cards";
import { DashboardCharts } from "@/components/dashboard/dashboard-charts";
import { FileTable } from "@/components/dashboard/file-table";
import { DependencyGraph3D } from "@/components/dashboard/dependency-graph-3d";
import type { ScanResult } from "@/lib/api";

export default function Home() {
  const [result, setResult] = useState<ScanResult | null>(null);

  return (
    <div className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 px-6 py-10">
      <header className="flex flex-col gap-1">
        <h1 className="text-2xl font-semibold tracking-tight">DevMetrics</h1>
        <p className="text-sm text-muted-foreground">
          Scan a codebase for complexity, duplication, and dependency health.
        </p>
      </header>

      <ScanForm onScanComplete={setResult} />

      {result ? (
        <>
          <OverviewCards scan={result.scan} files={result.files} />
          <DashboardCharts files={result.files} />
          <DependencyGraph3D files={result.files} dependencies={result.dependencies} />
          <FileTable files={result.files} />
        </>
      ) : (
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle>No scan yet</CardTitle>
            <CardDescription>
              Enter an absolute folder path above and run a scan to see its metrics.
            </CardDescription>
          </CardHeader>
        </Card>
      )}
    </div>
  );
}
