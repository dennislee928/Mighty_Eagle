# Product overview

## Problem statements

1. **Trust deficit in dating**: users face scams, fake personas, and unverifiable claims.
2. **Account enforcement opacity**: users report bans without clear explanation; they lack structured evidence for appeals.
3. **Privacy and compliance**: any solution must avoid collecting/processing third-party personal data without consent.

## Solution concept

Use the blockchain as a **notary / timestamp anchor**:
- Store **only commitments** (hashes / Merkle roots) on-chain.
- Store sensitive content **locally encrypted** (default) or off-chain with deletion controls.

## Non-negotiable constraints

- **No on-chain PII** (direct or indirect).
- Default: **no automated extraction** of platform content; user-entered events only.
- No features intended to **evade** platform detection or bans.

## MVP modules

### B. Local Dating CRM
- Personal notes and event summaries
- Reminders and lightweight risk flags (user-controlled)

### C. Consent & Action Traceability
- Event-based consent records (opt-in)
- Optional counterparty co-sign via QR/token (future)

### Ban Protection (Account Health)
- Proof-of-compliance operation logs (metadata only)
- Daily rollup anchor on-chain
- Exportable "Appeal Package"

## Out of scope (MVP)

- Scraping Tinder profiles/chats
- Automated actions on Tinder (auto-like, auto-message, etc.)
- KYC/issuer VC partnerships (A module is later)
