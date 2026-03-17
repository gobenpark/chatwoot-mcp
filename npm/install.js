#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");

const VERSION = require("./package.json").version;
const REPO = "gobenpark/chatwoot-mcp";

function getPlatform() {
  const osMap = { darwin: "darwin", linux: "linux", win32: "windows" };
  const archMap = { x64: "amd64", arm64: "arm64" };

  const os = osMap[process.platform];
  const cpu = archMap[process.arch];

  if (!os || !cpu) {
    throw new Error(`Unsupported platform: ${process.platform}/${process.arch}`);
  }
  return { os, cpu };
}

function download(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (res) => {
      if (res.statusCode === 302 || res.statusCode === 301) {
        return download(res.headers.location).then(resolve).catch(reject);
      }
      if (res.statusCode !== 200) {
        return reject(new Error(`HTTP ${res.statusCode} from ${url}`));
      }
      const chunks = [];
      res.on("data", (chunk) => chunks.push(chunk));
      res.on("end", () => resolve(Buffer.concat(chunks)));
      res.on("error", reject);
    }).on("error", reject);
  });
}

async function main() {
  const { os, cpu } = getPlatform();
  const binDir = path.join(__dirname, "bin");
  const binName = os === "windows" ? "chatwoot-mcp.exe" : "chatwoot-mcp";
  const binPath = path.join(binDir, binName);

  if (fs.existsSync(binPath)) {
    return;
  }

  fs.mkdirSync(binDir, { recursive: true });

  const ext = os === "windows" ? "zip" : "tar.gz";
  const assetName = `chatwoot-mcp_${os}_${cpu}.${ext}`;
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${assetName}`;

  console.log(`Downloading chatwoot-mcp v${VERSION} for ${os}/${cpu}...`);

  const data = await download(url);
  const tmpFile = path.join(binDir, assetName);
  fs.writeFileSync(tmpFile, data);

  if (ext === "tar.gz") {
    execSync(`tar xzf "${tmpFile}" -C "${binDir}"`);
  } else {
    execSync(`powershell -command "Expand-Archive -Path '${tmpFile}' -DestinationPath '${binDir}'"`);
  }

  fs.unlinkSync(tmpFile);

  if (os !== "windows") {
    fs.chmodSync(binPath, 0o755);
  }

  console.log(`chatwoot-mcp v${VERSION} installed successfully!`);
}

main().catch((err) => {
  console.error(`Failed to install chatwoot-mcp: ${err.message}`);
  process.exit(1);
});
