import { Badge } from "@/components/ui/badge";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import type { FileMetrics } from "@/lib/api";

interface FileDetailSheetProps {
  file: FileMetrics | null;
  onOpenChange: (open: boolean) => void;
}

function healthLabel(score: number): { label: string; variant: "default" | "secondary" | "destructive" } {
  if (score >= 75) return { label: "Good", variant: "default" };
  if (score >= 50) return { label: "Fair", variant: "secondary" };
  return { label: "Needs attention", variant: "destructive" };
}

export function FileDetailSheet({ file, onOpenChange }: FileDetailSheetProps) {
  const health = file ? healthLabel(file.healthScore) : null;

  return (
    <Sheet open={file !== null} onOpenChange={onOpenChange}>
      <SheetContent>
        {file && (
          <>
            <SheetHeader>
              <SheetTitle className="break-all font-mono text-sm">{file.path}</SheetTitle>
              <SheetDescription>File-level metrics from the most recent scan.</SheetDescription>
            </SheetHeader>
            <div className="flex flex-col gap-3 px-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Health score</span>
                <div className="flex items-center gap-2">
                  <span className="font-medium tabular-nums">{file.healthScore.toFixed(0)}</span>
                  {health && <Badge variant={health.variant}>{health.label}</Badge>}
                </div>
              </div>
              <DetailRow label="Extension" value={file.extension} />
              <DetailRow label="Size" value={`${file.sizeBytes.toLocaleString()} bytes`} />
              <DetailRow label="Lines of code" value={file.linesOfCode.toLocaleString()} />
              <DetailRow label="Comment lines" value={file.commentLines.toLocaleString()} />
              <DetailRow label="Comment ratio" value={`${(file.commentRatio * 100).toFixed(1)}%`} />
              <DetailRow label="Cyclomatic complexity" value={file.complexity.toLocaleString()} />
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Duplicate</span>
                <Badge variant={file.isDuplicate ? "destructive" : "secondary"}>
                  {file.isDuplicate ? "Yes" : "No"}
                </Badge>
              </div>
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className="font-medium tabular-nums">{value}</span>
    </div>
  );
}
