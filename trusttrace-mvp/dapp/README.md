# dApp (submit commitments)

This is a minimal wallet-enabled page to submit daily roots to `DailyRootRegistry`.

## Run
```bash
npm install
npm run dev
```

## Usage
1. Deploy `DailyRootRegistry` (see `contracts/`).
2. In the extension, export `commitment.json` and copy:
   - `dayInt` (YYYYMMDD)
   - `merkleRoot` (bytes32)
3. Paste into this page and submit.

## Safety
Never paste personal data into the dApp; only bytes32 roots.
