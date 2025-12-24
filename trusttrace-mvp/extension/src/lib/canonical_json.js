/**
 * Canonical JSON: stable key ordering for deterministic hashing.
 * - Arrays keep order.
 * - Objects keys sorted lexicographically.
 */
export function canonicalize(value) {
  if (value === null || typeof value !== "object") return value;
  if (Array.isArray(value)) return value.map(canonicalize);

  const out = {};
  for (const k of Object.keys(value).sort()) {
    out[k] = canonicalize(value[k]);
  }
  return out;
}

export function canonicalJSONString(value) {
  return JSON.stringify(canonicalize(value));
}
