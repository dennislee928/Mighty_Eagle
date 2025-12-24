import { getLocal, setLocal } from "./storage.js";
import { randomBytes, deriveKeyPBKDF2, encryptJSON, decryptJSON, leafHashV1, merkleRootV1, hex } from "./crypto.js";

const SETTINGS_KEY = "tt_settings_v1";
const EVENTS_KEY = "tt_events_v1";
const BUNDLES_KEY = "tt_bundles_v1";

/**
 * Settings schema:
 * {
 *   kdf: { salt16Hex, iterations },
 *   vaultId: string
 * }
 */
export async function getSettings() {
  return (await getLocal(SETTINGS_KEY)) ?? null;
}

export async function setSettings(s) {
  await setLocal(SETTINGS_KEY, s);
}

export async function setPassphrase(passphrase) {
  let s = await getSettings();
  if (!s) {
    const salt16 = randomBytes(16);
    s = { kdf: { salt16Hex: hex(salt16), iterations: 310000 }, vaultId: crypto.randomUUID() };
  }
  // Store a verifier to detect wrong passphrase: encrypt a known value.
  const key = await getVaultKey(passphrase, s);
  const verifier = await encryptJSON(key, { v: "tt_verifier_v1", t: Date.now() });
  s.verifier = verifier;
  await setSettings(s);
}

export async function getVaultKey(passphrase, settings = null) {
  const s = settings ?? await getSettings();
  if (!s?.kdf?.salt16Hex) throw new Error("Passphrase not set (open Options).");
  const salt = fromHex(s.kdf.salt16Hex);
  return await deriveKeyPBKDF2(passphrase, salt, s.kdf.iterations ?? 310000);
}

function fromHex(h) {
  const s = h.startsWith("0x") ? h.slice(2) : h;
  const bytes = new Uint8Array(s.length / 2);
  for (let i = 0; i < bytes.length; i++) bytes[i] = parseInt(s.slice(i * 2, i * 2 + 2), 16);
  return bytes;
}

async function loadEvents() {
  return (await getLocal(EVENTS_KEY)) ?? [];
}
async function saveEvents(events) {
  await setLocal(EVENTS_KEY, events);
}

async function loadBundles() {
  return (await getLocal(BUNDLES_KEY)) ?? {};
}
async function saveBundles(b) {
  await setLocal(BUNDLES_KEY, b);
}

export async function addEncryptedEvent(passphrase, { type, payload }) {
  const settings = await getSettings();
  const key = await getVaultKey(passphrase, settings);

  const now = new Date();
  const createdAt = now.toISOString();
  const day = createdAt.slice(0, 10);
  const salt32 = randomBytes(32);

  const encPayload = await encryptJSON(key, { payload, salt32: Array.from(salt32), createdAt, type, day });

  const events = await loadEvents();
  events.push({
    id: crypto.randomUUID(),
    type,
    createdAt,
    day,
    payloadSchema: "tt.event.v1",
    enc: encPayload,
  });
  await saveEvents(events);
}

export async function listEventsForDay(day) {
  const events = await loadEvents();
  return events.filter(e => e.day === day).slice(-20).reverse();
}

export async function decryptEventPreview(passphrase, encEvent) {
  const settings = await getSettings();
  const key = await getVaultKey(passphrase, settings);
  const dec = await decryptJSON(key, encEvent.enc);

  const payload = dec.payload;
  const preview = (payload?.text ?? payload?.terms ?? JSON.stringify(payload)).slice(0, 80);
  return { type: encEvent.type, createdAt: encEvent.createdAt, preview };
}

export async function buildDailyBundle(day, passphrase) {
  const settings = await getSettings();
  const key = await getVaultKey(passphrase, settings);
  const events = await listEventsForDay(day);

  const leaves = [];
  for (const e of events) {
    const dec = await decryptJSON(key, e.enc);
    const salt32 = new Uint8Array(dec.salt32);
    const leaf = await leafHashV1({
      day,
      type: e.type,
      createdAt: e.createdAt,
      salt32,
      payload: dec.payload
    });
    leaves.push(leaf);
  }

  const rootBytes = await merkleRootV1(leaves);
  const bundle = {
    day,
    schemaVersion: "tt.bundle.v1",
    leafCount: leaves.length,
    merkleRoot: hex(rootBytes),
    createdAt: new Date().toISOString(),
    vaultId: settings?.vaultId ?? "unknown",
    // leaf hashes can be included without payload
    leafHashes: leaves.map(l => hex(l)),
  };

  const bundles = await loadBundles();
  bundles[day] = bundle;
  await saveBundles(bundles);
  return bundle;
}

export async function getBundle(day) {
  const bundles = await loadBundles();
  return bundles[day] ?? null;
}
