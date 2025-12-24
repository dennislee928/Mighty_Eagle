import { canonicalJSONString } from "./canonical_json.js";

export function utf8(s) {
  return new TextEncoder().encode(s);
}

export function hex(buf) {
  const b = new Uint8Array(buf);
  return "0x" + Array.from(b).map(x => x.toString(16).padStart(2, "0")).join("");
}

export async function sha256Bytes(data) {
  const buf = (data instanceof ArrayBuffer) ? data : data.buffer ?? data;
  return await crypto.subtle.digest("SHA-256", buf);
}

export async function sha256Hex(data) {
  return hex(await sha256Bytes(data));
}

export function concatBytes(...arrs) {
  const total = arrs.reduce((s, a) => s + a.length, 0);
  const out = new Uint8Array(total);
  let off = 0;
  for (const a of arrs) {
    out.set(a, off);
    off += a.length;
  }
  return out;
}

export function randomBytes(n) {
  const b = new Uint8Array(n);
  crypto.getRandomValues(b);
  return b;
}

export async function payloadHash(payload) {
  const s = canonicalJSONString(payload);
  return new Uint8Array(await sha256Bytes(utf8(s)));
}

export async function leafHashV1({ day, type, createdAt, salt32, payload }) {
  const pHash = await payloadHash(payload);
  const prefix = utf8("TT|leaf|v1|");
  const bytes = concatBytes(
    prefix,
    utf8(day), utf8("|"),
    utf8(type), utf8("|"),
    utf8(createdAt), utf8("|"),
    salt32, utf8("|"),
    pHash
  );
  const h = await sha256Bytes(bytes);
  return new Uint8Array(h);
}

export async function merkleRootV1(leaves /* Uint8Array[] */) {
  if (!leaves.length) {
    return new Uint8Array(await sha256Bytes(utf8("TT|empty|v1")));
  }
  let layer = leaves.map(l => new Uint8Array(l));
  while (layer.length > 1) {
    const next = [];
    for (let i = 0; i < layer.length; i += 2) {
      const left = layer[i];
      const right = layer[i + 1] ?? layer[i]; // duplicate last if odd
      const parent = new Uint8Array(await sha256Bytes(concatBytes(left, right)));
      next.push(parent);
    }
    layer = next;
  }
  return layer[0];
}

// --- Encryption (AES-GCM) ---

export async function deriveKeyPBKDF2(passphrase, salt16, iterations = 310000) {
  const baseKey = await crypto.subtle.importKey(
    "raw",
    utf8(passphrase),
    "PBKDF2",
    false,
    ["deriveKey"]
  );
  return await crypto.subtle.deriveKey(
    { name: "PBKDF2", hash: "SHA-256", salt: salt16, iterations },
    baseKey,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"]
  );
}

export async function encryptJSON(key, obj) {
  const iv = randomBytes(12);
  const plaintext = utf8(canonicalJSONString(obj));
  const ciphertext = await crypto.subtle.encrypt({ name: "AES-GCM", iv }, key, plaintext);
  return { iv: hex(iv), ciphertext: hex(ciphertext) };
}

export async function decryptJSON(key, enc) {
  const iv = fromHex(enc.iv);
  const ct = fromHex(enc.ciphertext);
  const pt = await crypto.subtle.decrypt({ name: "AES-GCM", iv }, key, ct);
  const s = new TextDecoder().decode(pt);
  return JSON.parse(s);
}

export function fromHex(h) {
  const s = h.startsWith("0x") ? h.slice(2) : h;
  const bytes = new Uint8Array(s.length / 2);
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = parseInt(s.slice(i * 2, i * 2 + 2), 16);
  }
  return bytes;
}
