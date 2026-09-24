import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function Home() {
  return (
    <div className="flex flex-1 items-center justify-center bg-zinc-50 px-6 dark:bg-black">
      <Card className="max-w-md">
        <CardHeader>
          <CardTitle>DevMetrics</CardTitle>
          <CardDescription>
            Scaffold is up — dashboard, file table, and dependency graph land in upcoming commits.
          </CardDescription>
        </CardHeader>
      </Card>
    </div>
  );
}
