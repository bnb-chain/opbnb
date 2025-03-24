// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { LibClone } from "@solady/utils/LibClone.sol";
import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { ISemver } from "src/universal/ISemver.sol";

import { IDisputeGame } from "src/dispute/interfaces/IDisputeGame.sol";
import { IDisputeGameFactory } from "src/dispute/interfaces/IDisputeGameFactory.sol";

import "src/dispute/lib/Types.sol";
import "src/dispute/lib/Errors.sol";

/// @title DisputeGameFactory
/// @notice A factory contract for creating `IDisputeGame` contracts. All created dispute games are stored in both a
///         mapping and an append only array. The timestamp of the creation time of the dispute game is packed tightly
///         into the storage slot with the address of the dispute game to make offchain discoverability of playable
///         dispute games easier.
contract DisputeGameFactory is OwnableUpgradeable, IDisputeGameFactory, ISemver {
    /// @dev Allows for the creation of clone proxies with immutable arguments.
    using LibClone for address;

    /// @notice Semantic version.
    /// @custom:semver 1.0.0
    string public constant version = "1.0.0";

    /// @inheritdoc IDisputeGameFactory
    mapping(GameType => IDisputeGame) public gameImpls;

    /// @inheritdoc IDisputeGameFactory
    mapping(GameType => uint256) public initBonds;

    /// @notice Mapping of a hash of `gameType || rootClaim || extraData` to the deployed `IDisputeGame` clone (where
    //          `||` denotes concatenation).
    mapping(Hash => GameId) internal _disputeGames;

    /// @notice An append-only array of disputeGames that have been created. Used by offchain game solvers to
    ///         efficiently track dispute games.
    GameId[] internal _disputeGameList;

    /// @notice Constructs a new DisputeGameFactory contract.
    constructor() OwnableUpgradeable() {
        initialize(address(0));
    }

    /// @notice Initializes the contract.
    /// @param _owner The owner of the contract.
    function initialize(address _owner) public initializer {
        __Ownable_init();
        _transferOwnership(_owner);
    }

    /// @inheritdoc IDisputeGameFactory
    function gameCount() external view returns (uint256 gameCount_) {
        gameCount_ = _disputeGameList.length;
    }

    /// @inheritdoc IDisputeGameFactory
    function games(
        GameType _gameType,
        Claim _rootClaim,
        bytes calldata _extraData
    )
        external
        view
        returns (IDisputeGame proxy_, Timestamp timestamp_)
    {
        Hash uuid = getGameUUID(_gameType, _rootClaim, _extraData);
        (, Timestamp timestamp, address proxy) = _disputeGames[uuid].unpack();
        (proxy_, timestamp_) = (IDisputeGame(proxy), timestamp);
    }

    /// @inheritdoc IDisputeGameFactory
    function gameAtIndex(uint256 _index)
        external
        view
        returns (GameType gameType_, Timestamp timestamp_, IDisputeGame proxy_)
    {
        (GameType gameType, Timestamp timestamp, address proxy) = _disputeGameList[_index].unpack();
        (gameType_, timestamp_, proxy_) = (gameType, timestamp, IDisputeGame(proxy));
    }

    /// @inheritdoc IDisputeGameFactory
    function create(
        GameType _gameType,
        Claim _rootClaim,
        bytes calldata _extraData
    )
        external
        payable
        returns (IDisputeGame proxy_)
    {
        // Grab the implementation contract for the given `GameType`.
        IDisputeGame impl = gameImpls[_gameType];

        // If there is no implementation to clone for the given `GameType`, revert.
        if (address(impl) == address(0)) revert NoImplementation(_gameType);

        // If the required initialization bond is not met, revert.
        if (msg.value != initBonds[_gameType]) revert IncorrectBondAmount();

        // Get the hash of the parent block.
        bytes32 parentHash = blockhash(block.number - 1);

        // Clone the implementation contract and initialize it with the given parameters.
        //
        // CWIA Calldata Layout:
        // ┌──────────────┬────────────────────────────────────┐
        // │    Bytes     │            Description             │
        // ├──────────────┼────────────────────────────────────┤
        // │ [0, 20)      │ Game creator address               │
        // │ [20, 52)     │ Root claim                         │
        // │ [52, 84)     │ Parent block hash at creation time │
        // │ [84, 84 + n) │ Extra data (opaque)                │
        // └──────────────┴────────────────────────────────────┘
        proxy_ = IDisputeGame(address(impl).clone(abi.encodePacked(msg.sender, _rootClaim, parentHash, _extraData)));
        proxy_.initialize{ value: msg.value }();

        // Compute the unique identifier for the dispute game.
        Hash uuid = getGameUUID(_gameType, _rootClaim, _extraData);

        // If a dispute game with the same UUID already exists, revert.
        if (GameId.unwrap(_disputeGames[uuid]) != bytes32(0)) revert GameAlreadyExists(uuid);

        // Pack the game ID.
        GameId id = LibGameId.pack(_gameType, Timestamp.wrap(uint64(block.timestamp)), address(proxy_));

        // Store the dispute game id in the mapping & emit the `DisputeGameCreated` event.
        _disputeGames[uuid] = id;
        _disputeGameList.push(id);
        emit DisputeGameCreated(address(proxy_), _gameType, _rootClaim);
    }

    function _prepareExtraData(
        uint256 l2BlockNumber,
        Claim[] memory claims,
        address parentProxy,
        bytes memory extraData
    ) private pure returns (bytes memory) {
        // calculate the hash of all _claims
        bytes32 claimsHash = keccak256(abi.encodePacked(claims));
        return abi.encodePacked(
            l2BlockNumber,
            claimsHash,
            uint32(claims.length),
            parentProxy,
            extraData
        );
    }

    function _validateGameCreation(
        GameType _gameType,
        Claim[] calldata _claims,
        uint64 _parentGameIndex,
        IDisputeGame impl
    ) private view returns (IDisputeGame parentProxy) {
        // If there is no implementation to clone for the given `GameType`, revert.
        if (address(impl) == address(0)) revert NoImplementation(_gameType);
        // If the required initialization bond is not met, revert.
        if (msg.value != initBonds[_gameType]) revert IncorrectBondAmount();
        if (_claims.length == 0) revert NoClaims();
        // When the parent game index is the maximum value of a uint64, it means that there is no parent game.
        // And use the state from anchor state registry.
        if (_parentGameIndex != type(uint64).max && _parentGameIndex >= _disputeGameList.length) {
            revert InvalidParentGameIndex();
        }
        if (_parentGameIndex != type(uint64).max) {
            (GameType parentGameType, , address proxy) = _disputeGameList[_parentGameIndex].unpack();
            parentProxy = IDisputeGame(proxy);
            // The parent game must be of the same type as the child game.
            if (parentGameType.raw() != _gameType.raw()) revert InvalidParentGameType();
            if (parentProxy.status() == GameStatus.CHALLENGER_WINS) {
                revert InvalidParentGameStatus();
            }
        } else {
            parentProxy = IDisputeGame(address(0));
        }
    }

    function createZkFaultDisputeGame(
        GameType _gameType,
        Claim[] calldata _claims,
        uint64 _parentGameIndex,
        uint64 _l2BlockNumber,
        bytes calldata _extraData
    ) external
        payable
        returns (IDisputeGame proxy_)
    {
        // Grab the implementation contract for the given `GameType`.
        IDisputeGame impl = gameImpls[_gameType];
        IDisputeGame parentProxy = _validateGameCreation(_gameType, _claims, _parentGameIndex, impl);
        // Clone the implementation contract and initialize it with the given parameters.
        //
        // CWIA Calldata Layout:
        // ┌───────────────────────────┬────────────────────────────────────┐
        // │    Bytes                  │            Description             │
        // ├───────────────────────────┼────────────────────────────────────┤
        // │ [0, 20)                   │ Game creator address               │
        // | [20, 52)                  | Root claim                         |
        // │ [52, 84)                  │ Parent block hash at creation time │
        // │ [84, 84+n)                │ Extra data                         │
        // |    [84, 116)              | L2 block number                    |
        // |    [116, 148)             | Hash of all claims                 |
        // |    [148, 152)             | The length of claims               |
        // |    [152, 172)             | Parent game address                |
        // └───────────────────────────┴────────────────────────────────────┘
        Claim rootClaim = _claims[_claims.length-1];
        bytes memory extraData = _prepareExtraData(
            _l2BlockNumber,
            _claims,
            address(parentProxy),
            _extraData
        );

        proxy_ = IDisputeGame(address(impl).clone(abi.encodePacked(msg.sender, rootClaim, blockhash(block.number - 1), extraData)));

        // Compute the unique identifier for the dispute game.
        Hash uuid = getGameUUID(_gameType, rootClaim, extraData);
        // If a dispute game with the same UUID already exists, revert.
        if (GameId.unwrap(_disputeGames[uuid]) != bytes32(0)) revert GameAlreadyExists(uuid);

        // Pack the game ID.
        GameId id = LibGameId.pack(_gameType, Timestamp.wrap(uint64(block.timestamp)), address(proxy_));

        // Store the dispute game id in the mapping & emit the `DisputeGameCreated` event.
        _disputeGames[uuid] = id;
        _disputeGameList.push(id);

        proxy_.initialize{ value: msg.value }();

        emit DisputeGameCreated(address(proxy_), _gameType, rootClaim);
        emit ZkDisputeGameIndexUpdated(_disputeGameList.length - 1);

    }
    /// @inheritdoc IDisputeGameFactory
    function getGameUUID(
        GameType _gameType,
        Claim _rootClaim,
        bytes memory _extraData
    )
        public
        pure
        returns (Hash uuid_)
    {
        uuid_ = Hash.wrap(keccak256(abi.encode(_gameType, _rootClaim, _extraData)));
    }

    /// @inheritdoc IDisputeGameFactory
    function findLatestGames(
        GameType _gameType,
        uint256 _start,
        uint256 _n
    )
        external
        view
        returns (GameSearchResult[] memory games_)
    {
        // If the `_start` index is greater than or equal to the game array length or `_n == 0`, return an empty array.
        if (_start >= _disputeGameList.length || _n == 0) return games_;

        // Allocate enough memory for the full array, but start the array's length at `0`. We may not use all of the
        // memory allocated, but we don't know ahead of time the final size of the array.
        assembly {
            games_ := mload(0x40)
            mstore(0x40, add(games_, add(0x20, shl(0x05, _n))))
        }

        // Perform a reverse linear search for the `_n` most recent games of type `_gameType`.
        for (uint256 i = _start; i >= 0 && i <= _start;) {
            GameId id = _disputeGameList[i];
            (GameType gameType, Timestamp timestamp, address proxy) = id.unpack();

            if (gameType.raw() == _gameType.raw()) {
                // Increase the size of the `games_` array by 1.
                // SAFETY: We can safely lazily allocate memory here because we pre-allocated enough memory for the max
                //         possible size of the array.
                assembly {
                    mstore(games_, add(mload(games_), 0x01))
                }

                bytes memory extraData = IDisputeGame(proxy).extraData();
                Claim rootClaim = IDisputeGame(proxy).rootClaim();
                games_[games_.length - 1] = GameSearchResult({
                    index: i,
                    metadata: id,
                    timestamp: timestamp,
                    rootClaim: rootClaim,
                    extraData: extraData
                });
                if (games_.length >= _n) break;
            }

            unchecked {
                i--;
            }
        }
    }

    /// @inheritdoc IDisputeGameFactory
    function setImplementation(GameType _gameType, IDisputeGame _impl) external onlyOwner {
        gameImpls[_gameType] = _impl;
        emit ImplementationSet(address(_impl), _gameType);
    }

    /// @inheritdoc IDisputeGameFactory
    function setInitBond(GameType _gameType, uint256 _initBond) external onlyOwner {
        initBonds[_gameType] = _initBond;
        emit InitBondUpdated(_gameType, _initBond);
    }
}
