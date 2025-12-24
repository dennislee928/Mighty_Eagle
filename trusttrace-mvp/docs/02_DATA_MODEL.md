# Local data model

All persisted data is encrypted at rest.

## Storage choices

- **IndexedDB** for event records and bundles (recommended)
- `chrome.storage.local` for small settings (e.g., encryption params)

MVP starter uses `chrome.storage.local` for simplicity, with an upgrade path to IndexedDB.

## Entities

### EventRecord (encrypted payload)

- `id`: uuid v4
- `createdAt`: ISO8601
- `day`: YYYY-MM-DD (local timezone)
- `type`: enum
  - `NOTE`
  - `CONSENT`
  - `ACTION_LOG`  (ban-protection metadata)
- `payload`: encrypted blob (AES-GCM)
- `payloadSchema`: e.g., `tt.event.v1`

### LeafCommitment (derived, not stored necessarily)

- `leafHash`: 32 bytes
- `eventId`
- `day`

### DailyBundle

- `day`
- `schemaVersion`: `tt.bundle.v1`
- `leafHashes`: array[bytes32] (computed)
- `merkleRoot`: bytes32
- `createdAt`
- `status`: `DRAFT | EXPORTED | ANCHORED`
- `anchorTx`: optional

## Event payload schemas (v1)

### NOTE payload
- `title` (string, optional)
- `text` (string)
- `tags` (string[])
- `riskFlags` (string[], optional; user-chosen)

### CONSENT payload
- `consentType`: `CONTACT_EXCHANGE | MEETUP | OTHER`
- `terms`: string (user-entered; avoid identity)
- `counterpartyProof`: optional (future)
- `revocable`: boolean

### ACTION_LOG payload (metadata only)
- `action`: `LOGIN | SWIPE | MESSAGE_SENT | REPORT | BLOCK | OTHER`
- `count`: integer (aggregated count, not per-target)
- `appContext`: `TINDER | OTHER`
- `notes`: optional
