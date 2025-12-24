# Threat model (STRIDE-lite)

## Assets
- User vault plaintext
- Vault encryption key / passphrase
- Salts for commitments
- Daily bundle root
- Exported packages

## Threats & mitigations

### Spoofing
- Attacker pretends to be the user.
  - Mitigation: local passphrase; optional OS credential store (future).

### Tampering
- Modify records to create fake history.
  - Mitigation: on-chain timestamp anchor for daily root; local audit trail.

### Repudiation
- User denies making a claim.
  - Mitigation: voluntary anchoring; optional wallet signature on commitment.

### Information disclosure
- Leakage of PII via stored content or commitment.
  - Mitigation: encryption at rest; salted leaf hashing; strict schema; redaction warnings.

### Denial of service
- Platform changes or blocks extension pages.
  - Mitigation: no dependency on platform DOM or APIs.

### Elevation of privilege
- Malicious web pages trying to exfiltrate data.
  - Mitigation: no content script data extraction; least privilege; CSP (future).

## Abuse cases (product)
- Using the tool to "doxx" or store other people's identity on-chain.
  - Mitigation: guardrails + UX warnings + schema design (no fields for identity).
