# Extension specification (MV3)

## Purpose
Provide a local encrypted "vault" for user-entered dating notes and compliance logs, then produce a daily commitment.

## Permissions (MVP)
- `storage`
- `alarms`
- `activeTab` (optional; for UX context only)

Host permissions (optional, minimal):
- `https://tinder.com/*` (only to detect current site and show contextual UI; **no extraction**)

## UI surfaces
- Popup:
  - Add note
  - Add consent event
  - Add action log (aggregated)
  - Build daily rollup
  - Export commitment.json
  - Export Appeal Package (zip)
- Options:
  - Set / change passphrase
  - KDF parameters (advanced)
  - Vault export/import (future)

## Background service worker
- Handles alarms
- Maintains lightweight state
- Provides message API for popup

## Data handling
- All payloads are encrypted.
- Leaf salts are stored encrypted alongside payload.
- "Daily bundle" contains only hashes and can be stored unencrypted.

## Anti-abuse guardrails
- No auto-capture of profile details.
- No messaging automation.
- No "ban evasion" instructions or features.

## Acceptance criteria (MVP)
1. User can set passphrase and create encrypted vault.
2. User can add events and retrieve them (decrypt in popup).
3. User can generate daily Merkle root and export commitment.
4. User can export an appeal package including readable summary + root.
