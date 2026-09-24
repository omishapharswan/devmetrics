"use client";

import { useState } from "react";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScanForm } from "@/components/dashboard/scan-form";
import { OverviewCards } from "@/components/dashboard/overview-cards";
import { DashboardCharts } from "@/components/dashboard/dashboard-charts";
import { FileTable } from "@/components/dashboard/file-table";
import { DependencyGraph3D } from "@/components/dashboard/dependency-graph-3d";
import { ScanHistory } from "@/components/dashboard/scan-history";
import { ThemeToggle } from "@/components/theme-toggle";
import type { ScanResult } from "@/lib/api";

export default function Home() {
  const [result, setResult] = useState<ScanResult | null>(null);
  const [activeTab, setActiveTab] = useState("scan");

  function handleViewScan(scan: ScanResult) {
    setResult(scan);
    setActiveTab("scan");
  }

  return (
    <div className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 px-4 py-8 sm:px-6 sm:py-10">
      <header className="flex items-start justify-between gap-4">
        <div className="flex flex-col gap-1">
          <h1 className="text-2xl font-semibold tracking-tight">DevMetrics</h1>
          <p className="text-sm text-muted-foreground">
            Scan a codebase for complexity, duplication, and dependency health.
          </p>
        </div>
        <ThemeToggle />
      </header>

      <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as string)}>
        <TabsList>
          <TabsTrigger value="scan">Scan</TabsTrigger>
          <TabsTrigger value="history">History</TabsTrigger>
        </TabsList>

        <TabsContent value="scan" className="flex flex-col gap-6">
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
        </TabsContent>

        <TabsContent value="history">
          <ScanHistory onViewScan={handleViewScan} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
