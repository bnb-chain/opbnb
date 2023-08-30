import { predeploys } from '@eth-optimism/core-utils'
import { ethers } from 'ethers'
import portalArtifactsMainnet from '@eth-optimism/contracts-bedrock/deployments/mainnet/OptimismPortalProxy.json'
import portalArtifactsGoerli from '@eth-optimism/contracts-bedrock/deployments/goerli/OptimismPortalProxy.json'
import l2OutputOracleArtifactsMainnet from '@eth-optimism/contracts-bedrock/deployments/mainnet/L2OutputOracleProxy.json'
import l2OutputOracleArtifactsGoerli from '@eth-optimism/contracts-bedrock/deployments/goerli/L2OutputOracleProxy.json'
import addressManagerArtifactMainnet from '@eth-optimism/contracts-bedrock/deployments/mainnet/AddressManager.json'
import addressManagerArtifactGoerli from '@eth-optimism/contracts-bedrock/deployments/goerli/AddressManager.json'
import l1StandardBridgeArtifactMainnet from '@eth-optimism/contracts-bedrock/deployments/mainnet/L1StandardBridgeProxy.json'
import l1StandardBridgeArtifactGoerli from '@eth-optimism/contracts-bedrock/deployments/goerli/L1StandardBridgeProxy.json'
import l1CrossDomainMessengerArtifactMainnet from '@eth-optimism/contracts-bedrock/deployments/mainnet/L1CrossDomainMessengerProxy.json'
import l1CrossDomainMessengerArtifactGoerli from '@eth-optimism/contracts-bedrock/deployments/goerli/L1CrossDomainMessengerProxy.json'

const portalAddresses = {
  mainnet: portalArtifactsMainnet.address,
  goerli: portalArtifactsGoerli.address,
}

const l2OutputOracleAddresses = {
  mainnet: l2OutputOracleArtifactsMainnet.address,
  goerli: l2OutputOracleArtifactsGoerli.address,
}

const addressManagerAddresses = {
  mainnet: addressManagerArtifactMainnet.address,
  goerli: addressManagerArtifactGoerli.address,
}

const l1StandardBridgeAddresses = {
  mainnet: l1StandardBridgeArtifactMainnet.address,
  goerli: l1StandardBridgeArtifactGoerli.address,
}

const l1CrossDomainMessengerAddresses = {
  mainnet: l1CrossDomainMessengerArtifactMainnet.address,
  goerli: l1CrossDomainMessengerArtifactGoerli.address,
}

// legacy
const stateCommitmentChainAddresses = {
  mainnet: '0xBe5dAb4A2e9cd0F27300dB4aB94BeE3A233AEB19',
  goerli: '0x9c945aC97Baf48cB784AbBB61399beB71aF7A378',
}

// legacy
const canonicalTransactionChainAddresses = {
  mainnet: '0x5E4e65926BA27467555EB562121fac00D24E9dD2',
  goerli: '0x607F755149cFEB3a14E1Dc3A4E2450Cde7dfb04D',
}

import {
  L1ChainID,
  L2ChainID,
  OEContractsLike,
  OEL1ContractsLike,
  OEL2ContractsLike,
  BridgeAdapterData,
} from '../interfaces'
import {
  StandardBridgeAdapter,
  DAIBridgeAdapter,
  ECOBridgeAdapter,
} from '../adapters'

export const DEPOSIT_CONFIRMATION_BLOCKS: {
  [ChainID in L2ChainID]: number
} = {
  [L2ChainID.OPTIMISM]: 50 as const,
  [L2ChainID.OPTIMISM_GOERLI]: 12 as const,
  [L2ChainID.OPTIMISM_HARDHAT_LOCAL]: 2 as const,
  [L2ChainID.OPTIMISM_HARDHAT_DEVNET]: 2 as const,
  [L2ChainID.OPTIMISM_BEDROCK_LOCAL_DEVNET]: 2 as const,
  [L2ChainID.OPTIMISM_BEDROCK_ALPHA_TESTNET]: 12 as const,
  [L2ChainID.BASE_GOERLI]: 12 as const,
  [L2ChainID.ZORA_GOERLI]: 12 as const,
  [L2ChainID.ZORA_MAINNET]: 50 as const,
}

export const CHAIN_BLOCK_TIMES: {
  [ChainID in L1ChainID]: number
} = {
  [L1ChainID.MAINNET]: 13 as const,
  [L1ChainID.GOERLI]: 15 as const,
  [L1ChainID.HARDHAT_LOCAL]: 1 as const,
  [L1ChainID.BEDROCK_LOCAL_DEVNET]: 15 as const,
}

/**
 * Full list of default L2 contract addresses.
 */
export const DEFAULT_L2_CONTRACT_ADDRESSES: OEL2ContractsLike = {
  L2CrossDomainMessenger: predeploys.L2CrossDomainMessenger,
  L2ToL1MessagePasser: predeploys.L2ToL1MessagePasser,
  L2StandardBridge: predeploys.L2StandardBridge,
  OVM_L1BlockNumber: predeploys.L1BlockNumber,
  OVM_L2ToL1MessagePasser: predeploys.L2ToL1MessagePasser,
  OVM_DeployerWhitelist: predeploys.DeployerWhitelist,
  OVM_ETH: predeploys.LegacyERC20ETH,
  OVM_GasPriceOracle: predeploys.GasPriceOracle,
  OVM_SequencerFeeVault: predeploys.SequencerFeeVault,
  WETH: predeploys.WETH9,
  BedrockMessagePasser: predeploys.L2ToL1MessagePasser,
}

/**
 * Loads the L1 contracts for a given network by the network name.
 *
 * @param network The name of the network to load the contracts for.
 * @returns The L1 contracts for the given network.
 */
const getL1ContractsByNetworkName = (network: string): OEL1ContractsLike => {
  return {
    AddressManager: addressManagerAddresses[network],
    L1CrossDomainMessenger: l1CrossDomainMessengerAddresses[network],
    L1StandardBridge: l1StandardBridgeAddresses[network],
    StateCommitmentChain: stateCommitmentChainAddresses[network],
    CanonicalTransactionChain: canonicalTransactionChainAddresses[network],
    BondManager: ethers.constants.AddressZero,
    OptimismPortal: portalAddresses[network],
    L2OutputOracle: l2OutputOracleAddresses[network],
  }
}

/**
 * Mapping of L1 chain IDs to the appropriate contract addresses for the OE deployments to the
 * given network. Simplifies the process of getting the correct contract addresses for a given
 * contract name.
 */
export const CONTRACT_ADDRESSES: {
  [ChainID in L2ChainID]: OEContractsLike
} = {
  [L2ChainID.OPTIMISM]: {
    l1: getL1ContractsByNetworkName('mainnet'),
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.OPTIMISM_GOERLI]: {
    l1: getL1ContractsByNetworkName('goerli'),
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.OPTIMISM_HARDHAT_LOCAL]: {
    l1: {
      AddressManager: '0x5FbDB2315678afecb367f032d93F642f64180aa3' as const,
      L1CrossDomainMessenger:
        '0x8A791620dd6260079BF849Dc5567aDC3F2FdC318' as const,
      L1StandardBridge: '0x610178dA211FEF7D417bC0e6FeD39F05609AD788' as const,
      StateCommitmentChain:
        '0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9' as const,
      CanonicalTransactionChain:
        '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9' as const,
      BondManager: '0x5FC8d32690cc91D4c39d9d3abcBD16989F875707' as const,
      // FIXME
      OptimismPortal: '0x0000000000000000000000000000000000000000' as const,
      L2OutputOracle: '0x0000000000000000000000000000000000000000' as const,
    },
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.OPTIMISM_HARDHAT_DEVNET]: {
    l1: {
      AddressManager: '0x5FbDB2315678afecb367f032d93F642f64180aa3' as const,
      L1CrossDomainMessenger:
        '0x8A791620dd6260079BF849Dc5567aDC3F2FdC318' as const,
      L1StandardBridge: '0x610178dA211FEF7D417bC0e6FeD39F05609AD788' as const,
      StateCommitmentChain:
        '0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9' as const,
      CanonicalTransactionChain:
        '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9' as const,
      BondManager: '0x5FC8d32690cc91D4c39d9d3abcBD16989F875707' as const,
      OptimismPortal: '0x0000000000000000000000000000000000000000' as const,
      L2OutputOracle: '0x0000000000000000000000000000000000000000' as const,
    },
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.OPTIMISM_BEDROCK_LOCAL_DEVNET]: {
    l1: {
      AddressManager: '0x6900000000000000000000000000000000000005' as const,
      L1CrossDomainMessenger:
        '0x6900000000000000000000000000000000000002' as const,
      L1StandardBridge: '0x6900000000000000000000000000000000000003' as const,
      StateCommitmentChain:
        '0x0000000000000000000000000000000000000000' as const,
      CanonicalTransactionChain:
        '0x0000000000000000000000000000000000000000' as const,
      BondManager: '0x0000000000000000000000000000000000000000' as const,
      OptimismPortal: '0x6900000000000000000000000000000000000001' as const,
      L2OutputOracle: '0x6900000000000000000000000000000000000000' as const,
    },
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.OPTIMISM_BEDROCK_ALPHA_TESTNET]: {
    l1: {
      AddressManager: '0xb4e08DcE1F323608229265c9d4125E22a4B9dbAF' as const,
      L1CrossDomainMessenger:
        '0x838a6DC4E37CA45D4Ef05bb776bf05eEf50798De' as const,
      L1StandardBridge: '0xFf94B6C486350aD92561Ba09bad3a59df764Da92' as const,
      StateCommitmentChain:
        '0x0000000000000000000000000000000000000000' as const,
      CanonicalTransactionChain:
        '0x0000000000000000000000000000000000000000' as const,
      BondManager: '0x0000000000000000000000000000000000000000' as const,
      OptimismPortal: '0xA581Ca3353DB73115C4625FFC7aDF5dB379434A8' as const,
      L2OutputOracle: '0x3A234299a14De50027eA65dCdf1c0DaC729e04A6' as const,
    },
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.BASE_GOERLI]: {
    l1: {
      AddressManager: '0x4Cf6b56b14c6CFcB72A75611080514F94624c54e' as const,
      L1CrossDomainMessenger:
        '0x8e5693140eA606bcEB98761d9beB1BC87383706D' as const,
      L1StandardBridge: '0xfA6D8Ee5BE770F84FC001D098C4bD604Fe01284a' as const,
      StateCommitmentChain:
        '0x0000000000000000000000000000000000000000' as const,
      CanonicalTransactionChain:
        '0x0000000000000000000000000000000000000000' as const,
      BondManager: '0x0000000000000000000000000000000000000000' as const,
      OptimismPortal: '0xe93c8cD0D409341205A592f8c4Ac1A5fe5585cfA' as const,
      L2OutputOracle: '0x2A35891ff30313CcFa6CE88dcf3858bb075A2298' as const,
    },
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  // Zora Goerli
  [L2ChainID.ZORA_GOERLI]: {
    l1: {
      AddressManager: '0x54f4676203dEDA6C08E0D40557A119c602bFA246' as const,
      L1CrossDomainMessenger:
        '0xD87342e16352D33170557A7dA1e5fB966a60FafC' as const,
      L1StandardBridge: '0x7CC09AC2452D6555d5e0C213Ab9E2d44eFbFc956' as const,
      StateCommitmentChain:
        '0x0000000000000000000000000000000000000000' as const,
      CanonicalTransactionChain:
        '0x0000000000000000000000000000000000000000' as const,
      BondManager: '0x0000000000000000000000000000000000000000' as const,
      OptimismPortal: '0xDb9F51790365e7dc196e7D072728df39Be958ACe' as const,
      L2OutputOracle: '0xdD292C9eEd00f6A32Ff5245d0BCd7f2a15f24e00' as const,
    },
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.ZORA_MAINNET]: {
    l1: {
      AddressManager: '0xEF8115F2733fb2033a7c756402Fc1deaa56550Ef' as const,
      L1CrossDomainMessenger:
        '0xdC40a14d9abd6F410226f1E6de71aE03441ca506' as const,
      L1StandardBridge: '0x3e2Ea9B92B7E48A52296fD261dc26fd995284631' as const,
      StateCommitmentChain:
        '0x0000000000000000000000000000000000000000' as const,
      CanonicalTransactionChain:
        '0x0000000000000000000000000000000000000000' as const,
      BondManager: '0x0000000000000000000000000000000000000000' as const,
      OptimismPortal: '0x1a0ad011913A150f69f6A19DF447A0CfD9551054' as const,
      L2OutputOracle: '0x9E6204F750cD866b299594e2aC9eA824E2e5f95c' as const,
    },
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
}

/**
 * Mapping of L1 chain IDs to the list of custom bridge addresses for each chain.
 */
export const BRIDGE_ADAPTER_DATA: {
  [ChainID in L2ChainID]?: BridgeAdapterData
} = {
  [L2ChainID.OPTIMISM]: {
    wstETH: {
      Adapter: DAIBridgeAdapter,
      l1Bridge: '0x76943C0D61395d8F2edF9060e1533529cAe05dE6' as const,
      l2Bridge: '0x8E01013243a96601a86eb3153F0d9Fa4fbFb6957' as const,
    },
    BitBTC: {
      Adapter: StandardBridgeAdapter,
      l1Bridge: '0xaBA2c5F108F7E820C049D5Af70B16ac266c8f128' as const,
      l2Bridge: '0x158F513096923fF2d3aab2BcF4478536de6725e2' as const,
    },
    DAI: {
      Adapter: DAIBridgeAdapter,
      l1Bridge: '0x10E6593CDda8c58a1d0f14C5164B376352a55f2F' as const,
      l2Bridge: '0x467194771dAe2967Aef3ECbEDD3Bf9a310C76C65' as const,
    },
    ECO: {
      Adapter: ECOBridgeAdapter,
      l1Bridge: '0xAa029BbdC947F5205fBa0F3C11b592420B58f824' as const,
      l2Bridge: '0xAa029BbdC947F5205fBa0F3C11b592420B58f824' as const,
    },
  },
  [L2ChainID.OPTIMISM_GOERLI]: {
    DAI: {
      Adapter: DAIBridgeAdapter,
      l1Bridge: '0x05a388Db09C2D44ec0b00Ee188cD42365c42Df23' as const,
      l2Bridge: '0x467194771dAe2967Aef3ECbEDD3Bf9a310C76C65' as const,
    },
    ECO: {
      Adapter: ECOBridgeAdapter,
      l1Bridge: '0x9A4464D6bFE006715382D39D183AAf66c952a3e0' as const,
      l2Bridge: '0x6aA809bAeA2e4C057b3994127cB165119c6fc3B2' as const,
    },
  },
}
