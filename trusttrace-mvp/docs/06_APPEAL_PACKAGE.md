# Appeal Package format (MVP)

## Purpose
Provide a structured, human-readable and verifiable set of evidence for account appeals, without exposing third-party personal data.

## Contents
- `appeal_summary.md`
- `commitment.json`
- `bundle_<day>.json` (leaf hashes only, no payload)
- `environment.json` (browser + extension version)
- `README_verify.md` (how to verify hashes locally)

## appeal_summary.md (example fields)
- Time window: YYYY-MM-DD to YYYY-MM-DD
- High-level behavior summary:
  - total logins
  - total swipes (aggregated)
  - total messages sent (aggregated)
- User statement (free text)
- Links to on-chain tx (added after submission)

## Verification notes
- Without salts and payloads, third parties cannot reconstruct contents.
- The user may optionally disclose select decrypted entries during appeal, at their discretion.
