"use client";

import { useMemo, useState } from "react";
import { Canvas } from "@react-three/fiber";
import { Html, Line, OrbitControls } from "@react-three/drei";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import type { Dependency, FileMetrics } from "@/lib/api";
import { computeForceLayout3D, type Vec3 } from "@/lib/graph-layout";

interface DependencyGraph3DProps {
  files: FileMetrics[];
  dependencies: Dependency[];
}

function healthColor(score: number): string {
  if (score >= 75) return "#22c55e";
  if (score >= 50) return "#eab308";
  return "#ef4444";
}

function nodeRadius(complexity: number): number {
  return 0.25 + Math.min(complexity, 40) / 40;
}

export function DependencyGraph3D({ files, dependencies }: DependencyGraph3DProps) {
  const [hovered, setHovered] = useState<FileMetrics | null>(null);

  const indexById = useMemo(() => {
    const map = new Map<number, number>();
    files.forEach((f, i) => map.set(f.id, i));
    return map;
  }, [files]);

  const edges = useMemo(() => {
    return dependencies
      .map((d) => ({ from: indexById.get(d.fromFileId), to: indexById.get(d.toFileId) }))
      .filter((e): e is { from: number; to: number } => e.from !== undefined && e.to !== undefined);
  }, [dependencies, indexById]);

  const positions = useMemo(() => computeForceLayout3D(files.length, edges), [files.length, edges]);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Dependency Graph</CardTitle>
        <CardDescription>
          {edges.length > 0
            ? `${files.length} files, ${edges.length} resolved import edges. Drag to rotate, scroll to zoom.`
            : `${files.length} files. No resolvable import edges were found for this scan's languages.`}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {files.length === 0 ? (
          <p className="text-sm text-muted-foreground">No files to display.</p>
        ) : (
          <div className="h-[420px] w-full overflow-hidden rounded-lg bg-black/90">
            <Canvas camera={{ position: [0, 0, 26], fov: 50 }}>
              <ambientLight intensity={0.6} />
              <pointLight position={[15, 15, 15]} intensity={1} />
              <OrbitControls enableDamping dampingFactor={0.1} />
              {edges.map((edge, i) => (
                <Line
                  key={i}
                  points={[positions[edge.from], positions[edge.to]]}
                  color="#64748b"
                  lineWidth={1}
                  transparent
                  opacity={0.5}
                />
              ))}
              {files.map((file, i) => (
                <GraphNode
                  key={file.id}
                  file={file}
                  position={positions[i]}
                  isHovered={hovered?.id === file.id}
                  onHoverChange={(hovering) => setHovered(hovering ? file : null)}
                />
              ))}
            </Canvas>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function GraphNode({
  file,
  position,
  isHovered,
  onHoverChange,
}: {
  file: FileMetrics;
  position: Vec3;
  isHovered: boolean;
  onHoverChange: (hovering: boolean) => void;
}) {
  return (
    <mesh
      position={position}
      onPointerOver={(e) => {
        e.stopPropagation();
        onHoverChange(true);
      }}
      onPointerOut={() => onHoverChange(false)}
    >
      <sphereGeometry args={[nodeRadius(file.complexity), 16, 16]} />
      <meshStandardMaterial color={healthColor(file.healthScore)} />
      {isHovered && (
        <Html distanceFactor={12}>
          <div className="pointer-events-none rounded-md bg-popover px-2 py-1 text-xs whitespace-nowrap text-popover-foreground shadow-md ring-1 ring-foreground/10">
            {file.path}
            <span className="ml-1.5 text-muted-foreground">
              · health {file.healthScore.toFixed(0)} · complexity {file.complexity}
            </span>
          </div>
        </Html>
      )}
    </mesh>
  );
}
