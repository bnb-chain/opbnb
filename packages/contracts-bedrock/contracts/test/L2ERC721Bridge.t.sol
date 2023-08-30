// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import { Messenger_Initializer } from "./CommonTest.t.sol";

// Target contract dependencies
import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import { L1ERC721Bridge } from "../L1/L1ERC721Bridge.sol";
import { OptimismMintableERC721 } from "../universal/OptimismMintableERC721.sol";

// Target contract
import { L2ERC721Bridge } from "../L2/L2ERC721Bridge.sol";

contract TestERC721 is ERC721 {
    constructor() ERC721("Test", "TST") {}

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }
}

contract TestMintableERC721 is OptimismMintableERC721 {
    constructor(address _bridge, address _remoteToken)
        OptimismMintableERC721(_bridge, 1, _remoteToken, "Test", "TST")
    {}

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }
}

contract L2ERC721Bridge_Test is Messenger_Initializer {
    TestMintableERC721 internal localToken;
    TestERC721 internal remoteToken;
    L2ERC721Bridge internal bridge;
    address internal constant otherBridge = address(0x3456);
    uint256 internal constant tokenId = 1;

    event ERC721BridgeInitiated(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 tokenId,
        bytes extraData
    );

    event ERC721BridgeFinalized(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 tokenId,
        bytes extraData
    );

    /// @dev Sets up the test suite.
    function setUp() public override {
        super.setUp();

        // Create necessary contracts.
        bridge = new L2ERC721Bridge(address(L2Messenger), otherBridge);
        remoteToken = new TestERC721();
        localToken = new TestMintableERC721(address(bridge), address(remoteToken));

        // Label the bridge so we get nice traces.
        vm.label(address(bridge), "L2ERC721Bridge");

        // Mint alice a token.
        localToken.mint(alice, tokenId);

        // Approve the bridge to transfer the token.
        vm.prank(alice);
        localToken.approve(address(bridge), tokenId);
    }

    /// @dev Tests that the constructor sets the correct variables.
    function test_constructor_succeeds() public {
        assertEq(address(bridge.MESSENGER()), address(L2Messenger));
        assertEq(address(bridge.OTHER_BRIDGE()), otherBridge);
        assertEq(address(bridge.messenger()), address(L2Messenger));
        assertEq(address(bridge.otherBridge()), otherBridge);
    }

    /// @dev Tests that `bridgeERC721` correctly bridges a token and
    ///      burns it on the origin chain.
    function test_bridgeERC721_succeeds() public {
        // Expect a call to the messenger.
        vm.expectCall(
            address(L2Messenger),
            abi.encodeCall(
                L2Messenger.sendMessage,
                (
                    address(otherBridge),
                    abi.encodeCall(
                        L2ERC721Bridge.finalizeBridgeERC721,
                        (
                            address(remoteToken),
                            address(localToken),
                            alice,
                            alice,
                            tokenId,
                            hex"5678"
                        )
                    ),
                    1234
                )
            )
        );

        // Expect an event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit ERC721BridgeInitiated(
            address(localToken),
            address(remoteToken),
            alice,
            alice,
            tokenId,
            hex"5678"
        );

        // Bridge the token.
        vm.prank(alice);
        bridge.bridgeERC721(address(localToken), address(remoteToken), tokenId, 1234, hex"5678");

        // Token is burned.
        vm.expectRevert("ERC721: invalid token ID");
        localToken.ownerOf(tokenId);
    }

    /// @dev Tests that `bridgeERC721` reverts if the owner is not an EOA.
    function test_bridgeERC721_fromContract_reverts() external {
        // Bridge the token.
        vm.etch(alice, hex"01");
        vm.prank(alice);
        vm.expectRevert("ERC721Bridge: account is not externally owned");
        bridge.bridgeERC721(address(localToken), address(remoteToken), tokenId, 1234, hex"5678");

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `bridgeERC721` reverts if the local token is the zero address.
    function test_bridgeERC721_localTokenZeroAddress_reverts() external {
        // Bridge the token.
        vm.prank(alice);
        vm.expectRevert();
        bridge.bridgeERC721(address(0), address(remoteToken), tokenId, 1234, hex"5678");

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `bridgeERC721` reverts if the remote token is the zero address.
    function test_bridgeERC721_remoteTokenZeroAddress_reverts() external {
        // Bridge the token.
        vm.prank(alice);
        vm.expectRevert("L2ERC721Bridge: remote token cannot be address(0)");
        bridge.bridgeERC721(address(localToken), address(0), tokenId, 1234, hex"5678");

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `bridgeERC721` reverts if the caller is not the token owner.
    function test_bridgeERC721_wrongOwner_reverts() external {
        // Bridge the token.
        vm.prank(bob);
        vm.expectRevert("L2ERC721Bridge: Withdrawal is not being initiated by NFT owner");
        bridge.bridgeERC721(address(localToken), address(remoteToken), tokenId, 1234, hex"5678");

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `bridgeERC721To` correctly bridges a token
    ///      and burns it on the origin chain.
    function test_bridgeERC721To_succeeds() external {
        // Expect a call to the messenger.
        vm.expectCall(
            address(L2Messenger),
            abi.encodeCall(
                L2Messenger.sendMessage,
                (
                    address(otherBridge),
                    abi.encodeCall(
                        L1ERC721Bridge.finalizeBridgeERC721,
                        (address(remoteToken), address(localToken), alice, bob, tokenId, hex"5678")
                    ),
                    1234
                )
            )
        );

        // Expect an event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit ERC721BridgeInitiated(
            address(localToken),
            address(remoteToken),
            alice,
            bob,
            tokenId,
            hex"5678"
        );

        // Bridge the token.
        vm.prank(alice);
        bridge.bridgeERC721To(
            address(localToken),
            address(remoteToken),
            bob,
            tokenId,
            1234,
            hex"5678"
        );

        // Token is burned.
        vm.expectRevert("ERC721: invalid token ID");
        localToken.ownerOf(tokenId);
    }

    /// @dev Tests that `bridgeERC721To` reverts if the local token is the zero address.
    function test_bridgeERC721To_localTokenZeroAddress_reverts() external {
        // Bridge the token.
        vm.prank(alice);
        vm.expectRevert();
        bridge.bridgeERC721To(address(0), address(remoteToken), bob, tokenId, 1234, hex"5678");

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `bridgeERC721To` reverts if the remote token is the zero address.
    function test_bridgeERC721To_remoteTokenZeroAddress_reverts() external {
        // Bridge the token.
        vm.prank(alice);
        vm.expectRevert("L2ERC721Bridge: remote token cannot be address(0)");
        bridge.bridgeERC721To(address(localToken), address(0), bob, tokenId, 1234, hex"5678");

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `bridgeERC721To` reverts if the caller is not the token owner.
    function test_bridgeERC721To_wrongOwner_reverts() external {
        // Bridge the token.
        vm.prank(bob);
        vm.expectRevert("L2ERC721Bridge: Withdrawal is not being initiated by NFT owner");
        bridge.bridgeERC721To(
            address(localToken),
            address(remoteToken),
            bob,
            tokenId,
            1234,
            hex"5678"
        );

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `finalizeBridgeERC721` correctly finalizes a bridged token.
    function test_finalizeBridgeERC721_succeeds() external {
        // Bridge the token.
        vm.prank(alice);
        bridge.bridgeERC721(address(localToken), address(remoteToken), tokenId, 1234, hex"5678");

        // Expect an event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit ERC721BridgeFinalized(
            address(localToken),
            address(remoteToken),
            alice,
            alice,
            tokenId,
            hex"5678"
        );

        // Finalize a withdrawal.
        vm.mockCall(
            address(L2Messenger),
            abi.encodeWithSelector(L2Messenger.xDomainMessageSender.selector),
            abi.encode(otherBridge)
        );
        vm.prank(address(L2Messenger));
        bridge.finalizeBridgeERC721(
            address(localToken),
            address(remoteToken),
            alice,
            alice,
            tokenId,
            hex"5678"
        );

        // Token is not locked in the bridge.
        assertEq(localToken.ownerOf(tokenId), alice);
    }

    /// @dev Tests that `finalizeBridgeERC721` reverts if the token is not compliant
    ///      with the `IOptimismMintableERC721` interface.
    function test_finalizeBridgeERC721_interfaceNotCompliant_reverts() external {
        // Create a non-compliant token
        NonCompliantERC721 nonCompliantToken = new NonCompliantERC721(alice);

        // Bridge the non-compliant token.
        vm.prank(alice);
        bridge.bridgeERC721(address(nonCompliantToken), address(0x01), tokenId, 1234, hex"5678");

        // Attempt to finalize the withdrawal. Should revert because the token does not claim
        // to be compliant with the `IOptimismMintableERC721` interface.
        vm.mockCall(
            address(L2Messenger),
            abi.encodeWithSelector(L2Messenger.xDomainMessageSender.selector),
            abi.encode(otherBridge)
        );
        vm.prank(address(L2Messenger));
        vm.expectRevert("L2ERC721Bridge: local token interface is not compliant");
        bridge.finalizeBridgeERC721(
            address(address(nonCompliantToken)),
            address(address(0x01)),
            alice,
            alice,
            tokenId,
            hex"5678"
        );
    }

    /// @dev Tests that `finalizeBridgeERC721` reverts when not called by the remote bridge.
    function test_finalizeBridgeERC721_notViaLocalMessenger_reverts() external {
        // Finalize a withdrawal.
        vm.prank(alice);
        vm.expectRevert("ERC721Bridge: function can only be called from the other bridge");
        bridge.finalizeBridgeERC721(
            address(localToken),
            address(remoteToken),
            alice,
            alice,
            tokenId,
            hex"5678"
        );
    }

    /// @dev Tests that `finalizeBridgeERC721` reverts when not called by the remote bridge.
    function test_finalizeBridgeERC721_notFromRemoteMessenger_reverts() external {
        // Finalize a withdrawal.
        vm.mockCall(
            address(L2Messenger),
            abi.encodeWithSelector(L2Messenger.xDomainMessageSender.selector),
            abi.encode(alice)
        );
        vm.prank(address(L2Messenger));
        vm.expectRevert("ERC721Bridge: function can only be called from the other bridge");
        bridge.finalizeBridgeERC721(
            address(localToken),
            address(remoteToken),
            alice,
            alice,
            tokenId,
            hex"5678"
        );
    }

    /// @dev Tests that `finalizeBridgeERC721` reverts when the local token is the
    ///      address of the bridge itself.
    function test_finalizeBridgeERC721_selfToken_reverts() external {
        // Finalize a withdrawal.
        vm.mockCall(
            address(L2Messenger),
            abi.encodeWithSelector(L2Messenger.xDomainMessageSender.selector),
            abi.encode(otherBridge)
        );
        vm.prank(address(L2Messenger));
        vm.expectRevert("L2ERC721Bridge: local token cannot be self");
        bridge.finalizeBridgeERC721(
            address(bridge),
            address(remoteToken),
            alice,
            alice,
            tokenId,
            hex"5678"
        );
    }

    /// @dev Tests that `finalizeBridgeERC721` reverts when already finalized.
    function test_finalizeBridgeERC721_alreadyExists_reverts() external {
        // Finalize a withdrawal.
        vm.mockCall(
            address(L2Messenger),
            abi.encodeWithSelector(L2Messenger.xDomainMessageSender.selector),
            abi.encode(otherBridge)
        );
        vm.prank(address(L2Messenger));
        vm.expectRevert("ERC721: token already minted");
        bridge.finalizeBridgeERC721(
            address(localToken),
            address(remoteToken),
            alice,
            alice,
            tokenId,
            hex"5678"
        );
    }
}

/// @dev A non-compliant ERC721 token that does not implement the full ERC721 interface.
///      This is used to test that the bridge will revert if the token does not claim to
///      support the ERC721 interface.
contract NonCompliantERC721 {
    address internal immutable owner;

    constructor(address _owner) {
        owner = _owner;
    }

    function ownerOf(uint256) external view returns (address) {
        return owner;
    }

    function remoteToken() external pure returns (address) {
        return address(0x01);
    }

    function burn(address, uint256) external {
        // Do nothing.
    }

    function supportsInterface(bytes4) external pure returns (bool) {
        return false;
    }
}
