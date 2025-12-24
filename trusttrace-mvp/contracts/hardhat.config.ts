import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.24",
    settings: { optimizer: { enabled: true, runs: 200 } },
  },
  networks: {
    hardhat: {},
    // sepolia: { url: process.env.SEPOLIA_RPC_URL!, accounts: [process.env.DEPLOYER_KEY!] },
  },
};

export default config;
