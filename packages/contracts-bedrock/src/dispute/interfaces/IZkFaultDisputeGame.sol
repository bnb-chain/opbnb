// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { IDisputeGame } from "./IDisputeGame.sol";

import "src/dispute/lib/Types.sol";

/// @title IZkFaultDisputeGame
/// @notice The interface for a zk fault proof backed dispute game.
interface IZkFaultDisputeGame is IDisputeGame {

    /// @notice Parameters to initialize the contract.
    struct InitParams {
        uint256 chainId;
        bytes32 aggregationVkey;
        bytes32 rangeVkeyCommitment;
        address verifierGateway;
        bytes32 startingOutputRoot;
        address owner;
        bytes32 rollupConfigHash;
    }

    /// @notice The public values committed to for an OP Succinct aggregation program.
    struct AggregationOutputs {
        bytes32 l1Head;
        bytes32 l2PreRoot;
        bytes32 claimRoot;
        uint256 claimBlockNum;
        uint256 chainId;
        bytes32 rollupConfigHash;
        bytes32 rangeVkeyCommitment;
    }

    /// @notice Challengers use this function to challenge a claim by providing a proof.
    function challengeByProof(uint256 _disputeClaimIndex, Claim _expectedClaim,
                            Claim[] calldata _originalClaims, bytes calldata _proof) external;

    /// @notice Challengers use this function with a bond to signal that they want to challenge a claim.
    function challengeBySignal(uint256 _disputeClaimIndex) payable external;

    /// @notice Proposers use this function to submit a proof for a signal.
    function submitProofForSignal(uint256 _disputeClaimIndex, Claim[] calldata _originalClaims,
                                    bytes calldata _proof) external;

    /// @notice Challengers user this function to resolve challenge after the generation proof period has passed.
    function resolveClaim() external;

    /// @notice Starting output root and block number of the game.
    function startingOutputRoot() external view returns (Hash startingRoot_, uint256 l2BlockNumber_);

    /// @notice Only the starting block number of the game.
    function startingBlockNumber() external view returns (uint256 startingBlockNumber_);

    /// @notice Only the starting output root of the game.
    function startingRootHash() external view returns (Hash startingRootHash_);

    function gameWinner() external view returns (address gameWinner_);
}
