import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import App from "./App";

vi.mock("../services/env", () => ({
  loadConfig: vi.fn().mockResolvedValue({
    liffId: "",
    apiBase: "",
    useSharedToken: false,
    sharedToken: "",
    oidcEnabled: false,
    env: "dev",
  }),
}));
vi.mock("../services/liff", () => ({
  initLiff: vi.fn().mockResolvedValue(undefined),
  getGroupIdOrThrow: vi.fn().mockResolvedValue("gid"),
  getProfileSafe: vi.fn().mockResolvedValue({ displayName: "", userId: "" }),
}));
vi.mock("../api/client", () => ({ postScan: vi.fn().mockResolvedValue({}) }));
vi.mock("../ui/dom", () => ({ showToast: vi.fn() }));

describe("App", () => {
  it("renders header", async () => {
    render(<App />);
    expect(
      await screen.findByRole("heading", { name: /qr scanner/i }),
    ).toBeInTheDocument();
  });
});
