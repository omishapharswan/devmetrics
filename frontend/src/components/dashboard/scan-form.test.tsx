import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ScanForm } from "@/components/dashboard/scan-form";
import type { ScanResult } from "@/lib/api";

const sampleResult: ScanResult = {
  scan: {
    id: 1,
    folderPath: "/proj",
    scannedAt: "2026-01-01T00:00:00Z",
    totalFiles: 1,
    avgComplexity: 1,
    avgHealthScore: 100,
  },
  files: [],
  dependencies: [],
};

describe("ScanForm", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("disables the submit button until a folder path is entered", async () => {
    render(<ScanForm onScanComplete={vi.fn()} />);
    expect(screen.getByRole("button", { name: /run scan/i })).toBeDisabled();

    await userEvent.type(screen.getByLabelText(/folder path/i), "/proj");
    expect(screen.getByRole("button", { name: /run scan/i })).toBeEnabled();
  });

  it("submits the folder path and reports the parsed result", async () => {
    const fetchMock = fetch as unknown as ReturnType<typeof vi.fn>;
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: async () => sampleResult,
    });

    const onScanComplete = vi.fn();
    render(<ScanForm onScanComplete={onScanComplete} />);

    await userEvent.type(screen.getByLabelText(/folder path/i), "/proj");
    await userEvent.click(screen.getByRole("button", { name: /run scan/i }));

    await waitFor(() => expect(onScanComplete).toHaveBeenCalledWith(sampleResult));

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining("/api/scan"),
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ folderPath: "/proj" }),
      }),
    );
  });

  it("shows the server error message and does not call onScanComplete when the scan fails", async () => {
    const fetchMock = fetch as unknown as ReturnType<typeof vi.fn>;
    fetchMock.mockResolvedValueOnce({
      ok: false,
      status: 400,
      json: async () => ({ error: "folder not found" }),
    });

    const onScanComplete = vi.fn();
    render(<ScanForm onScanComplete={onScanComplete} />);

    await userEvent.type(screen.getByLabelText(/folder path/i), "/missing");
    await userEvent.click(screen.getByRole("button", { name: /run scan/i }));

    expect(await screen.findByText("folder not found")).toBeInTheDocument();
    expect(onScanComplete).not.toHaveBeenCalled();
  });
});
