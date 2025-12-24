import { ethers } from "hardhat";

async function main() {
  const Registry = await ethers.getContractFactory("DailyRootRegistry");
  const registry = await Registry.deploy();
  await registry.waitForDeployment();

  const addr = await registry.getAddress();
  console.log("DailyRootRegistry deployed to:", addr);
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});
