// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title DailyRootRegistry
/// @notice Stores one bytes32 root per (signer, day) with timestamp anchors.
/// @dev Do NOT store any payload or personal data on-chain. Root must be a commitment only.
contract DailyRootRegistry {
    mapping(address => mapping(uint256 => bytes32)) private _roots;

    event RootCommitted(
        address indexed signer,
        uint256 indexed day,
        bytes32 root,
        string schemaVersion,
        bytes32 metaHash,
        uint256 createdAt
    );

    event RootSuperseded(
        address indexed signer,
        uint256 indexed day,
        bytes32 oldRoot,
        bytes32 newRoot,
        string schemaVersion,
        bytes32 metaHash,
        uint256 createdAt
    );

    function getRoot(address signer, uint256 day) external view returns (bytes32) {
        return _roots[signer][day];
    }

    /// @param day YYYYMMDD, e.g., 20251224
    /// @param root Merkle root (bytes32), commitment only
    /// @param schemaVersion e.g., "tt.bundle.v1"
    /// @param metaHash optional hash of non-sensitive metadata (or bytes32(0))
    function commitRoot(uint256 day, bytes32 root, string calldata schemaVersion, bytes32 metaHash) external {
        require(day > 19700101 && day < 30000101, "invalid day");
        require(root != bytes32(0), "root required");
        require(_roots[msg.sender][day] == bytes32(0), "already committed");
        _roots[msg.sender][day] = root;
        emit RootCommitted(msg.sender, day, root, schemaVersion, metaHash, block.timestamp);
    }

    function supersedeRoot(uint256 day, bytes32 newRoot, string calldata schemaVersion, bytes32 metaHash) external {
        require(day > 19700101 && day < 30000101, "invalid day");
        require(newRoot != bytes32(0), "root required");
        bytes32 old = _roots[msg.sender][day];
        require(old != bytes32(0), "no root");
        require(old != newRoot, "same root");
        _roots[msg.sender][day] = newRoot;
        emit RootSuperseded(msg.sender, day, old, newRoot, schemaVersion, metaHash, block.timestamp);
    }
}
