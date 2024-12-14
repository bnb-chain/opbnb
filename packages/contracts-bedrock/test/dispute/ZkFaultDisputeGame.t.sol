// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test } from "forge-std/Test.sol";
import { Vm } from "forge-std/Vm.sol";

import { DisputeGameFactory_Init } from "test/dispute/DisputeGameFactory.t.sol";
import { ZkFaultDisputeGame, IDisputeGame} from "src/dispute/ZkFaultDisputeGame.sol";

import "src/dispute/lib/Types.sol";
import "src/dispute/lib/Errors.sol";
import { console2 as console} from "forge-std/Test.sol";


contract ZkFaultDisputeGame_Init is DisputeGameFactory_Init {
    /// @dev The type of the game being tested.
    GameType internal constant GAME_TYPE = GameType.wrap(3);

    /// @dev The implementation of the game being tested.
    ZkFaultDisputeGame internal gameImpl;
    /// @dev The `Clone` proxy of the game.
    ZkFaultDisputeGame internal gameProxy;

    // address immutable sp1VerifierGateway = 0x51d3960c929B27Db3f041eA3c3aD4fF3c2A121C7;
    // bytes32 immutable rollupConfigHash = 0x0b9b35ba1f4265979a10dea49f5501f81b3729b2856165574b1661323678e778;
    // bytes32 immutable aggregationVkey = 0x0006a81df67f2d5e48048edd5c051a5be0ef9720a2ce130f12f7021256160e73;
    // bytes32 immutable rangeVkeyCommitment = 0x5030974a2d74c494158e4af45836d72e2e0acae55f0f22d73c22bde90c1d6d98;

    bytes32 immutable l1BlockHash = 0x08837589aa817404bd09e35dd23b9efca9f3021809fa73c861d5eb9c2f41a06d;
    /// @dev The extra data passed to the game for initialization.
    bytes internal extraData;

    function init(Claim[] memory _claims, uint64 _parentGameIndex, uint64 _l2BlockNumber, bytes memory _extraData) public {

        gameImpl = new ZkFaultDisputeGame({
            _gameType: GAME_TYPE,
            _maxGenerateProofDuration: Duration.wrap(100),
            _maxDetectFaultDuration: Duration.wrap(100),
            _weth: delayedWeth,
            _anchorStateRegistry: anchorStateRegistry,
            _config: zkFaultProofConfig,
            _l2ChainId: 5611 // opbnb testnet chain id
        });

        // Register the game implementation.
        disputeGameFactory.setImplementation(GAME_TYPE, gameImpl);
        // Create a new game
        vm.setBlockhash(vm.getBlockNumber()-1, l1BlockHash);
        gameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame(
            GAME_TYPE, _claims, _parentGameIndex, _l2BlockNumber, _extraData)));

        // Check immutables
        assertEq(gameProxy.gameType().raw(), GAME_TYPE.raw());
        assertEq(gameProxy.maxGenerateProofDuration().raw(), 100);
        assertEq(gameProxy.maxDetectFaultDuration().raw(), 100);
        assertEq(gameProxy.maxClockDuration().raw(), 200);
        assertEq(address(gameProxy.anchorStateRegistry()), address(anchorStateRegistry));
        assertEq(address(gameProxy.config()), address(zkFaultProofConfig));

        // Label the proxy
        vm.label(address(gameProxy), "ZkFaultDisputeGame_Clone");
    }

}

/*
*** Use opbnb testnet block 47528965-47528980 to generate zk proof
*** https://testnet.opbnbscan.com/tx?block=47528956
*/

contract ZkFaultDisputeGame_Test is ZkFaultDisputeGame_Init {
    Claim[] internal claims;
    function setUp() public override {
        // uint256 id = vm.createFork("https://ancient-responsive-seed.bsc-testnet.quiknode.pro/989bbd40f106bddfe2e7f9c74e49644f28a6451e", 46571120);
        // vm.selectFork(id);
        super.setUp();
        OutputRoot memory outputRoot = OutputRoot({
            root: Hash.wrap(bytes32(0xeb184356c1e393bf4ab709254068cbb11ed723e45a5ff832b7394636786f52e2)),
            l2BlockNumber: 47528965
        });
        anchorStateRegistry.setAnchorState(GAME_TYPE, outputRoot);
        claims.push(Claim.wrap(bytes32(0xab7649c15939b35ace364954fd6b3859363784d3751d735c40581ade07808dc5)));
        for (uint256 i = 1; i < 5; i++) {
            claims.push(Claim.wrap(0));
        }
        bytes memory extraData;
        uint64 claimBlockNumber = 47528980;
        super.init(claims, type(uint64).max, claimBlockNumber, extraData);
    }

    function testChallengeBySignal() public {
        // Challenge the game
        uint256 disputeClaimIndex = 0;
        gameProxy.challengeBySignal{ value: 1 ether}(disputeClaimIndex);
        // Check the game status
        assertEq(gameProxy.challengedClaimIndexes(0), 0);

        // proof for block 47528965(agreed claim)-47528968(dispute claim)
        bytes memory _proof = hex"8653506623244edee8024efd99f19b806a8527b78291d9728516364ff8ee36728cd582500c68ca9edb5681591fb4ad0796ed6b983f8796c3b75f8bd5e360d61baa65191a1d63351599586537a801f816ae02b809be4a7b4076df97d36e48390b0ccce0e700dfe1e2be1ce304ed5e91b84f4914ee7b4a9c1f0cbb28b5dbb7ab7e1efaa5e4264449cf0efcf92eeeb51763f76c0232c47853e2d929a6af64754b0d23b90d791931dec9b23c02c0671d872a47d97c4c99050f78a48be7ab65aad140356b850210f70475b1388d112a18d7be71c39fba19efc10310d46f5d3f8ac1d4d2c118fe210d6483910c4ed03381ac18211238fe59f00b57d37c2e7146436dab82c544442c8188f22adba2be9fc50316a4a8888f03493b6e63a917c7742fbc53b764aca00e22058ecc7e4865eeac88007aa1d387a630522d5d1333e64023a76b71c8635d304fa2d0b95e5cb3aab490e215c086103290fd62ae9c3599bf244a15212f45fd0ff30621ac7a1d311695ec737e2fb306332a69fe73d9db110cf73a1d23630caf09c649fea0964cdd2855239803bb0421843eac61073090bbf456f7ad315ab6ea274fa6aa6375cada850c0081f30163c8ba95b433c6f53148a656fb934870ec192052b7c898af154eb1e51fa77361765c0244c285f667ce2c3297e7c67920060f2e823f1e94b5919cb386362ad0b752ba232fe734c982f83cfd6e1234d11dbe5b2cc612f3a915721531c4c5ae327605d73a999c43867d6aae6fc955f5de40b78f0bf6f436f3908dd11de20030558502d0f84be4abaf8058067f4ac439410751b22fdc89bcb749833f5ac39baaa29de3b8a58d660f19e5885775a693596f2adfec0b5f33b5113e56bf1409dc12ffc2698aae41a9515c98e5660de179ba463046951ea791ea82af0135da8f759b1c41b6098e07e9563a0263969472a7abb3f5f38424a44a135fe04846da2609ce77fedea0aa11bef7cbc1d505728df5d010b365890d507299095b152518fd4b00f7315d011aebb65b0893392b49c64d7fb782df9c1f9f04fbb757cfb9e0b72306e8162a45246b5f947093f99f4122ad64fc87d9131e7b6a212d7b8f19479ff8f0dd3edd8b7961783d0fcf009a085dab597504182e11610ca08ee877d657a768584dcbf8a366052a7d418703aa8e4744b14d8b905b1e2552638af0be86f0500d4a8bc563482d22c2b8fd9c4121898fc54415c6773b";
        gameProxy.submitProofForSignal(disputeClaimIndex, claims, _proof);
    }

    // function testChallengeByProof() public {
    //     bytes memory publicValue = hex"08837589aa817404bd09e35dd23b9efca9f3021809fa73c861d5eb9c2f41a06deb184356c1e393bf4ab709254068cbb11ed723e45a5ff832b7394636786f52e2cebef2969b201f8ec59652ca3fc07691251d685d3a6eeed153cd8473b92040e30000000000000000000000000000000000000000000000000000000002d53c1400000000000000000000000000000000000000000000000000000000000015eb0b9b35ba1f4265979a10dea49f5501f81b3729b2856165574b1661323678e7785030974a2d74c494158e4af45836d72e2e0acae55f0f22d73c22bde90c1d6d98";
    //     bytes memory proof = hex"feb5e54e0d9689d86160ae9a23b98aed031b5ec90623996a38e75379559262836d9a716b06694be58afee3bb8ec564db9ba6e2c93a66bfad2ce7b0a3f17c91c2e2a33967230bcadedf88d76675307b2d011b30ff948a54b48fd4e0bdf319cef399215a6d074b4e05cd9f847bed4f99772735e91cb2e8671aee3fd8580a06ea7e823367f8212463916c16a0e37038867ca437caf1d90974fa087741fac8668f70c72472fd23c10323826586613d0e4f00ba1580785e675cc6e294a76b50647b1d74ebc48d0c5b0d9e066770da0b3c83c65d2a952afcd90cf70bda18ba20f87a7e6762cce208d4a3ee267fc8856d21c0995d065214a2d336be77a0898ab629dbdc507b72ba";
    //     bytes32 memory
    // }
}