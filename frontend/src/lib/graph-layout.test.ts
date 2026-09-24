import { describe, expect, it } from "vitest";
import { computeForceLayout3D } from "@/lib/graph-layout";

describe("computeForceLayout3D", () => {
  it("returns an empty array for zero nodes", () => {
    expect(computeForceLayout3D(0, [])).toEqual([]);
  });

  it("returns one finite position per node", () => {
    const positions = computeForceLayout3D(12, [
      { from: 0, to: 1 },
      { from: 1, to: 2 },
    ]);
    expect(positions).toHaveLength(12);
    for (const [x, y, z] of positions) {
      expect(Number.isFinite(x)).toBe(true);
      expect(Number.isFinite(y)).toBe(true);
      expect(Number.isFinite(z)).toBe(true);
    }
  });

  it("spreads unconnected nodes apart rather than collapsing to one point", () => {
    const positions = computeForceLayout3D(5, []);
    const [x0, y0, z0] = positions[0];
    const [x1, y1, z1] = positions[1];
    const dist = Math.hypot(x1 - x0, y1 - y0, z1 - z0);
    expect(dist).toBeGreaterThan(0.5);
  });

  it("ignores edges referencing out-of-range indices instead of throwing", () => {
    expect(() => computeForceLayout3D(3, [{ from: 0, to: 99 }])).not.toThrow();
  });
});
