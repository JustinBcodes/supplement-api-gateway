import { chromium } from "playwright";
import { mkdir } from "node:fs/promises";

const browser = await chromium.launch({
  headless: true,
  executablePath:
    process.env.CHROME_PATH ||
    "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
});
const page = await browser.newPage({
  viewport: { width: 1440, height: 900 },
  deviceScaleFactor: 1,
});
await mkdir("../screenshots", { recursive: true });
await page.goto(process.env.APP_URL || "http://localhost:4173", {
  waitUntil: "networkidle",
});
await page.getByRole("heading", { name: /Overview/ }).waitFor();
await page.screenshot({
  path: "../screenshots/control-overview.png",
  fullPage: true,
});
await page.getByRole("button", { name: "Traffic analytics" }).click();
await page.screenshot({
  path: "../screenshots/traffic-analytics.png",
  fullPage: true,
});
await page.getByRole("button", { name: "Reliability" }).click();
await page.screenshot({
  path: "../screenshots/reliability.png",
  fullPage: true,
});
await page.getByRole("button", { name: "Routes & policies" }).click();
await page.screenshot({ path: "../screenshots/routes.png", fullPage: true });
await browser.close();
