"use client";

import { Bar, BarChart, CartesianGrid, XAxis } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import type { FileMetrics } from "@/lib/api";
import {
  complexityDistribution,
  extensionBreakdown,
  healthDistribution,
} from "@/lib/dashboard-metrics";

const countConfig: ChartConfig = {
  count: { label: "Files", color: "var(--chart-1)" },
};

interface DashboardChartsProps {
  files: FileMetrics[];
}

export function DashboardCharts({ files }: DashboardChartsProps) {
  const complexity = complexityDistribution(files);
  const health = healthDistribution(files);
  const extensions = extensionBreakdown(files).slice(0, 8);

  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
      <BucketChart title="Complexity Distribution" data={complexity} />
      <BucketChart title="Health Score Distribution" data={health} />
      <BucketChart title="Files by Extension" data={extensions} />
    </div>
  );
}

function BucketChart({
  title,
  data,
}: {
  title: string;
  data: { bucket: string; count: number }[];
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <ChartContainer config={countConfig} className="aspect-auto h-[220px] w-full">
          <BarChart data={data} margin={{ left: 0, right: 8 }}>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="bucket"
              tickLine={false}
              axisLine={false}
              interval={0}
              angle={-20}
              textAnchor="end"
              height={50}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <Bar dataKey="count" fill="var(--color-count)" radius={4} />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
