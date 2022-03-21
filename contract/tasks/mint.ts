import { task } from "hardhat/config";
import "@nomiclabs/hardhat-waffle";

import type { MockToken } from "../typechain/MockToken";
import { MockToken__factory } from "../typechain/factories/MockToken__factory";

task("mint", "mint tokens to accounts")
  .addParam("token", "token address")
  .setAction(async (taskArgs, { ethers }) => {
    const tokenAddress = taskArgs.token;

    if ((await ethers.provider.getCode(tokenAddress)) === "0x") {
      console.error("You need to deploy your contract first");
      return;
    }

    const accounts = await ethers.getSigners();
    const token: MockToken = MockToken__factory.connect(
      tokenAddress,
      ethers.provider
    );
    for (const account of accounts) {
      // console.log(account.address);
      await token
        .connect(accounts[0])
        .mint(account.address, ethers.utils.parseEther("100"));
    }
    for (const account of accounts) {
      // console.log(account.address);
      const balance = await token.balanceOf(account.address);
      console.log(`${account.address}: ${balance}`);
    }
  });
