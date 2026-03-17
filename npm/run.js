#!/usr/bin/env node
// @ts-check

const { execFileSync } = require("child_process");
const path = require("path");

const binName =
  process.platform === "win32" ? "chatwoot-mcp.exe" : "chatwoot-mcp";
const binPath = path.join(__dirname, "bin", binName);

try {
  execFileSync(binPath, process.argv.slice(2), {
    stdio: "inherit",
    env: process.env,
  });
} catch (err) {
  if (err.status !== null) {
    process.exit(err.status);
  }
  throw err;
}
