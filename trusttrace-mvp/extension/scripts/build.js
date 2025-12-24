/**
 * MV3 build script using esbuild (dependency-light).
 * Bundles src/ into dist/.
 */
import * as esbuild from "esbuild";
import fs from "fs";
import path from "path";
import url from "url";

const __dirname = path.dirname(url.fileURLToPath(import.meta.url));
const root = path.join(__dirname, "..");
const src = path.join(root, "src");
const dist = path.join(root, "dist");

const watch = process.argv.includes("--watch");

if (!fs.existsSync(dist)) fs.mkdirSync(dist, { recursive: true });

// Copy static files
for (const file of ["manifest.json", "popup.html", "options.html"]) {
  fs.copyFileSync(path.join(src, file), path.join(dist, file));
}

const common = {
  bundle: true,
  format: "esm",
  target: ["chrome114"],
  sourcemap: true,
  outdir: dist,
};

const entryPoints = [
  path.join(src, "background.js"),
  path.join(src, "popup.js"),
  path.join(src, "options.js"),
  path.join(src, "content.js"),
];

if (watch) {
  const ctx = await esbuild.context({ ...common, entryPoints });
  await ctx.watch();
  console.log("Watching...");
} else {
  await esbuild.build({ ...common, entryPoints });
  console.log("Built to dist/");
}
