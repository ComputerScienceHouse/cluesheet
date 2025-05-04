import { mockCluesheetUUID } from "./lib/data";
import { test, expect } from "./mock/test";

const unitTestTimeout = 10000;

test("happy view cluesheet", async ({ page }) => {
  test.setTimeout(unitTestTimeout);
  await page.goto(`/cluesheet/${mockCluesheetUUID}`);
  await page.waitForTimeout(5000);
});
