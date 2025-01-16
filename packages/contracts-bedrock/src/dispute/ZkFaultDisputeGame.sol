// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { FixedPointMathLib } from "@solady/utils/FixedPointMathLib.sol";

import { IDelayedWETH } from "src/dispute/interfaces/IDelayedWETH.sol";
import { IDisputeGame } from "src/dispute/interfaces/IDisputeGame.sol";
import { IZkFaultDisputeGame } from "src/dispute/interfaces/IZkFaultDisputeGame.sol";
import { ZkFaultProofConfig } from "src/dispute/ZkFaultProofConfig.sol";
import { SP1VerifierGateway } from "@sp1-contracts/src/SP1VerifierGateway.sol";
import { IInitializable } from "src/dispute/interfaces/IInitializable.sol";
import { IBigStepper, IPreimageOracle } from "src/dispute/interfaces/IBigStepper.sol";
import { IAnchorStateRegistry } from "src/dispute/interfaces/IAnchorStateRegistry.sol";
import { Clone } from "@solady/utils/Clone.sol";
import { ISemver } from "src/universal/ISemver.sol";
import "src/dispute/lib/Types.sol";
import "src/dispute/lib/Errors.sol";

/// @title FaultDisputeGame
/// @notice An implementation of the `IFaultDisputeGame` interface.
contract ZkFaultDisputeGame is IZkFaultDisputeGame, Clone, ISemver {
    ////////////////////////////////////////////////////////////////
    //                         State Vars                         //
    ////////////////////////////////////////////////////////////////

    /// @notice The maximum duration that may accumulate on a team's chess clock before they may no longer respond.
    Duration internal immutable MAX_CLOCK_DURATION;

    Duration internal immutable MAX_GENERATE_PROOF_DURATION;
    Duration internal immutable MAX_DETECT_FAULT_DURATION;

    /// @notice The game type ID.
    GameType internal immutable GAME_TYPE;

    /// @notice WETH contract for holding ETH.
    IDelayedWETH internal immutable WETH;

    /// @notice The anchor state registry.
    IAnchorStateRegistry internal immutable ANCHOR_STATE_REGISTRY;

    /// @notice The configuration contract for the zk fault proof system.
    ZkFaultProofConfig internal immutable CONFIG;

    /// @notice The chain ID of the L2 network this contract argues about.
    uint256 internal immutable L2_CHAIN_ID;

    /// @notice Semantic version.
    /// @custom:semver 1.0.0
    string public constant version = "1.0.0";

    /// @notice The starting timestamp of the game
    Timestamp public createdAt;

    /// @notice The timestamp of the game's global resolution.
    Timestamp public resolvedAt;

    /// @inheritdoc IDisputeGame
    GameStatus public status;

    /// @notice Flag for the `initialize` function to prevent re-initialization.
    bool internal initialized;

    /// @notice Credited balances for winning participants.
    mapping(address => uint256) public credit;

    /// @notice The mapping of claims which have been challenged by signal.
    mapping(uint256 => bool) public challengedClaims;
    /// @notice The indexes list of the challenged claims.
    uint256[] public challengedClaimIndexes;
    /// @notice The mapping of challenged claims and their timestamp.
    mapping(uint256 => uint64) public challengedClaimsTimestamp;
    /// @notice The mapping of claims and their challengers.
    mapping(uint256 => address payable) public challengers;
    /// @notice The fault prover who submitted the valid fault proof.
    address payable public faultProofProver;

    /// @notice The mapping of claims and their validity proof provers.
    mapping(uint256 => address payable) public validityProofProvers;
    /// @notice The mapping of claims which have been challenged by signal and proven to be invalid.
    mapping(uint256 => bool) public invalidChallengeClaims;
    /// @notice The indexes list of the invalid challenged claims.
    uint256[] public invalidChallengeClaimIndexes;

    /// @notice The flag for the challenge success.
    bool public isChallengeSuccess;
    /// @notice The index of the successful challenge.
    uint256 public successfulChallengeIndex;

    /// @notice The latest finalized output root, serving as the anchor for output bisection.
    OutputRoot public startingOutputRoot;

    uint256 public constant PERCENTAGE_DIVISOR = 10000;
    uint256 public immutable PROPOSER_BOND;
    uint256 public immutable CHALLENGER_BOND;
    address payable public immutable FEE_VAULT_ADDRESS;
    uint256 public immutable CHALLENGER_REWARD_PERCENTAGE;
    uint256 public immutable PROVER_REWARD_PERCENTAGE;

    /// @param _gameType The type ID of the game.
    /// @param _maxGenerateProofDuration The maximum amount of time that a validity prover has to generate a proof.
    /// @param _maxDetectFaultDuration The maximum amount of time that a challenger has to detect a fault and sunmit a signal.
    /// @param _weth WETH contract for holding ETH.
    /// @param _anchorStateRegistry The contract that stores the anchor state for each game type.
    /// @param _config The configuration contract for the zk fault proof system.
    /// @param _l2ChainId Chain ID of the L2 network this contract argues about.
    constructor(
        GameType _gameType,
        Duration _maxGenerateProofDuration,
        Duration _maxDetectFaultDuration,
        uint256 _PROPOSER_BOND,
        uint256 _CHALLENGER_BOND,
        address payable _FEE_VAULT_ADDRESS,
        uint256 _CHALLENGER_REWARD_PERCENTAGE,
        uint256 _PROVER_REWARD_PERCENTAGE,
        IDelayedWETH _weth,
        IAnchorStateRegistry _anchorStateRegistry,
        ZkFaultProofConfig _config,
        uint256 _l2ChainId

    ) {
        GAME_TYPE = _gameType;
        MAX_GENERATE_PROOF_DURATION = _maxGenerateProofDuration;
        MAX_DETECT_FAULT_DURATION = _maxDetectFaultDuration;
        PROPOSER_BOND = _PROPOSER_BOND;
        CHALLENGER_BOND = _CHALLENGER_BOND;
        FEE_VAULT_ADDRESS = _FEE_VAULT_ADDRESS;
        CHALLENGER_REWARD_PERCENTAGE = _CHALLENGER_REWARD_PERCENTAGE;
        PROVER_REWARD_PERCENTAGE = _PROVER_REWARD_PERCENTAGE;
        if (CHALLENGER_REWARD_PERCENTAGE + PROVER_REWARD_PERCENTAGE > PERCENTAGE_DIVISOR) {
            revert InvalidRewardPercentage();
        }
        MAX_CLOCK_DURATION = Duration.wrap(MAX_GENERATE_PROOF_DURATION.raw() + MAX_DETECT_FAULT_DURATION.raw());
        WETH = _weth;
        ANCHOR_STATE_REGISTRY = _anchorStateRegistry;
        CONFIG = _config;
        L2_CHAIN_ID = _l2ChainId;
    }

    /// @inheritdoc IInitializable
    function initialize() public payable virtual {
        // SAFETY: Any revert in this function will bubble up to the DisputeGameFactory and
        // prevent the game from being created.
        //
        // Implicit assumptions:
        // - The `gameStatus` state variable defaults to 0, which is `GameStatus.IN_PROGRESS`
        // - The dispute game factory will enforce the required bond to initialize the game.
        //
        // Explicit checks:
        // - The game must not have already been initialized.
        // - An output root cannot be proposed at or before the starting block number.

        // INVARIANT: The game must not have already been initialized.
        if (initialized) revert AlreadyInitialized();
        // if (msg.value != PROPOSER_BOND) revert IncorrectBondAmount();
        IZkFaultDisputeGame parentGameContract = parentGameProxy();
        // Grab the latest anchor root.
        (Hash root, uint256 rootBlockNumber) = ANCHOR_STATE_REGISTRY.anchors(GAME_TYPE);

        if (address(parentGameContract) != address(0)) {
            root = Hash.wrap(parentGameContract.rootClaim().raw());
            rootBlockNumber = parentGameContract.l2BlockNumber();
            if (parentGameStatus() == GameStatus.CHALLENGER_WINS) {
                revert ParentGameIsInvalid();
            }
        }

        // Should only happen if this is a new game type that hasn't been set up yet.
        if (root.raw() == bytes32(0)) revert AnchorRootNotFound();

        // Set the starting output root.
        startingOutputRoot = OutputRoot({ l2BlockNumber: rootBlockNumber, root: root });

        // Revert if the calldata size is not the expected length.
        //
        // This is to prevent adding extra or omitting bytes from to `extraData` that result in a different game UUID
        // in the factory, but are not used by the game, which would allow for multiple dispute games for the same
        // output proposal to be created.
        //
        // Expected length: 0xb2
        // - 0x04 selector
        // - 0x14 creator address
        // - 0x20 root claim
        // - 0x20 l1 head
        // - 0x58 extraData
            // - 0x20 l2 block number
            // - 0x20 claims hash
            // - 0x04 claims length
            // - 0x14 parent game contract address
        // - 0x02 CWIA bytes
        assembly {
            if iszero(eq(calldatasize(), 0xb2)) {
                // Store the selector for `BadExtraData()` & revert
                mstore(0x00, 0x9824bdab)
                revert(0x1C, 0x04)
            }
        }

        // Do not allow the game to be initialized if the root claim corresponds to a block at or before the
        // configured starting block number.
        uint256 currentL2BlockNumber = l2BlockNumber();
        if (currentL2BlockNumber <= rootBlockNumber) revert UnexpectedRootClaim(rootClaim());

        // This statement also ensures the correctness of currentL2BlockNumber, The reason is as follows:
        // If the currentL2BlockNumber is wrong, then the proposer must provide a wrong list of claims.
        // The challenger can detect the mismatch between the l2 block number and claims.
        // This mismatch can be proved in the zk proof.
        if (CONFIG.blockDistance() * claimsLength() != currentL2BlockNumber - rootBlockNumber) {
            revert InvalidClaimsLength();
        }

        // Set the game as initialized.
        initialized = true;

        // Deposit the bond.
        WETH.deposit{ value: msg.value }();

        // Set the game's starting timestamp
        createdAt = Timestamp.wrap(uint64(block.timestamp));
    }

    function requireNotExpired(Duration _maxDuration, Timestamp _startDuration) internal view {
        if (block.timestamp - _startDuration.raw() > _maxDuration.raw()) {
            revert ClockTimeExceeded();
        }
    }

    function challengeByProof(uint256 _disputeClaimIndex, Claim _expectedClaim, Claim[] calldata _originalClaims, bytes calldata _proof) external override {
        if (status != GameStatus.IN_PROGRESS) revert GameNotInProgress();
        if (parentGameStatus() == GameStatus.CHALLENGER_WINS) revert ParentGameIsInvalid();
        requireNotExpired(MAX_CLOCK_DURATION, createdAt);
        if (isChallengeSuccess) revert GameChallengeSucceeded();

        // allow direct challenge even there is a signal challenge
        // check the validity of origin claims
        if (_originalClaims.length != claimsLength()) {
            revert InvalidClaimsLength();
        }
        if (keccak256(abi.encodePacked(_originalClaims)) != claimsHash().raw()) {
            revert InvalidOriginClaims();
        }

        // check the validity of _disputeClaimIndex
        if (_disputeClaimIndex >= _originalClaims.length) {
            revert InvalidDisputeClaimIndex();
        }

        if (_expectedClaim.raw() == _originalClaims[_disputeClaimIndex].raw()) {
            revert InvalidExpectedClaim();
        }
        Claim agreedClaim;
        if (_disputeClaimIndex == 0) {
            agreedClaim = Claim.wrap(startingOutputRoot.root.raw());
        } else {
            agreedClaim = _originalClaims[_disputeClaimIndex - 1];
        }

        uint256 claimBlockNumber = startingOutputRoot.l2BlockNumber + (_disputeClaimIndex+1) * CONFIG.blockDistance();
        AggregationOutputs memory publicValues = AggregationOutputs({
            l1Head: l1Head().raw(),
            l2PreRoot: agreedClaim.raw(),
            claimRoot: _expectedClaim.raw(),
            claimBlockNum: claimBlockNumber,
            chainId: L2_CHAIN_ID,
            rollupConfigHash: CONFIG.rollupConfigHash(),
            rangeVkeyCommitment: CONFIG.rangeVkeyCommitment()
        });

        SP1VerifierGateway verifierGateway = SP1VerifierGateway(CONFIG.verifierGateway());
        verifierGateway.verifyProof(CONFIG.aggregationVkey(), abi.encode(publicValues), _proof);

        // update dispute claim index status
        isChallengeSuccess = true;
        successfulChallengeIndex = _disputeClaimIndex;
        faultProofProver = payable(msg.sender);
    }

    function challengeBySignal(uint256 _disputeClaimIndex) payable external override {
        if (msg.value != CHALLENGER_BOND) revert IncorrectBondAmount();
        WETH.deposit{ value: msg.value }();
        if (status != GameStatus.IN_PROGRESS) revert GameNotInProgress();
        if (parentGameStatus() == GameStatus.CHALLENGER_WINS) revert ParentGameIsInvalid();
        requireNotExpired(MAX_DETECT_FAULT_DURATION, createdAt);
        if (isChallengeSuccess) revert GameChallengeSucceeded();
        if (challengedClaims[_disputeClaimIndex]) revert ClaimAlreadyChallenged();
        if (_disputeClaimIndex >= claimsLength()) {
            revert InvalidDisputeClaimIndex();
        }
        challengedClaims[_disputeClaimIndex] = true;
        challengedClaimsTimestamp[_disputeClaimIndex] = uint64(block.timestamp);
        challengers[_disputeClaimIndex] = payable(msg.sender);
        challengedClaimIndexes.push(_disputeClaimIndex);
    }

    function submitProofForSignal(uint256 _disputeClaimIndex, Claim[] calldata _originalClaims, bytes calldata _proof) external override {
        if (status != GameStatus.IN_PROGRESS) revert GameNotInProgress();
        if (parentGameStatus() == GameStatus.CHALLENGER_WINS) revert ParentGameIsInvalid();
        // if the dispute claim index is already proven to be valid, revert
        if (isChallengeSuccess) revert GameChallengeSucceeded();
        // if there is no signal challenge, revert
        if (!challengedClaims[_disputeClaimIndex]) revert ClaimNotChallenged();
        // if the dispute claim index is already proven to be invalid, revert
        if (invalidChallengeClaims[_disputeClaimIndex]) revert ChallengeAlreadyInvalid();
        // if the challenge signal is submitted after the MAX_GENERATE_PROOF_DURATION, revert
        requireNotExpired(MAX_GENERATE_PROOF_DURATION, Timestamp.wrap(challengedClaimsTimestamp[_disputeClaimIndex]));

        // check the validity of origin claims
        if (_originalClaims.length != claimsLength()) {
            revert InvalidClaimsLength();
        }
        if (keccak256(abi.encodePacked(_originalClaims)) != claimsHash().raw()) {
            revert InvalidOriginClaims();
        }

        // check the validity of _disputeClaimIndex
        if (_disputeClaimIndex >= _originalClaims.length) {
            revert InvalidDisputeClaimIndex();
        }

        Claim agreedClaim;
        if (_disputeClaimIndex == 0) {
            agreedClaim = Claim.wrap(startingOutputRoot.root.raw());
        } else {
            agreedClaim = _originalClaims[_disputeClaimIndex - 1];
        }
        uint256 claimBlockNumber = startingOutputRoot.l2BlockNumber + (_disputeClaimIndex+1) * CONFIG.blockDistance();
        AggregationOutputs memory publicValues = AggregationOutputs({
            l1Head: l1Head().raw(),
            l2PreRoot: agreedClaim.raw(),
            claimRoot: _originalClaims[_disputeClaimIndex].raw(),
            claimBlockNum: claimBlockNumber,
            chainId: L2_CHAIN_ID,
            rollupConfigHash: CONFIG.rollupConfigHash(),
            rangeVkeyCommitment: CONFIG.rangeVkeyCommitment()
        });

        SP1VerifierGateway verifierGateway = SP1VerifierGateway(CONFIG.verifierGateway());
        verifierGateway.verifyProof(CONFIG.aggregationVkey(), abi.encode(publicValues), _proof);

        validityProofProvers[_disputeClaimIndex] = payable(msg.sender);
        invalidChallengeClaims[_disputeClaimIndex] = true;
        invalidChallengeClaimIndexes.push(_disputeClaimIndex);
    }

    function resolveClaim() external override {
        if (status != GameStatus.IN_PROGRESS) revert GameNotInProgress();
        if (parentGameStatus() == GameStatus.CHALLENGER_WINS) revert ParentGameIsInvalid();
        if (isChallengeSuccess) revert GameChallengeSucceeded();

        bool findValidChallenge = false;
        for (uint256 i = 0; i < challengedClaimIndexes.length; i++) {
            if (invalidChallengeClaims[challengedClaimIndexes[i]]) {
                continue;
            }
            if (block.timestamp - challengedClaimsTimestamp[challengedClaimIndexes[i]] > MAX_GENERATE_PROOF_DURATION.raw()) {
                isChallengeSuccess = true;
                successfulChallengeIndex = challengedClaimIndexes[i];
                faultProofProver = challengers[successfulChallengeIndex];
                findValidChallenge = true;
                break;
            }
        }
        if (!findValidChallenge) {
            revert NoExpiredChallenges();
        }
        return;
    }

    /// @inheritdoc IZkFaultDisputeGame
    function startingBlockNumber() external view returns (uint256 startingBlockNumber_) {
        startingBlockNumber_ = startingOutputRoot.l2BlockNumber;
    }

    /// @inheritdoc IZkFaultDisputeGame
    function startingRootHash() external view returns (Hash startingRootHash_) {
        startingRootHash_ = startingOutputRoot.root;
    }

    ////////////////////////////////////////////////////////////////
    //                    `IDisputeGame` impl                     //
    ////////////////////////////////////////////////////////////////

    /// @inheritdoc IDisputeGame
    function resolve() external returns (GameStatus status_) {
        // INVARIANT: Resolution cannot occur unless the game is currently in progress.
        if (status != GameStatus.IN_PROGRESS) revert GameNotInProgress();
        GameStatus parentStatus = parentGameStatus();
        // parent game must be resolved
        if (parentStatus == GameStatus.IN_PROGRESS) {
            revert ParentGameNotResolved();
        }
        if (parentStatus == GameStatus.CHALLENGER_WINS) {
            status_ = GameStatus.CHALLENGER_WINS;
            // re-update fault proof prover
            faultProofProver = payable(parentGameProxy().gameWinner());
        }
        if (isChallengeSuccess) {
            status_ = GameStatus.CHALLENGER_WINS;
        }

        if (status_ != GameStatus.CHALLENGER_WINS) {
            // first check if the challenge window is expired
            if (block.timestamp - createdAt.raw() <= MAX_CLOCK_DURATION.raw()) {
                revert ClockNotExpired();
            }
            // then check if there is any remaining challenges
            for (uint256 i = 0; i < challengedClaimIndexes.length; i++) {
                if (!invalidChallengeClaims[challengedClaimIndexes[i]]) {
                    revert UnresolvedChallenges();
                }
            }
            status_ = GameStatus.DEFENDER_WINS;
        }

        uint256 currentContractBalance = WETH.balanceOf(address(this));
        if (status_ == GameStatus.CHALLENGER_WINS) {
            // refund valid challengers if there is any
            for (uint256 i = 0; i < challengedClaimIndexes.length; i++) {
                if (!invalidChallengeClaims[challengedClaimIndexes[i]]) {
                    // refund the bond
                    _distributeBond(challengers[challengedClaimIndexes[i]], CHALLENGER_BOND);
                    currentContractBalance = currentContractBalance - CHALLENGER_BOND;
                }
            }
            // TODO reward part of challengers bond to validity provers, current reward is zero
            uint256 initialBalance = currentContractBalance;
            // reward the special challenger who submitted the signal which is proven to be valid
            // 1. someone submitted a valid fault proof corresponding to the challenge index; or
            // 2. the generate proof window is expired and no one submitted a validity proof
            // If isChallengeSuccess is true, then the challenger exists;
            // If isChallengeSuccess is false, then it indicates that the parent game is CHALLENGER_WINS and
            // there is no successful challenge in the current game.
            if (isChallengeSuccess && parentStatus != GameStatus.CHALLENGER_WINS) {
                // there is a challenger who submmitted the dispute claim index by `challengeBySignal`
                uint256 challengerBond = (currentContractBalance * CHALLENGER_REWARD_PERCENTAGE) / PERCENTAGE_DIVISOR;
                if (challengedClaims[successfulChallengeIndex]) {
                    _distributeBond(challengers[successfulChallengeIndex], challengerBond);
                } else {
                    // if there is no challenger, then the challenger is the fault proof prover self
                    _distributeBond(faultProofProver, challengerBond);
                }
                currentContractBalance = currentContractBalance - challengerBond;
            }
            // reward the fault proof prover
            uint256 proverBond = (initialBalance * PROVER_REWARD_PERCENTAGE) / PERCENTAGE_DIVISOR;
            _distributeBond(faultProofProver, proverBond);
            currentContractBalance = currentContractBalance - proverBond;
        } else if (status_ == GameStatus.DEFENDER_WINS) {
            // reward part of challengers bond to validity provers
            for (uint256 i = 0; i < invalidChallengeClaimIndexes.length; i++) {
                uint256 proverBond = (CHALLENGER_BOND * PROVER_REWARD_PERCENTAGE) / PERCENTAGE_DIVISOR;
                _distributeBond(validityProofProvers[invalidChallengeClaimIndexes[i]], proverBond);
                currentContractBalance = currentContractBalance - proverBond;
            }
            // refund the bond to proposer
            _distributeBond(gameCreator(), PROPOSER_BOND);
            currentContractBalance = currentContractBalance - PROPOSER_BOND;
        } else {
            // sanity check
            revert InvalidGameStatus();
        }
        // transfer the rest
        _distributeBond(FEE_VAULT_ADDRESS, currentContractBalance);

        resolvedAt = Timestamp.wrap(uint64(block.timestamp));

        // Update the status and emit the resolved event, note that we're performing an assignment here.
        emit Resolved(status = status_);

        // Try to update the anchor state, this should not revert.
        ANCHOR_STATE_REGISTRY.tryUpdateAnchorState();
    }

    /// @inheritdoc IDisputeGame
    function gameType() public view override returns (GameType gameType_) {
        gameType_ = GAME_TYPE;
    }

    /// @inheritdoc IDisputeGame
    function gameCreator() public pure returns (address creator_) {
        creator_ = _getArgAddress(0x00);
    }

    /// @inheritdoc IDisputeGame
    function rootClaim() public pure returns (Claim rootClaim_) {
        rootClaim_ = Claim.wrap(_getArgBytes32(0x14));
    }

    /// @inheritdoc IDisputeGame
    function l1Head() public pure returns (Hash l1Head_) {
        l1Head_ = Hash.wrap(_getArgBytes32(0x34));
    }

    function claimsHash() public pure returns (Hash claimsHash_) {
        claimsHash_ = Hash.wrap(_getArgBytes32(0x74));
    }

    function claimsLength() public pure returns (uint256 claimsLength_) {
        claimsLength_ = uint256(_getArgUint32(0x94));
    }

    function parentGameProxy() public pure returns (IZkFaultDisputeGame parentGameProxy_) {
        parentGameProxy_ = IZkFaultDisputeGame(_getArgAddress(0x98));
    }

    /// @inheritdoc IDisputeGame
    function l2BlockNumber() public pure returns (uint256 l2BlockNumber_) {
        l2BlockNumber_ = uint256(_getArgUint256(0x54));
    }

    /// @inheritdoc IDisputeGame
    function extraData() public pure returns (bytes memory extraData_) {
        // The extra data starts at the second word within the cwia calldata and
        // is 60 bytes long.
        extraData_ = _getArgBytes(0x54, 0x58);
    }

    /// @inheritdoc IDisputeGame
    function gameData() external view returns (GameType gameType_, Claim rootClaim_, bytes memory extraData_) {
        gameType_ = gameType();
        rootClaim_ = rootClaim();
        extraData_ = extraData();
    }

    ////////////////////////////////////////////////////////////////
    //                          HELPERS                           //
    ////////////////////////////////////////////////////////////////

    /// @notice Pays out the bond of a claim to a given recipient.
    /// @param _recipient The recipient of the bond.
    /// @param _bond The bond to pay out.
    function _distributeBond(address _recipient, uint256 _bond) internal {
        // Increase the recipient's credit.
        credit[_recipient] += _bond;

        // Unlock the bond.
        WETH.unlock(_recipient, _bond);
    }

    ////////////////////////////////////////////////////////////////
    //                       MISC EXTERNAL                        //
    ////////////////////////////////////////////////////////////////

    /// @notice Claim the credit belonging to the recipient address.
    /// @param _recipient The owner and recipient of the credit.
    function claimCredit(address _recipient) external {
        // Remove the credit from the recipient prior to performing the external call.
        uint256 recipientCredit = credit[_recipient];
        credit[_recipient] = 0;

        // Revert if the recipient has no credit to claim.
        if (recipientCredit == 0) revert NoCreditToClaim();

        // Try to withdraw the WETH amount so it can be used here.
        WETH.withdraw(_recipient, recipientCredit);

        // Transfer the credit to the recipient.
        (bool success,) = _recipient.call{ value: recipientCredit }(hex"");
        if (!success) revert BondTransferFailed();
    }

    /// @notice Returns the max clock duration.
    function maxClockDuration() external view returns (Duration maxClockDuration_) {
        maxClockDuration_ = MAX_CLOCK_DURATION;
    }

    /// @notice Returns the maximum duration allowed for generating a proof.
    function maxGenerateProofDuration() external view returns (Duration maxGenerateProofDuration_) {
        maxGenerateProofDuration_ = MAX_GENERATE_PROOF_DURATION;
    }

    /// @notice Returns the maximum duration allowed for detecting a fault.
    function maxDetectFaultDuration() external view returns (Duration maxDetectFaultDuration_) {
        maxDetectFaultDuration_ = MAX_DETECT_FAULT_DURATION;
    }

    /// @notice Returns the WETH contract for holding ETH.
    function weth() external view returns (IDelayedWETH weth_) {
        weth_ = WETH;
    }

    /// @notice Returns the anchor state registry contract.
    function anchorStateRegistry() external view returns (IAnchorStateRegistry registry_) {
        registry_ = ANCHOR_STATE_REGISTRY;
    }

    /// @notice Returns Config contract for holding zk fault proof configuration.
    function config() external view returns (ZkFaultProofConfig config_) {
        config_ = CONFIG;
    }

    /// @notice Returns the chain ID of the L2 network this contract argues about.
    function l2ChainId() external view returns (uint256 l2ChainId_) {
        l2ChainId_ = L2_CHAIN_ID;
    }

    /// @notice Returns the fault prover for the game
    function gameWinner() external view returns (address gameWinner_) {
        if (status != GameStatus.CHALLENGER_WINS) {
            gameWinner_ = address(0);
        } else {
            gameWinner_ = faultProofProver;
        }
    }

    /// @notice Returns the parent game status
    /// if the parent game is 0x00...00, then the parent game is from anchor state registry
    /// and its status is always DEFENDER_WINS
    function parentGameStatus() internal view returns (GameStatus parentGameStatus_) {
        if (address(parentGameProxy()) == address(0)) {
            parentGameStatus_ = GameStatus.DEFENDER_WINS;
        } else {
            parentGameStatus_ = parentGameProxy().status();
        }
    }

}
