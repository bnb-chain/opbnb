package bsc

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum-optimism/optimism/op-node/eth"
)

var DefaultBaseFee = big.NewInt(3000000000)
var DefaultOPBNBTestnetBaseFee = big.NewInt(5000000000)
var OPBNBTestnet = big.NewInt(5611)

type BlockInfoBSCWrapper struct {
	eth.BlockInfo
	baseFee *big.Int
}

var _ eth.BlockInfo = (*BlockInfoBSCWrapper)(nil)

func NewBlockInfoBSCWrapper(info eth.BlockInfo, baseFee *big.Int) *BlockInfoBSCWrapper {
	return &BlockInfoBSCWrapper{
		BlockInfo: info,
		baseFee:   baseFee,
	}
}

func (b *BlockInfoBSCWrapper) BaseFee() *big.Int {
	return b.baseFee
}

// BaseFeeByTransactions calculates the average gas price of the non-zero-gas-price transactions in the block.
// If there is no non-zero-gas-price transaction in the block, it returns DefaultBaseFee.
func BaseFeeByTransactions(transactions types.Transactions) *big.Int {
	nonZeroTxsCnt := big.NewInt(0)
	nonZeroTxsSum := big.NewInt(0)
	for _, tx := range transactions {
		if tx.GasPrice().Cmp(common.Big0) > 0 {
			nonZeroTxsCnt.Add(nonZeroTxsCnt, big.NewInt(1))
			nonZeroTxsSum.Add(nonZeroTxsSum, tx.GasPrice())
		}
	}

	if nonZeroTxsCnt.Cmp(big.NewInt(0)) == 0 {
		return DefaultBaseFee
	}
	return nonZeroTxsSum.Div(nonZeroTxsSum, nonZeroTxsCnt)
}

// BaseFeeByNetworks set l1 gas price by network.
func BaseFeeByNetworks(chainId *big.Int) *big.Int {
	if chainId.Cmp(OPBNBTestnet) == 0 {
		return DefaultOPBNBTestnetBaseFee
	} else {
		return DefaultBaseFee
	}
}

func ToLegacyTx(dynTx *types.DynamicFeeTx) types.TxData {
	return &types.LegacyTx{
		Nonce:    dynTx.Nonce,
		GasPrice: dynTx.GasFeeCap,
		Gas:      dynTx.Gas,
		To:       dynTx.To,
		Value:    dynTx.Value,
		Data:     dynTx.Data,
	}
}

func ToLegacyCallMsg(callMsg ethereum.CallMsg) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:     callMsg.From,
		To:       callMsg.To,
		GasPrice: callMsg.GasFeeCap,
		Data:     callMsg.Data,
	}
}
