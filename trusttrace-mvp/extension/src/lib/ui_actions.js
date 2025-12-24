import { addEncryptedEvent, decryptEventPreview, buildDailyBundle as buildBundle } from "./vault.js";
import { getLocal, setLocal } from "./storage.js";

const UI_CACHE_KEY = "tt_ui_cache_v1";

async function promptPassphrase() {
  // MVP: prompt each session. Improve with session memory + auto-lock.
  const cached = await getLocal(UI_CACHE_KEY);
  if (cached?.passphrase) return cached.passphrase;
  const pass = prompt("Enter TrustTrace passphrase (stored only for this session).");
  if (!pass) throw new Error("passphrase required");
  await setLocal(UI_CACHE_KEY, { passphrase: pass });
  return pass;
}

export async function addEvent({ type, text }) {
  const pass = await promptPassphrase();

  // Minimal PII heuristic: warn on phone/email patterns.
  const pii = /(@\w+\.\w+)|(\+?\d[\d\s\-]{8,})/;
  if (pii.test(text)) {
    const ok = confirm("Warning: text looks like it may contain personal identifiers. Continue?");
    if (!ok) return;
  }

  const day = new Date().toISOString().slice(0, 10);
  let payload;
  if (type === "NOTE") payload = { text, tags: [], day };
  else if (type === "CONSENT") payload = { consentType: "OTHER", terms: text, revocable: true, day };
  else payload = { action: "OTHER", count: 1, appContext: "TINDER", notes: text, day };

  await addEncryptedEvent(pass, { type, payload });
}

export async function decryptPreviewEvents() {
  const pass = await promptPassphrase();
  const today = new Date().toISOString().slice(0, 10);
  const resp = await chrome.runtime.sendMessage({ type: "TT_LIST_EVENTS_TODAY" });
  if (!resp.ok) return [];

  const previews = [];
  for (const e of resp.data) {
    previews.push(await decryptEventPreview(pass, e));
  }
  return previews;
}

export async function buildDailyBundle() {
  const pass = await promptPassphrase();
  const today = new Date().toISOString().slice(0, 10);
  return await buildBundle(today, pass);
}

export async function exportCommitmentJson(bundle) {
  const commitment = {
    day: bundle.day,
    dayInt: parseInt(bundle.day.replaceAll("-", ""), 10),
    schemaVersion: bundle.schemaVersion,
    merkleRoot: bundle.merkleRoot,
    leafCount: bundle.leafCount,
    createdAt: bundle.createdAt,
    vaultId: bundle.vaultId,
    metaHash: "0x" + "00".repeat(32),
  };

  const blob = new Blob([JSON.stringify(commitment, null, 2)], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  await chrome.downloads?.download?.({ url, filename: `trusttrace_commitment_${bundle.day}.json`, saveAs: true })
    .catch(() => {
      // fallback: open in new tab
      window.open(url);
    });
}

export async function exportAppealZip(bundle) {
  // Dependency-light MVP: create a pseudo-zip via JSON-only export instructions.
  // For a real zip, integrate a zip library (e.g., fflate) in Sprint 3.
  const summary = `# Appeal Summary\n\nDay: ${bundle.day}\nRoot: ${bundle.merkleRoot}\nLeafCount: ${bundle.leafCount}\n\nThis package contains no third-party personal data.\n`;
  const env = {
    userAgent: navigator.userAgent,
    createdAt: new Date().toISOString(),
    extension: { name: "TrustTrace MVP", version: "0.1.0" }
  };

  const pkg = {
    appeal_summary_md: summary,
    commitment: {
      day: bundle.day,
      dayInt: parseInt(bundle.day.replaceAll("-", ""), 10),
      schemaVersion: bundle.schemaVersion,
      merkleRoot: bundle.merkleRoot,
      leafCount: bundle.leafCount,
      createdAt: bundle.createdAt,
      vaultId: bundle.vaultId,
    },
    bundle_day: { day: bundle.day, leafHashes: bundle.leafHashes, merkleRoot: bundle.merkleRoot },
    environment: env
  };

  const blob = new Blob([JSON.stringify(pkg, null, 2)], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  await chrome.downloads?.download?.({ url, filename: `trusttrace_appeal_package_${bundle.day}.json`, saveAs: true })
    .catch(() => window.open(url));
}
