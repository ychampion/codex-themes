#!/usr/bin/env node
import { spawn } from "node:child_process";
import { existsSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const binary = process.platform === "win32" ? join(here, "codex-theme.exe") : join(here, "codex-theme");

if (!existsSync(binary)) {
  console.error(`[codex-theme] Missing native binary at ${binary}. Try reinstalling the package.`);
  process.exit(1);
}

const child = spawn(binary, process.argv.slice(2), { stdio: "inherit" });
child.on("exit", (code) => process.exit(code ?? 1));
child.on("error", (error) => {
  console.error(`[codex-theme] Failed to launch ${binary}: ${error.message}`);
  process.exit(1);
});
