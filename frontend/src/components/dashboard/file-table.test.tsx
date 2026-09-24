import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";
import { FileTable } from "@/components/dashboard/file-table";
import type { FileMetrics } from "@/lib/api";

function makeFile(overrides: Partial<FileMetrics>): FileMetrics {
  return {
    id: 1,
    scanId: 1,
    path: "a.ts",
    extension: ".ts",
    sizeBytes: 100,
    linesOfCode: 10,
    commentLines: 1,
    commentRatio: 0.1,
    complexity: 1,
    isDuplicate: false,
    healthScore: 100,
    ...overrides,
  };
}

const files: FileMetrics[] = [
  makeFile({ id: 1, path: "b-medium.go", extension: ".go", complexity: 10, isDuplicate: false }),
  makeFile({ id: 2, path: "a-high.ts", extension: ".ts", complexity: 30, isDuplicate: true }),
  makeFile({ id: 3, path: "c-low.ts", extension: ".ts", complexity: 2, isDuplicate: false }),
];

function pathCells() {
  return screen.getAllByTitle(/\.(ts|go)$/).map((el) => el.textContent);
}

describe("FileTable", () => {
  it("renders every file initially", () => {
    render(<FileTable files={files} />);
    expect(pathCells()).toEqual(["b-medium.go", "a-high.ts", "c-low.ts"]);
  });

  it("sorts by a column when its header is clicked, toggling direction on a second click", async () => {
    render(<FileTable files={files} />);

    await userEvent.click(screen.getByRole("button", { name: /complexity/i }));
    expect(pathCells()).toEqual(["c-low.ts", "b-medium.go", "a-high.ts"]);

    await userEvent.click(screen.getByRole("button", { name: /complexity/i }));
    expect(pathCells()).toEqual(["a-high.ts", "b-medium.go", "c-low.ts"]);
  });

  it("filters rows by path search text", async () => {
    render(<FileTable files={files} />);

    await userEvent.type(screen.getByPlaceholderText(/filter by path/i), "high");
    expect(pathCells()).toEqual(["a-high.ts"]);
  });

  it("filters to duplicates only when the badge is toggled", async () => {
    render(<FileTable files={files} />);

    await userEvent.click(screen.getByText("Duplicates only"));
    expect(pathCells()).toEqual(["a-high.ts"]);
  });

  it("opens a detail sheet with the file's metrics when a row is clicked", async () => {
    render(<FileTable files={files} />);

    await userEvent.click(screen.getByTitle("a-high.ts"));

    const dialog = await screen.findByRole("dialog");
    expect(within(dialog).getByText("a-high.ts")).toBeInTheDocument();
    expect(within(dialog).getByText("30")).toBeInTheDocument(); // complexity
  });
});
