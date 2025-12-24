# Contract specification: DailyRootRegistry

## Purpose
Provide an immutable, publicly verifiable timestamp anchor for a daily Merkle root.

## Data stored
- `(signer, day, root, schemaVersion, metaHash, createdAt)`

No payload data is stored on-chain.

## Interface (v1)

- `commitRoot(uint256 day, bytes32 root, string schemaVersion, bytes32 metaHash)`
  - `day` is `YYYYMMDD` as integer (e.g., 20251224)
  - emits `RootCommitted`

- `getRoot(address signer, uint256 day) -> bytes32`

Optional:
- `supersedeRoot(uint256 day, bytes32 newRoot, string schemaVersion, bytes32 metaHash)`
  - emits `RootSuperseded`
  - old root remains discoverable via events

## Events

- `event RootCommitted(address indexed signer, uint256 indexed day, bytes32 root, string schemaVersion, bytes32 metaHash, uint256 createdAt);`
- `event RootSuperseded(address indexed signer, uint256 indexed day, bytes32 oldRoot, bytes32 newRoot, string schemaVersion, bytes32 metaHash, uint256 createdAt);`

## Threat considerations

- Malicious users can commit arbitrary roots; verification requires local vault + salts.
- Replay: committing same root is allowed; supersede is explicit.
- Privacy: day and signer are public; recommend a dedicated wallet address per user.

## Networks

- MVP testing: Sepolia or local Hardhat network.
- Production: L2 (e.g., Base, Arbitrum, Optimism) to reduce fees.

## Gas considerations

- Storing one `bytes32` per day per address is minimal.
- Strings cost more; schemaVersion should be short (e.g., `tt.bundle.v1`).
