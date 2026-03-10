import { createWriteStream, existsSync, chmodSync, unlinkSync } from "node:fs";
import { mkdir } from "node:fs/promises";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { get } from "node:https";
import AdmZip from "adm-zip";
import tar from "tar";

const rootDir = dirname(fileURLToPath(import.meta.url));
const packageJson = await import("./package.json", { with: { type: "json" } });
const version = packageJson.default.version;
const binDir = join(rootDir, "bin");
const binaryName = process.platform === "win32" ? "codex-theme.exe" : "codex-theme";
const binaryPath = join(binDir, binaryName);

if (existsSync(binaryPath)) {
  process.exit(0);
}

await mkdir(binDir, { recursive: true });

const asset = releaseAssetName(process.platform, process.arch);
if (!asset) {
  console.warn(`[codex-theme] No prebuilt release is available for ${process.platform}/${process.arch}.`);
  process.exit(0);
}

const url = `https://github.com/ychampion/codex-themes/releases/download/v${version}/${asset.archive}`;
const archivePath = join(binDir, asset.archive);

await download(url, archivePath);

if (asset.archive.endsWith(".gz")) {
  await extractSingleTarGz(archivePath, binaryPath);
} else if (asset.archive.endsWith(".zip")) {
  await extractZip(archivePath, binDir);
}

if (existsSync(binaryPath) && process.platform !== "win32") {
  chmodSync(binaryPath, 0o755);
}

if (existsSync(archivePath)) {
  unlinkSync(archivePath);
}

function releaseAssetName(platform, arch) {
  const archMap = { x64: "x64", arm64: "arm64" };
  const mappedArch = archMap[arch];
  if (!mappedArch) return null;
  if (platform === "linux" || platform === "darwin") {
    return { archive: `codex-theme_${platform}_${mappedArch}.tar.gz` };
  }
  if (platform === "win32") {
    return { archive: `codex-theme_windows_${mappedArch}.zip` };
  }
  return null;
}

function download(url, destination) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(destination);

    const request = (currentUrl) => {
      get(currentUrl, (response) => {
        if (response.statusCode && response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
          request(response.headers.location);
          return;
        }

        if (response.statusCode && response.statusCode >= 400) {
          reject(new Error(`download failed: ${response.statusCode} ${currentUrl}`));
          return;
        }

        response.pipe(file);
        file.on("finish", () => file.close(resolve));
      }).on("error", reject);
    };

    request(url);
  });
}

async function extractSingleTarGz(archivePath, destinationPath) {
  await tar.x({ file: archivePath, cwd: dirname(destinationPath) });
}

async function extractZip(archivePath, destinationDir) {
  const zip = new AdmZip(archivePath);
  zip.extractAllTo(destinationDir, true);
}

