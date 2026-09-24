"use client";

import { useMemo, useState } from "react";
import { ArrowDown, ArrowUp, ArrowUpDown } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { FileMetrics } from "@/lib/api";
import {
  filterFiles,
  sortFiles,
  uniqueExtensions,
  type SortColumn,
  type SortDirection,
} from "@/lib/file-table";
import { FileDetailSheet } from "@/components/dashboard/file-detail-sheet";

interface FileTableProps {
  files: FileMetrics[];
}

const COLUMNS: { key: SortColumn; label: string }[] = [
  { key: "path", label: "Path" },
  { key: "extension", label: "Ext" },
  { key: "linesOfCode", label: "LOC" },
  { key: "complexity", label: "Complexity" },
  { key: "commentRatio", label: "Comments" },
  { key: "healthScore", label: "Health" },
];

export function FileTable({ files }: FileTableProps) {
  const [query, setQuery] = useState("");
  const [extension, setExtension] = useState<string>("all");
  const [duplicatesOnly, setDuplicatesOnly] = useState(false);
  const [sortColumn, setSortColumn] = useState<SortColumn>("healthScore");
  const [sortDirection, setSortDirection] = useState<SortDirection>("asc");
  const [selectedFile, setSelectedFile] = useState<FileMetrics | null>(null);

  const extensions = useMemo(() => uniqueExtensions(files), [files]);

  const visibleFiles = useMemo(() => {
    const filtered = filterFiles(files, { query, extension, duplicatesOnly });
    return sortFiles(filtered, sortColumn, sortDirection);
  }, [files, query, extension, duplicatesOnly, sortColumn, sortDirection]);

  function toggleSort(column: SortColumn) {
    if (column === sortColumn) {
      setSortDirection((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortColumn(column);
      setSortDirection("asc");
    }
  }

  return (
    <Card>
      <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <CardTitle className="text-base">Files ({visibleFiles.length})</CardTitle>
        <div className="flex flex-wrap items-center gap-2">
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Filter by path..."
            className="w-48"
          />
          <Select value={extension} onValueChange={(value) => setExtension(value ?? "all")}>
            <SelectTrigger size="sm" className="w-28">
              <SelectValue placeholder="Extension" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All extensions</SelectItem>
              {extensions.map((ext) => (
                <SelectItem key={ext} value={ext}>
                  {ext}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <button
            type="button"
            onClick={() => setDuplicatesOnly((v) => !v)}
            className="focus-visible:outline-none"
          >
            <Badge variant={duplicatesOnly ? "destructive" : "outline"} className="cursor-pointer">
              Duplicates only
            </Badge>
          </button>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              {COLUMNS.map((col) => (
                <TableHead key={col.key}>
                  <button
                    type="button"
                    onClick={() => toggleSort(col.key)}
                    className="flex items-center gap-1 hover:text-foreground"
                  >
                    {col.label}
                    <SortIcon active={sortColumn === col.key} direction={sortDirection} />
                  </button>
                </TableHead>
              ))}
              <TableHead>Duplicate</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {visibleFiles.map((file) => (
              <TableRow
                key={file.id}
                onClick={() => setSelectedFile(file)}
                className="cursor-pointer"
              >
                <TableCell className="max-w-xs truncate font-mono text-xs" title={file.path}>
                  {file.path}
                </TableCell>
                <TableCell>{file.extension}</TableCell>
                <TableCell className="tabular-nums">{file.linesOfCode}</TableCell>
                <TableCell className="tabular-nums">{file.complexity}</TableCell>
                <TableCell className="tabular-nums">{(file.commentRatio * 100).toFixed(0)}%</TableCell>
                <TableCell className="tabular-nums">{file.healthScore.toFixed(0)}</TableCell>
                <TableCell>
                  {file.isDuplicate && <Badge variant="destructive">Yes</Badge>}
                </TableCell>
              </TableRow>
            ))}
            {visibleFiles.length === 0 && (
              <TableRow>
                <TableCell colSpan={COLUMNS.length + 1} className="text-center text-muted-foreground">
                  No files match the current filters.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </CardContent>
      <FileDetailSheet file={selectedFile} onOpenChange={(open) => !open && setSelectedFile(null)} />
    </Card>
  );
}

function SortIcon({ active, direction }: { active: boolean; direction: SortDirection }) {
  if (!active) return <ArrowUpDown className="size-3.5 text-muted-foreground" />;
  return direction === "asc" ? (
    <ArrowUp className="size-3.5" />
  ) : (
    <ArrowDown className="size-3.5" />
  );
}
