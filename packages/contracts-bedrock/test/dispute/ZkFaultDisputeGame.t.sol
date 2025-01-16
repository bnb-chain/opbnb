// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test } from "forge-std/Test.sol";
import { Vm } from "forge-std/Vm.sol";

import { DisputeGameFactory_Init } from "test/dispute/DisputeGameFactory.t.sol";
import { ZkFaultDisputeGame, IDisputeGame} from "src/dispute/ZkFaultDisputeGame.sol";
import { AnchorStateRegistry } from "src/dispute/AnchorStateRegistry.sol";

import "src/dispute/lib/Types.sol";
import "src/dispute/lib/Errors.sol";
import { console2 as console} from "forge-std/Test.sol";


contract ZkFaultDisputeGame_Init is DisputeGameFactory_Init {
    /// @dev The type of the game being tested.
    GameType internal constant GAME_TYPE = GameType.wrap(3);
    Duration internal maxGenerateProofDuration = Duration.wrap(100);
    Duration internal maxDetectFaultDuration = Duration.wrap(100);
    uint256 internal constant challengerRewardPercentage = 1000;
    uint256 internal constant proverRewardPercentage = 5000;
    uint256 internal constant percentageDivisor = 10000;
    uint256 internal constant challengerBond = 1 ether;
    uint256 internal constant proposerBond = 1 ether;
    address payable internal constant feeVaultAddress = payable(address(0));
    /// @dev The implementation of the game being tested.
    ZkFaultDisputeGame internal gameImpl;

    // address immutable sp1VerifierGateway = 0x51d3960c929B27Db3f041eA3c3aD4fF3c2A121C7;
    // bytes32 immutable rollupConfigHash = 0x0b9b35ba1f4265979a10dea49f5501f81b3729b2856165574b1661323678e778;
    // bytes32 immutable aggregationVkey = 0x0006a81df67f2d5e48048edd5c051a5be0ef9720a2ce130f12f7021256160e73;
    // bytes32 immutable rangeVkeyCommitment = 0x5030974a2d74c494158e4af45836d72e2e0acae55f0f22d73c22bde90c1d6d98;

    // This hash is the block 46562356 hash of BSC testnet
    bytes32 immutable l1BlockHash = 0x08837589aa817404bd09e35dd23b9efca9f3021809fa73c861d5eb9c2f41a06d;
    /// @dev The extra data passed to the game for initialization.
    bytes internal extraData;

    function init() public {

        gameImpl = new ZkFaultDisputeGame({
            _gameType: GAME_TYPE,
            _maxGenerateProofDuration: maxGenerateProofDuration,
            _maxDetectFaultDuration: maxDetectFaultDuration,
            _PROPOSER_BOND: proposerBond,
            _CHALLENGER_BOND: challengerBond,
            _FEE_VAULT_ADDRESS: feeVaultAddress,
            _CHALLENGER_REWARD_PERCENTAGE: challengerRewardPercentage,
            _PROVER_REWARD_PERCENTAGE: proverRewardPercentage,
            _weth: delayedWeth,
            _anchorStateRegistry: anchorStateRegistry,
            _config: zkFaultProofConfig,
            _l2ChainId: 5611 // opbnb testnet chain id
        });

        // Register the game implementation.
        disputeGameFactory.setImplementation(GAME_TYPE, gameImpl);
        disputeGameFactory.setInitBond(GAME_TYPE, proposerBond);
    }

}

/*
*** Use opbnb testnet block 47528965-47528980 to generate zk proof
*** Two Games:
*** 1. Parent Game: 47528965-47528971
*** 2. Child Game : 47528971-47528980
*** https://testnet.opbnbscan.com/tx?block=47528956
*** (47528965, 47528968] proof: 8653506623244edee8024efd99f19b806a8527b78291d9728516364ff8ee36728cd582500c68ca9edb5681591fb4ad0796ed6b983f8796c3b75f8bd5e360d61baa65191a1d63351599586537a801f816ae02b809be4a7b4076df97d36e48390b0ccce0e700dfe1e2be1ce304ed5e91b84f4914ee7b4a9c1f0cbb28b5dbb7ab7e1efaa5e4264449cf0efcf92eeeb51763f76c0232c47853e2d929a6af64754b0d23b90d791931dec9b23c02c0671d872a47d97c4c99050f78a48be7ab65aad140356b850210f70475b1388d112a18d7be71c39fba19efc10310d46f5d3f8ac1d4d2c118fe210d6483910c4ed03381ac18211238fe59f00b57d37c2e7146436dab82c544442c8188f22adba2be9fc50316a4a8888f03493b6e63a917c7742fbc53b764aca00e22058ecc7e4865eeac88007aa1d387a630522d5d1333e64023a76b71c8635d304fa2d0b95e5cb3aab490e215c086103290fd62ae9c3599bf244a15212f45fd0ff30621ac7a1d311695ec737e2fb306332a69fe73d9db110cf73a1d23630caf09c649fea0964cdd2855239803bb0421843eac61073090bbf456f7ad315ab6ea274fa6aa6375cada850c0081f30163c8ba95b433c6f53148a656fb934870ec192052b7c898af154eb1e51fa77361765c0244c285f667ce2c3297e7c67920060f2e823f1e94b5919cb386362ad0b752ba232fe734c982f83cfd6e1234d11dbe5b2cc612f3a915721531c4c5ae327605d73a999c43867d6aae6fc955f5de40b78f0bf6f436f3908dd11de20030558502d0f84be4abaf8058067f4ac439410751b22fdc89bcb749833f5ac39baaa29de3b8a58d660f19e5885775a693596f2adfec0b5f33b5113e56bf1409dc12ffc2698aae41a9515c98e5660de179ba463046951ea791ea82af0135da8f759b1c41b6098e07e9563a0263969472a7abb3f5f38424a44a135fe04846da2609ce77fedea0aa11bef7cbc1d505728df5d010b365890d507299095b152518fd4b00f7315d011aebb65b0893392b49c64d7fb782df9c1f9f04fbb757cfb9e0b72306e8162a45246b5f947093f99f4122ad64fc87d9131e7b6a212d7b8f19479ff8f0dd3edd8b7961783d0fcf009a085dab597504182e11610ca08ee877d657a768584dcbf8a366052a7d418703aa8e4744b14d8b905b1e2552638af0be86f0500d4a8bc563482d22c2b8fd9c4121898fc54415c6773b
*** (47528968, 47528971] proof: 865350662799a0f21694ab08e1b809bc2a4e29bfced51b8c874c15854824230bb654e8090c79aad0adb9933f7bf7a9ae25835dc9f23036eeb91d2f07ac2b9db7bd0da148155d54c08125cd34d56cc3400fdfabe4869aa99f6ea55ce81abbbe960c88ee78290b65352cdfe26641ea2a473cb5034ba117b8a1771ccc60cd9ce7b2efc6eb5422f204c12f24beae751cc14214d485b16bb637416d98a4092726ddc47c7276f6215713698141591b79d3c978753442e6c8019c5e224d8a6d63bf7853a072d8a71871b64bf3f96508dfa74b4dd5416f1c4acd0a9780a907e00feec9b494e300e00fa3e79754d4574a279104e6f652afd4f329a467f5a4bfcd47ee8df07ef87eea19bcced91d1f4ddd4fa71af478c9aba42e219e3b27fbd3fe3dfd2eabbe703dd21a60fbeacccca280425c9cf9793d729db4023c7b5c58e23bc4b923789806fac82d272f19950ae0eccc3bc4e993b4319c15f69b69809ede167d1fb03b3c2e8629126e55f56fedddfbff41eab7b03365f1302afda1921a20b69778c1df0eb2987d0d11d70a01305b249a1806c0ebd2aca0b22b91ceb6a2cb4dc363a97edf216eac1d32e95fb9ab56349e82bccca0ae9ce6b7fb4aadfcab024e08d7585ddd875b782445f35f345683324c9a27db07f2e2fe0a55b4a187272f45aa76838e007ad8ac2ae51a6eb3448a6396d5420c2667f5d2f06e92107e894cb767c800ac1260972728870b1c346be18921b7c8768a71a759a09dcaa9f57a5a484f2b04afee6689570364d3285028eb7deebaac404cdf9ef700dcb7c50e3c7284222929ff881121ff1f419490b47d63403854b628483f51d649ef6ce8cef27a004eebcf04cc6f675c17fdee4478a23a95d8e38f429b2aa2fe24f44e0ad8d87c26a899e3f2757281782e51658e5a3724a2393dd715038bd0876ca78fe0474d730454e99d7f5f33780d1d85c555d6ed79a6cc645bc70968f81bedb7fa9ceea85a96b7229dd417e802d82271fae8f776bc05bb72f73f4e6a0256cacde3f9aa7e3ea3a9665b2c34abf0a80ef13896de49152e936346be7f8dd1536d8fe57878ae8864890323dbe79ff58a29d7f2632dd5a05d9d2589d28e2c0f99b3793533869b355f404f98c12a8bf84f1777426713c0e79ef5772b7416cf144c1a80b915b8da94d6b9dff3ee15dfbe1505a31c41ee062bb6b15a3324cf2fb929735178e59c29f65a16fda981204d4bb3
*** (47528971, 47528974] proof: 8653506626d363b85c9a5d9ca3ff305f93621efcad98edb97278bf91e3965f79221758e6305469b9f769d5b4f1af6fc36d9b906e125a6bbf0ca8429fc22875f98b12402d2d57ad93b081e9bdab5aea15d86d98077b833fa2527b0d71650520d56db9039f1a62897864cd400ec08b1cfe3f5e55cc1eb28f86f971e296dd4358bcb374683b232e1cc975181cf2d3e224095edf916f35d37922fa7e62a4af5ae564929feec62700f8435383ae63e5a2f4c8996d3578277251f6dd29978a10297b8f5b3b4bcb1e32a58fc3ed0b76092ed72e5c4c51deb31839555798a248e86175457840eb9f14226aa60c4380991dedaaa3026b8165b216e4d6a6b7ff330dbd906c200bc95c0296f130f46504cbd4f6fad92df66564b251482fcd7fefa37949e76d554163131e1f22cc5926ae46d652f545ad95520c0c92f7e36bacb03223fda191ee67d1432fb918dc85f0bc8ee7c6af391e95bbeea6df71161ef4153b74ea9bbbf9fdcf9927dedfca2d5896531b669966af2f351b1ad5b20260d94106b014b58ca37534cf1cd38d37f4765830b84a9b9460104555813decea22fe06d2100078708f5d54a00f6101a1a7c8dea39d6e5d97fa23c88c77919d0ad9999d5e8adf63c2d6a396290833ccd6d754481abbba5fa9c08506679bf9e38361e2ab3cd9e2986e2d9c9ec803c57a7d357535339fcebfac5fe320cf05ed86a61858b880005f4bd68fb5433207e97b0309d712116e80dafcac389670a3cc02d2012191ecaa459c030f96840a12f02fc9b06c9d4c8d4cb5fd3b0e7863f55db5862faf1e9e86040a5028541ac227b7e932794832756cc293efc8ee35feb66f25bcbf1201b3adfe03b4770de18d0c9040ec5e989285f9fe53761496b497a5dfbe822fc57ff24b13d4f1bc11712d0938a63efe5985bcdbfbae3b5f9d5f929655f24271b5523682814faeeb5a5dbf089a02906313ec148e4c91e0ab65e55f9b93e7ee244739d733b685461231eb46152c72f2cdd5dec91bfe510110c97bbf2f0ae9e61f6f4b707c894b6ba1c09224010f54aebd98da4473e93d91aef8aad009699c3307db44a3e465654567eeb2cc0fdb561bc244019b2cd044b1e04cc34f936ae1df9a396c1315f8307c1ed88df218572d4a53493b02286cd9ede0b59a5720b62706d4d7dcfa7b4730f099d314cd1142131a9daa5750e9c85b25f76f4f1b722ae62a457ceeb7e9cccb06d9b7f314
*** (47528974, 47528977] proof: 86535066010ba13283ec6fa24f4441ad332130b3ed51e587618909945a0b30d0313142b5168d69b26926f3d44b451ccc52f07f3bd7378ecb3c8fcf01d4a22e7f1813104e1c40aa55e8efcd32b7ada67eed9a6f1e140df00f4a71c9f04ab4b816147c9d75053be92c211dd62f3377faa5b043f8c62c13f18ab2109ea3b96fecc02b935990165bd5bf9970add163ed805137592182220f43d3b4c70233d8afe68151acef32154fd1d79b12ba7052316ef203a6a60eb9a7dc38084f220667def722f697ea582e8902c29ab8d5aa34d31e3e5c2c928e7dbb966af95af939b5ca855513def9062ea4ea31a831f552eb73a3e0d765f3e7f21e26f5ef484e6f4746d963dd24744506ff063b48dd2695f5d275fc0071dec47573b65ff291bea5d6e4f98b82d65eb821e1c0ee113cf95a8165518d9edc9f2590dad503bd04d86ca637b2f8742b02710fe99faa8879ca9035a64c6dfd9d5dc0d077ead0076d0844842b517aaac43ccb2257df9c2a90dfe3d580b661d0b3ea4b37ee77ab3138bd2d3d1ea8f026b24de51393de27de7a5715f8cee55ddab20713674191f7982b47399dc23520c7d36b98166bfdeb11347fc275a411d84618353e125631cfa2a7da21142b4a964b705f39222d66a0804574cd0734dfde380c42a1002f6133cca715e2ffe964dacd95839404d05be96d16a2c3a71c9c47de68bd4c378eb2ee0acd79e57e7206e5ee7ec01011ce2bd7e4f58c6777cdec5184688428e947808cac761c09897e0cd8c357766d28f7213e8cb4f7bcde040050617a08215343611fb30bcf05ddbec2c4d9bf6cd5181a3d3819353399cea935c1d6c21afb6925c9e57f945e886d8220183d1d3b562ef5838efead233faccb0665ec13bbb36d27c3497a6d52f8565b61d759315e072d233237a60c0cfafb822a3eeba2d1811aacf7188c6200b90b675bfd8b1fedcd2d32094de1a5c2509c2b67d3196edbedce7f2882d2af3467ce679588edc095231fa286b955f629b54846ee2d3e4a642a48d5715acc4f5462d140b237d4ab9449144e7e1ff6e2d07ffcb3a128a579aba7299e7a396735f7cc6e37e542f645cc3e152023dc3ac23f9d3aeaa048ae89b803ed1849d14efa6363d554931fddb0d61c14a302f38cb419d5cbaaa1145c379856910ec40e8bd67d1b617082988e01c48f270a6672603fd93ffbd92e18beba7145dfd29545db35e834d41a77bf7463da23
*** (47528977, 47528980] proof: 8653506606e0d6cd8bea8879f9eeecf333ee4f3ba9dfdef37e2f262a2ad58f2c616b54fd24da0dd7f6a9bec26ea8d1e8fb73af4efae6def6c9fb49d6e4900cad51b0c7ba21507acc7a5c1f680299c82812ec0d0f047fbdd2adbf39d3a24632e6dd18cef91e9b318af91f7ce31fba5391af8a9b2d37c897eb336de1e444662f96e5b419f00c3f494ea32d8d631c2bb032e06cfaf840745294f7ad71c4a0e31a0bc46e1b4d073978db4fc150b592c227ad36f651699f038b3558cbf8d195dfc541eb051f5618f209f176b4b0d76f184dc0a26c333ddc88d9705c504c9a91ed0745130d9ad2270f26db6a25bfd1e8d32003511f2e41b6eae60dd5a0aa5d014e137877187bef1fedca4f3bd1b7faa40a159ab2ead1dc44fc19aa97852fef867fd70ea1484c46053a3a86ecf11876875648c6bec10fad0074b47f8f8678f2934c49632e53624714b0dad48d4295c38e69fce3a6b75a17ff294acccf555444483e5ddd6f9ca59201fdf4528a2bec9efae36a4f918410cd142d5745286c22037d1259f41a90e3302f5c41e89a830fe9646dfd510b4324eeca5a574d59cad37043f436f8974c4c250dff99a6217df65135081b4cde00a52d334f575109976061cba4b662d34c08a317563edd216829caa5ff7c3a81beb236c794ca07dc8aff32826f0a6b27da6bd50737a94b0596047bf07f6d68c33f45b97fb90ded3395bcc6979d23392b62ddcb03026aacaf52035a157228c6b04c91c1ec3bc46d5647fede5dd390b95c4095002e8766cc9c205c7464ae68164015d0e962f2837ce827c6e4779e96c8f586c5252d20092677eb6afc9464a9d35ba6859289d06c2416189fd1ec67f9c6bf9195742c34d759fbcab908da1d30620cd30371f69ca2ae2d4fb0b03d518e8d31ed90231eefcc085ea31c34af5cf6a976e81f4d512a3cd71a64524e97949511c7d2e0cc00e99c5d6f93233506c38f85500b6bc1ad9b9acf413f2a232e469b38dd7fd1af00dfa8bd318a22945da83bb4bda7a08c5cd52b3172f5388077c49f3df7471ad01d90204ec131348a2996f628bc7c568a12b1488bf99653b12e6d72365303589a08bfb045101ba8c4b204e58895a5f92a8332811dff8edf51f33a86a549077c8a0ecbf0a099f8a910d906e7c753deb521e1cd9d16a1a030ecf09658ddacd7d12a0f8e7693f8fb76b23af6b1ed25a9373e09b935b298cfb51411170acb7956e3fe
*/

contract ZkFaultDisputeGame_Test is ZkFaultDisputeGame_Init {
    Claim[] internal parentClaims;
    uint64 parentClaimBlockNumber = 47528971;
    uint256 parentCreateTimestamp = 1734663238;
    Claim[] internal childClaims;
    uint64 childClaimBlockNumber = 47528980;
    uint256 childCreateTimestamp = 1734666838;
    /// @dev The `Clone` proxy of the game.
    ZkFaultDisputeGame internal parentGameProxy;
    ZkFaultDisputeGame internal childGameProxy;
    bytes[] parentProofs;
    bytes[] childProofs;

    // used address
    address parentFaultProver;
    address parentValidityProver;
    address parentChallenger;
    address childValidityProver;
    address childFaultProver;
    address childChallenger;

    function setUp() public override {
        // uint256 id = vm.createFork("https://ancient-responsive-seed.bsc-testnet.quiknode.pro/989bbd40f106bddfe2e7f9c74e49644f28a6451e", 46571120);
        // vm.selectFork(id);
        super.setUp();
        OutputRoot memory outputRoot = OutputRoot({
            root: Hash.wrap(bytes32(0xeb184356c1e393bf4ab709254068cbb11ed723e45a5ff832b7394636786f52e2)),
            l2BlockNumber: 47528965
        });
        AnchorStateRegistry.StartingAnchorRoot[] memory _startingAnchorRoots = new AnchorStateRegistry.StartingAnchorRoot[](1);
        _startingAnchorRoots[0] = AnchorStateRegistry.StartingAnchorRoot({
            gameType: GAME_TYPE,
            outputRoot: outputRoot
        });
        vm.store(address(anchorStateRegistry),
            bytes32(0x0000000000000000000000000000000000000000000000000000000000000000),
            bytes32(0x0000000000000000000000000000000000000000000000000000000000000000));
        // anchorStateRegistry.setAnchorState(GAME_TYPE, outputRoot);
        anchorStateRegistry.initialize(_startingAnchorRoots);
        // block 47528965
        parentClaims.push(Claim.wrap(bytes32(0xab7649c15939b35ace364954fd6b3859363784d3751d735c40581ade07808dc5)));
        parentClaims.push(Claim.wrap(bytes32(0xb9b9a9fe5b39804effc888ec9acb23b2f3c5772c368a6297966df44e53c2ebae)));
        parentProofs.push(hex"8653506623244edee8024efd99f19b806a8527b78291d9728516364ff8ee36728cd582500c68ca9edb5681591fb4ad0796ed6b983f8796c3b75f8bd5e360d61baa65191a1d63351599586537a801f816ae02b809be4a7b4076df97d36e48390b0ccce0e700dfe1e2be1ce304ed5e91b84f4914ee7b4a9c1f0cbb28b5dbb7ab7e1efaa5e4264449cf0efcf92eeeb51763f76c0232c47853e2d929a6af64754b0d23b90d791931dec9b23c02c0671d872a47d97c4c99050f78a48be7ab65aad140356b850210f70475b1388d112a18d7be71c39fba19efc10310d46f5d3f8ac1d4d2c118fe210d6483910c4ed03381ac18211238fe59f00b57d37c2e7146436dab82c544442c8188f22adba2be9fc50316a4a8888f03493b6e63a917c7742fbc53b764aca00e22058ecc7e4865eeac88007aa1d387a630522d5d1333e64023a76b71c8635d304fa2d0b95e5cb3aab490e215c086103290fd62ae9c3599bf244a15212f45fd0ff30621ac7a1d311695ec737e2fb306332a69fe73d9db110cf73a1d23630caf09c649fea0964cdd2855239803bb0421843eac61073090bbf456f7ad315ab6ea274fa6aa6375cada850c0081f30163c8ba95b433c6f53148a656fb934870ec192052b7c898af154eb1e51fa77361765c0244c285f667ce2c3297e7c67920060f2e823f1e94b5919cb386362ad0b752ba232fe734c982f83cfd6e1234d11dbe5b2cc612f3a915721531c4c5ae327605d73a999c43867d6aae6fc955f5de40b78f0bf6f436f3908dd11de20030558502d0f84be4abaf8058067f4ac439410751b22fdc89bcb749833f5ac39baaa29de3b8a58d660f19e5885775a693596f2adfec0b5f33b5113e56bf1409dc12ffc2698aae41a9515c98e5660de179ba463046951ea791ea82af0135da8f759b1c41b6098e07e9563a0263969472a7abb3f5f38424a44a135fe04846da2609ce77fedea0aa11bef7cbc1d505728df5d010b365890d507299095b152518fd4b00f7315d011aebb65b0893392b49c64d7fb782df9c1f9f04fbb757cfb9e0b72306e8162a45246b5f947093f99f4122ad64fc87d9131e7b6a212d7b8f19479ff8f0dd3edd8b7961783d0fcf009a085dab597504182e11610ca08ee877d657a768584dcbf8a366052a7d418703aa8e4744b14d8b905b1e2552638af0be86f0500d4a8bc563482d22c2b8fd9c4121898fc54415c6773b");
        parentProofs.push(hex"865350662799a0f21694ab08e1b809bc2a4e29bfced51b8c874c15854824230bb654e8090c79aad0adb9933f7bf7a9ae25835dc9f23036eeb91d2f07ac2b9db7bd0da148155d54c08125cd34d56cc3400fdfabe4869aa99f6ea55ce81abbbe960c88ee78290b65352cdfe26641ea2a473cb5034ba117b8a1771ccc60cd9ce7b2efc6eb5422f204c12f24beae751cc14214d485b16bb637416d98a4092726ddc47c7276f6215713698141591b79d3c978753442e6c8019c5e224d8a6d63bf7853a072d8a71871b64bf3f96508dfa74b4dd5416f1c4acd0a9780a907e00feec9b494e300e00fa3e79754d4574a279104e6f652afd4f329a467f5a4bfcd47ee8df07ef87eea19bcced91d1f4ddd4fa71af478c9aba42e219e3b27fbd3fe3dfd2eabbe703dd21a60fbeacccca280425c9cf9793d729db4023c7b5c58e23bc4b923789806fac82d272f19950ae0eccc3bc4e993b4319c15f69b69809ede167d1fb03b3c2e8629126e55f56fedddfbff41eab7b03365f1302afda1921a20b69778c1df0eb2987d0d11d70a01305b249a1806c0ebd2aca0b22b91ceb6a2cb4dc363a97edf216eac1d32e95fb9ab56349e82bccca0ae9ce6b7fb4aadfcab024e08d7585ddd875b782445f35f345683324c9a27db07f2e2fe0a55b4a187272f45aa76838e007ad8ac2ae51a6eb3448a6396d5420c2667f5d2f06e92107e894cb767c800ac1260972728870b1c346be18921b7c8768a71a759a09dcaa9f57a5a484f2b04afee6689570364d3285028eb7deebaac404cdf9ef700dcb7c50e3c7284222929ff881121ff1f419490b47d63403854b628483f51d649ef6ce8cef27a004eebcf04cc6f675c17fdee4478a23a95d8e38f429b2aa2fe24f44e0ad8d87c26a899e3f2757281782e51658e5a3724a2393dd715038bd0876ca78fe0474d730454e99d7f5f33780d1d85c555d6ed79a6cc645bc70968f81bedb7fa9ceea85a96b7229dd417e802d82271fae8f776bc05bb72f73f4e6a0256cacde3f9aa7e3ea3a9665b2c34abf0a80ef13896de49152e936346be7f8dd1536d8fe57878ae8864890323dbe79ff58a29d7f2632dd5a05d9d2589d28e2c0f99b3793533869b355f404f98c12a8bf84f1777426713c0e79ef5772b7416cf144c1a80b915b8da94d6b9dff3ee15dfbe1505a31c41ee062bb6b15a3324cf2fb929735178e59c29f65a16fda981204d4bb3");

        childClaims.push(Claim.wrap(bytes32(0xd6c969048a04b51295d8bab9e4739e6684e0cb2c3ca8cc04ec21fb7b7e55dcc2)));
        childClaims.push(Claim.wrap(bytes32(0xc9797bac31da11798dd5b7c1b3b8e00c5b59d19fc830011dcc527532745a23ca)));
        childClaims.push(Claim.wrap(bytes32(0xcebef2969b201f8ec59652ca3fc07691251d685d3a6eeed153cd8473b92040e3)));
        childProofs.push(hex"8653506626d363b85c9a5d9ca3ff305f93621efcad98edb97278bf91e3965f79221758e6305469b9f769d5b4f1af6fc36d9b906e125a6bbf0ca8429fc22875f98b12402d2d57ad93b081e9bdab5aea15d86d98077b833fa2527b0d71650520d56db9039f1a62897864cd400ec08b1cfe3f5e55cc1eb28f86f971e296dd4358bcb374683b232e1cc975181cf2d3e224095edf916f35d37922fa7e62a4af5ae564929feec62700f8435383ae63e5a2f4c8996d3578277251f6dd29978a10297b8f5b3b4bcb1e32a58fc3ed0b76092ed72e5c4c51deb31839555798a248e86175457840eb9f14226aa60c4380991dedaaa3026b8165b216e4d6a6b7ff330dbd906c200bc95c0296f130f46504cbd4f6fad92df66564b251482fcd7fefa37949e76d554163131e1f22cc5926ae46d652f545ad95520c0c92f7e36bacb03223fda191ee67d1432fb918dc85f0bc8ee7c6af391e95bbeea6df71161ef4153b74ea9bbbf9fdcf9927dedfca2d5896531b669966af2f351b1ad5b20260d94106b014b58ca37534cf1cd38d37f4765830b84a9b9460104555813decea22fe06d2100078708f5d54a00f6101a1a7c8dea39d6e5d97fa23c88c77919d0ad9999d5e8adf63c2d6a396290833ccd6d754481abbba5fa9c08506679bf9e38361e2ab3cd9e2986e2d9c9ec803c57a7d357535339fcebfac5fe320cf05ed86a61858b880005f4bd68fb5433207e97b0309d712116e80dafcac389670a3cc02d2012191ecaa459c030f96840a12f02fc9b06c9d4c8d4cb5fd3b0e7863f55db5862faf1e9e86040a5028541ac227b7e932794832756cc293efc8ee35feb66f25bcbf1201b3adfe03b4770de18d0c9040ec5e989285f9fe53761496b497a5dfbe822fc57ff24b13d4f1bc11712d0938a63efe5985bcdbfbae3b5f9d5f929655f24271b5523682814faeeb5a5dbf089a02906313ec148e4c91e0ab65e55f9b93e7ee244739d733b685461231eb46152c72f2cdd5dec91bfe510110c97bbf2f0ae9e61f6f4b707c894b6ba1c09224010f54aebd98da4473e93d91aef8aad009699c3307db44a3e465654567eeb2cc0fdb561bc244019b2cd044b1e04cc34f936ae1df9a396c1315f8307c1ed88df218572d4a53493b02286cd9ede0b59a5720b62706d4d7dcfa7b4730f099d314cd1142131a9daa5750e9c85b25f76f4f1b722ae62a457ceeb7e9cccb06d9b7f314");
        childProofs.push(hex"86535066010ba13283ec6fa24f4441ad332130b3ed51e587618909945a0b30d0313142b5168d69b26926f3d44b451ccc52f07f3bd7378ecb3c8fcf01d4a22e7f1813104e1c40aa55e8efcd32b7ada67eed9a6f1e140df00f4a71c9f04ab4b816147c9d75053be92c211dd62f3377faa5b043f8c62c13f18ab2109ea3b96fecc02b935990165bd5bf9970add163ed805137592182220f43d3b4c70233d8afe68151acef32154fd1d79b12ba7052316ef203a6a60eb9a7dc38084f220667def722f697ea582e8902c29ab8d5aa34d31e3e5c2c928e7dbb966af95af939b5ca855513def9062ea4ea31a831f552eb73a3e0d765f3e7f21e26f5ef484e6f4746d963dd24744506ff063b48dd2695f5d275fc0071dec47573b65ff291bea5d6e4f98b82d65eb821e1c0ee113cf95a8165518d9edc9f2590dad503bd04d86ca637b2f8742b02710fe99faa8879ca9035a64c6dfd9d5dc0d077ead0076d0844842b517aaac43ccb2257df9c2a90dfe3d580b661d0b3ea4b37ee77ab3138bd2d3d1ea8f026b24de51393de27de7a5715f8cee55ddab20713674191f7982b47399dc23520c7d36b98166bfdeb11347fc275a411d84618353e125631cfa2a7da21142b4a964b705f39222d66a0804574cd0734dfde380c42a1002f6133cca715e2ffe964dacd95839404d05be96d16a2c3a71c9c47de68bd4c378eb2ee0acd79e57e7206e5ee7ec01011ce2bd7e4f58c6777cdec5184688428e947808cac761c09897e0cd8c357766d28f7213e8cb4f7bcde040050617a08215343611fb30bcf05ddbec2c4d9bf6cd5181a3d3819353399cea935c1d6c21afb6925c9e57f945e886d8220183d1d3b562ef5838efead233faccb0665ec13bbb36d27c3497a6d52f8565b61d759315e072d233237a60c0cfafb822a3eeba2d1811aacf7188c6200b90b675bfd8b1fedcd2d32094de1a5c2509c2b67d3196edbedce7f2882d2af3467ce679588edc095231fa286b955f629b54846ee2d3e4a642a48d5715acc4f5462d140b237d4ab9449144e7e1ff6e2d07ffcb3a128a579aba7299e7a396735f7cc6e37e542f645cc3e152023dc3ac23f9d3aeaa048ae89b803ed1849d14efa6363d554931fddb0d61c14a302f38cb419d5cbaaa1145c379856910ec40e8bd67d1b617082988e01c48f270a6672603fd93ffbd92e18beba7145dfd29545db35e834d41a77bf7463da23");
        childProofs.push(hex"8653506606e0d6cd8bea8879f9eeecf333ee4f3ba9dfdef37e2f262a2ad58f2c616b54fd24da0dd7f6a9bec26ea8d1e8fb73af4efae6def6c9fb49d6e4900cad51b0c7ba21507acc7a5c1f680299c82812ec0d0f047fbdd2adbf39d3a24632e6dd18cef91e9b318af91f7ce31fba5391af8a9b2d37c897eb336de1e444662f96e5b419f00c3f494ea32d8d631c2bb032e06cfaf840745294f7ad71c4a0e31a0bc46e1b4d073978db4fc150b592c227ad36f651699f038b3558cbf8d195dfc541eb051f5618f209f176b4b0d76f184dc0a26c333ddc88d9705c504c9a91ed0745130d9ad2270f26db6a25bfd1e8d32003511f2e41b6eae60dd5a0aa5d014e137877187bef1fedca4f3bd1b7faa40a159ab2ead1dc44fc19aa97852fef867fd70ea1484c46053a3a86ecf11876875648c6bec10fad0074b47f8f8678f2934c49632e53624714b0dad48d4295c38e69fce3a6b75a17ff294acccf555444483e5ddd6f9ca59201fdf4528a2bec9efae36a4f918410cd142d5745286c22037d1259f41a90e3302f5c41e89a830fe9646dfd510b4324eeca5a574d59cad37043f436f8974c4c250dff99a6217df65135081b4cde00a52d334f575109976061cba4b662d34c08a317563edd216829caa5ff7c3a81beb236c794ca07dc8aff32826f0a6b27da6bd50737a94b0596047bf07f6d68c33f45b97fb90ded3395bcc6979d23392b62ddcb03026aacaf52035a157228c6b04c91c1ec3bc46d5647fede5dd390b95c4095002e8766cc9c205c7464ae68164015d0e962f2837ce827c6e4779e96c8f586c5252d20092677eb6afc9464a9d35ba6859289d06c2416189fd1ec67f9c6bf9195742c34d759fbcab908da1d30620cd30371f69ca2ae2d4fb0b03d518e8d31ed90231eefcc085ea31c34af5cf6a976e81f4d512a3cd71a64524e97949511c7d2e0cc00e99c5d6f93233506c38f85500b6bc1ad9b9acf413f2a232e469b38dd7fd1af00dfa8bd318a22945da83bb4bda7a08c5cd52b3172f5388077c49f3df7471ad01d90204ec131348a2996f628bc7c568a12b1488bf99653b12e6d72365303589a08bfb045101ba8c4b204e58895a5f92a8332811dff8edf51f33a86a549077c8a0ecbf0a099f8a910d906e7c753deb521e1cd9d16a1a030ecf09658ddacd7d12a0f8e7693f8fb76b23af6b1ed25a9373e09b935b298cfb51411170acb7956e3fe");

        bytes memory extraData;
        super.init();
        // Create a new game
        address proposer = makeAddr("proposer");
        vm.deal(proposer, 100 ether);
        vm.startPrank(proposer);
        vm.setBlockhash(vm.getBlockNumber()-1, l1BlockHash);
        vm.warp(parentCreateTimestamp);
        parentGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond} (
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));
        vm.stopPrank();
        // Check immutables
        assertEq(parentGameProxy.gameType().raw(), GAME_TYPE.raw());
        assertEq(parentGameProxy.maxGenerateProofDuration().raw(), 100);
        assertEq(parentGameProxy.maxDetectFaultDuration().raw(), 100);
        assertEq(parentGameProxy.maxClockDuration().raw(), 200);
        assertEq(parentGameProxy.CHALLENGER_BOND(), 1 ether);
        assertEq(parentGameProxy.PROPOSER_BOND(), 1 ether);
        assertEq(parentGameProxy.FEE_VAULT_ADDRESS(), address(0));
        assertEq(parentGameProxy.CHALLENGER_REWARD_PERCENTAGE(), 1000);
        assertEq(parentGameProxy.PROVER_REWARD_PERCENTAGE(), 5000);
        assertEq(address(parentGameProxy.anchorStateRegistry()), address(anchorStateRegistry));
        assertEq(address(parentGameProxy.config()), address(zkFaultProofConfig));
        // Label the proxy
        vm.label(address(parentGameProxy), "ParentZkFaultDisputeGame_Clone");

        // Create a new game
        vm.startPrank(proposer);
        vm.warp(childCreateTimestamp);
        childGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond} (
            GAME_TYPE, childClaims, 0, childClaimBlockNumber, extraData)));
        vm.label(address(childGameProxy), "ChildZkFaultDisputeGame_Clone");
        vm.stopPrank();

        // used address
        parentFaultProver = makeAddr("parentFaultProver");
        parentChallenger = makeAddr("parentChallenger");
        parentValidityProver = makeAddr("parentValidityProver");
        childValidityProver = makeAddr("childValidityProver");
        childFaultProver = makeAddr("childFaultProver");
        childChallenger = makeAddr("challenger");

        vm.deal(parentFaultProver, 100 ether);
        vm.deal(parentChallenger, 100 ether);
        vm.deal(parentValidityProver, 100 ether);
        vm.deal(childValidityProver, 100 ether);
        vm.deal(childFaultProver, 100 ether);
        vm.deal(childChallenger, 100 ether);
    }

    function testChallengeBySignal() public {
        vm.warp(parentCreateTimestamp + 1);
        // Challenge the game
        uint256 disputeClaimIndex = 1;
        parentGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);
        // Check the game status
        assertEq(parentGameProxy.challengedClaimIndexes(0), disputeClaimIndex);

        // proof for block 47528968(agreed claim)-47528971(dispute claim)
        bytes memory proof = parentProofs[disputeClaimIndex];
        parentGameProxy.submitProofForSignal(disputeClaimIndex, parentClaims, proof);
        // Check the game status
        assertEq(parentGameProxy.validityProofProvers(disputeClaimIndex), address(this));
        assertEq(parentGameProxy.invalidChallengeClaims(disputeClaimIndex), true);
        assertEq(parentGameProxy.invalidChallengeClaimIndexes(0), disputeClaimIndex);

        vm.warp(childCreateTimestamp + 1);
        childGameProxy.challengeBySignal{ value: challengerBond}(disputeClaimIndex);
        // Check the game status
        assertEq(childGameProxy.challengedClaimIndexes(0), disputeClaimIndex);

    }

    function testChallengeBySignal_FailCases() public {
        uint256 disputeClaimIndex = 1;
        // Challenge the game with not enough bond
        vm.expectRevert(abi.encodeWithSelector(IncorrectBondAmount.selector));
        parentGameProxy.challengeBySignal{ value: challengerBond - 1 }(disputeClaimIndex);

        // Challenge the game after the max detect fault duration
        vm.warp(parentCreateTimestamp + parentGameProxy.maxDetectFaultDuration().raw() + 1);
        vm.expectRevert(abi.encodeWithSelector(ClockTimeExceeded.selector));
        parentGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);

        // Challenge the same dispute claim index again
        vm.warp(parentCreateTimestamp + 1);
        parentGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);
        vm.expectRevert(abi.encodeWithSelector(ClaimAlreadyChallenged.selector));
        parentGameProxy.challengeBySignal{ value: challengerBond}(disputeClaimIndex);

        // Challenge the game with invalid claim index
        vm.expectRevert(abi.encodeWithSelector(InvalidDisputeClaimIndex.selector));
        parentGameProxy.challengeBySignal{ value: challengerBond }(parentClaims.length);

        // construct invalid parent game
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        bytes memory extraData;
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000003));
        ZkFaultDisputeGame invalidParentGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));

        ZkFaultDisputeGame invalidChildGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, childClaims, uint64(disputeGameFactory.gameCount()-1), childClaimBlockNumber, extraData)));
        // parent game is resolved to CHALLENGER_WINS
        vm.startPrank(parentFaultProver);
        invalidParentGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, parentProofs[disputeClaimIndex]);
        invalidParentGameProxy.resolve();
        vm.stopPrank();

        // Challenge the game which is already resolved
        vm.expectRevert(abi.encodeWithSelector(GameNotInProgress.selector));
        invalidParentGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);

        // Challenge the game with invalid parent game
        vm.expectRevert(abi.encodeWithSelector(ParentGameIsInvalid.selector));
        invalidChildGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);

    }

    function testChallengeByProof() public {
        // construct invalid games
        uint256 disputeClaimIndex = 1;
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000001));
        bytes memory extraData;
        ZkFaultDisputeGame invalidGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));
        // Challenge the game
        bytes memory proof = parentProofs[disputeClaimIndex];
        vm.prank(msg.sender);
        invalidGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, proof);
        assertEq(invalidGameProxy.isChallengeSuccess(), true);
        assertEq(invalidGameProxy.successfulChallengeIndex(), disputeClaimIndex);
        assertEq(invalidGameProxy.faultProofProver(), msg.sender);

        uint64 parentGameIndex = uint64(disputeGameFactory.gameCount()-1);
        console.log("parent game index is %d", parentGameIndex);
        // create a new child game
        ZkFaultDisputeGame invalidChildGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, childClaims, parentGameIndex, childClaimBlockNumber, extraData)));

        // resolve the parent game
        invalidGameProxy.resolve();
        assertEq(uint256(invalidGameProxy.status()), uint256(GameStatus.CHALLENGER_WINS));
        vm.expectRevert(abi.encodeWithSelector(ParentGameIsInvalid.selector));
        invalidChildGameProxy.resolveClaim();
    }

    function testChallengeByProof_Child_FirstClaim() public {
        uint256 disputeClaimIndex = 0;
        // The valid parent game index which is created in the setup function
        uint64 parentGameIndex = 0;
        Claim expectedClaim = childClaims[disputeClaimIndex];
        childClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000001));
        ZkFaultDisputeGame invalidChildProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, childClaims, parentGameIndex, childClaimBlockNumber, new bytes(0))));
        // block (47528974, 47528977] proof
        bytes memory proof = childProofs[disputeClaimIndex];
        vm.startPrank(childFaultProver);
        invalidChildProxy.challengeByProof(disputeClaimIndex, expectedClaim, childClaims, proof);
        assertEq(invalidChildProxy.isChallengeSuccess(), true);
        assertEq(invalidChildProxy.successfulChallengeIndex(), disputeClaimIndex);
        assertEq(invalidChildProxy.faultProofProver(), childFaultProver);

        // another dispute claim index: first challenge by signal, then challenge by proof
        uint256 secondDisputeClaimIndex = 2;
        Claim secondExpectedClaim = childClaims[secondDisputeClaimIndex];
        childClaims[secondDisputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000001));
        ZkFaultDisputeGame secondInvalidChildProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, childClaims, parentGameIndex, childClaimBlockNumber, extraData)));
        secondInvalidChildProxy.challengeBySignal{ value: challengerBond }(secondDisputeClaimIndex);
        secondInvalidChildProxy.challengeByProof(secondDisputeClaimIndex, secondExpectedClaim, childClaims, childProofs[secondDisputeClaimIndex]);
        assertEq(secondInvalidChildProxy.isChallengeSuccess(), true);
        assertEq(secondInvalidChildProxy.successfulChallengeIndex(), secondDisputeClaimIndex);
        assertEq(secondInvalidChildProxy.faultProofProver(), childFaultProver);
        vm.stopPrank();
    }

    function testChallengeByProof_FailCases() public {
        uint256 disputeClaimIndex = 1;
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        Claim[] memory invalidClaims = new Claim[](parentClaims.length + 1);
        for (uint256 i = 0; i < parentClaims.length; i++) {
            invalidClaims[i] = parentClaims[i];
        }
        bytes memory extraData;
        vm.warp(parentCreateTimestamp + 1);
        vm.expectRevert(abi.encodeWithSelector(InvalidClaimsLength.selector));
        parentGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, invalidClaims, parentProofs[disputeClaimIndex]);

        Claim[] memory invalidClaims2 = parentClaims;
        invalidClaims2[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000005));
        vm.expectRevert(abi.encodeWithSelector(InvalidOriginClaims.selector));
        parentGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, invalidClaims2, parentProofs[disputeClaimIndex]);

        vm.expectRevert(abi.encodeWithSelector(InvalidDisputeClaimIndex.selector));
        parentGameProxy.challengeByProof(parentClaims.length, expectedClaim, parentClaims, parentProofs[disputeClaimIndex]);

        vm.expectRevert(abi.encodeWithSelector(InvalidExpectedClaim.selector));
        parentGameProxy.challengeByProof(disputeClaimIndex, parentClaims[disputeClaimIndex], parentClaims, parentProofs[disputeClaimIndex]);

        vm.warp(parentCreateTimestamp + parentGameProxy.maxClockDuration().raw() + 1);
        vm.expectRevert(abi.encodeWithSelector(ClockTimeExceeded.selector));
        parentGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, parentProofs[disputeClaimIndex]);

        vm.warp(parentCreateTimestamp + 1);
        // 0x09bde339 is sig of InvalidProof()
        vm.expectRevert(abi.encodeWithSelector(0x09bde339));
        parentGameProxy.challengeByProof(disputeClaimIndex, invalidClaims2[disputeClaimIndex], parentClaims, parentProofs[disputeClaimIndex]);

        // construct invalid parent game
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000005));
        ZkFaultDisputeGame invalidGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));

        ZkFaultDisputeGame invalidChildGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, childClaims, uint64(disputeGameFactory.gameCount()-1), childClaimBlockNumber, extraData)));

        vm.startPrank(parentFaultProver);
        invalidGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, parentProofs[disputeClaimIndex]);
        invalidGameProxy.resolve();
        vm.stopPrank();

        vm.expectRevert(abi.encodeWithSelector(ParentGameIsInvalid.selector));
        invalidChildGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, childClaims, childProofs[disputeClaimIndex]);
    }

    function testResolveClaim() public {
        uint256 disputeClaimIndex = 1;
        uint256 disputeTimestamp = parentCreateTimestamp + 10;
        vm.warp(disputeTimestamp);
        vm.prank(msg.sender);
        parentGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);

        uint256 secondDisputeClaimIndex = 0;
        uint256 secondDisputeTimestamp = parentCreateTimestamp + 20;
        vm.warp(secondDisputeTimestamp);
        parentGameProxy.challengeBySignal{ value: challengerBond }(secondDisputeClaimIndex);

        Duration maxGenerateProofDuration = parentGameProxy.maxGenerateProofDuration();
        vm.warp(disputeTimestamp + maxGenerateProofDuration.raw());
        vm.expectRevert(abi.encodeWithSelector(NoExpiredChallenges.selector));
        parentGameProxy.resolveClaim();

        vm.warp(disputeTimestamp + maxGenerateProofDuration.raw() + 1);
        parentGameProxy.resolveClaim();
        assertEq(parentGameProxy.isChallengeSuccess(), true);
        assertEq(parentGameProxy.successfulChallengeIndex(), disputeClaimIndex);
        assertEq(parentGameProxy.faultProofProver(), msg.sender);
    }

    function testResolveClaim_FailCases() public {
        uint256 disputeClaimIndex = 1;
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000005));
        vm.warp(parentCreateTimestamp);
        ZkFaultDisputeGame invalidGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame { value: proposerBond } (
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));

        uint256 disputeTimestamp = parentCreateTimestamp + 10;
        vm.warp(disputeTimestamp);
        vm.prank(msg.sender);
        parentGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);

        Duration maxGenerateProofDuration = parentGameProxy.maxGenerateProofDuration();
        vm.warp(disputeTimestamp + maxGenerateProofDuration.raw());
        vm.expectRevert(abi.encodeWithSelector(NoExpiredChallenges.selector));
        parentGameProxy.resolveClaim();

        vm.startPrank(parentFaultProver);
        invalidGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, parentProofs[disputeClaimIndex]);
        vm.expectRevert(abi.encodeWithSelector(GameChallengeSucceeded.selector));
        invalidGameProxy.resolveClaim();
        vm.stopPrank();

        invalidGameProxy.resolve();
        vm.expectRevert(abi.encodeWithSelector(GameNotInProgress.selector));
        invalidGameProxy.resolveClaim();
    }

    function testResolve_CHALLENGER_WINS_ByProof() public {
        address faultProver = makeAddr("faultProver");
        // construct invalid games
        uint256 disputeClaimIndex = 0;
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000001));
        bytes memory extraData;
        vm.warp(parentCreateTimestamp);
        ZkFaultDisputeGame invalidGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame{ value: proposerBond }(
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));
        // Challenge the game
        // proof for block 47528971(agreed claim)-47528974(dispute claim)
        bytes memory proof = parentProofs[disputeClaimIndex];
        vm.deal(faultProver, 100 ether);
        vm.startPrank(faultProver);
        uint256 proverInitialBalance = faultProver.balance;
        uint256 contractBalance = delayedWeth.balanceOf(address(invalidGameProxy));
        assertEq(contractBalance, proposerBond);
        invalidGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, proof);
        invalidGameProxy.resolve();
        uint256 reward = contractBalance * (proverRewardPercentage + challengerRewardPercentage) / percentageDivisor;
        uint256 faultGameWithdrawalDelay = 604800;
        vm.warp(parentCreateTimestamp + faultGameWithdrawalDelay + 1);
        invalidGameProxy.claimCredit(faultProver);
        uint256 proverAfterBalance = faultProver.balance;
        assertEq(proverAfterBalance - proverInitialBalance, reward);
        assertEq(uint256(invalidGameProxy.status()), uint256(GameStatus.CHALLENGER_WINS));
        vm.stopPrank();
    }

    // The challenger and the prover are two different accounts
    // two challengers and one prover
    function testResolve_CHALLENGER_WINS_BySignalAndProof() public {
        address faultProver = makeAddr("faultProver");
        address challenger = makeAddr("challenger");
        address challenger2 = makeAddr("challenger2");
        vm.deal(challenger, 100 ether);
        vm.deal(challenger2, 100 ether);
        vm.deal(faultProver, 100 ether);
        uint256 disputeClaimIndex = 0;
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000002));
        bytes memory proof = parentProofs[disputeClaimIndex];
        vm.warp(parentCreateTimestamp);
        ZkFaultDisputeGame invalidGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame{ value: proposerBond }(
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));

        vm.startPrank(challenger);
        invalidGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);
        vm.stopPrank();

        vm.startPrank(challenger2);
        uint256 challenger2DisputeIndex = 1;
        invalidGameProxy.challengeBySignal{ value: challengerBond }(challenger2DisputeIndex);
        vm.stopPrank();

        vm.startPrank(faultProver);
        invalidGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, proof);
        vm.stopPrank();

        uint256 contractBalance = delayedWeth.balanceOf(address(invalidGameProxy));
        assertEq(contractBalance, proposerBond + challengerBond * 2);

        uint256 challengerInitialBalance = challenger.balance;
        uint256 challenger2InitialBalance = challenger2.balance;
        uint256 proverInitialBalance = faultProver.balance;
        invalidGameProxy.resolve();
        uint256 faultGameWithdrawalDelay = 604800;
        vm.warp(parentCreateTimestamp + faultGameWithdrawalDelay + 1);
        invalidGameProxy.claimCredit(faultProver);
        uint256 proverAfterBalance = faultProver.balance;
        assertEq(uint256(invalidGameProxy.status()), uint256(GameStatus.CHALLENGER_WINS));
        // First refund to challenger and challenger2
        contractBalance = contractBalance - challengerBond * 2;
        invalidGameProxy.claimCredit(challenger2);
        assertEq(challenger2.balance - challenger2InitialBalance, challengerBond);
        // Then calculate the challenger reward and the prover reward
        uint256 challengeReward = contractBalance * challengerRewardPercentage / percentageDivisor;
        invalidGameProxy.claimCredit(challenger);
        assertEq(challenger.balance - challengerInitialBalance, challengeReward + challengerBond);
        uint256 proverReward = contractBalance * proverRewardPercentage / percentageDivisor;
        assertEq(proverAfterBalance - proverInitialBalance, proverReward);
    }

    // Two games: invalid parent game and one valid child game
    // The valid child game has validity prover and two challengers
    // expect result:
    // a. the validity prover doesn't get any reward
    // b. the challenger who submit the challenge which is proved invalid by validity prover get slashed;
    // c. the challenger who submit the challenge which isn't proved invalid by validity prover get refund;
    function testResolve_CHALLENGER_WINS_ParentGameIsInvalid_S1() public {
        // first construct an invalid parent game and an invalid child game
        uint256 disputeClaimIndex = 0;
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000002));
        bytes memory proof = parentProofs[disputeClaimIndex];
        ZkFaultDisputeGame invalidGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame{ value: proposerBond }(
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));

        uint64 currentParentGameIndex = uint64(disputeGameFactory.gameCount() - 1);
        // construct a valid child game whose parent game is invalid
        ZkFaultDisputeGame childGameProxy2 = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame{ value: proposerBond }(
            GAME_TYPE, childClaims, currentParentGameIndex, childClaimBlockNumber, extraData)));

        // Challenge the parent game
        vm.startPrank(parentFaultProver);
        invalidGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, proof);
        vm.stopPrank();

        vm.startPrank(childChallenger);
        uint256 challengerInitialBalance = childChallenger.balance;
        childGameProxy2.challengeBySignal{ value: challengerBond }(disputeClaimIndex);
        childGameProxy2.challengeBySignal{ value: challengerBond }(disputeClaimIndex + 1);
        vm.stopPrank();

        vm.startPrank(childValidityProver);
        uint256 validityProverInitialBalance = childValidityProver.balance;
        proof = childProofs[disputeClaimIndex+1];
        childGameProxy2.submitProofForSignal(disputeClaimIndex + 1, childClaims, proof);
        vm.stopPrank();

        uint256 proverInitialBalance = parentFaultProver.balance;
        invalidGameProxy.resolve();
        assertEq(uint256(invalidGameProxy.status()), uint256(GameStatus.CHALLENGER_WINS));
        uint256 proverAfterBalance = parentFaultProver.balance;
        assertEq(proverAfterBalance - proverInitialBalance, proposerBond * (proverRewardPercentage + challengerRewardPercentage) / percentageDivisor);
        proverInitialBalance = proverAfterBalance;

        childGameProxy2.resolve();
        assertEq(uint256(childGameProxy2.status()), uint256(GameStatus.CHALLENGER_WINS));
        uint256 challengerAfterBalance = childChallenger.balance;
        uint256 validityProverAfterBalance = childValidityProver.balance;
        proverAfterBalance = parentFaultProver.balance;
        // the validity prover can't get the reward because of the parent game is invalid
        assertEq(validityProverAfterBalance, validityProverInitialBalance);
        // the challenger can only get one bond refund, because one of the two challenges is invalid
        assertEq(challengerInitialBalance - challengerAfterBalance, challengerBond);
        // the total reward is proposerBond + challengerBond
        assertEq(proverAfterBalance - proverInitialBalance, (proposerBond + challengerBond) * proverRewardPercentage / percentageDivisor);
    }

    // Two games: invalid parent game and one inValid child game
    // The invalid child game has two challengers, one child fault prover and one validity prover
    // expect result:
    // a. the child fault prover doesn't get any reward
    // b. the challenger who submit the challenger which is proved valid by fault prover get reward;
    // c. the child validity prover doesn't get any reward
    // d. the challenger who submit the challenge which is proved invalid by validity prover get slashed;
    function testResolve_CHALLENGER_WINS_ParentGameIsInvalid_S2() public {
        // first construct an invalid parent game and an invalid child game
        uint256 disputeClaimIndex = 0;
        Claim expectedClaim = parentClaims[disputeClaimIndex];
        parentClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000002));
        bytes memory proof = parentProofs[disputeClaimIndex];
        ZkFaultDisputeGame invalidGameProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame{ value: proposerBond }(
            GAME_TYPE, parentClaims, type(uint64).max, parentClaimBlockNumber, extraData)));

        Claim childExpectedClaim = childClaims[disputeClaimIndex];
        childClaims[disputeClaimIndex] = Claim.wrap(bytes32(0x0000000000000000000000000000000000000000000000000000000000000002));
        // construct a valid child game whose parent game is invalid
        ZkFaultDisputeGame invalidChildProxy = ZkFaultDisputeGame(address(disputeGameFactory.createZkFaultDisputeGame{ value: proposerBond }(
            GAME_TYPE, childClaims, uint64(disputeGameFactory.gameCount() - 1), childClaimBlockNumber, extraData)));

        // Challenge the parent game

        vm.startPrank(parentChallenger);
        uint256 parentChallengerInitialBalance = parentChallenger.balance;
        invalidGameProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);
        uint256 afterBalance = parentChallenger.balance;
        assertEq(parentChallengerInitialBalance - afterBalance, challengerBond);
        vm.stopPrank();

        vm.startPrank(parentFaultProver);
        uint256 parentFaultProverInitialBalance = parentFaultProver.balance;
        invalidGameProxy.challengeByProof(disputeClaimIndex, expectedClaim, parentClaims, proof);
        assertEq(invalidGameProxy.isChallengeSuccess(), true);
        vm.stopPrank();

        vm.startPrank(childChallenger);
        uint256 childChallengerInitialBalance = childChallenger.balance;
        invalidChildProxy.challengeBySignal{ value: challengerBond }(2);
        invalidChildProxy.challengeBySignal{ value: challengerBond }(disputeClaimIndex);
        vm.stopPrank();
        proof = childProofs[2];

        vm.startPrank(childValidityProver);
        uint256 childValidityProverInitialBalance = childValidityProver.balance;
        invalidChildProxy.submitProofForSignal(2, childClaims, proof);
        assertEq(invalidChildProxy.invalidChallengeClaims(2), true);
        vm.stopPrank();

        vm.startPrank(childFaultProver);
        uint256 childFaultProverInitialBalance = childFaultProver.balance;
        proof = childProofs[disputeClaimIndex];
        invalidChildProxy.challengeByProof(disputeClaimIndex, childExpectedClaim, childClaims, proof);
        assertEq(invalidChildProxy.isChallengeSuccess(), true);
        vm.stopPrank();

        invalidGameProxy.resolve();
        assertEq(uint256(invalidGameProxy.status()), uint256(GameStatus.CHALLENGER_WINS));
        afterBalance = parentFaultProver.balance;
        assertEq(afterBalance - parentFaultProverInitialBalance, proposerBond * proverRewardPercentage / percentageDivisor);
        parentFaultProverInitialBalance = afterBalance;

        afterBalance = parentChallenger.balance;
        assertEq(afterBalance - parentChallengerInitialBalance, proposerBond * challengerRewardPercentage / percentageDivisor);

        invalidChildProxy.resolve();
        assertEq(uint256(invalidChildProxy.status()), uint256(GameStatus.CHALLENGER_WINS));
        afterBalance = childFaultProver.balance;
        assertEq(afterBalance - childFaultProverInitialBalance, 0);

        afterBalance = childValidityProver.balance;
        assertEq(afterBalance - childValidityProverInitialBalance, 0);

        afterBalance = childChallenger.balance;
        assertEq(childChallengerInitialBalance - afterBalance, challengerBond);

        afterBalance = parentFaultProver.balance;
        assertEq(afterBalance - parentFaultProverInitialBalance, (proposerBond + challengerBond) * proverRewardPercentage / percentageDivisor);
    }

    //
    function testResolve_DEFENDER_WINS() public {
        vm.startPrank(parentChallenger);
        vm.warp(parentCreateTimestamp + 1);
        uint256 parentChallengerInitialBalance = parentChallenger.balance;
        for (uint256 i = 0; i < parentClaims.length; i++) {
            parentGameProxy.challengeBySignal{ value: challengerBond }(i);
        }
        vm.stopPrank();
        vm.startPrank(parentValidityProver);
        uint256 validityProverInitialBalance = parentValidityProver.balance;
        for (uint256 i = 0; i < parentClaims.length; i++) {
            parentGameProxy.submitProofForSignal(i, parentClaims, parentProofs[i]);
        }

        uint256 proposerInitialBalance = parentGameProxy.gameCreator().balance;

        assertEq(address(parentGameProxy).balance, proposerBond + challengerBond * parentClaims.length);
        vm.stopPrank();
        vm.expectRevert(abi.encodeWithSelector(ClockNotExpired.selector));
        parentGameProxy.resolve();
        vm.warp(parentCreateTimestamp + parentGameProxy.maxClockDuration().raw() + 1);
        parentGameProxy.resolve();
        assertEq(uint256(parentGameProxy.status()), uint256(GameStatus.DEFENDER_WINS));

        uint256 parentChallengerAfterBalance = parentChallenger.balance;
        assertEq(parentChallengerInitialBalance - parentChallengerAfterBalance, challengerBond * parentClaims.length);

        uint256 proposerAfterBalance = parentGameProxy.gameCreator().balance;
        assertEq(proposerAfterBalance - proposerInitialBalance, proposerBond);

        uint256 validityProverAfterBalance = parentValidityProver.balance;
        assertEq(validityProverAfterBalance - validityProverInitialBalance, challengerBond * parentClaims.length * proverRewardPercentage / percentageDivisor);

        assertEq(address(parentGameProxy).balance, 0);
    }

}
