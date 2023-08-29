const abi_1 = require("@ethersproject/abi");
var __importDefault = function (mod) {
  return (mod && mod.__esModule) ? mod : { "default": mod };
};
import * as foundryup from '@foundry-rs/easy-foundryup'
import {task, types} from "hardhat/config";
const axios_1 = __importDefault(require("axios"));
const qs_1 = __importDefault(require("qs"));
const { execSync } = require('child_process');


task('opbnb-verify',
  `Verify contracts for opbnb.
  Before use, please add SCAN_API and API_KEY two environment variables in the .env file
  1. If you have developments and want to verify all contracts, please execute: npx hardhat opbnb-verify --network your-network-name
  2. If you have developments, but only want to verify one of the contracts, please execute: npx hardhat opbnb-verify --network your-network-name yourContractName
  3. If you do not have developments, then you can only pass in the necessary parameters (contract name, contract address, contract constructor input parameters): npx hardhat opbnb-verify --contract-address 0xc0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d30007 L2CrossDomainMessenger 0xd506952e78eecd5d4424b1990a0c99b1568e7c2c
  The default compiler version is: v0.8.15+commit.e14f2714, you can override it with the --compiler-version parameter`)
  .addOptionalPositionalParam(
    'contractName',
    'Name of the contract to verify',
    '',
    types.string
  )
  .addOptionalParam(
    'contractAddress',
    'Address of the contract to verify',
    '',
    types.string
  )
  .addOptionalParam(
    'compilerVersion',
    'compiler version of the contract to verify',
    'v0.8.15+commit.e14f2714',
    types.string
  )
  .addOptionalVariadicPositionalParam(
    'constructorArguments',
    'arguments of constructor',
  )
  .setAction(async (args, hre) => {
    if (!args.contractName) {
      const all=await hre.deployments.all()
      //handle all contract in deployment
      for (let fileName in all ) {
        const one=all[fileName]
        const metadataString = one.metadata
        if (metadataString) {
          const metadata = JSON.parse(metadataString);
          var _a;
          const compilationTarget = (_a = metadata.settings) === null || _a === void 0 ? void 0 : _a.compilationTarget;
          let contractName;
          if (compilationTarget) {
            contractName = compilationTarget[Object.keys(compilationTarget)[0]];
          }else{
            console.log(`can not find compilationTarget in metadata in ${fileName}!,skip`)
          }
          if (!contractName) {
            console.log(`Failed to extract contract name from metadata.settings.compilationTarget in ${fileName}. Skipping.`);
          }
          console.log(`will handle contract:${contractName}`)
          const result = await oneContractVerify({contractName:contractName},hre)
          if (result) {
            console.log(`success`)
          } else {
            console.log(`fail`)
          }
        } else {
          console.log(`Contract deployment file ${fileName} was deployed without saving metadata. Cannot submit to etherscan, skipping.`);
        }
      }
    }else{
      const result = await oneContractVerify(args,hre)
      if (result) {
        console.log(`success`)
      } else {
        console.log(`fail`)
      }
    }
  })

async function oneContractVerify(args, hre):Promise<boolean>{
  let contractName = args.contractName
  let contractAddress = args.contractAddress
  let compilerVersion = args.compilerVersion

  let contractNamePath;
  let constructorArguments;

  //if deployment exist,use it firstly
  try{
    const all=await hre.deployments.all()
    let deployment
    for (let fileName in all) {
      const oneDeployment = all[fileName]
      const {metadata: metadataString } = oneDeployment;

      if (metadataString) {
        const metadata = JSON.parse(metadataString);
        var _a;
        const compilationTarget = (_a = metadata.settings) === null || _a === void 0 ? void 0 : _a.compilationTarget;
        let contractNameInDev;
        if (compilationTarget) {
          contractNameInDev = compilationTarget[Object.keys(compilationTarget)[0]];
        }
        if (contractNameInDev&&contractNameInDev===contractName){
          deployment = oneDeployment
          break
        }
      }
    }

    if (deployment.args) {
      const constructor = deployment.abi.find((v) => v.type === 'constructor');
      if (constructor) {
        constructorArguments = abi_1.defaultAbiCoder
          .encode(constructor.inputs, deployment.args)
          .slice(2);
      }
    } else {
      console.log(`no args found, assuming empty constructor...`);
    }
    const { address, metadata: metadataString } = deployment;
    contractAddress = address

    if (metadataString) {
      const metadata = JSON.parse(metadataString);
      var _b;
      const compilationTarget = (_b = metadata.settings) === null || _b === void 0 ? void 0 : _b.compilationTarget;
      let contractFilepath;
      if (compilationTarget) {
        contractFilepath = Object.keys(compilationTarget)[0];
      }
      if (contractFilepath) {
        contractNamePath = `${contractFilepath}:${contractName}`;
      } else {
        console.log(`Failed to extract contract fully qualified name from metadata.settings.compilationTarget for ${contractName}. Skipping.`);
      }
      compilerVersion=`v${metadata.compiler.version}`
    } else {
      console.log(`Contract ${contractName} was deployed without saving metadata. Cannot submit to etherscan, skipping.`);
    }
  } catch (e) {
    console.log(`deployment not found for:${contractName},err:${e.toString()},will verify without deployment`)
  }

  if (!contractAddress||contractAddress=='') {
    console.log(`miss contractAddress!`)
    return false
  }

  if (!contractNamePath) {
    try {
      let artifact = hre.artifacts.readArtifactSync(contractName);
      if (!artifact){
        throw Error(`artifact not found for ${contractName}`)
      }
      contractNamePath = `${artifact.sourceName}:${artifact.contractName}`
    } catch (e){
      const msg = (e as Error).message;
      console.log(`find artifact for ${contractName} fail:${msg}`)
      throw e;
    }
  }

  if (args.constructorArguments && args.constructorArguments.length>0) {
    try {
      let artifact = hre.artifacts.readArtifactSync(contractName);
      if (!artifact){
        throw Error(`artifact not found for ${contractName}`)
      }
      const constructor = artifact.abi.find((v) => v.type === 'constructor');
      if (constructor) {
        constructorArguments = abi_1.defaultAbiCoder
          .encode(constructor.inputs, args.constructorArguments)
          .slice(2);
      }
    } catch (e){
      const msg = (e as Error).message;
      console.log(`find artifact for ${contractName} fail:${msg}`)
      throw e;
    }
  }
  //all params ready
  console.log(`will verify.name:${contractName},address:${contractAddress},compiler:${compilerVersion},contractNamePath:${contractNamePath},constructorArguments:${constructorArguments}`)

  const verified = await isContractVerified(contractAddress)
  if (verified){
    console.log(`contract already be verified`)
    return true
  }

  //foundry
  const forgeCmd = await foundryup.getForgeCommand()
  let fArgs=['verify-contract', '--num-of-optimizations','999999','--compiler-version',compilerVersion,'--show-standard-json-input',contractAddress,contractNamePath]
  const solcInput=execSync(`${forgeCmd} ${fArgs.join(' ')}`,{encoding:'utf8'})

  // console.log(`result:${solcInput}`)
  const postData = {
    apikey: process.env.API_KEY,
    module: 'contract',
    action: 'verifysourcecode',
    contractaddress: contractAddress,
    sourceCode: solcInput,
    codeformat: 'solidity-standard-json-input',
    contractname: contractNamePath,
    compilerversion: compilerVersion,
    constructorArguements:constructorArguments,
    licenseType:3,
  };
  // console.log(`verify submit data:${postData}`)
  const formDataAsString = qs_1.default.stringify(postData);
  const submissionResponse = await axios_1.default.request({
    url: `${process.env.SCAN_API}`,
    method: 'POST',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    data: formDataAsString,
  });
  const { data: submissionData } = submissionResponse;
  let guid;
  if (submissionData.status === '1') {
    guid = submissionData.result;
  }else {
    console.log(`contract ${postData.contractname} failed to submit : "${submissionData.message}" : "${submissionData.result}"`, submissionData);
    return false;
  }
  if (!guid) {
    console.log(`contract submission for ${postData.contractname} failed to return a guid`);
    return false;
  }

  console.log(`submit success,guid:${guid},checking result status...`)

  let result;
  while (!result) {
    await new Promise((resolve) => setTimeout(resolve, 10 * 1000));
    result = await checkStatus(contractName,guid);
  }
  if (result === 'success') {
    console.log(` => contract ${postData.contractname} is now verified`);
    return true
  }
  if (result === 'failure') {
    console.log(`fail`)
    return false
  }
  return false
}

async function isContractVerified(address):Promise<boolean>{
  const abiResponse = await axios_1.default.get(`${process.env.SCAN_API}?module=contract&action=getabi&address=${address}&apikey=${process.env.API_KEY}`);
  const { data: abiData } = abiResponse;
  let contractABI;
  if (abiData.status !== '0') {
    try {
      if (typeof abiData.result === 'string'){
        contractABI = JSON.parse(abiData.result);
      }else{
        contractABI = abiData.result
      }
    }catch (e) {
      console.log(`abi data err,result:${JSON.stringify(abiData.result)},${e.message}`);
      throw e
    }
  }
  if (contractABI && contractABI !== '') {
    console.log(`already verified: (${address}), skipping.`);
    return true;
  }
  return false
}

async function checkStatus(contractName,guid) {
  const statusResponse = await axios_1.default.get(`${process.env.SCAN_API}?apikey=${process.env.API_KEY}`, {
    params: {
      guid,
      module: 'contract',
      action: 'checkverifystatus',
    },
  });
  const { data: statusData } = statusResponse;
  // blockscout seems to return status == 1 in case of failure
  // so we check string first
  if (statusData.result === 'Pending in queue') {
    return undefined;
  }
  if (statusData.result !== 'Fail - Unable to verify') {
    if (statusData.status === 1||statusData.status==='1') {
      return 'success';
    }
  }
  console.log(`Failed to verify contract ${contractName}: ${statusData.message}, ${statusData.result}`);
  return 'failure';
}
