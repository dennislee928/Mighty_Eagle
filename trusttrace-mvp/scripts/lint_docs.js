/**
 * Minimal docs lint placeholder. Replace with markdownlint if desired.
 */
const fs = require("fs");
const path = require("path");

function walk(dir, files = []) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, entry.name);
    if (entry.isDirectory()) walk(p, files);
    else if (p.endsWith(".md")) files.push(p);
  }
  return files;
}

const docsDir = path.join(__dirname, "..", "docs");
if (!fs.existsSync(docsDir)) process.exit(0);

const files = walk(docsDir);
let ok = true;

for (const f of files) {
  const content = fs.readFileSync(f, "utf8");
  if (content.includes("\t")) {
    console.error(`[lint] Tabs found in ${f}`);
    ok = false;
  }
}
process.exit(ok ? 0 : 1);
