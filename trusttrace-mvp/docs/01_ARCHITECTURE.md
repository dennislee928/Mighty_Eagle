# Architecture

## High-level components

1. **Chrome Extension (MV3)**
   - Local encrypted vault
   - Event capture (user-entered)
   - Daily rollup (Merkle root)
   - Proof export (Appeal Package)

2. **On-chain: DailyRootRegistry**
   - Stores `(signer, day, root, schemaVersion, metaHash)`
   - Emits events for auditability
   - Does not store any payload

3. **Web dApp (wallet-enabled)**
   - Used to submit daily roots
   - Reads commitment from copy/paste or file import
   - Uses injected wallet (MetaMask) on a normal web origin

## Data flow (MVP)

1. User creates events in extension (notes / consent / actions).
2. Extension encrypts and stores events locally.
3. Extension builds daily leaf hashes and Merkle root.
4. Extension exports:
   - `commitment.json` (root + day + schema)
   - `appeal-package.zip` (optional)
5. User submits the root via the dApp to the on-chain registry.

## Trust boundaries

- The extension must be safe even if the platform webpage is malicious.
- The chain is public; only commitments should be published.
- The dApp is only used for transaction submission.

## Availability assumptions

- No backend required for MVP.
- Optional future backend for B2B risk APIs (not in MVP).
