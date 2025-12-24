# TrustTrace MVP (Commitment / Proof Anchoring)

This repository is a **build-ready engineering spec + starter code** for an MVP that:
- Keeps **all sensitive data local** (encrypted on the user's device)
- Anchors **only non-PII commitments** (hashes / Merkle roots) on-chain
- Supports:
  - **B (Local Dating CRM)**: private interaction summaries and reminders
  - **C (Consent & Action Traceability)**: user-controlled consent events (opt-in)
  - **Ban Protection**: proof-of-compliance activity logs + verifiable timestamp anchors

> Hard constraint: **Never put personal data on-chain** (including other people's data, photos, chat, identifiers, URLs, etc.).

## Monorepo layout

- `docs/` — product + security + privacy engineering specifications
- `extension/` — Chrome Extension (Manifest V3) starter (local vault + event rollups)
- `contracts/` — Solidity contract (Daily Root Registry) + Hardhat project
- `dapp/` — simple web dApp to submit commitments on-chain (wallet-injected environment)

## Quick start

### 1) Contracts (Hardhat)

```bash
cd contracts
npm install
npx hardhat test
npx hardhat compile
```

### 2) Extension (MV3)

```bash
cd extension
npm install
npm run build
```

Load unpacked extension:
- Chrome → `chrome://extensions` → Enable *Developer mode* → *Load unpacked* → select `extension/dist`

### 3) dApp (submit commitments)

```bash
cd dapp
npm install
npm run dev
```

Open the local dApp in the browser (Metamask available on this page), paste a commitment (bytes32 hex), and submit.

## MVP scope

- **Local encrypted vault** for event records
- **Daily rollup**: Merkle root over event leaves
- **Anchor**: on-chain commit of daily root (bytes32)
- **Export**: Appeal Package (human-readable + proof references)
- **No scraping / no automated interaction with Tinder**: only user-entered events, with optional minimal domain detection for UX.

## License

MIT
