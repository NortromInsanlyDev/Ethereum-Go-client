import "@nomiclabs/hardhat-waffle";
import { task, types } from "hardhat/config";
import { MockToken__factory } from "../typechain/factories/MockToken__factory";
import { Vault__factory } from "../typechain/factories/Vault__factory";
import type { MockToken } from "../typechain/MockToken";
import type { Vault } from "../typechain/Vault";

task("withdraw", "withdraw tokens from vault")
  .addParam("vault", "vault address")
  .addParam("token", "token address")
  .addParam("index", "index of account to withdraw", 0, types.int)
  .addParam("amount", "amount to withdraw")
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
    await vault.connect(user).withdraw(amount);
    console.log(
      `user ${user.address} has withdrawn ${amount} of tokens from vault`
    );
  });
