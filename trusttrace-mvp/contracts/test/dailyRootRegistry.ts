import { expect } from "chai";
import { ethers } from "hardhat";

describe("DailyRootRegistry", function () {
  it("commits and reads a root", async function () {
    const [alice] = await ethers.getSigners();
    const Registry = await ethers.getContractFactory("DailyRootRegistry");
    const reg = await Registry.deploy();
    await reg.waitForDeployment();

    const day = 20251224;
    const root = "0x" + "11".repeat(32);
    await expect(reg.connect(alice).commitRoot(day, root, "tt.bundle.v1", ethers.ZeroHash))
      .to.emit(reg, "RootCommitted");

    expect(await reg.getRoot(alice.address, day)).to.equal(root);
  });

  it("prevents double commit, allows supersede", async function () {
    const [alice] = await ethers.getSigners();
    const Registry = await ethers.getContractFactory("DailyRootRegistry");
    const reg = await Registry.deploy();
    await reg.waitForDeployment();

    const day = 20251224;
    const r1 = "0x" + "11".repeat(32);
    const r2 = "0x" + "22".repeat(32);

    await reg.commitRoot(day, r1, "tt.bundle.v1", ethers.ZeroHash);
    await expect(reg.commitRoot(day, r2, "tt.bundle.v1", ethers.ZeroHash)).to.be.revertedWith("already committed");

    await expect(reg.supersedeRoot(day, r2, "tt.bundle.v1", ethers.ZeroHash))
      .to.emit(reg, "RootSuperseded");
    expect(await reg.getRoot(alice.address, day)).to.equal(r2);
  });
});
