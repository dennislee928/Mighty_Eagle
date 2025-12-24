# Contracts

## Deploy
1. Copy `.env.example` → `.env` and fill your RPC + key (optional for testnet).
2. Run:

```bash
npm install
npx hardhat test
npx hardhat run scripts/deploy.ts --network hardhat
```

## Notes
- Use a dedicated address per user for privacy.
- The contract stores only bytes32 roots.
