// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
import hre from "hardhat";

async function main() {
  await oneToken('0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee', 'Binance USD-fake2', 'BUSDf2');
  process.exit();
  return

}

async function oneToken(anotherToken,name,symbol) {
  const factoryContract = await hre.ethers.getContractAt("OptimismMintableERC20Factory", '0x4200000000000000000000000000000000000012');
  const tx = await factoryContract.createOptimismMintableERC20(anotherToken,name,symbol)
  console.log(`${tx.hash}`)

  // 2. Listen for the events

  const tokenAddress = await getTokenAddress(factoryContract); // Await the Promise to get the value
  console.log("Token address:", tokenAddress);

  try {
    await hre.run("opbnb-verify", {
      contractName: "OptimismMintableERC20",
      contractAddress: tokenAddress,
      constructorArguments: ['0x4200000000000000000000000000000000000010', anotherToken,name,symbol]
    })
  } catch (e) {
    const msg = (e as Error).message;
    console.log(`verify err ${msg}`);
    throw e;
  }
}
function getTokenAddress(factoryContract): Promise<any> {
  return new Promise((resolve) => {
    factoryContract.on("OptimismMintableERC20Created", (localToken, remoteToken, deployer) => {
      console.log("Event received:", localToken, remoteToken, deployer);
      console.log(`deployed to ${localToken}`);
      resolve(localToken); // Resolve the Promise with the updated value
    });
  });
}
// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
