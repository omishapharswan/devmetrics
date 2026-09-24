export interface GraphEdge {
  from: number;
  to: number;
}

export type Vec3 = [number, number, number];

const REPULSION = 18;
const ATTRACTION = 0.02;
const DAMPING = 0.85;
const ITERATIONS = 200;
const SPREAD = 12;

/**
 * A minimal force-directed layout: nodes repel each other, connected nodes
 * attract along their edge, iterated until roughly settled. Positions are
 * indexed the same as the input node count (index i -> layout[i]).
 */
export function computeForceLayout3D(nodeCount: number, edges: GraphEdge[]): Vec3[] {
  if (nodeCount === 0) return [];

  const positions: Vec3[] = Array.from({ length: nodeCount }, (_, i) => {
    const phi = Math.acos(1 - (2 * (i + 0.5)) / nodeCount);
    const theta = Math.PI * (1 + Math.sqrt(5)) * i;
    return [
      SPREAD * Math.sin(phi) * Math.cos(theta),
      SPREAD * Math.sin(phi) * Math.sin(theta),
      SPREAD * Math.cos(phi),
    ];
  });

  const velocities: Vec3[] = Array.from({ length: nodeCount }, () => [0, 0, 0]);

  for (let iter = 0; iter < ITERATIONS; iter++) {
    const forces: Vec3[] = Array.from({ length: nodeCount }, () => [0, 0, 0]);

    for (let i = 0; i < nodeCount; i++) {
      for (let j = i + 1; j < nodeCount; j++) {
        const dx = positions[i][0] - positions[j][0];
        const dy = positions[i][1] - positions[j][1];
        const dz = positions[i][2] - positions[j][2];
        const distSq = Math.max(dx * dx + dy * dy + dz * dz, 0.01);
        const force = REPULSION / distSq;
        const dist = Math.sqrt(distSq);
        const fx = (dx / dist) * force;
        const fy = (dy / dist) * force;
        const fz = (dz / dist) * force;
        forces[i][0] += fx;
        forces[i][1] += fy;
        forces[i][2] += fz;
        forces[j][0] -= fx;
        forces[j][1] -= fy;
        forces[j][2] -= fz;
      }
    }

    for (const edge of edges) {
      if (edge.from === edge.to) continue;
      const a = positions[edge.from];
      const b = positions[edge.to];
      if (!a || !b) continue;
      const dx = b[0] - a[0];
      const dy = b[1] - a[1];
      const dz = b[2] - a[2];
      forces[edge.from][0] += dx * ATTRACTION;
      forces[edge.from][1] += dy * ATTRACTION;
      forces[edge.from][2] += dz * ATTRACTION;
      forces[edge.to][0] -= dx * ATTRACTION;
      forces[edge.to][1] -= dy * ATTRACTION;
      forces[edge.to][2] -= dz * ATTRACTION;
    }

    for (let i = 0; i < nodeCount; i++) {
      velocities[i][0] = (velocities[i][0] + forces[i][0]) * DAMPING;
      velocities[i][1] = (velocities[i][1] + forces[i][1]) * DAMPING;
      velocities[i][2] = (velocities[i][2] + forces[i][2]) * DAMPING;
      positions[i][0] += velocities[i][0];
      positions[i][1] += velocities[i][1];
      positions[i][2] += velocities[i][2];
    }
  }

  return positions;
}
