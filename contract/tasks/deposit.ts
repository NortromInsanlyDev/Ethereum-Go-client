import "@nomiclabs/hardhat-waffle";
import { task, types } from "hardhat/config";
import { MockToken__factory } from "../typechain/factories/MockToken__factory";
import { Vault__factory } from "../typechain/factories/Vault__factory";
import type { MockToken } from "../typechain/MockToken";
import type { Vault } from "../typechain/Vault";

task("deposit", "deposit tokens to vault")
  .addParam("vault", "vault address")
  .addParam("token", "token address")
  .addParam("index", "index account to deposit", 0, types.int)
  .addParam("amount", "amount to deposit")
  .setAction(async (taskArgs, { ethers }) => {
    const vaultAddress = taskArgs.vault;
    const tokenAddress = taskArgs.token;
    const amount = ethers.utils.parseEther(taskArgs.amount);
    if ((await ethers.provider.getCode(vaultAddress)) === "0x") {
      console.error("You need to deploy your contract first");
      return;
    }

    const accounts = await ethers.getSigners();
    const user = accounts[taskArgs.index];
    const token: MockToken = MockToken__factory.connect(
      tokenAddress,
      ethers.provider
    );

    const vault: Vault = Vault__factory.connect(vaultAddress, ethers.provider);
    const balance = await token.balanceOf(user.address);
    if (balance < amount) {
      console.error("You need to charge tokens first!");
      return;
    }
    await token.connect(user).approve(vault.address, amount);
    await vault.connect(user).deposit(amount);
    console.log(`user ${user.address} deposited ${amount} of tokens to vault`);
  });
