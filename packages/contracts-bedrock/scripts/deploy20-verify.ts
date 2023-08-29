// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
import hre from "hardhat";

async function main() {
  // const factoryContract = await hre.ethers.getContractAt("OptimismMintableERC20Factory",'0x4200000000000000000000000000000000000012');
  // const tx = await factoryContract.createOptimismMintableERC20('0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee','Binance USD-fake','BUSDf')
  // console.log(`${tx.hash}`)
  // let optimismMintableERC20 = await hre.ethers.getContractFactory("OptimismMintableERC20");
  // let tokenContract = optimismMintableERC20.attach(erc20ContractAddress)
  // const name = await tokenContract.name()
  // const symbol = await tokenContract.symbol()
  //
  // console.log(
  //   `deployed to ${erc20ContractAddress}, ${name}, ${symbol}`
  // );
  try {
    await hre.run("verify:verify", {
      address:"0xe4f27b04cc7729901876b44f4eaa5102ec150265",
      constructorArguments: ["0x4200000000000000000000000000000000000010","0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee","Binance USD-fake","BUSDf"],
    });
  } catch (e) {
    const msg = (e as Error).message;
    if (msg.includes('Contract source code already verified')) {
      console.log(`Contract [0xe4f27b04cc7729901876b44f4eaa5102ec150265] source code already verified`);
    } else {
      throw e;
    }
  }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
