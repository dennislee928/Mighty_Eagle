# Cryptography & commitments

## Goals

- Local confidentiality (encrypted vault)
- Public verifiability of "I had these records by time T"
- No leakage of PII via commitments

## Hashing

- Use **SHA-256** for:
  - event leaf hash
  - Merkle root

SHA-256 output is 32 bytes → fits EVM `bytes32`.

## Leaf commitment construction (v1)

We hash only **non-PII metadata + a salt + payloadHash**.

1. `payloadHash = SHA256( canonicalJSON(payload) )`
2. `leaf = SHA256( "TT|leaf|v1" || day || type || createdAt || salt32 || payloadHash )`

Where:
- `salt32` is cryptographically random 32 bytes per event
- `canonicalJSON` uses stable key ordering and UTF-8 encoding

Rationale:
- Even if someone guesses the payload, they still need `salt32` to match the leaf.

## Merkle tree (v1)

- Build a binary Merkle tree over `leaf[]`
- If odd number of leaves, duplicate the last leaf (common approach)
- Parent = SHA256(left || right)
- Root = bytes32

## Local encryption

- Use `AES-GCM 256` with a per-record random 12-byte IV.
- Key derivation: `PBKDF2(SHA-256)` from a user passphrase.
  - `salt`: 16 bytes random (stored in settings)
  - `iterations`: 310,000 (tune for UX; target ~200ms on desktop)
- Store:
  - `kdfParams` (non-secret)
  - `ciphertext` + `iv` + `tag` (tag included in WebCrypto output)

## Export proofs

`commitment.json`:
- `day`
- `schemaVersion`
- `merkleRoot` (0x...)
- `leafCount`
- `createdAt`
- `vaultId` (random id; not user identity)
- `metaHash` (optional; hash of non-sensitive metadata)
