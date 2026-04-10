#!/usr/bin/env node

const fs = require("node:fs");
const fsp = require("node:fs/promises");
const https = require("node:https");
const path = require("node:path");
const { spawnSync } = require("node:child_process");

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

function getDownloadUrl(archiveName) {
  const customUrl = process.env.QMAI_CLI_DOWNLOAD_URL;
  if (customUrl) {
    return customUrl;
  }

  const baseUrl =
    process.env.QMAI_CLI_RELEASE_BASE_URL ||
    "https://github.com/feixiao629/qmai-cli-public/releases/download";
  return `${baseUrl}/v${packageJson.version}/${archiveName}`;
}

async function downloadFile(url, outputPath) {
  await new Promise((resolve, reject) => {
    const file = fs.createWriteStream(outputPath);
    const request = https.get(
      url,
      {
        headers: {
          "User-Agent": "qmai-cli-public-installer",
          Accept: "application/octet-stream",
        },
      },
      (response) => {
        if (
          response.statusCode &&
          response.statusCode >= 300 &&
          response.statusCode < 400 &&
          response.headers.location
        ) {
          file.close();
          fs.unlinkSync(outputPath);
          downloadFile(response.headers.location, outputPath).then(resolve).catch(reject);
          return;
        }

        if (response.statusCode !== 200) {
          file.close();
          fs.unlinkSync(outputPath);
          reject(
            new Error(`download failed: ${response.statusCode} ${response.statusMessage || ""}`.trim())
          );
          return;
        }

        response.pipe(file);
        file.on("finish", () => file.close(resolve));
      }
    );

    request.on("error", (error) => {
      file.close();
      if (fs.existsSync(outputPath)) {
        fs.unlinkSync(outputPath);
      }
      reject(error);
    });
  });
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
    const url = getDownloadUrl(archiveName);
    console.log(`下载 qmai 二进制：${url}`);
    await downloadFile(url, archivePath);
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
