#!/usr/bin/env node

const fs = require("node:fs");
const fsp = require("node:fs/promises");
const http = require("node:http");
const https = require("node:https");
const path = require("node:path");
const { spawnSync } = require("node:child_process");
const { URL } = require("node:url");

const packageRoot = path.resolve(__dirname, "..");
const distDir = path.join(packageRoot, "npm", "dist");
const packageJson = require(path.join(packageRoot, "package.json"));

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const archMap = {
  x64: "amd64",
  arm64: "arm64",
};

const DOWNLOAD_TIMEOUT_MS = parseInt(process.env.QMAI_CLI_DOWNLOAD_TIMEOUT_MS || "30000", 10);
const DOWNLOAD_RETRIES = parseInt(process.env.QMAI_CLI_DOWNLOAD_RETRIES || "3", 10);
const MAX_REDIRECTS = parseInt(process.env.QMAI_CLI_DOWNLOAD_MAX_REDIRECTS || "5", 10);

function fail(message) {
  console.error(message);
  process.exit(1);
}

function runCommand(command, args) {
  const result = spawnSync(command, args, { stdio: "inherit" });
  if (result.error || result.status !== 0) {
    fail(`qmai-cli-public 安装失败：执行 ${command} ${args.join(" ")} 失败`);
  }
}

function getPlatform() {
  const mapped = platformMap[process.platform];
  if (!mapped) {
    fail(`qmai-cli-public 当前不支持平台 ${process.platform}`);
  }
  return mapped;
}

function getArch() {
  const mapped = archMap[process.arch];
  if (!mapped) {
    fail(`qmai-cli-public 当前不支持架构 ${process.arch}`);
  }
  return mapped;
}

function getAssetInfo() {
  const version = packageJson.version;
  const platform = getPlatform();
  const arch = getArch();
  const baseName = `qmai_${version}_${platform}_${arch}`;
  const archiveName = platform === "windows" ? `${baseName}.zip` : `${baseName}.tar.gz`;
  const binaryName = platform === "windows" ? "qmai.exe" : "qmai";
  return { version, platform, arch, baseName, archiveName, binaryName };
}

function parseList(value) {
  return (value || "")
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean);
}

function getDownloadCandidates(archiveName) {
  const candidates = [];
  const seen = new Set();

  const add = (url) => {
    if (!url || seen.has(url)) {
      return;
    }
    seen.add(url);
    candidates.push(url);
  };

  add(process.env.QMAI_CLI_DOWNLOAD_URL);
  parseList(process.env.QMAI_CLI_FALLBACK_DOWNLOAD_URLS).forEach(add);

  const baseUrls = [
    process.env.QMAI_CLI_RELEASE_BASE_URL || "https://github.com/feixiao629/qmai-cli-public/releases/download",
    ...parseList(process.env.QMAI_CLI_FALLBACK_RELEASE_BASE_URLS),
  ];
  for (const baseUrl of baseUrls) {
    const normalized = baseUrl.replace(/\/+$/, "");
    add(`${normalized}/v${packageJson.version}/${archiveName}`);
  }

  return candidates;
}

function getHttpClient(url) {
  const protocol = new URL(url).protocol;
  if (protocol === "https:") {
    return https;
  }
  if (protocol === "http:") {
    return http;
  }
  throw new Error(`unsupported protocol: ${protocol}`);
}

async function removeIfExists(filePath) {
  await fsp.rm(filePath, { force: true });
}

async function downloadFileOnce(url, outputPath, redirectCount = 0) {
  if (redirectCount > MAX_REDIRECTS) {
    throw new Error(`redirect limit exceeded (${MAX_REDIRECTS})`);
  }

  await new Promise((resolve, reject) => {
    const file = fs.createWriteStream(outputPath);
    const client = getHttpClient(url);

    const cleanup = (callback) => {
      file.destroy();
      removeIfExists(outputPath).then(() => callback()).catch(callback);
    };

    const request = client.get(
      url,
      {
        headers: {
          "User-Agent": "qmai-cli-public-installer",
          Accept: "application/octet-stream",
        },
      },
      (response) => {
        const { statusCode = 0, statusMessage = "", headers } = response;

        if (statusCode >= 300 && statusCode < 400 && headers.location) {
          const nextUrl = new URL(headers.location, url).toString();
          response.resume();
          file.close(() => {
            removeIfExists(outputPath)
              .then(() => downloadFileOnce(nextUrl, outputPath, redirectCount + 1))
              .then(resolve)
              .catch(reject);
          });
          return;
        }

        if (statusCode !== 200) {
          response.resume();
          cleanup(() =>
            reject(new Error(`download failed: ${statusCode} ${statusMessage}`.trim()))
          );
          return;
        }

        response.pipe(file);
        file.on("finish", () => file.close(resolve));
        response.on("error", (error) => {
          cleanup(() => reject(error));
        });
      }
    );

    request.setTimeout(DOWNLOAD_TIMEOUT_MS, () => {
      request.destroy(new Error(`download timeout after ${DOWNLOAD_TIMEOUT_MS}ms`));
    });

    request.on("error", (error) => {
      cleanup(() => reject(error));
    });

    file.on("error", (error) => {
      request.destroy(error);
    });
  });
}

async function downloadFile(urls, outputPath) {
  const attempts = [];

  for (const url of urls) {
    for (let attempt = 1; attempt <= DOWNLOAD_RETRIES; attempt += 1) {
      try {
        console.log(`下载 qmai 二进制（${attempt}/${DOWNLOAD_RETRIES}）：${url}`);
        await downloadFileOnce(url, outputPath);
        return;
      } catch (error) {
        attempts.push(`${url} [${attempt}/${DOWNLOAD_RETRIES}]: ${error.message}`);
        console.warn(`下载失败：${error.message}`);
        await removeIfExists(outputPath);
      }
    }
  }

  throw new Error(
    [
      "无法下载 qmai 二进制。",
      `已尝试 ${urls.length} 个下载源，每个最多重试 ${DOWNLOAD_RETRIES} 次。`,
      `超时设置：${DOWNLOAD_TIMEOUT_MS}ms。`,
      "可设置以下环境变量重试：",
      "- QMAI_CLI_DOWNLOAD_TIMEOUT_MS=60000",
      "- QMAI_CLI_DOWNLOAD_RETRIES=5",
      "- QMAI_CLI_RELEASE_BASE_URL=https://<your-mirror>/releases/download",
      "- QMAI_CLI_FALLBACK_RELEASE_BASE_URLS=https://mirror-a/releases/download,https://mirror-b/releases/download",
      "- QMAI_CLI_DOWNLOAD_URL=https://.../qmai_<version>_<os>_<arch>.tar.gz",
      "- QMAI_CLI_LOCAL_ASSET_DIR=/path/to/pre-downloaded-assets",
      "失败明细：",
      ...attempts.map((item) => `  - ${item}`),
    ].join("\n")
  );
}

async function extractArchive(archivePath, targetDir, platform) {
  if (platform === "windows") {
    runCommand("powershell.exe", [
      "-NoProfile",
      "-Command",
      `Expand-Archive -LiteralPath '${archivePath}' -DestinationPath '${targetDir}' -Force`,
    ]);
    return;
  }

  runCommand("tar", ["-xzf", archivePath, "-C", targetDir]);
}

async function copyLocalAsset(localAssetDir, archiveName, outputPath) {
  const sourcePath = path.join(localAssetDir, archiveName);
  await fsp.copyFile(sourcePath, outputPath);
}

async function main() {
  const { archiveName, baseName, binaryName, platform } = getAssetInfo();
  const archivePath = path.join(distDir, archiveName);
  const extractedDir = path.join(distDir, baseName);
  const finalBinaryPath = path.join(distDir, binaryName);
  const shimPath = path.join(packageRoot, "npm", "bin", "qmai.js");

  await fsp.mkdir(distDir, { recursive: true });
  await fsp.rm(archivePath, { force: true });
  await fsp.rm(extractedDir, { recursive: true, force: true });
  await fsp.rm(finalBinaryPath, { force: true });

  const localAssetDir = process.env.QMAI_CLI_LOCAL_ASSET_DIR;
  if (localAssetDir) {
    await copyLocalAsset(localAssetDir, archiveName, archivePath);
  } else {
    const urls = getDownloadCandidates(archiveName);
    await downloadFile(urls, archivePath);
  }

  await extractArchive(archivePath, distDir, platform);

  const extractedBinaryPath = path.join(extractedDir, binaryName);
  if (!fs.existsSync(extractedBinaryPath)) {
    fail(`qmai-cli-public 安装失败：未在压缩包中找到 ${binaryName}`);
  }

  await fsp.copyFile(extractedBinaryPath, finalBinaryPath);
  if (platform !== "windows") {
    await fsp.chmod(finalBinaryPath, 0o755);
    await fsp.chmod(shimPath, 0o755);
  }

  await fsp.rm(archivePath, { force: true });
  await fsp.rm(extractedDir, { recursive: true, force: true });

  console.log(`qmai 已安装到 ${finalBinaryPath}`);
}

main().catch((error) => {
  fail(`qmai-cli-public 安装失败：${error.message}`);
});
