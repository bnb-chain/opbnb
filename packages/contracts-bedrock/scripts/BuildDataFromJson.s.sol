// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {IGnosisSafe, Enum} from "scripts/interfaces/IGnosisSafe.sol";
import {IMulticall3} from "forge-std/interfaces/IMulticall3.sol";
import {stdJson} from "forge-std/StdJson.sol";
import {console} from "forge-std/console.sol";
import {Vm, VmSafe} from "forge-std/Vm.sol";
import {CommonBase} from "forge-std/Base.sol";

contract BuildDataFromJson is CommonBase {

    string json;

    IMulticall3 internal constant multicall = IMulticall3(MULTICALL3_ADDRESS);

    function _loadJson(string memory _path) internal {
        console.log("Reading transaction bundle %s", _path);
        json = vm.readFile(_path);
    }

    function _buildCalls() internal view returns (IMulticall3.Call3[] memory) {
        return _buildCallsFromJson(json);
    }

    function _buildCallsFromJson(string memory jsonContent) internal pure returns (IMulticall3.Call3[] memory) {
        // A hacky way to get the total number of elements in a JSON
        // object array because Forge does not support this natively.
        uint256 MAX_LENGTH_SUPPORTED = 999;
        uint256 transaction_count = MAX_LENGTH_SUPPORTED;
        for (uint256 i = 0; transaction_count == MAX_LENGTH_SUPPORTED; i++) {
            require(
                i < MAX_LENGTH_SUPPORTED,
                "Transaction list longer than MAX_LENGTH_SUPPORTED is not "
                "supported, to support it, simply bump the value of " "MAX_LENGTH_SUPPORTED to a bigger one."
            );
            try vm.parseJsonAddress(jsonContent, string(abi.encodePacked("$.transactions[", vm.toString(i), "].to")))
            returns (address) {} catch {
                transaction_count = i;
            }
        }

        IMulticall3.Call3[] memory calls = new IMulticall3.Call3[](transaction_count);

        for (uint256 i = 0; i < transaction_count; i++) {
            calls[i] = IMulticall3.Call3({
                target: stdJson.readAddress(
                    jsonContent, string(abi.encodePacked("$.transactions[", vm.toString(i), "].to"))
                ),
                allowFailure: false,
                callData: stdJson.readBytes(
                    jsonContent, string(abi.encodePacked("$.transactions[", vm.toString(i), "].data"))
                )
            });
        }

        return calls;
    }

    function rawInputData(string memory _path) public {
        _loadJson(_path);
        IMulticall3.Call3[] memory calls = _buildCalls();
        bytes memory data = abi.encodeCall(IMulticall3.aggregate3, (calls));
        bytes memory txData = abi.encodeCall(IGnosisSafe.execTransaction,
            (
                address(multicall),
                0,
                data,
                Enum.Operation.DelegateCall,
                0,
                0,
                0,
                address(0),
                payable(address(0)),
                prevalidatedSignature(msg.sender)
            )
        );
        console.log("raw input data:");
        console.log(vm.toString(txData));
    }

    function prevalidatedSignature(address _address) internal pure returns (bytes memory) {
        uint8 v = 1;
        bytes32 s = bytes32(0);
        bytes32 r = bytes32(uint256(uint160(_address)));
        return abi.encodePacked(r, s, v);
    }
}
