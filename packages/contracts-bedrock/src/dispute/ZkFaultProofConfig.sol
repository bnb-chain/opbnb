// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { SP1VerifierGateway } from "@sp1-contracts/src/SP1VerifierGateway.sol";
import { Initializable } from "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import { ISemver } from "src/universal/ISemver.sol";

contract ZkFaultProofConfig is OwnableUpgradeable, Initializable, ISemver {

    /// @notice Semantic version.
    /// @custom:semver 1.0.0
    string public constant version = "1.0.0";

    /// @notice the block interval of claims proposed by proposers
    uint256 public blockDistance;

    /// @notice the hash of the rollup configuration
    bytes32 public rollupConfigHash;

    /// @notice The chain ID of the L2 chain.
    uint256 public chainId;

    /// @notice The verification key of the aggregation SP1 program.
    bytes32 public aggregationVkey;

    /// @notice The 32 byte commitment to the BabyBear representation of the verification key of the range SP1 program. Specifically,
    /// this verification is the output of converting the [u32; 8] range BabyBear verification key to a [u8; 32] array.
    bytes32 public rangeVkeyCommitment;

    /// @notice The deployed SP1VerifierGateway contract to request proofs from.
    SP1VerifierGateway public verifierGateway;

    /// @notice A trusted mapping of block numbers to block hashes.
    mapping(uint256 => bytes32) public historicBlockHashes;

    /// @notice Emitted when the aggregation vkey is updated.
    /// @param oldVkey The old aggregation vkey.
    /// @param newVkey The new aggregation vkey.
    event UpdatedAggregationVKey(bytes32 indexed oldVkey, bytes32 indexed newVkey);

    /// @notice Emitted when the range vkey commitment is updated.
    /// @param oldRangeVkeyCommitment The old range vkey commitment.
    /// @param newRangeVkeyCommitment The new range vkey commitment.
    event UpdatedRangeVkeyCommitment(bytes32 indexed oldRangeVkeyCommitment, bytes32 indexed newRangeVkeyCommitment);

    /// @notice Emitted when the verifier gateway is updated.
    /// @param oldVerifierGateway The old verifier gateway.
    /// @param newVerifierGateway The new verifier gateway.
    event UpdatedVerifierGateway(address indexed oldVerifierGateway, address indexed newVerifierGateway);

    /// @notice Emitted when the rollup config hash is updated.
    /// @param oldRollupConfigHash The old rollup config hash.
    /// @param newRollupConfigHash The new rollup config hash.
    event UpdatedRollupConfigHash(bytes32 indexed oldRollupConfigHash, bytes32 indexed newRollupConfigHash);

    constructor() {}

    function initialize (
        address _owner,
        uint256 _blockDistance,
        uint256 _chainId,
        bytes32 _aggregationVkey,
        bytes32 _rangeVkeyCommitment,
        address _verifierGateway,
        bytes32 _rollupConfigHash
    ) public initializer {
        __Ownable_init();
        transferOwnership(_owner);

        blockDistance = _blockDistance;
        chainId = _chainId;
        aggregationVkey = _aggregationVkey;
        rangeVkeyCommitment = _rangeVkeyCommitment;
        verifierGateway = SP1VerifierGateway(_verifierGateway);
        rollupConfigHash = _rollupConfigHash;
    }

    function updateAggregationVKey(bytes32 _aggregationVKey) external onlyOwner {
        aggregationVkey = _aggregationVKey;
        emit UpdatedAggregationVKey(aggregationVkey, _aggregationVKey);
    }

    function updateRangeVkeyCommitment(bytes32 _rangeVkeyCommitment) external onlyOwner {
        rangeVkeyCommitment = _rangeVkeyCommitment;
        emit UpdatedRangeVkeyCommitment(rangeVkeyCommitment, _rangeVkeyCommitment);
    }

    function updateVerifierGateway(address _verifierGateway) external onlyOwner {
        verifierGateway = SP1VerifierGateway(_verifierGateway);
        emit UpdatedVerifierGateway(address(verifierGateway), _verifierGateway);
    }

    function updateRollupConfigHash(bytes32 _rollupConfigHash) external onlyOwner {
        rollupConfigHash = _rollupConfigHash;
        emit UpdatedRollupConfigHash(rollupConfigHash, _rollupConfigHash);
    }

}
