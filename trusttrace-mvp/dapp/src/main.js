import { ethers } from "ethers";

const ABI = [
  "function commitRoot(uint256 day, bytes32 root, string schemaVersion, bytes32 metaHash) external",
  "function getRoot(address signer, uint256 day) external view returns (bytes32)",
  "event RootCommitted(address indexed signer, uint256 indexed day, bytes32 root, string schemaVersion, bytes32 metaHash, uint256 createdAt)"
];

const $ = (id) => document.getElementById(id);
let signer = null;

function setStatus(s) { $("status").textContent = s; }

$("connect").addEventListener("click", async () => {
  if (!window.ethereum) {
    setStatus("No injected wallet found. Install MetaMask.");
    return;
  }
  const provider = new ethers.BrowserProvider(window.ethereum);
  await provider.send("eth_requestAccounts", []);
  signer = await provider.getSigner();
  $("submit").disabled = false;
  setStatus("Connected: " + await signer.getAddress());
});

$("submit").addEventListener("click", async () => {
  try {
    const addr = $("registry").value.trim();
    const day = BigInt($("day").value.trim());
    const root = $("root").value.trim();
    const schema = $("schema").value.trim();

    if (!ethers.isAddress(addr)) throw new Error("Invalid registry address");
    if (!/^0x[0-9a-fA-F]{64}$/.test(root)) throw new Error("Root must be 32-byte hex");
    if (schema.length > 32) throw new Error("Schema too long (<= 32 chars recommended)");

    const reg = new ethers.Contract(addr, ABI, signer);
    setStatus("Submitting transaction...");
    const tx = await reg.commitRoot(day, root, schema, ethers.ZeroHash);
    setStatus("Tx sent: " + tx.hash);
    const receipt = await tx.wait();
    setStatus("Confirmed in block " + receipt.blockNumber + " | tx " + tx.hash);
  } catch (e) {
    setStatus("Error: " + (e?.message ?? String(e)));
  }
});
